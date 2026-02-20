package repo

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type conf map[string]map[string]string

// Gitrepo struct represents a git repository
// Worktree is the path to the working directory
// Gitdir is the path to the .git directory
// Conf is the configuration of the repository
type Gitrepo struct {
	Worktree string
	Gitdir   string
	Conf     conf
}

// writeStringFile writes the content to the file at path
func writeStringFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// writeConfigFile writes the config to the file at path
func writeConfigFile(path string, c conf) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for section, kv := range c {
		if _, err := fmt.Fprintf(f, "[%s]\n", section); err != nil {
			return err
		}
		for key, value := range kv {
			if _, err := fmt.Fprintf(f, "\t%s = %s\n", key, value); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(f); err != nil {
			return err
		}
	}
	return nil
}

// isDirEmpty checks if the directory at path is empty
func isDirEmpty(path string) (bool, error) {
	isPath, isDir := pathExist(path)

	if !isPath && !isDir {
		return true, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil // directory is empty
	}
	if err != nil {
		return false, err // some other error
	}

	return false, nil // directory has at least one entry
}

// parseConfig parses the config file data and returns a conf map
func parseConfig(data []byte) (conf, error) {
	c := conf{} // initialize your map
	section := ""
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		//[] #
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSpace(line[1 : len(line)-1])
			c[section] = make(map[string]string)
			continue
		}
		pair := strings.SplitN(line, "=", 2)
		if len(pair) == 2 {
			key := strings.TrimSpace(pair[0])
			value := strings.TrimSpace(pair[1])
			c[section][key] = value
		}

	}

	return c, nil
}

// NewGitrepo creates a new Gitrepo struct
func NewGitrepo(path string, force bool) (*Gitrepo, error) {
	repo := &Gitrepo{}
	repo.Worktree = path
	repo.Gitdir = filepath.Join(repo.Worktree, ".tit")
	cf, err := RepoFile(repo, false, "config")
	if err != nil {
		return nil, err
	}
	isPath, _ := pathExist(cf)
	if isPath {
		data, err := os.ReadFile(cf)
		if err != nil {
			return nil, err
		}
		conf, err := parseConfig(data)
		if err != nil {
			return nil, err
		}
		repo.Conf = conf
	} else if !force {
		return nil, errors.New("No config file in this repo")
	}
	if !force {
		ver, exist := repo.Conf["core"]["repoformatversion"]
		if !exist || ver != "0" {
			return nil, errors.New("Unsupported version")
		}
	}
	return repo, nil
}

// pathExist checks if the path exists and if it is a directory
func pathExist(path string) (exists bool, isDir bool) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, false
		}
		return false, false // or return err
	}
	return true, info.IsDir()
}

// RepoPath returns the path to the file in the .git directory calls recursive
func RepoPath(repo *Gitrepo, paths ...string) string {
	all := append([]string{repo.Gitdir}, paths...)
	return filepath.Join(all...)
}

// RepoDir returns the path to the directory in the .git directory
func RepoDir(repo *Gitrepo, mkdir bool, paths ...string) (string, error) {

	path := RepoPath(repo, paths...)

	isPath, isDir := pathExist(path)

	if isPath {

		if isDir {
			return path, nil
		} else {
			return "", errors.New("path exists but is not a directory")
		}
	}

	if mkdir {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return "", err
		}
		return path, nil
	}

	return "", nil
}

// RepoFile returns the path to the file in the .git directory
func RepoFile(repo *Gitrepo, mkdir bool, paths ...string) (string, error) {
	if len(paths) == 0 {
		return "", errors.New("no paths provided")
	}
	_, err := RepoDir(repo, mkdir, paths[:len(paths)-1]...)
	if err != nil {
		return "", err
	}
	return RepoPath(repo, paths...), nil
}

// RepoCreate creates a new git repository
func RepoCreate(path string) (*Gitrepo, error) {
	repo, err := NewGitrepo(path, true)
	if err != nil {
		return nil, err
	}
	isWorkTree, isWorkDir := pathExist(repo.Worktree)

	if isWorkTree {
		if !isWorkDir {
			return nil, errors.New("There is no directory")
		}

		isGit, _ := isDirEmpty(repo.Gitdir)

		if !isGit {
			return nil, errors.New("Git dir is not empty")

		}
	} else {
		err = os.Mkdir(path, 0755)
		if err != nil {
			return nil, err
		}
	}
	dirs := [][]string{
		{"branches"},
		{"objects"},
		{"refs", "tags"},
		{"refs", "heads"},
	}

	for _, d := range dirs {
		if _, err := RepoDir(repo, true, d...); err != nil {
			return nil, err
		}
	}
	if err := writeStringFile(RepoPath(repo, "description"),
		"Unnamed Repo; Change this description file to make a repository.\n"); err != nil {
		return nil, err
	}

	// HEAD
	if err := writeStringFile(RepoPath(repo, "HEAD"),
		"ref: refs/heads/master\n"); err != nil {
		return nil, err
	}

	// config
	if err := writeConfigFile(RepoPath(repo, "config"), getDefaultConfig()); err != nil {
		return nil, err
	}

	return repo, nil
}

// getDefaultConfig returns the default configuration for a new repository
func getDefaultConfig() conf {
	Conf := conf{}
	addSection(Conf, "core")
	setSection(Conf, "core", "repoformatversion", "0")
	setSection(Conf, "core", "bare", "false")
	return Conf
}

// addSection adds a new section to the config map
func addSection(c conf, section string) {
	if _, ok := c[section]; !ok {
		c[section] = make(map[string]string)
	}
}

// setSection sets the value for a key in a section
func setSection(c conf, section, key, value string) {
	addSection(c, section)
	c[section][key] = value
}

// RepoFind finds the root of the git repository
func RepoFind(path string, req bool) (*Gitrepo, error) {
	path, err := filepath.Abs(path)


	if err != nil {
		return nil, err
	}


	_, isDir := pathExist(filepath.Join(path, ".tit"))

	if isDir {
		return NewGitrepo(path, false)
	}
	parentPath := filepath.Dir(path)
	if parentPath == path {
		if req {
			return nil, errors.New("No tit found")

		} else {
			return nil, nil
		}
	}

	return RepoFind(parentPath, req)
}

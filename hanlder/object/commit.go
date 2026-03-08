package object

import "bytes"



type Commit struct {
	Data CommitData
	Fmt  []byte
}

type CommitData struct{
	Header map[string][]string
	Message []byte
}
func (c *Commit) Type() string {
	return "commit"
}
func (c *Commit) Deserialize(raw []byte)error{
	kvlm:=map[string][]string{}
	header,message,err:=kvlmParse(raw,0,kvlm)
	if err != nil {
		return err
	}
	c.Data.Header=header
	c.Data.Message=message
	c.Fmt=[]byte("commit")
	return nil
}
func (c *Commit) Serialize() ([]byte,error) {
	var buf bytes.Buffer

	// Write headers
	for key, values := range c.Data.Header {
		for _, value := range values {

			
			lines := bytes.Split([]byte(value), []byte("\n"))


			buf.WriteString(key)
			buf.WriteByte(' ')
			buf.Write(lines[0])
			buf.WriteByte('\n')


			for _, line := range lines[1:] {
				buf.WriteByte(' ')
				buf.Write(line)
				buf.WriteByte('\n')
			}
		}
	}


	buf.WriteByte('\n')


	buf.Write(c.Data.Message)

	return buf.Bytes(),nil
}
func kvlmParse(raw []byte, start int, kvlm map[string][]string) (map[string][]string, []byte, error) {
	spc := bytes.Index(raw[start:], []byte(" "))
	nl := bytes.Index(raw[start:], []byte("\n"))
	if spc == -1 || nl == -1 || spc > nl {
		return kvlm, raw[start+1:], nil
	}
	spc += start
	nl += start

	key := raw[start:spc]
	end := nl

	for end+1 < len(raw) && raw[end+1] == ' ' {
		nextNl := bytes.Index(raw[end+1:], []byte("\n"))
		if nextNl == -1 {
			end = len(raw)
			break
		}
		end += nextNl + 1
	}
	value := bytes.ReplaceAll(raw[spc+1:end], []byte("\n "), []byte("\n"))
	v, ok := kvlm[string(key)]
	if ok {
		kvlm[string(key)] = append(v, string(value))
	} else {
		kvlm[string(key)] = []string{string(value)}
	}

	return kvlmParse(raw, end+1, kvlm)
}

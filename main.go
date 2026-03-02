package main

import (
	"bytes"
	"fmt"

	"github.com/Blue-Onion/pygo/hanlder/object"
)

func main() {
	fmt.Println("========== COMMIT OBJECT TEST ==========")
	fmt.Println()

	// --- Test 1: Deserialize a raw commit ---
	fmt.Println("--- Test 1: Deserialize ---")

	raw := []byte("tree abc123\nparent def456\nparent fedcba\nauthor Blue Onion\n <blue@onion.com>\ncommitter Blue Onion <blue@onion.com>\n\nThis is the commit message\nWith multiple lines\nAnd even more lines\n")

	commit := &object.Commit{}
	err := commit.Deserialize(raw)
	if err != nil {
		fmt.Println("  [FAIL] Deserialize error:", err)
		return
	}
	fmt.Println("  [PASS] Deserialize succeeded")

	// Check Type()
	if commit.Type() != "commit" {
		fmt.Printf("  [FAIL] Type() = %q, want \"commit\"\n", commit.Type())
	} else {
		fmt.Println("  [PASS] Type() = \"commit\"")
	}

	// Print parsed headers
	fmt.Println()
	fmt.Println("  Parsed Headers:")
	for key, values := range commit.Data.Header {
		for _, v := range values {
			fmt.Printf("    %s: %s\n", key, v)
		}
	}

	// Print parsed message
	fmt.Println()
	fmt.Println("  Parsed Message:")
	fmt.Printf("    %q\n", string(commit.Data.Message))

	// Check specific header values
	fmt.Println()
	if tree, ok := commit.Data.Header["tree"]; ok && len(tree) > 0 && tree[0] == "abc123" {
		fmt.Println("  [PASS] Header 'tree' = \"abc123\"")
	} else {
		fmt.Printf("  [FAIL] Header 'tree' = %v, want [\"abc123\"]\n", commit.Data.Header["tree"])
	}

	if parents, ok := commit.Data.Header["parent"]; ok && len(parents) == 2 {
		fmt.Println("  [PASS] Header 'parent' has 2 values")
		if parents[0] == "def456" {
			fmt.Println("  [PASS] parent[0] = \"def456\"")
		} else {
			fmt.Printf("  [FAIL] parent[0] = %q, want \"def456\"\n", parents[0])
		}
		if parents[1] == "fedcba" {
			fmt.Println("  [PASS] parent[1] = \"fedcba\"")
		} else {
			fmt.Printf("  [FAIL] parent[1] = %q, want \"fedcba\"\n", parents[1])
		}
	} else {
		fmt.Printf("  [FAIL] Header 'parent' = %v, want 2 values\n", commit.Data.Header["parent"])
	}

	// --- Test 2: Serialize back ---
	fmt.Println()
	fmt.Println("--- Test 2: Serialize ---")

	serialized, err := commit.Serialize()
	if err != nil {
		fmt.Println("  [FAIL] Serialize error:", err)
		return
	}
	fmt.Println("  [PASS] Serialize succeeded")
	fmt.Println()
	fmt.Println("  Serialized output:")
	fmt.Println("  ---")
	fmt.Print("  " + string(bytes.ReplaceAll(serialized, []byte("\n"), []byte("\n  "))))
	fmt.Println()
	fmt.Println("  ---")

	// --- Test 3: Round-trip (Deserialize the serialized output) ---
	fmt.Println()
	fmt.Println("--- Test 3: Round-trip (Deserialize → Serialize → Deserialize) ---")

	commit2 := &object.Commit{}
	err = commit2.Deserialize(serialized)
	if err != nil {
		fmt.Println("  [FAIL] Round-trip Deserialize error:", err)
		return
	}
	fmt.Println("  [PASS] Round-trip Deserialize succeeded")

	// Compare headers
	match := true
	for key, vals := range commit.Data.Header {
		vals2, ok := commit2.Data.Header[key]
		if !ok || len(vals) != len(vals2) {
			match = false
			fmt.Printf("  [FAIL] Header mismatch for key %q\n", key)
			break
		}
		for i := range vals {
			if vals[i] != vals2[i] {
				match = false
				fmt.Printf("  [FAIL] Header value mismatch for key %q at index %d: %q vs %q\n", key, i, vals[i], vals2[i])
				break
			}
		}
	}
	if match {
		fmt.Println("  [PASS] Headers match after round-trip")
	}

	// Compare messages
	if bytes.Equal(commit.Data.Message, commit2.Data.Message) {
		fmt.Println("  [PASS] Message matches after round-trip")
	} else {
		fmt.Printf("  [FAIL] Message mismatch:\n    original: %q\n    roundtrip: %q\n",
			string(commit.Data.Message), string(commit2.Data.Message))
	}

	// --- Test 4: Multi-line header value (author with continuation line) ---
	fmt.Println()
	fmt.Println("--- Test 4: Multi-line header (author) ---")
	if author, ok := commit.Data.Header["author"]; ok && len(author) > 0 {
		fmt.Printf("  [PASS] Author parsed: %q\n", author[0])
		// The continuation line " <blue@onion.com>" should be joined
		if bytes.Contains([]byte(author[0]), []byte("<blue@onion.com>")) {
			fmt.Println("  [PASS] Continuation line merged into author value")
		} else {
			fmt.Println("  [FAIL] Continuation line NOT found in author value")
		}
	} else {
		fmt.Println("  [FAIL] No 'author' header found")
	}

	fmt.Println()
	fmt.Println("========== ALL TESTS DONE ==========")
}

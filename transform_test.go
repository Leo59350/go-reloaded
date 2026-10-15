package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTransform(t *testing.T) {
	cases := []struct{ name, input, want string }{
		{"hex", "1E (hex) files were added", "30 files were added"},
		{"bin", "It has been 10 (bin) years", "It has been 2 years"},
		{"up", "Ready, set, go (up) !", "Ready, set, GO!"},
		{"low", "I should stop SHOUTING (low)", "I should stop shouting"},
		{"cap", "Welcome to the Brooklyn bRIDGE (cap)", "Welcome to the Brooklyn Bridge"},
		{"multiple", "This is so exciting (up, 2)", "This is SO EXCITING"},
		{"count beyond start", "one two (cap, 6)", "One Two"},
		{"punctuation", "I was sitting over there ,and then BAMM !!", "I was sitting over there, and then BAMM!!"},
		{"groups", "thinking . . . right ! ?", "thinking... right!?"},
		{"quotes", "describe me: ' awesome '", "describe me: 'awesome'"},
		{"phrase quotes", "He said: ' I am the most famous person ' .", "He said: 'I am the most famous person'."},
		{"several quotes", "' hello ' and ' world '", "'hello' and 'world'"},
		{"contractions", "don't stop John's work", "don't stop John's work"},
		{"articles", "A Amazing rock and a house and a book", "An Amazing rock and an house and a book"},
		{"quoted article", "a ' apple '", "an 'apple'"},
		{"combined", "hello (cap) ! I have 10 (bin) apples and a orange .", "Hello! I have 2 apples and an orange."},
		{"commands in sequence", "HELLO (low) (cap),world (up)!", "Hello, WORLD!"},
		{"punctuation before command", "hello, world! (up, 2)", "HELLO, WORLD!"},
		{"newlines", "hello (up) \n\n a apple .\n", "HELLO\n\nan apple.\n"},
		{"windows newlines", "one\r\ntwo", "one\r\ntwo"},
		{"unicode", "éCOLE (cap) été (up)", "École ÉTÉ"},
		{"empty", "", ""},
		{"whitespace", "  \t ", ""},
		{"numeric mix", "Simply add 42 (hex) and 10 (bin) and you will see the result is 68.", "Simply add 66 and 2 and you will see the result is 68."},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Transform(tc.input)
			if err != nil || got != tc.want {
				t.Fatalf("Transform(%q) = %q, %v ; attendu %q", tc.input, got, err, tc.want)
			}
		})
	}
}

func TestInvalidCommands(t *testing.T) {
	for _, input := range []string{"(up)", "hello (up, 0)", "hello (cap, -2)", "hello (low, abc)", "XYZ (hex)", "102 (bin)", "FFFFFFFFFFFFFFFFFFFF (hex)", "10 (hex, 2)"} {
		if _, err := Transform(input); err == nil {
			t.Errorf("erreur attendue pour %q", input)
		}
	}
}

func TestRun(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "input.txt")
	out := filepath.Join(dir, "output.txt")
	if err := os.WriteFile(in, []byte("hello (up) !"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{in, out}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(out)
	if err != nil || string(data) != "HELLO!" {
		t.Fatalf("sortie = %q, erreur = %v", data, err)
	}
	for _, args := range [][]string{nil, {in}, {in, out, "extra"}, {filepath.Join(dir, "missing"), out}, {in, filepath.Join(dir, "missing", "out")}} {
		if err := run(args); err == nil {
			t.Errorf("erreur attendue pour %v", args)
		}
	}
	if err := os.WriteFile(in, []byte("bad (hex) 102 (bin)"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{in, out}); err == nil {
		t.Fatal("erreur de conversion attendue")
	}
	data, _ = os.ReadFile(out)
	if string(data) != "HELLO!" {
		t.Fatal("une erreur ne doit pas écraser la sortie")
	}
}

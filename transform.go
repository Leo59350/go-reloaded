package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var tokenPattern = regexp.MustCompile(`\((?:hex|bin|up|low|cap)(?:\s*,\s*[^)]*)?\)|\r\n|\n|\r|[\p{L}\p{N}_]+(?:['’\-][\p{L}\p{N}_]+)*|[.,!?:;]+|[^\s]`)
var commandPattern = regexp.MustCompile(`^\((hex|bin|up|low|cap)(?:\s*,\s*([^)]*))?\)$`)

type token struct {
	text string
	kind string
}

func Transform(input string) (string, error) {
	var tokens []token
	for _, text := range tokenPattern.FindAllString(input, -1) {
		if match := commandPattern.FindStringSubmatch(text); match != nil {
			if err := applyCommand(tokens, match[1], match[2]); err != nil {
				return "", fmt.Errorf("%s : %w", text, err)
			}
			continue
		}
		kind := "word"
		switch {
		case text == "\n" || text == "\r" || text == "\r\n":
			kind = "newline"
		case text == "'":
			kind = "quote"
		case strings.ContainsRune(".,!?:;", []rune(text)[0]):
			kind = "punct"
		}
		tokens = append(tokens, token{text, kind})
	}
	fixArticles(tokens)
	return render(tokens), nil
}

func applyCommand(tokens []token, name, countText string) error {
	count := 1
	if countText != "" {
		if name == "hex" || name == "bin" {
			return fmt.Errorf("cette conversion ne prend pas de nombre de mots")
		}
		var err error
		count, err = strconv.Atoi(strings.TrimSpace(countText))
		if err != nil || count < 1 {
			return fmt.Errorf("nombre de mots invalide")
		}
	}
	changed := 0
	for i := len(tokens) - 1; i >= 0 && changed < count; i-- {
		if tokens[i].kind != "word" {
			continue
		}
		word := tokens[i].text
		switch name {
		case "up":
			word = strings.ToUpper(word)
		case "low":
			word = strings.ToLower(word)
		case "cap":
			letters := []rune(strings.ToLower(word))
			letters[0] = unicode.ToUpper(letters[0])
			word = string(letters)
		case "hex", "bin":
			base := 16
			if name == "bin" {
				base = 2
			}
			number, err := strconv.ParseInt(word, base, 64)
			if err != nil {
				return fmt.Errorf("nombre %q invalide ou trop grand : %w", word, err)
			}
			word = strconv.FormatInt(number, 10)
		}
		tokens[i].text = word
		changed++
	}
	if changed == 0 {
		return fmt.Errorf("aucun mot précédent")
	}
	// Si N dépasse le nombre de mots disponibles, tous sont transformés.
	return nil
}

func fixArticles(tokens []token) {
	for i := range tokens {
		if tokens[i].text != "a" && tokens[i].text != "A" {
			continue
		}
		j := i + 1
		for j < len(tokens) && (tokens[j].kind == "quote" || tokens[j].kind == "newline") {
			j++
		}
		if j < len(tokens) && tokens[j].kind == "word" {
			first := unicode.ToLower([]rune(tokens[j].text)[0])
			if strings.ContainsRune("aeiouh", first) {
				tokens[i].text += "n"
			}
		}
	}
}

func render(tokens []token) string {
	var output strings.Builder
	previous := ""
	insideQuote := false
	for _, t := range tokens {
		if t.kind == "newline" {
			output.WriteString(t.text)
			previous = "newline"
			continue
		}
		kind := t.kind
		if kind == "quote" {
			if insideQuote {
				kind = "closeQuote"
			} else {
				kind = "openQuote"
			}
			insideQuote = !insideQuote
		}
		space := output.Len() > 0 && previous != "newline" && previous != "openQuote" && kind != "punct" && kind != "closeQuote"
		if space {
			output.WriteByte(' ')
		}
		output.WriteString(t.text)
		previous = kind
	}
	return output.String()
}

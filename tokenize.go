package tim2

import (
	_ "fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/thanhpk/ascii"
)

type MatchLiteral struct {
	Str  string
	Psrc []int
}

var replacer = strings.NewReplacer("/", " ", "\"", " ", "/", " ", "_", " ", "'", " ", "{", " ", "}", " ",
	"(", " ", ")", " ", "[", " ", "]", " ", "&", " ", "?", " ", "!", " ", "=", " ", ">", " ", "<", " ")

// var replacerLiteral = strings.NewReplacer("-", " ", "_", " ")

func splitSentence(r rune) bool {
	if r >= '0' && r <= '9' {
		return false
	}
	return r < 'A' || r > 'z'
}

func Tokenize(str string) []string {
	str = strings.TrimSpace(strings.ToLower(str))
	tokenM := make(map[string]bool)

	// phones and emails
	var tokens []string
	emails := Email_regexp.FindAllString(str, -1)
	for _, email := range emails {
		tokens = append(tokens, strings.ToLower(email))
	}
	tokens = append(tokens, findPersonalPhoneNumber(str)...)
	for _, t := range tokens {
		tokenM[t] = true
	}

	str = ascii.Convert(str)

	tokens = tokenizeFilename(str)
	for _, t := range tokens {
		tokenM[t] = true
	}

	// remove space and weird characters
	str = replacer.Replace(str)
	str = strings.Join(strings.Fields(str), " ")

	strs := strings.FieldsFunc(str, splitSentence)

	lines := [][]string{}
	line := []string{}
	for _, str := range strs {
		if len(str) == 0 {
			continue
		}

		line = append(line, str)
		if len(str) < 2 || utf8.RuneCountInString(str) < 2 {
			lines = append(lines, line)
			line = []string{}
			continue
		}
	}
	if len(line) > 0 {
		lines = append(lines, line)
	}
	for _, line := range lines {
		for i, word := range line {
			if len(word) > Token_max_len {
				continue
			}

			if !stopWordM[word] && len(word) > Token_min_len {
				tokenM[word] = true
			}

			// could have bi-word
			if i < len(line)-1 {
				if len(word) > 9 || len(line[i+1]) > 9 /* we dont want be-word to long */ {
					continue
				}

				tokenM[word+"-"+line[i+1]] = true
			}
		}
	}

	tokens = []string{}
	for k := range tokenM {
		tokens = append(tokens, k)
	}
	return tokens
}

const Email_regex = `([a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+)`

var Email_regexp = regexp.MustCompile(Email_regex)

// +84 2473.021.368
const PersonalPhoneNumber_digit = "0123456789"
const PersonalPhoneNumber_format = "+() .-"
const PersonalPhoneNumber_min_len = 7

var PersonalPhoneNumber_regexp = regexp.MustCompile(RegexPhone)
var PersonalPhoneNumber_norm_map map[rune]rune

func findPersonalPhoneNumber(str string) []string {
	if len(str) < PersonalPhoneNumber_min_len {
		return nil
	}
	str += ","
	phoneNumbers := make([]string, 0)
	rarr := []rune(str)
	from, to := 0, 0
	for i, r := range rarr {
		if _, has := PersonalPhoneNumber_norm_map[r]; has {
			to++
			continue
		}

		token := string(rarr[from:to])
		if len(token) >= PersonalPhoneNumber_min_len {
			normtoken := make([]rune, 0)
			for _, tr := range token {
				if '0' <= tr && tr <= '9' {
					normtoken = append(normtoken, tr)
				}
			}
			if len(normtoken) >= PersonalPhoneNumber_min_len {
				phoneNumbers = append(phoneNumbers, string(normtoken))
			}
		}

		from, to = i+1, i+1
	}
	return phoneNumbers
}

func tokenizeFilename(str string) []string {
	strs := strings.FieldsFunc(str, func(r rune) bool {
		return r == ' ' || r == '\t'
	})
	out := []string{}
	for _, str := range strs {
		if !strings.Contains(str, ".") {
			continue
		}

		if strings.HasSuffix(str, ".") {
			continue
		}

		if len(str) > 51 {
			str = str[:50]
		}

		out = append(out, str)
	}

	return out
}

const Token_min_len = 2
const Token_max_len = 45

var Norm_map map[rune]rune

const RegexPhone = `([0-9._-]{3,})`

var Regexp_phone = regexp.MustCompile(RegexPhone)

// see http://www.clc.hcmus.edu.vn/?page_id=1507
var stopWordM = map[string]bool{}

func initPhoneKit() {
	PersonalPhoneNumber_norm_map = make(map[rune]rune)
	for _, r := range PersonalPhoneNumber_digit {
		PersonalPhoneNumber_norm_map[r] = r
	}
	for _, r := range PersonalPhoneNumber_format {
		PersonalPhoneNumber_norm_map[r] = r
	}
}

func init() {
	initPhoneKit()
	for _, stopword := range stopwords {
		stopword = strings.TrimSpace(stopword)
		if len(stopword) > 0 {
			stopWordM[stopword] = true
		}
	}

}

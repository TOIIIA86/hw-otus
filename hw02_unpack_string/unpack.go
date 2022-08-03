//go:build ignore
// +build ignore

package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}
	if !utf8.ValidString(str) {
		return "", ErrInvalidString
	}
	var res []rune
	prevDigit := false
	backslash := false
	for i, r := range str {
		switch {
		case unicode.IsLetter(r) || unicode.IsSpace(r):
			if backslash {
				return "", ErrInvalidString
			}
			prevDigit = false
			res = append(res, r)

		case unicode.IsDigit(r):
			if i == 0 || prevDigit {
				return "", ErrInvalidString
			}
			if backslash {
				backslash = false
				res = append(res, r)
				continue
			}
			prevDigit = true
			q, err := strconv.Atoi(string(r))
			if err != nil {
				return "", ErrInvalidString
			}
			res = cloneRune(res, q)

		case r == 92:
			prevDigit = false
			if backslash {
				backslash = false
				res = append(res, r)
			} else {
				backslash = true
			}

		default:
			return "", ErrInvalidString
		}
	}
	return string(res), nil
}

func cloneRune(r []rune, q int) []rune {
	if q == 0 {
		r = r[:len(r)-1]
	} else {
		q--
		lastChar := r[len(r)-1]
		for j := 0; j < q; j++ {
			r = append(r, lastChar)
		}
	}
	return r
}

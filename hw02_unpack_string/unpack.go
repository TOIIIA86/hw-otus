package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}
	converted, err := convertToSlice(str)
	if err != nil {
		return "", err
	}
	ret := &strings.Builder{}
	var lastSym string
	for i := 0; i < len(converted); i++ {
		sym := converted[i]
		var prevSym, nextSym rune
		if i != 0 {
			prevSym = converted[i-1]
		}
		if i+1 < len(converted) {
			nextSym = converted[i+1]
		}
		if unicode.IsDigit(sym) {
			if unicode.IsDigit(prevSym) && !unicode.IsDigit(rune(lastSym[0])) {
				return "", ErrInvalidString
			}
			num, _ := strconv.Atoi(string(sym))
			if num == 0 {
				reduceBuilder(ret, prevSym)
			} else {
				ret.WriteString(strings.Repeat(lastSym, num-1))
			}
		}
		if !unicode.IsDigit(sym) {
			if sym == '\\' {
				lastSym = slashProcessing(nextSym)
				ret.WriteString(lastSym)
				if i+1 >= len(converted) {
					return "", ErrInvalidString
				}
				i++
			}
			if sym != '\\' {
				ret.WriteRune(sym)
				lastSym = string(sym)
			}
		}
	}
	return ret.String(), nil
}

func slashProcessing(nextSym rune) string {
	lastSym := ""
	if unicode.IsDigit(nextSym) {
		lastSym = string(nextSym)
	}
	if nextSym == '\\' {
		lastSym = string(nextSym)
	}
	if nextSym == 'n' {
		lastSym = `\n`
	}
	return lastSym
}

func reduceBuilder(builder *strings.Builder, prevSym rune) {
	temp := builder.String()
	temp = temp[:len(temp)-utf8.RuneLen(prevSym)]
	builder.Reset()
	builder.WriteString(temp)
}

func convertToSlice(str string) ([]rune, error) {
	converted := []rune{}
	for p, sym := range str {
		if p == 0 && unicode.IsDigit(sym) {
			return nil, ErrInvalidString
		}
		converted = append(converted, sym)
	}
	return converted, nil
}

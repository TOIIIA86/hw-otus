package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type kv struct {
	key string
	val int
}

const maxResultLen = 10

var regexpTextSplit = regexp.MustCompile(`\s+`)

func Top10(in string) []string {
	if len(in) == 0 {
		return []string{}
	}
	dict := map[string]int{}

	words := regexpTextSplit.Split(in, -1)
	for _, word := range words {
		dict[word]++
	}

	if len(dict) == 1 {
		result := make([]string, 0, 1)
		for key := range dict {
			result = append(result, key)
		}
		return result
	}

	dictSlice := make([]kv, 0, len(dict))
	for k, v := range dict {
		dictSlice = append(dictSlice, kv{k, v})
	}

	sort.SliceStable(dictSlice, func(i, j int) bool {
		if dictSlice[i].val > dictSlice[j].val {
			return true
		} else if dictSlice[i].val < dictSlice[j].val {
			return false
		}
		switch strings.Compare(dictSlice[i].key, dictSlice[j].key) {
		case -1:
			return true
		case 1:
			return false
		default:
			return false
		}
	})

	resultLen := maxResultLen
	if len(dictSlice) < maxResultLen {
		resultLen = len(dictSlice)
	}
	result := make([]string, 0, resultLen)
	for _, val := range dictSlice {
		result = append(result, val.key)
		if len(result) == maxResultLen {
			break
		}
	}

	return result
}

package garkov

import (
	"github.com/mickuehl/garkov/dictionary"
)

func wordsToPrefixString(prefix []dictionary.Word) string {
	k := ""
	for i := range prefix {
		k = k + prefix[i].Word
	}

	return k
}

func indexToPrefixString(prefix []int, dict *dictionary.Dictionary) string {
	k := ""
	for i := range prefix {
		k = k + dict.V[prefix[i]]
	}

	return k
}

func wordsToIndexArray(prefix []dictionary.Word) []int {
	idx := make([]int, len(prefix))

	for i := range prefix {
		idx[i] = prefix[i].Idx
	}

	return idx
}

func wordsToSentence(sentence []dictionary.Word) string {
	k := ""
	for i := range sentence {
		if sentence[i].Type < 10 {
			k = k + " " + sentence[i].Word
		} else {
			k = k + sentence[i].Word
		}
	}

	return k
}

func isStopToken(t rune) bool {
	if t == 46 {
		return true
	}
	return false
}

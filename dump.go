package garkov

import (
	"fmt"

	"github.com/mickuehl/garkov/dictionary"
)

// Debug prints the model for debugging
func (m *Markov) Debug() {

	fmt.Println("\nDumping model ...\n")

	// the word vector
	fmt.Println("Words:")
	fmt.Println(m.Dict.V)
	fmt.Println("")

	// the dictionary itself
	fmt.Println("Dictionary:")
	for w := range m.Dict.Words {
		fmt.Println(prettyPrintWord(m.Dict.Words[w]))
	}
	fmt.Println("")

	fmt.Println("Start prefixes:")
	fmt.Println(m.Start)
	fmt.Println("")

	fmt.Println("Word Chains:")
	for _, suffix := range m.Chain {
		fmt.Println(suffix.PrettyPrintChain(m.Dict))
	}

	fmt.Println("")
}

func prettyPrintWord(w dictionary.Word) string {
	return fmt.Sprintf("%v: %v[%v,%v]", w.Idx, w.Word, w.Type, w.Count)
}

func (c *WordChain) PrettyPrintChain(d *dictionary.Dictionary) string {
	_prefix := ""
	for i := range c.Prefix {
		_prefix = _prefix + d.V[c.Prefix[i]] + " "
	}

	_suffix := ""
	for w := range c.Words {
		_suffix = _suffix + w + " "
	}
	return fmt.Sprintf("%v -> %v", _prefix, _suffix)
}

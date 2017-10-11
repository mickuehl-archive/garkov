package garkov

import (
	"fmt"
	"io/ioutil"

	"github.com/jdkato/prose/tokenize"
	"github.com/mickuehl/garkov/dictionary"
)

// Build reads an input file and updates the markov model with its content.
func (m *Markov) Build(fileName string) {

	// read the file line-by-line and create an array of words
	var tokens []dictionary.Word

	tokenizer := tokenize.NewTreebankWordTokenizer()
	sentenizer, _ := tokenize.NewPragmaticSegmenter("en")

	all, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Print(err)
	}

	// split the text into complete sentences fist, regardless of the individual lines.
	content := string(all)
	for _, sentence := range sentenizer.Tokenize(content) {
		if len(sentence) > 0 {
			var word dictionary.Word
			prefix := make([]int, m.Depth)
			_prefix := 0

			//fmt.Println("Sentence: " + sentence)

			// now split the sentence into words
			for _, w := range tokenizer.Tokenize(sentence) {
				//fmt.Println(w)
				word = m.Dict.Add(w)
				tokens = append(tokens, word)

				// build the start index vector
				if _prefix < m.Depth {
					prefix[_prefix] = word.Idx
					_prefix = _prefix + 1
				}
			}

			// add the prefix to the index
			m.Start = append(m.Start, prefix)

			// check if the sentence ends with a STOP token and add one if not
			if word.Type != dictionary.SENTENCE_END {
				tokens = append(tokens, m.Dict.Add(dictionary.SENTENCE_END_TOKEN))
			}

		}
	}

	if len(tokens) > m.Depth+1 {
		pos := 0

		// only so far as there are tuples + a word
		for pos < len(tokens)-(m.Depth) {
			prefix := make([]dictionary.Word, m.Depth)

			// read the prefix
			i := 0
			for i < m.Depth {
				prefix[i] = tokens[pos+i]
				i = i + 1
			}

			// the word following the prefix
			suffix := tokens[pos+m.Depth]

			// update the chain
			m.Update(prefix, suffix)
			pos = pos + 1
		}
	}

}

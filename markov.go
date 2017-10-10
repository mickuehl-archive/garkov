package garkov

import (
	"bufio"
	"log"
	"os"
	"strings"
	"text/scanner"

	"github.com/mickuehl/garkov/dictionary"
)

// WordCount the number of occurences of a word from the word vector
type WordCount struct {
	Idx   int
	Count int
}

// WordChain is the main structure of the model. It represents a prefix and all its suffixes.
type WordChain struct {
	Prefix []int                // arrary of words forming the prefix. Index into the dictionaries word vector
	Type   int                  // the chains position, i.e. start, middle or end of sentence
	Words  map[string]WordCount // the collection of suffixes and their count
}

type Markov struct {
	Name  string                 // name of the model
	Depth int                    // prefix size
	Chain map[string]WordChain   // the prefixes mapped to the word chains
	Dict  *dictionary.Dictionary // the dictionary used in the model
}

// New creates an empty markov model.
func New(name string, depth int) *Markov {

	m := Markov{
		Name:  name,
		Depth: depth,
		Chain: make(map[string]WordChain),
		Dict:  dictionary.New(name),
	}

	return &m
}

// Train reads an input file and updates the markov model with its content.
func (m *Markov) Train(fileName string) {

	// open and read the file
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read the file line-by-line and create an array of words
	var tokens []dictionary.Word
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		tokens = m.StringToWords(line, tokens)
	}

	// analyze the array of words
	if len(tokens) > m.Depth+1 {
		state := SENTENCE_START
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
			state = m.Update(prefix, suffix, state)
			pos = pos + 1
		}
	}

}

// Update adds a prefix + suffix to the markov model
func (m *Markov) Update(prefix []dictionary.Word, suffix dictionary.Word, state int) int {

	_prefix := wordsToPrefixString(prefix)
	chain, found := m.Chain[_prefix]

	if !found {
		chain = WordChain{
			Prefix: wordsToIndexArray(prefix),
			Type:   state,
			Words:  make(map[string]WordCount),
		}
	}

	// add the word to the sequence
	chain.AddWord(suffix)

	//fmt.Println(seq)

	// update the model
	m.Chain[_prefix] = chain

	return state
}

// Close writes the model to disc
func (m *Markov) Close() {
	m.Dict.Close()

}

// StringToWords parse a sentence into an array of words
func (m *Markov) StringToWords(sentence string, tokens []dictionary.Word) []dictionary.Word {

	var sc scanner.Scanner
	sc.Init(strings.NewReader(sentence))

	var tok rune
	for tok != scanner.EOF {
		tok = sc.Scan()

		if tok != scanner.EOF {

			if tok == SINGLE_QUOTE || tok == DOUBLE_QUOTE {

				// resolve a quote to a sequence of tokens, recursively.

				// open quote
				word := m.Dict.Add("QUOTE_BEGIN", QUOTE_START_RUNE)
				tokens = append(tokens, word)

				// sentence without quotes
				l := sc.TokenText()
				tokens = m.StringToWords(l[1:len(l)-1], tokens)

				// close quote
				word = m.Dict.Add("QUOTE_END", QUOTE_END_RUNE)
				tokens = append(tokens, word)
			} else {
				word := m.Dict.Add(sc.TokenText(), tok)
				tokens = append(tokens, word)
			}
		}
	}

	return tokens
}

// AddWord updates a word chain
func (s *WordChain) AddWord(w dictionary.Word) {
	words, found := s.Words[w.Word]
	if found {
		words.Count = words.Count + 1
	} else {
		words = WordCount{
			Idx:   w.Idx,
			Count: 1,
		}
	}
	// update
	s.Words[w.Word] = words
}

func wordsToPrefixString(prefix []dictionary.Word) string {
	k := ""
	for i := range prefix {
		k = k + prefix[i].Word
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

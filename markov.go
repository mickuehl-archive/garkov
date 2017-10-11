package garkov

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"
	"text/scanner"
	"time"

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

// Markov wraps all data of a markov-chain into one
type Markov struct {
	Name   string                 // name of the model
	Depth  int                    // prefix size
	Chain  map[string]WordChain   // the prefixes mapped to the word chains
	Dict   *dictionary.Dictionary // the dictionary used in the model
	Start  [][]int                // array of start prefixes
	Random *rand.Rand
}

// New creates an empty markov model.
func New(name string, depth int) *Markov {

	m := Markov{
		Name:   name,
		Depth:  depth,
		Chain:  make(map[string]WordChain),
		Dict:   dictionary.New(name),
		Start:  make([][]int, 0),
		Random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	return &m
}

// Sentence creates a new sentence based on the markov-chain
func (m *Markov) Sentence(minWords, maxWords int) string {

	sentence := make([]dictionary.Word, m.Depth)

	// select a first prefix to start with
	_prefix := m.Start[m.Random.Intn(len(m.Start))]
	for i := range _prefix {
		w, _ := m.Dict.GetAt(_prefix[i])
		sentence[i] = w
	}
	prefix := sentence

	n := 0
	for {
		// get the next word, until we get a STOP word
		suffix := m.SuffixFor(prefix)
		sentence = append(sentence, suffix)

		if suffix.Type == dictionary.STOP && n >= minWords {
			break
		}

		// new prefix
		prefix = sentence[len(sentence)-m.Depth:]
		n = n + 1

		if n > maxWords {
			break // emergency break
		}

	}

	return wordsToSentence(sentence)
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

		last := tokens[len(tokens)-1]
		if last.Type != dictionary.STOP {

		}
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

	// create an array of start prefixes
	m.Start = make([][]int, 0)
	for c := range m.Chain {
		prefix := m.Chain[c]
		if prefix.Prefix[0] == 0 { // assume that the START token is always the first entry in the vector, i.e. has index 0
			a := make([]int, m.Depth)
			var b []int
			a = prefix.Prefix[1:]

			for w := range prefix.Words {
				// we only expect one ...
				word := prefix.Words[w]
				b = append(a, word.Idx)
			}

			m.Start = append(m.Start, b)

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

	// update the model
	m.Chain[_prefix] = chain

	return state
}

// Close writes the model to disc
func (m *Markov) Close() {
	m.Dict.Close()

}

// SuffixFor returns a word that succeedes a given prefix
func (m *Markov) SuffixFor(prefix []dictionary.Word) dictionary.Word {

	// lookup the word chain
	_prefix := wordsToPrefixString(prefix)
	chain, found := m.Chain[_prefix]

	if found {
		idx := 0
		max := m.Random.Intn(len(chain.Words))
		i := 0
		// FIXME poor implementation of a random lookup ... need a better way
		for p := range chain.Words {
			if i == max {
				idx = chain.Words[p].Idx
				break
			}
			i = i + 1
		}

		word, _ := m.Dict.GetAt(idx)
		return word
	}

	// FIXME we should never get here ...
	return dictionary.Word{}
}

// StringToWords parse a sentence into an array of words
func (m *Markov) StringToWords(sentence string, tokens []dictionary.Word) []dictionary.Word {

	// make sure that the first word is always a START token
	word := m.Dict.Add("START", SENTENCE_START_RUNE)
	tokens = append(tokens, word)

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

				if isStopToken(tok) {
					// Since the current token is a stop token, we have to insert an artificial start token.
					word := m.Dict.Add("START", SENTENCE_START_RUNE)
					tokens = append(tokens, word)
				}

			}
		}
	}

	// check that we do not end with a START token and the last one is a STOP token
	last := tokens[len(tokens)-1]
	if last.Type == SENTENCE_START {
		// cut off the last element
		tokens = tokens[:len(tokens)]
	}

	// make sure we have a proper STOP token
	last = tokens[len(tokens)-1]
	if last.Type != dictionary.STOP {
		word := m.Dict.Add(".", SENTENCE_END_RUNE)
		tokens = append(tokens, word)
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

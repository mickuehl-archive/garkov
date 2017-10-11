package garkov

import (
	"math/rand"
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

// Update adds a prefix + suffix to the markov model
func (m *Markov) Update(prefix []dictionary.Word, suffix dictionary.Word) {

	_prefix := wordsToPrefixString(prefix)
	chain, found := m.Chain[_prefix]

	if !found {
		chain = WordChain{
			Prefix: wordsToIndexArray(prefix),
			Words:  make(map[string]WordCount),
		}
	}

	// add the word to the sequence
	chain.AddWord(suffix)

	// update the model
	m.Chain[_prefix] = chain

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

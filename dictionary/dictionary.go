package dictionary

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	// token types
	WORD             int = 1
	STOP             int = 2
	COMMA            int = 3
	HYPHEN           int = 4
	COLON            int = 5
	SEMICOLON        int = 6
	INTEGER          int = 7
	FLOAT            int = 8
	EXCLAMATION_MARK int = 9
	QUESTION_MARK    int = 10
	QUOTE_BEGIN      int = 50
	QUOTE_END        int = 51
	UNKNOWN          int = -1
)

// Word the basic dictionary structure
type Word struct {
	Word  string
	Idx   int
	Type  int
	Count int
}

// Dictionary the collection of words
type Dictionary struct {
	Name  string          // name of the dictionary
	Size  int             // number of words in the dictionary
	Words map[string]Word // map of words and their stats
	V     []string        // the word vector
}

// New creates and initialize a new dictionary
func New(name string) *Dictionary {

	dict := Dictionary{
		Name:  name,
		Size:  0,
		Words: make(map[string]Word),
		V:     make([]string, 0),
	}

	return &dict

}

// Open creates a new dictionary and reads a persisted version from disc if available.
func Open(name string) *Dictionary {

	// new, empty dictionary
	dict := New(name)

	// try to open dictionary
	fileName := name + ".dict"
	file, err := os.Open(fileName)

	if err == nil {
		defer file.Close()
		// read an existing dictionary
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// parse a single line into a word
			w, word, _ := parseWord(scanner.Text())

			// update the dictionary
			dict.Words[w] = word
			dict.Size = dict.Size + 1
		}

		// initialize the word vector
		dict.V = make([]string, dict.Size)
		for w := range dict.Words {
			word := dict.Words[w]
			dict.V[word.Idx] = word.Word
		}
	}

	return dict

}

// Close persists the dictionary to disc.
func (d *Dictionary) Close() {

	fileName := d.Name + ".dict"
	f, err := os.Create(fileName)

	if err == nil {
		for w := range d.Words {
			word := d.Words[w]
			f.Write([]byte(word.ToS() + "\n"))
		}
		f.Close()
	}
}

// Add add a word to the dictionary
func (d *Dictionary) Add(w string, t rune) Word {

	word, found := d.Words[w]
	if found {
		// update the word count
		word.Count = word.Count + 1
		d.Words[w] = word
		return word
	}

	// add the word to the word vector
	d.V = append(d.V, w)

	// create a new entry
	word = Word{
		Word:  w,
		Count: 1,
		Type:  tokenType(t),
		Idx:   len(d.V) - 1,
	}

	// update the dictionary
	d.Words[w] = word
	d.Size = d.Size + 1

	// done
	return word

}

// Exists returns true if a word exists in the dictionary
func (d *Dictionary) Exists(w string) bool {
	_, found := d.Words[w]
	return found
}

// ToS dumps a word into a string
func (w *Word) ToS() string {
	return fmt.Sprintf("%v,%v,%v,%v", w.Word, w.Type, w.Count, w.Idx)
}

// parseWord parses a comma separated string into a word
func parseWord(s string) (string, Word, error) {
	// Format: word, type, count, ix
	// Example: one,1,1,1

	parts := strings.Split(s, ",")
	if len(parts) != 4 {
		return "", Word{}, errors.New("Insufficient number of parts")
	}

	// extract the parts
	t, _ := strconv.Atoi(parts[1])
	count, _ := strconv.Atoi(parts[2])
	idx, _ := strconv.Atoi(parts[3])

	w := Word{
		Word:  parts[0],
		Type:  t,
		Count: count,
		Idx:   idx,
	}

	return w.Word, w, nil
}

func tokenType(t rune) int {
	switch {
	case t == -2:
		return WORD
	case t == 46:
		return STOP
	case t == 44:
		return COMMA
	case t == 45:
		return HYPHEN
	case t == 58:
		return COLON
	case t == 59:
		return SEMICOLON
	case t == -3:
		return INTEGER
	case t == -4:
		return FLOAT
	case t == -50:
		return QUOTE_BEGIN
	case t == -51:
		return QUOTE_END
	case t == 33:
		return EXCLAMATION_MARK
	case t == 63:
		return QUESTION_MARK
	}

	return UNKNOWN
}

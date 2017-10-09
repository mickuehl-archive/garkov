package garkov

import (
	"bufio"
	"os"
)

// NewDictionary create and initialize a new dictionary
func NewDictionary(name string) *Dictionary {

	dict := Dictionary{
		Name:  name,
		Size:  0,
		Words: make(map[string]Word),
	}

	return &dict

}

// NewDictionary create and initialize a new dictionary
func OpenDictionary(name string) *Dictionary {

	dict := Dictionary{
		Name:  name,
		Size:  0,
		Words: make(map[string]Word),
	}

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
	}

	return &dict

}

// Add add a word to the dictionary
func (d *Dictionary) Add(w string, t int) Word {

	word, found := d.Words[w]
	if found {
		// update the dictionary count
		word.Count = word.Count + 1
		d.Words[w] = word
		return word
	}

	// create a new entry
	word = Word{
		Word:  w,
		Count: 1,
		Type:  t,
		Idx:   len(d.Words),
	}

	// update the dictionary
	d.Words[w] = word
	d.Size = d.Size + 1

	// all god
	return word

}

// Exists returns true if a word exists in the dictionary
func (d *Dictionary) Exists(w string) bool {
	_, found := d.Words[w]
	return found
}

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

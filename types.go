package garkov

// Word the basic dictionary structure
type Word struct {
	Word  string
	Idx   int
	Type  int
	Count int
}

// Dictionary the collection of words
type Dictionary struct {
	Name  string
	Size  int
	Words map[string]Word
	V     []string
}

type Markov struct {
	Name  string
	Depth int
	Dict  *Dictionary
}

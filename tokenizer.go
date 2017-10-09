package garkov

import (
	"strings"
	"text/scanner"
)

// AddSentence tokenizes the sentence and adds all the words to the dictionary
func (d *Dictionary) AddSentence(s string) {

	var sc scanner.Scanner
	sc.Init(strings.NewReader(s))

	var tok rune
	for tok != scanner.EOF {
		tok = sc.Scan()

		if tok != scanner.EOF {
			d.Add(sc.TokenText(), tokenType(tok))
		}
	}

}

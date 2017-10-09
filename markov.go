package garkov

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"text/scanner"
)

func OpenModel(name string) *Markov {

	m := Markov{
		Name:  name,
		Depth: 2,
	}

	m.Dict = OpenDictionary(name)
	fmt.Println(m.Dict.Words)
	fmt.Println(m.Dict.V)
	return &m
}

func (m *Markov) TrainModel(fileName string) {

	// open and read the file
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read the file line-by-line and tokenize it
	var tokens []Word
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		tokens = m.StringToTokens(line, tokens)
	}

	// ready to analyze now
	if len(tokens) > m.Depth+1 {
		pos := 0
		l := len(tokens) - (m.Depth + 1)
		for pos < l {

			// read the tupel
			w1 := tokens[pos]
			w2 := tokens[pos+1]

			// read the follower
			w3 := tokens[pos+2]

			msg := fmt.Sprintf("[%v,%v] -> %v", w1.Word, w2.Word, w3.Word)
			fmt.Println(msg)

			pos = pos + 1
		}
	}

}

func (m *Markov) Close() {
	m.Dict.Close()

}

func (m *Markov) StringToTokens(line string, tokens []Word) []Word {

	var sc scanner.Scanner
	sc.Init(strings.NewReader(line))

	var tok rune
	for tok != scanner.EOF {
		tok = sc.Scan()

		if tok != scanner.EOF {

			if tok == -5 || tok == -6 { // single/double quotes

				// open
				word := m.Dict.Add("QUOTE_BEGIN", QUOTE_BEGIN)
				tokens = append(tokens, word)

				// sentence without quotes
				l := sc.TokenText()
				tokens = m.StringToTokens(l[1:len(l)-1], tokens)

				//close
				word = m.Dict.Add("QUOTE_END", QUOTE_END)
				tokens = append(tokens, word)
			} else {
				word := m.Dict.Add(sc.TokenText(), tokenType(tok))
				tokens = append(tokens, word)
			}
		}
	}

	return tokens
}

package garkov

import (
	"fmt"
	"strconv"
	"strings"
)

/*
func token(t rune) string {
	switch {
	case t == 33:
		return "EXCLAMATION_MARK"
	case t == 46:
		return "STOP"
	case t == 44:
		return "COMMA"
	case t == 45:
		return "HYPHEN"
	case t == 58:
		return "COLON"
	case t == 59:
		return "SEMICOLON"
	case t == 63:
		return "QUESTION_MARK"
	case t == -2:
		return "WORD"
	case t == -3:
		return "INTEGER"
	case t == -4:
		return "FLOAT"
	case t == -5:
		return "QUOTE"
	case t == -6:
		return "QUOTE"
	case t == scanner.EOF:
		return "EOF"
	}
	return "???"
}
*/

func (w *Word) ToS() string {
	return fmt.Sprintf("%v,%v,%v,%v", w.Word, w.Type, w.Count, w.Idx)
}

func parseWord(s string) (string, Word, error) {
	// word, type, count, ix
	// Example: one,1,1,1

	parts := strings.Split(s, ",")
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
	case t == 33:
		return EXCLAMATION_MARK
	case t == 63:
		return QUESTION_MARK
	}

	return UNKNOWN
}

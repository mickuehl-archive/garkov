package garkov

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

	// tuples
	SENTENCE_START int = 100
	SENTENCE_MAIN  int = 200
	SENTENCE_END   int = 300

	SENTENCE_START_RUNE rune = -60
	SENTENCE_END_RUNE   rune = -61

	QUOTE_START_RUNE rune = -50
	QUOTE_END_RUNE   rune = -51

	SINGLE_QUOTE rune = -5
	DOUBLE_QUOTE rune = -6
)

package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + litlerals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators.
	ASSIGN = "="
	PLUS   = "+"

	// Delimiters.
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords.
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIdent(ident string) TokenType {
	if ttype, ok := keywords[ident]; ok {
		return ttype
	}
	return IDENT
}

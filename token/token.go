package token

import "fmt"

// Character and Literal types are useful while lexing for
// treating single characters and strings as the same.

type Character byte

func (c Character) String() string {
	return string(c)
}

func (c Character) IsLetter() bool {
	ch := byte(c)
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (c Character) IsDigit() bool {
	ch := byte(c)
	return '0' <= ch && ch <= '9'
}

func (c Character) IsChar(ch byte) bool {
	return byte(c) == ch
}

type Literal string

func (l Literal) String() string {
	return string(l)
}

// Token types are listed below. These types are used by the lexer
// in order to generate a list of tokens for the parser.

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + litlerals
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	// Operators.
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	MOD      = "%"
	BANG     = "!"

	// Relational operators.
	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters.
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords.
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

// A token can be create using the make function (not new since the returned
// value is not a pointer).

type Token struct {
	Type    TokenType
	Literal string
}

func Make(t TokenType, literal fmt.Stringer) Token {
	return Token{t, literal.String()}
}

// keywords is a map of reserved keywords in the language.
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent checks if the given identiier is a keyword.
// Returns the TokenType for the identifier if a match is found,
// otherwise, it returns the IDENT type.
func LookupIdent(ident string) TokenType {
	if ttype, ok := keywords[ident]; ok {
		return ttype
	}
	return IDENT
}

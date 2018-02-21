package lexer

import (
	"strings"

	"github.com/danielrs/monkey/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           token.Character
}

// Returns a pointer to a new Lexer.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Gets the next token in the stream.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Make(token.EQ, token.Literal("=="))
		} else {
			tok = token.Make(token.ASSIGN, l.ch)
		}
	case '+':
		tok = token.Make(token.PLUS, l.ch)
	case '-':
		tok = token.Make(token.MINUS, l.ch)
	case '*':
		tok = token.Make(token.ASTERISK, l.ch)
	case '/':
		tok = token.Make(token.SLASH, l.ch)
	case '%':
		tok = token.Make(token.MOD, l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Make(token.NOT_EQ, token.Literal("!="))
		} else {
			tok = token.Make(token.BANG, l.ch)
		}

	case '<':
		tok = token.Make(token.LT, l.ch)
	case '>':
		tok = token.Make(token.GT, l.ch)

	case ',':
		tok = token.Make(token.COMMA, l.ch)
	case ';':
		tok = token.Make(token.SEMICOLON, l.ch)
	case ':':
		tok = token.Make(token.COLON, l.ch)

	case '(':
		tok = token.Make(token.LPAREN, l.ch)
	case ')':
		tok = token.Make(token.RPAREN, l.ch)
	case '{':
		tok = token.Make(token.LBRACE, l.ch)
	case '}':
		tok = token.Make(token.RBRACE, l.ch)
	case '[':
		tok = token.Make(token.LBRACKET, l.ch)
	case ']':
		tok = token.Make(token.RBRACKET, l.ch)

	case 0:
		tok = token.Make(token.EOF, token.Literal(""))

	default:
		if l.ch.IsLetter() {
			tok.Literal = l.readWhile(token.Character.IsLetter)
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if l.ch.IsDigit() {
			tok.Literal = l.readWhile(token.Character.IsDigit)
			if len(tok.Literal) > 1 {
				tok.Literal = strings.TrimLeft(tok.Literal, "0")
			}
			tok.Type = token.INT
			return tok
		} else if l.ch.IsChar('"') {
			l.readChar()
			tok.Literal = l.readUntil(func(c token.Character) bool {
				return c == '"'
			})
			l.readChar()
			tok.Type = token.STRING
			return tok
		} else {
			tok = token.Make(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// readChar reads the next byte in the stream and advances the lexer position.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = token.Character(0)
	} else {
		l.ch = token.Character(l.input[l.readPosition])
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// readWhile acts like readChar, but continues reading until the given
// predicate is false.
func (l *Lexer) readWhile(predicate func(token.Character) bool) string {
	position := l.position
	for l.ch != 0 && predicate(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readUntil is like readWhile but stops when predicate is true.
func (l *Lexer) readUntil(predicate func(token.Character) bool) string {
	return l.readWhile(func(ch token.Character) bool {
		return !predicate(ch)
	})
}

// skipWhitespace ignores all the whitespacec starting in the current
// Lexer position.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

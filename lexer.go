package main

import (
	"bytes"
	"errors"
	"unicode/utf8"

	"github.com/danwakefield/gosh/char"
)

const EOFRune = -1

var (
	ErrQuotedString = errors.New("Unterminated quoted string")
)

type StateFn func(*Lexer) StateFn

type LexItem struct {
	Tok    Token
	Pos    int
	LineNo int
	Val    string
}

type Lexer struct {
	input       string
	inputLen    int
	lineNo      int
	lastPos     int
	pos         int
	backupWidth int
	state       StateFn
	itemChan    chan LexItem
	buf         bytes.Buffer
	backslash   bool
	quoted      bool

	CheckNewline bool
	CheckAlias   bool
	CheckKeyword bool
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:        input,
		inputLen:     len(input),
		itemChan:     make(chan LexItem),
		lineNo:       1,
		CheckNewline: false,
		CheckAlias:   true,
		CheckKeyword: true,
	}
	go l.run()
	return l
}

func (l *Lexer) emit(t Token) {
	l.itemChan <- LexItem{
		Tok:    t,
		Pos:    l.lastPos,
		LineNo: l.lineNo,
		Val:    l.buf.String(),
	}
	l.lastPos = l.pos
	l.buf.Reset()
}

func (l *Lexer) next() rune {
	if l.pos >= l.inputLen {
		l.pos++
		return EOFRune
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += w
	l.backupWidth = w
	return r
}

func (l *Lexer) backup() {
	l.pos -= l.backupWidth
	l.backupWidth = 0
}

func (l *Lexer) hasNext(r rune) bool {
	if r == l.next() {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) run() {
	for l.state = lexStart; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.itemChan)
}

func (l *Lexer) ignore() {
	l.lastPos = l.pos
}

func (l *Lexer) NextLexItem() LexItem {
	tok := <-l.itemChan

	if l.CheckNewline {
		for tok.Tok == TNewLine {
			tok = <-l.itemChan
		}
	}

	if tok.Tok != TWord || l.quoted {
		return tok
	}

	if l.CheckKeyword {
		t, found := KeywordLookup[tok.Val]
		if found {
			return LexItem{Tok: t, Pos: tok.Pos, LineNo: tok.LineNo, Val: tok.Val}
		}
	}

	return tok
}

func lexStart(l *Lexer) StateFn {
	for {
		c := l.next()

		switch c {
		default:
			l.backup()
			return lexWord
		case EOFRune:
			l.emit(TEOF)
			return nil
		case ' ', '\t': // Ignore Whitespace
			continue
		case '#': // Consume comments upto EOF or a newline
			for {
				c = l.next()
				if c == '\n' || c == EOFRune {
					l.ignore()
					l.backup()
					break
				}
			}
		case '\\': // Line Continuation or escaped character
			if l.hasNext('\n') {
				l.lineNo++
				continue
			}
			l.backup()
			l.backslash = true
			l.quoted = true
			return lexWord
		case '\n':
			l.emit(TNewLine)
			l.lineNo++
		case '&':
			if l.hasNext('&') {
				l.emit(TAnd)
			}
			l.emit(TBackground)
		case '|':
			if l.hasNext('|') {
				l.emit(TOr)
			}
			l.emit(TPipe)
		case ';':
			if l.hasNext(';') {
				l.emit(TEndCase)
			}
			l.emit(TSemicolon)
		case '(':
			l.emit(TLeftParen)
		case ')':
			l.emit(TRightParen)
		}
	}
}

func lexSubstitution(l *Lexer) StateFn {
	// Upon entering we have only read the '$'
	c := l.next()

	switch {
	default:
		// If $ is not followed by a valid variable character
		l.buf.WriteRune('$')
		l.backup()
		return lexWord
	case c == '(':
		// ARITH / SUBSHELL
	case c == '{':
		// VAR / VAROP
	case char.IsSpecial(c), char.IsFirstInVarName(c):
		//
	}

	return nil //TODO
}

func lexBackQuote(l *Lexer) StateFn {
	return nil
}
func lexWord(l *Lexer) StateFn {

OuterLoop:
	for {
		c := l.next()

		if l.backslash {
			if c == EOFRune {
				l.backup()
				l.buf.WriteRune('\\')
				break
			}
			l.buf.WriteRune(c)
			l.backslash = false
			continue
		}

		switch c {
		case '\n', ' ', EOFRune:
			l.backup()
			break OuterLoop
		case '\'':
			l.quoted = true
			return lexSingleQuote
		case '"':
			l.quoted = true
			return lexDoubleQuote
		case '`':
			return lexBackQuote
		case '$':
			return lexSubstitution
		default:
			l.buf.WriteRune(c)
		}
	}

	l.emit(TWord)
	return lexStart
}

func lexDoubleQuote(l *Lexer) StateFn {
	// XXX: Fix interpolation
	for {
		c := l.next()

		switch c {
		case EOFRune:
			panic(ErrQuotedString) //XXX: Dont make this panic
		case '\x01':
			continue
		case '"':
			return lexWord
		default:
			l.buf.WriteRune(c)
		}
	}
}

func lexSingleQuote(l *Lexer) StateFn {
	// We have consumed the first single quote before entering
	// this state.
	for {
		c := l.next()

		switch c {
		case EOFRune:
			panic(ErrQuotedString) //XXX: Dont make this panic
		case '\x01':
			continue
		case '\'':
			return lexWord
		default:
			l.buf.WriteRune(c)
		}
	}
}

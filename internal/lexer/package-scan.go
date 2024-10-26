package lexer

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/internal/common"
)


type packageLexer struct {
	*Lexer
}

func InitPackageLexer(src api.Source) *packageLexer {
	return &packageLexer{Init(src)}
}

func (lex *packageLexer) maybeOutputPackageKeyword(id string) token.Token {
	if id == "package" {
		return token.Package.MakeValued("package")
	} else if id == "require" {
		return token.Require.MakeValued("require")
	} else if id == "use" {
		return token.Use.MakeValued("use")
	}

	return token.Id.MakeValued(id)
}

func (lex *packageLexer) scan() (tok token.Token) {
	if lex.Line == 0 {
		lex.fixLineChar()
	}

	lex.skipWhitespace()
	c := lex.peek()

	// get rid of comment
	if k := lex.peek(); c == '-' && k == '-' {
		line, ok := lex.remainingLine()
		if !ok {
			panic("expected line content")
		}

		lex.Pos += len(line) // skip comment
		c = lex.peek()
	}

	// set position for whatever token is read
	defer func(start int) {
		tok.Start, tok.End = start, lex.Pos
	}(lex.Pos)

	if c == '\n' {
		lex.nextChar()
		lex.Line++
		return token.Newline.Make()
	} else if c == 0 {
		return token.EndOfTokens.Make()
	}

	switch c {
	case '(':
		lex.nextChar()
		return token.LeftParen.Make()
	case ')':
		lex.nextChar()
		return token.RightParen.Make()
	case ',':
		lex.nextChar()
		return token.Comma.Make()
	}

	line, ok := lex.remainingLine()
	if !ok {
		return lex.error(UnexpectedEndOfLine)
	}

	if res := common.Version.Match(line); res != nil {
		lex.Pos += len(*res)
		return token.Version.MakeValued(*res)
	} else if res = common.AlphanumericId.Match(line); res != nil {
		lex.Pos += len(*res)
		return lex.maybeOutputPackageKeyword(*res)
	} else if res = common.SimpleString.Match(line); res != nil {
		lex.Pos += len(*res)
		s := (*res)[1 : len(*res)-1] // remove leading and trailing '"'
		return token.StringValue.MakeValued(s)
	}

	return lex.error(UnexpectedSymbol)
}

func (lex *packageLexer) Scan() api.Token {
	return lex.scan()
}
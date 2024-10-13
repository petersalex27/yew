package token

import (
	"errors"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/util"
)

type Token struct {
	Value      string
	Typ        Type
	Start, End int
}

func (token Token) Pos() (int, int) {
	return token.Start, token.End
}

func (token Token) GetPos() api.Position {
	return api.MakePosition(token.Start, token.End)
}

func (a Token) Equals(b Token) bool {
	return a.Value == b.Value &&
		a.Typ == b.Typ &&
		a.Start == b.Start &&
		a.End == b.End
}

func (token Token) String() string {
	return token.Value
}

func (token Token) Debug() string {
	return util.ExposeToken(token)
}

func (token Token) Describe() (value string, alwaysNil []api.Node) {
	return token.Typ.String() + ": " + token.Value, nil
}

func (token Token) Type() api.NodeType { return token.Typ }

func (token Token) Error() error {
	if token.Typ == Error {
		return errors.New(token.Value)
	}
	return nil
}

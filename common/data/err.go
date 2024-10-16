package data

import "github.com/petersalex27/yew/api"

type (
	Ers = List[Err]

	Err struct {
		msg string
		fatal bool
		api.Position
	}

	EmbedsErr interface {
		api.DescribableNode
		~struct{ Err }
	}
)

func (e Err) Msg() string { return e.msg }

func (e Err) Fatal() bool { return e.fatal }

func MkErr(msg string, pos api.Positioned) Err {
	return Err{msg: msg, fatal: true, Position: pos.GetPos()}
}

func MkWarning(msg string, pos api.Positioned) Err {
	return Err{msg: msg, fatal: false, Position: pos.GetPos()}
}

func (e Err) Describe() (string, []api.Node) {
	if e.fatal {
		return "error: " + e.msg, []api.Node{} 
	}
	return "warning: " + e.msg, []api.Node{}
}

// use to apply changes to an Err
func ErrMap(f func(msg string, fatal bool) Err) func(Err) Err {
	return func(e Err) Err {
		return f(e.msg, e.fatal)
	}
} 
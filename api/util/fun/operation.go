package fun

type UnaryOp[a, b any] func(a) b
type BinaryOp[a, b, c any] func(a, b) c
type liquid[a, b any] struct{ x a }

func Pipe[out, in any](x in) liquid[in, out] { return liquid[in, out]{x} }

func (l liquid[in, out]) Into(pipeline func(in) out) out {
	return pipeline(l.x)
}

func MakeUnaryOp[a, b any](f func(a) b) UnaryOp[a, b] {
	return UnaryOp[a, b](f)
}

func MakeBinaryOp[a, b, c any](f func(a, b) c) BinaryOp[a, b, c] {
	return BinaryOp[a, b, c](f)
}

func Target[a any](f func(a)) func(a) a {
	return func(x a) a {
		f(x)
		return x 
	}
}

type CurriedBiFunc[a, b, c any] func(a) func(b) c

func (op CurriedBiFunc[a, b, c]) On(f func(a) a, g func(a) b) func(a) c {
	return func(x a) c {
		return op(f(x))(g(x))
	}
}

func (op BinaryOp[a, b, c]) On(f func(a) a, g func(a) b) UnaryOp[a, c] {
	return func(x a) c {
		return op(f(x), g(x))
	}
}

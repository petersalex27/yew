package fun

func Compose[a, b, c any](f func(a) b, g func(c) a) func(c) b {
	return func(x c) b { return f(g(x)) }
}

func ComposeOn[a, b, c, d any](op func(a, b) c, f func(d) a, g func(d) b) func(d) c {
	return func(x d) c {
		return op(f(x), g(x))
	}
}

// curry a binary function
func Curry[a, b, c any](f func(a, b) c) func(a) func(b) c {
	return func(x a) func(b) c {
		return func(y b) c {
			return f(x, y)
		}
	}
}

// uncurry a curried binary function
func Uncurry[a, b, c any](f func(a) func(b) c) func(a, b) c {
	return func(x a, y b) c { return f(x)(y) }
}

// curry a binary function w/ its arguments flipped
func CurryFlip[a, b, c any](f func(a, b) c) func(b) func(a) c {
	return func(y b) func(x a) c {
		return func(x a) c {
			return f(x, y)
		}
	}
}

// flip the arguments of a binary function
func Flip[a, b, c any](f func(a, b) c) func(b, a) c {
	return func(x b, y a) c { return f(y, x) }
}

// compose a unary function w/ a binary operation
func BinaryOpCompose[a, b, c, d any](f func(a) b, g func(c, d) a) BinaryOp[c, d, b] {
	return func(x c, y d) b { return f(g(x, y)) }
}

// compose a unary function w/ a binary function
func BiCompose[a, b, c, d any](f func(a) b, g func(c, d) a) func(c) func(d) b {
	return func(x c) func(d) b {
		return func(y d) b {
			return f(g(x, y))
		}
	}
}

// the identity function
func Id[a any](x a) a { return x }

// the Constant function
func Constant[b, a any](x a) func(b) a { return func(b) a { return x } }

func BiConstant[b, c, a any](f func(c) a) func(b, c) a {
	return func(_ b, x c) a { return f(x) }
}

func ComposeRightFlip[a, b, c, d any](f func(a, b) c, g func(d) b) func(a, d) c {
	return func(x a, y d) c {
		return f(x, g(y))
	}
}

func ComposeRightCurryFlip[a, b, c, d any](f func(a, b) c, g func(d) b) func(a) func(d) c {
	return func(x a) func(d d) c {
		return func(y d) c {
			return f(x, g(y))
		}
	}
}

func ComposeLeft[a, b, c, d any](f func(a, b) c, g func(d) a) func(d, b) c {
	return func(x d, y b) c {
		return f(g(x), y)
	}
}

func PairCurry[a, b, c, d any](f func(a, b, c) d) func(a, b) func(c) d {
	return func(x a, y b) func(c) d {
		return func(z c) d {
			return f(x, y, z)
		}
	}
}

func TripleCurry[a, b, c, d any](f func(a, b, c) d) func(a) func(b) func(c) d {
	return func(x a) func(y b) func(z c) d {
		return func(y b) func(z c) d {
			return func(z c) d {
				return f(x, y, z)
			}
		}
	}
}

func UnaryBind1st[a, b any](f func(a) b, x a) func() b {
	return func() b { return f(x) }
}

func BinBind1st_PairTarget[a, b, c, d any](f func(a, b) (c, d), y b) func(a) (c, d) {
	return func(x a) (c, d) { return f(x, y) }
}

func Bind1stOf2[a, b, c any](f func(a, b) c, y b) func(a) c {
	return func(x a) c { return f(x, y) }
}

func Bind2ndOf2[a, b, c any](f func(a, b) c, x a) func(b) c {
	return func(y b) c { return f(x, y) }
}

func Bind1stOf3[a, b, c, d any](f func(a, b, c) d, y b, z c) func(a) d {
	return func(x a) d { return f(x, y, z) }
}

func Bind2ndOf3[a, b, c, d any](f func(a, b, c) d, x a, z c) func(b) d {
	return func(y b) d { return f(x, y, z) }
}

func Bind3rdOf3[a, b, c, d any](f func(a, b, c) d, x a, y b) func(c) d {
	return func(z c) d { return f(x, y, z) }
}

// ((a, b, c) -> d) -> b -> (a -> c) -> (a -> d)
func Bind1stOn3rdOf3[a, b, c, d any](f func(a, b, c) d, y b, z func(a) c) func(a) d {
	return func(x a) d { return f(x, y, z(x)) }
}

// ((a, b, c, d) -> e) -> b -> c -> d -> (a -> e)
func Bind1stOf4[a, b, c, d, e any](f func(a, b, c, d) e, y b, z c, w d) func(a) e {
	return func(x a) e { return f(x, y, z, w) }
}

// the constant function that always returns `true`
func ConstTrue[a any](a) bool { return true }

// the constant function that always returns `false`
func ConstFalse[a any](a) bool { return false }

// the constant function that always returns `nil`
func ConstNil[a any](a) any { return nil }

// the constant function that always returns the zero value of type argument for `z`
func ConstZeroValue[z, a any](a) (zeroValue z) { return }

// performs a left fold over a list
func FoldOver[a, b any](f func(a, b) b, z b) func(...a) b {
	return func(xs ...a) b {
		acc := z
		for _, x := range xs {
			acc = f(x, acc)
		}
		return acc
	}
}

// performs a right fold over a list
func Fold[a, b any](f func(a, b) b, z b, xs []a) b {
	// right fold
	acc := z
	for i := len(xs) - 1; i >= 0; i-- {
		acc = f(xs[i], acc)
	}
	return acc
}

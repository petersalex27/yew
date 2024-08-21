## Infix Functions with Arity > 2
example 1, some kind of unnecessary function composition:
```
(<=:) : (a -> b) -> (c -> d) -> (d -> a) -> c -> b
f <=: g h = (\x => f (h (g x)))
```
example 2, ternary operator (conditional expression) w/o ':':
```
(?) : Bool -> a -> a -> a
True ? x _ = x
False ? _ x = x
```
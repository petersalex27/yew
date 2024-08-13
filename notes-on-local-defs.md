# Notes on Local (Inline) Definitions
Local (inline) definitions are those of the form:

```
fun : a -> b -> c := (\x, y => z)
```
Local definitions cannot have multiple definition clauses: For example, the following is illegal (reports a "split" definition error, ig? [this needs to be changed so it reports something more useful]):
```
div : Uint -> Uint -> Maybe Uint := (\x, y => Just (divUint x y))
div _ 0 = Nothing
```
Even if this was permitted, the local definition would cover all possible cases and `div _ 0` would never be used

They exist for two main reasons, one to resolve an ambiguity and the other to provide a more concise way of writing certain functions:

### Ambiguity Resolution
Sometimes writing both a declaration and a definition on two lines can feel like a waste of space, especially for a small let-binding, for example ...
```
let
  x : Int
  x = 3 in 
  1 + x 
```
It would be nice if we could just do this ...
```
let x : Int = 3 in 1 + x
```
Unfortunately, this causes an ambiguity between the equality type and the RHS of a binding--for example ...
```
let threeIsThree : 3 = 3 = Refl in threeIsThree
```
As a human reading this, we can figure it out, but the parser isn't quite that smart; it will struggle to understand where the type ends and the RHS of the binding begins.

#### The solution: introduce a new binding operator `:=`

The two expressions above become:
```
let x : Int := 3 in 1 + x
let threeIsThree : 3 = 3 := Refl in threeIsThree
```

# NOTE: the reason below is not implemented and will likely not be (gonna remove)
### Concise Type Term Availability

There are times we want a term of a product type to be available in the definition of the function that's given that type. One way to do this is to put it in the type twice. Once as an implicit parameter, and once as an explicit parameter. For example ...
```
TypeOf : {a : Type} -> (_ : a) -> Type
TypeOf _ = a
```
But this definition is unnecessarily verbose. The only reason for the implicit binding is to capture it for the definition.

#### The solution: allow local (inline) definitions to capture the terms of the product type that have been met

The declaration/definition pair above becomes¹ ...
```
TypeOf : a -> Type := const a
```
---
#### notes
¹ The function `const` is used above (defined in `/std/prelude/functions.yew`)
```
const : forall a . a -> _ -> a := \x, _ => x
```
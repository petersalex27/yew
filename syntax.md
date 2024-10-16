# Yew Language Syntax

Key:
- <code>name:</code> : element "name" must appear in its specific location
- <code>[name]</code> : optional elements "name"
- <code>{name}</code> : zero or more occurrences of elements "name"  

Quick, slightly ambiguous, surface view of a Yew source file
```
┌──────── Yew Source ────────┐
│ meta:                      │
│ - {annotation}             │
├┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┤
│ header:                    │
│ - [module]                 │
│ - {import}:                │
│   - {annotation}           │
│     - import statement     │
├┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┤
│ {body}:                    │
│ - {annotation}             │
│   - [data type definition] │
│   - [type alias]           │
│   - [spec definition]      │
│   - [inst definition]      │
│   - [function signature]   │
│   - [function definition]  │
│   - [syntax definition]    │
├┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┤
│ footer:                    │
│ - {annotations}            │
└──────── Source End ────────┘
```

Example file with all sections (labeled with comments):
```
[@infixl 5 (+)]                      -- meta

module example                       -- header, module

--@log Symbols                       -- import
import (
  "base" using _
  "builtin/lazy"
  "base/bool"
)

Nat : Type where (                   -- body, data type def.
  0 : Nat
  Succ : Nat -> Nat
)

alias Bool = bool.Bool               -- type alias

spec Summand sm where (              -- spec def.
  (+) : sm -> sm -> sm
)

inst Summand Nat where (             -- inst def.
  0 + y = y
  (Succ x) + y = Succ (x + y) 
)

[@inline]                            -- annotation
ifThenElse : Bool                    -- function type sig.
  -> lazy.Lazy a 
  -> lazy.Lazy a 
  -> a

ifThenElse True true _ = true        -- function def.
ifThenElse False _ false = false

syntax                               -- syntax def.
  `if` cnd `then` true `else` false  
  = ifThenElse cnd true false
                                     -- footer
```

## EBNF

NOTE: There might be slight inconsistencies with the *actual* grammar. The most accurate representation of Yew's grammar can be found in `./internal/parser/yew.ebnf`

```ebnf
(* root *)
yew source = {"\n"}, [header | body | header, then, body], footer ; 
(* header *)
header = 
    [[annotations_], module, {then, [annotations_], import statement}] 
    | [[annotations_], import statement, {then, [annotations_], import statement}] ;
  module = "module", {"\n"}, lower ident ;
  import statement = "import", {"\n"}, imports ;
    imports = 
        package import 
        | "(", {"\n"}, package import, {then, package import}, {"\n"}, ")" ;
    package import = import path, [{"\n"}, import specification] ;
      import specification = as clause | using clause ;
        as clause = "as", {"\n"}, module alias ;
        module alias = lower ident | "_" ;
        using clause = "using", {"\n"}, symbol selection group ;
        symbol selection group = 
            "_"
            | name 
            | "(", {"\n"}, name, {{"\n"}, ",", {"\n"}, name}, [{"\n"}, ","], {"\n"}, ")" ;
(* footer *)
footer = [annotations_], eof ;
(* body *)
body = [annotations_], body elem, {then, [annotations_], body elem} ;
  visibility = "public" | "open" ;
  modality = "once" | "erase" ;
  main elem = 
      def
      | spec def 
      | spec inst 
      | type def 
      | type alias 
      | typing 
      | syntax ; 
  body elem = def | visible body elem ;
  visible body elem = [visibility],
      ( spec def 
      | spec inst 
      | type def 
      | type alias 
      | typing 
      | syntax 
      ) ;
  syntax = "syntax", {"\n"}, syntax rule, {"\n"}, "=", {"\n"}, expr ;
    syntax rule = {syntax symbol, {"\n"}}, raw keyword, {{"\n"}, syntax symbol} ;
    syntax symbol = ident | "{", {"\n"}, ident, {"\n"}, "}" | raw keyword ;
    raw keyword = ? RAW STRING OF JUST A VALID ident OR symbol ? ;
  type def = typing, {"\n"}, "where", {"\n"}, type def body, [{"\n"}, deriving clause] ;
    type def body =
        "impossible"
        | [annotations_], type constructor
        | "(", {"\n"}, [annotations_], type constructor, {then, [annotations_], type constructor}, {"\n"}, ")" ;
    type constructor = constructor name seq, {"\n"}, ":", {"\n"}, type ;
    constructor name seq = constructor name, {{"\n"}, ",", {"\n"}, constructor name}, [{"\n"}, ","] ;
    deriving clause = "deriving", {"\n"}, deriving body ;
    deriving body = constrainer | "(", {"\n"}, constrainer, {{"\n"}, ",", {"\n"}, constrainer}, [{"\n"}, ","], {"\n"}, ")" ;
  type alias = "alias", {"\n"}, name, {"\n"}, "=", {"\n"}, type ;
  constrainer = upper ident, pattern | "(", {"\n"}, enc constrainer, {"\n"}, ")" ;
  enc constrainer = upper ident, {"\n"}, pattern ;
  spec head = [constraint, {"\n"}, "=>", {"\n"}], constrainer ;
  spec def = "spec", {"\n"}, spec head, [{"\n"}, spec dependency], {"\n"}, "where", {"\n"}, spec body, [{"\n"}, requiring clause] ;
    spec dependency = "from", {"\n"}, pattern ;
    spec body =
        spec member 
        | "(", {"\n"}, spec member, {then, spec member}, {"\n"}, ")" ;
    spec member = [annotations_], def | [annotations_], typing ;
    requiring clause = "requiring", {"\n"}, 
        ( [annotations_], def 
        | "(", {"\n"}, [annotations_], def, {then, [annotations_], def}, {"\n"}, ")" 
        ) ;
  spec inst = "inst", {"\n"}, spec head, {"\n"}, [spec inst target, {"\n"}], spec inst where clause ;
    spec inst target = "=", {"\n"}, constrainer ;
    spec inst where clause = "where", {"\n"}, spec inst member group ;
    spec inst member group = spec body ;
  constraint = "(", {"\n"}, constraint seq, {"\n"}, ")" | constrainer ;
    constraint seq = constraint elem, {{"\n"}, ",", {"\n"}, constraint elem}, [{"\n"}, ",", {"\n"}] ;
    constraint elem = {upper ident, {"\n"}, ",", {"\n"}}, enc constrainer ;
  typing = ["auto", {"\n"}], name, {"\n"}, ":", {"\n"}, type ;
  pattern typing = pattern, {"\n"}, ":", {"\n"}, type ;
  enc type = ["forall", {"\n"}, forall binders, {"\n"}, "in", {"\n"}], enc type tail | "(", {"\n"}, enc type, {"\n"}, ")" ;
  type = ["forall", {"\n"}, forall binders, {"\n"}, "in", {"\n"}], type tail | "(", {"\n"}, enc type, {"\n"}, ")" ;
  type tail = type term, {type term rhs}, [{"\n"}, function rhs] ;
  enc type tail = type term, {{"\n"}, type term rhs}, [{"\n"}, function rhs] ;
    type term rhs = type term | access ;
    function rhs = arrow, {"\n"}, type tail ;
    arrow = "->" | "=>" ;
    type term =
        expr atom
        | "_" | "()" | "="
        | "(", {"\n"}, enc type inner, [{"\n"}, enc typing end], {"\n"}, ")" 
        | "{", {"\n"}, enc type inner, [{"\n"}, enc typing end, [{"\n"}, default expr]], {"\n"}, "}" ;
    enc type inner = modality, {"\n"}, ident | inner type terms ;
    enc typing end = ":", {"\n"}, enc type ;
    inner type terms = enc type tail, [{{"\n"}, ",", {"\n"}, enc type tail}, [{"\n"}, ","]] ;
    constrained type = constraint, {"\n"}, "=>", {"\n"}, type ;
    default expr = ":=", {"\n"}, expr ;
    forall binders = ident, {ident} | "(", {"\n"}, ident, {{"\n"}, ident}, {"\n"}, ")" ;
  def = pattern, {"\n"}, def body ;
  def body thick arrow = def body thick arrow main, def body tail ;
  def body = def body main, def body tail ;
    def body main = with clause | "=", {"\n"}, expr ;
    def body thick arrow main = with clause | "=>", {"\n"}, expr ;
    def body tail = [{"\n"}, where clause] | "impossible" ;
    where clause = "where", {"\n"}, where body ;
    where body = main elem | "(", {"\n"}, main elem, {then, main elem}, {"\n"}, ")" ;
    with clause = "with", {"\n"}, pattern, {"\n"}, "of", {"\n"}, with clause arms ;
      with clause arms = 
          "(", {"\n"}, with clause arm, {then, with clause arm}, {"\n"}, ")" 
          | with clause arm ;
      with clause arm = [view refined pattern, {"\n"}], pattern, {"\n"}, def body thick arrow ;
      view refined pattern = pattern, {"\n"}, "|" ;
  pattern atom = literal | name | "[]" | hole ;
  pattern = pattern term, {pattern term rhs} ;
  access = ".", {"\n"}, name ;
  pattern term rhs = pattern term | access ;
  enc pattern = enc pattern term, {{"\n"}, enc pattern term rhs} ;
  enc pattern term = "=" | pattern term ; 
  enc pattern term rhs = enc pattern term | access ;
  pattern term = 
      pattern atom 
      | "_"
      | "(", {"\n"}, enc pattern seq, {"\n"}, ")" 
      | "{", {"\n"}, enc pattern seq, {"\n"}, "}" ;
    enc pattern seq = enc pattern, {{"\n"}, ",", {"\n"}, enc pattern}, [{"\n"}, ","] ;
  expr atom = pattern atom | lambda abstraction ;
  expr = expr term, {expr term rhs} ;
  enc expr = expr term, {{"\n"}, expr term rhs} ;
    expr term rhs = expr term | access ;
    expr term = 
        expr atom
        | "(", {"\n"}, enc expr, {"\n"}, ")"
        | let expr 
        | case expr ;
    binder = lower ident | upper ident | "(", {"\n"}, enc pattern, {"\n"}, ")" ;
  lambda abstraction = "\\", {"\n"}, lambda binders, {"\n"}, "=>", {"\n"}, expr ;
    lambda binders = lambda binder, {{"\n"}, ",", {"\n"}, lambda binder}, [{"\n"}, ","] ;
    lambda binder = binder | "_" ;
  let expr = "let", {"\n"}, let binding, {"\n"}, "in", {"\n"}, expr ;
    let binding = 
        binding group member 
        | "(", {"\n"}, binding group member, {then, binding group member}, {"\n"}, ")" ;
    binding group member =
        binder, {"\n"}, binding assignment
        | typing, [{"\n"}, binding assignment] ;
    binding assignment = ":=", {"\n"}, expr ;
  case expr = "case", {"\n"}, pattern, {"\n"}, "of", {"\n"}, case arms ;
    case arms = case arm | "(", {"\n"}, case arm, {then, case arm}, {"\n"}, ")" ;
    case arm = pattern, {"\n"}, def body thick arrow ;

(* annotations *)
annotations_ = annotation, {{"\n"}, annotation}, {"\n"} ;
  annotation = bound annotation | flat annotation ;
  flat annotation = ? REGEX "--[ \t]*@[ \t]*[a-z][A-Za-z0-9']*\b.*$" ? ;
  bound annotation = "[@", {"\n"}, ident, {? ANY OTHER RULE ?}, "]" ;

(* identifiers *)
lower ident = ? REGEX "[a-z][a-zA-Z0-9']*" ? ;
upper ident = ? REGEX "[A-Z][a-zA-Z0-9']*" ? ;
hole = ? REGEX "\?[A-Za-z][A-Za-z0-9']*" ? ;
infix lower ident = ? REGEX "\([a-z][a-zA-Z0-9']*\)" ? ;
infix upper ident = ? REGEX "\([A-Z][a-zA-Z0-9']*\)" ? ;
infix symbol = ? REGEX "\((?![-=]>\B|[.]{1,2}\B|:\B|\?\|)[-/*=<>!@#$%^&|~?+:.]+\)" ? ;
ident = lower ident | upper ident ;
symbol = ? REGEX "(?![-=]>\B|[.]{1,2}\B|:\B|\?\|)[-/*=<>!@#$%^&|~?+:.]+|\[\]|\(\)" ? ;
method symbol = ? REGEX "\([.]([a-z][A-Z0-9']*|[-/*=<>!@#$%^&|~?+:.]+|\[\]|\(\))\)" ? ;
infix name = infix lower ident | infix upper ident | infix symbol ;
constructor name = infix upper ident | upper ident | infix symbol | symbol ;
name = ident | symbol | infix name | method symbol ;
import path = ? REGEX "\"([a-z][a-zA-Z0-9']*)(/[a-z][a-zA-Z0-9']*)*\"" ? ;

(* literals *)
literal = string | integer | float | char ;
string = raw string | escapable string | prompted string | import path ;
raw string = ? REGEX "`.*`" ? ;
escapable string = ? REGEX "\"(\\[afnvtbr\\\"]|[^\a\f\n\v\t\b\r])*\"" ? ;
integer = decimal | hexadecimal | binary | octal ;
decimal = ? REGEX "[0-9]+(_[0-9]+)*" ? ;
hexadecimal = ? REGEX "0[xX][0-9a-fA-F]+(_[0-9a-fA-F]+)*" ? ;
binary = ? REGEX "0[bB][01]+(_[01]+)*" ? ;
octal = ? REGEX "0[oO][0-7]+(_[0-7]+)*" ? ;
float = ? REGEX "([0-9]+(_[0-9]+)*)?(\.[0-9]+(_[0-9]+)*)([eE][\+-][0-9]+(_[0-9]+)*)?" ? ;
char = ? REGEX "'(\\[afnvtbr\\']|[^\a\f\n\v\t\b\r])'" ? ;

(* misc *)
eof = ? END OF FILE ? ;
then = "\n", {"\n"} ;

(* reserved *)
reserved = 
    "deriving"
    | "module" | "import" | "as" | "using"
    | "public" | "open" 
    | "once" | "erase"
    | "where" | "with" | "requiring" 
    | "inst" | "spec" | "from" 
    | "alias" | "syntax"
    | "case" | "of" | "let" | "forall" | "in"
    | "=>" | "->" | ":" | ":=" | "=" | "|" | ","
    | "(" | ")" | "{" | "}" | "[" | "]" | "." | "_"
    | "impossible"
    | "auto" ;
pseudo keyword = nil | unit type | annotation scope head ;
unit type = "()" ;
nil = "[]" ;
annotation scope head = "[@" ;
comment = ? REGEX "--(?![ \t]*@[a-z]).*$" ? ;
comments = comment, {{"\n"}, comment} ;

(* reserved, but unused *)
ref modality = "ref" ;
dot dot = ".." ;
syntax specification = "pattern" | "term" ;
```

## Some unfortunate accepted expression of the grammar rules

- This should be removed eventually, but are documented here as a word of caution
- There are likely more, this is just the known examples that have not been dealt with

- removed `impossible` or mere expression termination ending def. bodies. Def. bodies now require one, and only one, of the previously mentioned parts AND, now, additionally require either a newline or the EOF
TODO
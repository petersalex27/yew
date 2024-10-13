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

TODO
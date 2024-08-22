# Yew Language Syntax

Key:
- <code>Name:</code> : file section "Name", must appear in its specific location
- <code>[name]</code> : optional elements "name", can appear anywhere within the top-level of the respective parent section

A yew file is broken down into the following sections:
```
┌──────────────────┐
│ Header:          │
│  Meta:           │
│  - [directives]  │
│  - [embeds]      │
│ - [module]       │
│ - [imports]      │
├┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┤
│ Body:            │
│ - [definitions]  │
│ - [directives]   │
│ - [mutuals]      │
│ - [syntaxes]     │
├┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┄ ┤
│ Footer:          │
│ - [directives]   │
└──────────────────┘
```

Example file with all sections:
```
-- Header Meta [directives]:
-- @build "example"
-- @excludeBase

-- Header Meta [embeds]
-- @embed "hello.h" "c" 
int fhello(FILE *);
-- @end embed "hello.h"
-- @embed "hello.c" "c"
#include <stdio.h>

int fhello(FILE *f) {
  return fprintf(f, "hello, world!");
}
-- @end embed "hello.c"

-- Header [module]:
module exampleMod

-- Header [imports]:
import example2/mod.Example

-- Body [mutual]:
mutual (
  -- Body [directives]:
  -- @infix Left 2 (foo)

  -- Body [definitions]:
  (foo) : Int -> Int
  0 foo = 0
  (Succ n) foo = bar n

  bar : Int -> Int
  bar 0 = 1
  bar (Succ n) = n foo
)

-- Body [definitions]:
Nat : Type where (
  Zero : Nat
  Succ : Nat -> Nat
)

syntax ifThenElse : (
  {`if`} Bool
    -> {`then`} Lazy a
    -> {`else`} Lazy a
    -> a
)
ifThenElse True t _ = t
ifThenElse False _ f = f

-- Footer [directives]:
-- @eof
```

## EBNF

TODO
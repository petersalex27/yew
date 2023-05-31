# Yew Grammar (W.I.P.)
## Definitions
_VALUE_ ∈ {x : x matches  }\
_ID_ ∈ {x : x matches [A-Za-z]+[A-Za-z0-9'_]*}\
## Operator Precedence and Evaluation Direction 
| Operator Class | Members | Direction |
|--|--|--|
| Composition | `.` | left to right |
| Application | - | left to right |
| Power | `^` | right to left |
| Postfix | `!` | left to right |
| Unary | `! - +` | right to left |
| Productive | `* /` | left to right |
| Addative | `+ - ++ :` | left to right |
| Order | `> >= < <=` | left to right |
| Equality | `== !=` | left to right |
| Logical And | `&` | left to right |
| Logical Or | <code>&#124;&#124;</code> | left to right |
## Grammar
Top ::= Package Modules\
MaybeAnotation ::= ɛ\
    &emsp;| `@` _ID_\
Modules ::= ɛ\
    &emsp;| Program Modules\
    &emsp;| Module Modules\
Module ::= MaybeAnotation `module` _ID_ `where` `{` Program `}`\
Package ::= MaybeAnotation `package` _ID_\
Program ::= Definition ProgramTail\
    &emsp;| Expression ProgramTail\
    &emsp;| `{` ProgramTail `}`\
ProgramTail ::= ɛ\
    &emsp;| Program\
Expression ::= Expression2 ExpressionTypeAnotation\
ExpressionTypeAnotation ::= ɛ\
    &emsp;| `::` Type\
Expression2 ::= `(` Expression `)`\
    &emsp;| Value\
    &emsp;| ExpressionSequence\
    &emsp;| Program\
    &emsp;| _ID_\
    &emsp;| Operation\
    &emsp;| UnaryOperation\
Operation ::= Expression BinaryOp Expression\
UnaryOperation ::= UnaryOp Expression\
BinaryOp ::= `+`\
    &emsp;| `-`\
    &emsp;| `*`\
    &emsp;| `/`\
    &emsp;| `^`\
    &emsp;| `:`\
    &emsp;| `++`\
    &emsp;| `.`\
    &emsp;| `>`\
    &emsp;| `>=`\
    &emsp;| `<`\
    &emsp;| `<=`\
    &emsp;| `&&`\
    &emsp;| `||`\
UnaryOp ::= 
ExpressionSequence ::= \
    ExpressionSequence `;` Expression\
Value ::= VALUE \
    &emsp;| List\
    &emsp;| Tuple\
    &emsp;| `True`\
    &emsp;| `False`\
Tuple ::= `(` ListElements `)`\
List ::= `[` ListElements `]`\
ListElements ::= ɛ \
    &emsp;| _VALUE_ ListTail\
ListTail ::= `,` ListElements \
Type ::= `(` TypeListInitial `)`\
    &emsp;| `(` Type `)`\
    &emsp;| FunctionType\
    &emsp;| `(` NamedTypes `)`\
    &emsp;| `[` Type ListSize `]`\
    &emsp;| `Int`\
    &emsp;| `Char`\
    &emsp;| `Bool`\
    &emsp;| `Float`\
    &emsp;| `String`\
    &emsp;| Kind\
FunctionType ::= TypeList `->` Type\
TypeListInitial ::= ɛ\
    &emsp;| TypeList\
TypeList ::= Type TypeListTail
TypeListTail ::= ɛ\
    &emsp;| `,` ɛ
    &emsp;| `,` TypeList\
Kind ::= VALUE\
    &emsp;| TypeConstructor\
TypeConstructor ::= _ID_ Type\
    &emsp;| _ID_ \
ListSize ::= ɛ\
    &emsp;| `;` _VALUE_\
NamedTypes ::= _ID_ Type NamedTypesTail\
NamedTypesTail ::= ɛ\
    &emsp;| `,` ɛ\
    &emsp;| `,` NamedTypes\
TypeAnotation ::= ɛ\
    &emsp;| Type\
DeclarationKind ::= `let`\
    &emsp;| `const`\
    &emsp;| `mut`\
Declaration ::= MaybeAnotation `let` _ID_ TypeAnotation\
    &emsp;| MaybeAnotation FunctionDeclaration\
FunctionDeclaration ::= _ID_ TypeAnotationList `->` Type\
TypeAnotationListInitial ::= `(` `)`\
    &emsp;| TypeAnotationList\
TypeAnotationList ::= \
    ID TypeAnotation TypeAnotationListTail\
    &emsp;| Pattern TypeAnotationTail\
TypeAnotationListTail ::= ɛ\
    &emsp;| `,` TypeAnotationList\
Pattern ::= \
Definition ::= Declaration `=` Assignment\
Assignment ::= Expression\
Application ::=\
Function ::=\
Class ::= `class` _ID_ `where` ClassStart\
ClassStart ::= `{` ClassBodies `}`\
    &emsp;| `{` ClassSequence `}`\
    &emsp;| ClassSequence\
ClassBodies ::= ClassBody ClassBodiesTail\
ClassBodiesTail ::= ɛ\
    &emsp;| ClassBodies\
ClassSequence ::= ClassBody ClassSequenceTail\
ClassSequenceTail ::= ɛ\
    &emsp;| `;` ClassSequence\
ClassBody ::= _ID_ FunctionType 
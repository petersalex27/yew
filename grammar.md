# Yew Grammar (W.I.P.)
Top ::= Package Modules\
MaybeAnotation ::= ɛ\
    &emsp;| `@` ID\
Modules ::= ɛ\
    &emsp;| Program Modules\
    &emsp;| Module Modules\
Module ::= MaybeAnotation `module` ID `where` `{` Program `}`\
Package ::= MaybeAnotation `package` ID\
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
    &emsp;| ID\
    &emsp;| Operation\
    &emsp;| UnaryOperation\
Operation ::=\
UnaryOperation ::=\
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
    &emsp;| VALUE ListTail\
ListTail ::= `,` ListElements \
Type ::= `(` TypeAnotation `)`\
    &emsp;| `(` NamedTypes `)`\
    &emsp;| `[` Type ListSize `]`\
    &emsp;| `Int`\
    &emsp;| `Char`\
    &emsp;| `Bool`\
    &emsp;| `Float`\
    &emsp;| `String`\
    &emsp;| Kind\
Kind ::= VALUE\
    &emsp;| TypeConstructor\
TypeConstructor ::= ID Type\
    &emsp;| ID \
ListSize ::= ɛ\
    &emsp;| `;` VALUE\
NamedTypes ::= ID Type NamedTypesTail\
NamedTypesTail ::= ɛ\
    &emsp;| `,` ɛ\
    &emsp;| `,` NamedTypes\
TypeAnotation ::= ɛ\
    &emsp;| Type\
Declaration ::= MaybeAnotation `let` ID TypeAnotation\
    &emsp;| MaybeAnotation FunctionDeclaration\
FunctionDeclaration ::= ID TypeAnotationList\
TypeAnotationListInitial ::= `(` `)`\
    &emsp;| TypeAnotationList\
TypeAnotationList ::= \
    ID TypeAnotation TypeAnotationListTail\
    &emsp;| Pattern TypeAnotationTail\
TypeAnotationListTail ::= ɛ\
    &emsp;| `,` TypeAnotationList_
Pattern ::= \
Definition ::= Declaration `=` Assignment\
Assignment ::= Expression\
Application ::=\
Function ::=\
Class ::= \
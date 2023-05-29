# Yew Grammar (W.I.P.)
Top ::= Package Modules
MaybeAnotation ::= ɛ
    | `@` ID
Modules ::= ɛ
    | Program Modules
    | Module Modules
Module ::= MaybeAnotation `module` ID `where` `{` Program `}`
Package ::= MaybeAnotation `package` ID
Program ::= Definition ProgramTail
    | Expression ProgramTail
    | `{` ProgramTail `}`
ProgramTail ::= ɛ
    | Program 
Expression ::= Expression2 ExpressionTypeAnotation
ExpressionTypeAnotation ::= ɛ
    | `::` Type
Expression2 ::= `(` Expression `)`
    | Value
    | ExpressionSequence
    | Program
    | ID
    | Operation
    | UnaryOperation 
Operation ::=
UnaryOperation ::=
ExpressionSequence ::= 
    ExpressionSequence `;` Expression
Value ::= VALUE 
    | List
    | Tuple
    | `True`
    | `False`
Tuple ::= `(` ListElements `)`
List ::= `[` ListElements `]`
ListElements ::= ɛ 
    | VALUE ListTail
ListTail ::= `,` ListElements 
Type ::= `(` TypeAnotation `)`
    | `(` NamedTypes `)`
    | `[` Type ListSize `]`
    | `Int`
    | `Char`
    | `Bool`
    | `Float`
    | `String`
    | Kind
Kind ::= VALUE
    | TypeConstructor
TypeConstructor ::= ID Type
    | ID 
ListSize ::= ɛ
    | `;` VALUE
NamedTypes ::= ID Type NamedTypesTail
NamedTypesTail ::= ɛ
    | `,` ɛ
    | `,` NamedTypes
TypeAnotation ::= ɛ
    | Type
Declaration ::= MaybeAnotation `let` ID TypeAnotation
    | MaybeAnotation FunctionDeclaration
FunctionDeclaration ::= ID TypeAnotationList
TypeAnotationListInitial ::= `(` `)`
    | TypeAnotationList
TypeAnotationList ::= 
    ID TypeAnotation TypeAnotationListTail
    | Pattern TypeAnotationTail
TypeAnotationListTail ::= ɛ
    | `,` TypeAnotationList
Pattern ::= 
Definition ::= Declaration `=` Assignment
Assignment ::= Expression
Application ::=
Function ::=
Class ::= 
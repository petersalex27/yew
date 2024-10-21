package symbol

const (
	OutsideLinearContext = "linear symbol cannot be accessed in a non-linear context" // outside-linear-context
	OutsideErasedContext = "erased symbol not accessible"                             // outside-erased-context
	InvalidAliasIdent    = "invalid form for alias identifier"                        // invalid-alias-ident
	AliasCycle           = "alias cannot be defined in terms of itself"               // alias-cycle
	AliasRedefined       = "alias cannot be redefined"                                // alias-redefined
	DuplicateDeclaration = "symbol already declared"                                  // duplicate-declaration
	UndeclaredSymbol     = "symbol not declared"                                      // undeclared-symbol
)

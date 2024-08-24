package parser

func (v literalValue) f()

func (w wildcard) f()

func (n name) f()

func (ps tuple[pattern]) f()

func (ps app[pattern]) f()

func (ps list[pattern]) f()

func (v literalValue) visitPatternApp(pp app[pattern], parser *Parser) maybe[app[pattern]] {
	pp.apply(v)
	return unit(pp)
}
func (v literalValue) visitPatternTuple(pp tuple[pattern], parser *Parser) maybe[tuple[pattern]]
func (v literalValue) visitPatternList(ps list[pattern], parser *Parser) maybe[list[pattern]]
func (v literalValue) visitNameAppPattern(pp nameAppPattern, parser *Parser) maybe[nameAppPattern]
func (v literalValue) visitTyp(t typ, parser *Parser) maybe[typ]
func (v literalValue) visitExplicitType(t explicitTyp, parser *Parser) maybe[explicitTyp]
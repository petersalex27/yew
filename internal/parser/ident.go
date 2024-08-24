package parser

// things that implement the identity node interface 'ident'

func (id lowerIdent) visitAnnotation(a annotation, parser *Parser) maybe[annotation] {
	a.id = id
	return unit(a)
}
func (id lowerIdent) visitIdentTyping(it identTyping, parser *Parser) maybe[identTyping] {
	it.ident = id
	return unit(it)
}
func (id lowerIdent) visitForall(fa forall, parser *Parser) maybe[forall] {
	fa.binding = append(fa.binding, id)
	return unit(fa)
}
func (id lowerIdent) visitLambdaBinder(lb lambdaBinder, parser *Parser) maybe[lambdaBinder] {
	lb.ident = id
	return unit(lb)
}
func (id lowerIdent) visitBindingAssignment(ba bindingAssignment, parser *Parser) maybe[bindingAssignment] {
	panic("TODO: implement")
}
func (id lowerIdent) visitPrefixedName(pn prefixedName, parser *Parser) maybe[prefixedName] {
	pn.ident = id.ident
	return unit(pn)
}

func (id upperIdent) visitAnnotation(a annotation, parser *Parser) maybe[annotation] {
	a.id = id
	return unit(a)
}
func (id upperIdent) visitIdentTyping(it identTyping, parser *Parser) maybe[identTyping] {
	it.ident = id
	return unit(it)
}
func (id upperIdent) visitForall(fa forall, parser *Parser) maybe[forall] {
	fa.binding = append(fa.binding, id)
	return unit(fa)
}
func (id upperIdent) visitLambdaBinder(lb lambdaBinder, parser *Parser) maybe[lambdaBinder] {
	lb.ident = id
	return unit(lb)
}
func (id upperIdent) visitBindingAssignment(ba bindingAssignment, parser *Parser) maybe[bindingAssignment] {
	panic("TODO: implement")
}
func (id upperIdent) visitPrefixedName(pn prefixedName, parser *Parser) maybe[prefixedName] {
	pn.ident = id.ident
	return unit(pn)
}
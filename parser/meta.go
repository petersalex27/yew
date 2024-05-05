package parser

type MetaData interface {
	LoadMetaData(parser *Parser, start, end int, args ...Term) (ok bool)
}
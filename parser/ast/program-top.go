/* interface for programs and modules */

package ast

import "yew/parser/parser"

type ProgramTop interface {
	parser.Ast
	GetProgram() Program
	GetNameSpace() Id
}
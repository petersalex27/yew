//go:build !debug
// +build !debug

package types

import "io"

func debug_log_step(io.Writer, Term, string, bool)         {}
func debug_log_conclusion(io.Writer, Term, string, string) {}
func debug_log_steps[T Term](io.Writer, []T, bool)         {}

func debug_log_Rule(Type, Type, Type) {}

func debug_log_Constrain(Type, Type, Pi) {}

func debug_log_Gen(Type, Type, Pi) {}

func debug_log_Forall([]Variable, Type, Forall) {}

func debug_log_App(Lambda, Type, Term, Type, Pi, Term, Type) {}

func debug_log_IAbs(Variable, Type, Term, Type, Pi, Lambda) {}

func debug_log_Abs(Variable, Type, Term, Type, Pi, Lambda) {}

func debug_log_IProd(Type, Sort, Variable, Type, Sort, Pi) {}

func debug_log_Prod(Type, Sort, Variable, Type, Sort, Pi) {}

func debug_log_Var(Variable, Type) {}

func debug_log_Red(Type, Type) {}

func debug_log_Con(Constant, Universe, Constant, Type, Term) {}

func debug_log_IApp(Lambda, Pi, Term, Type, Term, Type) {}

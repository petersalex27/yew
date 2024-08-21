//go:build debug
// +build debug

package types

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func debug_log_step(w io.Writer, t Term, addition string, assumption bool) {
	if !assumption {
		fmt.Fprintf(w, "%v%v\n", t, addition)
		return
	}
	fmt.Fprintf(w, "[%v%v]\n", t, addition)
}

func debug_log_steps[T Term](w io.Writer, ts []T, assumption bool) {
	for _, t := range ts {
		debug_log_step(w, t, "", assumption)
	}
}

func debug_log_conclusion(w io.Writer, t Term, addition, rule string) {
	conclusion := t.String() + addition
	line := strings.Repeat("-", len(conclusion))
	fmt.Fprintf(w, "%s\n%s\n\n", line + " " + rule, conclusion)
}

func debug_log_conclusion2(w io.Writer, assumption string, t Term, addition, rule string) {
	conclusion := t.String() + addition
	line := strings.Repeat("-", len(conclusion))
	fmt.Fprintf(w, "%s\n%s\n%s\n\n", line + " " + rule, assumption, conclusion)
}

func debug_log_Rule(s, t, u Type) {
	debug_log_step(os.Stderr, s, "", false)
	debug_log_step(os.Stderr, t, "", false)
	debug_log_conclusion(os.Stderr, u, "", "(Rule)")
}

func debug_log_Constrain(constraint, constrained Type, pi Pi) {
	debug_log_step(os.Stderr, constraint, "", false)
	debug_log_step(os.Stderr, constrained, "", false)
	debug_log_conclusion(os.Stderr, pi, "", "(Constrain)")
}

func debug_log_Gen(H, A Type, Hp Pi) {
	debug_log_step(os.Stderr, H, "", false)
	debug_log_step(os.Stderr, A, "", false)
	debug_log_conclusion(os.Stderr, Hp, "", "(Gen)")
}

func debug_log_Forall(as []Variable, A Type, f Forall) {
	debug_log_steps(os.Stderr, as, true)
	debug_log_step(os.Stderr, A, "", false)
	debug_log_conclusion(os.Stderr, f, "", "(Forall)")
}

func debug_log_App(f Lambda, F Type, a Term, Ap Type, H Pi, term Term, ty Type) {
	debug_log_step(os.Stderr, f, ":->>" + F.String(), false)
	debug_log_step(os.Stderr, a, ":" + Ap.String(), false)
	debug_log_step(os.Stderr, H.binderVar.Kind, " =β " + Ap.String(), false)
	debug_log_conclusion(os.Stderr, term, ":" + ty.String(), "(App)")
}

func debug_log_IAbs(x Variable, A Type, b Term, B Type, P Pi, lam Lambda) {
	debug_log_step(os.Stderr, x, ":" + A.String(), true)
	debug_log_step(os.Stderr, b, ":" + B.String(), false)
	_, kB := B.GetKind()
	debug_log_step(os.Stderr, P, ":" + kB.String(), false)
	debug_log_conclusion(os.Stderr, lam, ":" + P.String(), "(IAbs)")
}

func debug_log_Abs(x Variable, A Type, b Term, B Type, P Pi, lam Lambda) {
	debug_log_step(os.Stderr, x, ":" + A.String(), true)
	debug_log_step(os.Stderr, b, ":" + B.String(), false)
	_, kB := B.GetKind()
	debug_log_step(os.Stderr, P, ":" + kB.String(), false)
	debug_log_conclusion(os.Stderr, lam, ":" + P.String(), "(Abs)")
}

func debug_log_IProd(A Type, s Sort, x Variable, B Type, u Sort, pi Pi) {
	debug_log_step(os.Stderr, A, ":->>" + s.String(), false)
	debug_log_step(os.Stderr, x, ":" + A.String(), true)
	_, t := B.GetKind()
	debug_log_step(os.Stderr, B, ":->>" + t.String(), false)
	debug_log_step(os.Stderr, s, "~>" + t.String() + ":" + u.String(), false)
	debug_log_conclusion(os.Stderr, pi, "", "(IProd)")
}

func debug_log_Prod(A Type, s Sort, x Variable, B Type, u Sort, pi Pi) {
	debug_log_step(os.Stderr, A, ":->>" + s.String(), false)
	debug_log_step(os.Stderr, x, ":" + A.String(), true)
	_, t := B.GetKind()
	debug_log_step(os.Stderr, B, ":->>" + t.String(), false)
	debug_log_step(os.Stderr, s, "~>" + t.String() + ":" + u.String(), false)
	debug_log_conclusion(os.Stderr, pi, "", "(Prod)")
}

func debug_log_Var(x Variable, A Type) {
	debug_log_step(os.Stderr, x, ":" + A.String() + " ∈ 𝚪", false)
	debug_log_conclusion(os.Stderr, x, ":" + A.String(), "(Var)")
}

func debug_log_Red(A, rA Type) {
	debug_log_step(os.Stderr, A, "", false)
	debug_log_conclusion(os.Stderr, rA, "", "(Red)")
}

func debug_log_Con(Z Constant, Type_n Universe, C Constant, D Type, Cons Term) {
	fmt.Fprintf(os.Stderr, "free(%v) ∈ 𝚪", C)
	assumption := ":" + Type_n.String() + " .. where .."
	debug_log_step(os.Stderr, Z, assumption, true)
	assumption = Z.String() + assumption
	debug_log_step(os.Stderr, D, "", false)
	debug_log_conclusion2(os.Stderr, assumption, Cons, ":" + D.String(), "(Con)")
}

func debug_log_IApp(f Lambda, F Pi, a Term, Ap Type, term Term, ty Type) {
	debug_log_step(os.Stderr, f, ":->>" + F.String(), false)
	debug_log_step(os.Stderr, a, ":" + Ap.String(), false)
	debug_log_step(os.Stderr, F.binderVar.Kind, " =β " + Ap.String(), false)
	debug_log_conclusion(os.Stderr, term, ":" + ty.String(), "(IApp)")
}
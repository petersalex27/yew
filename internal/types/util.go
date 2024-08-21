// =================================================================================================
// Alex Peters - 2024
// =================================================================================================
package types

func GetConstant(t Term) Constant {
	switch a := t.(type) {
	case Constant:
		return a
	case Application:
		if len(a.terms) == 0 {
			return Constant{}
		}
		return GetConstant(a.terms[0])
	case Universe:
		return Constant{a.String(), 0, 0}
	}
	return Constant{}
}

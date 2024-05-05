// =============================================================================
// Author-Date: Alex Peters - November 19, 2023
//
// Content: methods associated w/ Context's error reports
//
// Notes: -
// =============================================================================
package inf

func (cxt *Context[N]) appendReport(report errorReport[N]) {
	cxt.reports = append(cxt.reports, report)
}

func (cxt *Context[N]) GetReports() []errorReport[N] {
	return cxt.reports
}

func (cxt *Context[N]) HasErrors() bool {
	return len(cxt.reports) != 0
}
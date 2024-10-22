package parser

import (
	"fmt"

	"github.com/petersalex27/yew/api"
	t "github.com/petersalex27/yew/internal/parser/typ"
)

// requirements implements packageBodyElement
func (reqs requirements) packageBodyElement() packageBodyElement { return reqs }

func (n language) Type() api.NodeType         { return t.PackageLanguage }
func (n requirements) Type() api.NodeType     { return t.YewPackageRequirements }
func (n versionCanonical) Type() api.NodeType { return t.YewPackageVersionCanonical }
func (n versionRelative) Type() api.NodeType  { return t.YewPackageVersionRelative }
func (n YewPackageSource) Type() api.NodeType { return t.YewPackageSource }

func (n language) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n requirements) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n versionCanonical) Describe() (string, []api.Node) {
	s := n.Type().String() + fmt.Sprintf(" : v%d.%d.%d", n.major, n.minor, n.patch)
	if n.milestone != "" {
		s += fmt.Sprintf("@%s.%d", n.milestone, n.milestoneNum)
	}
	return s, []api.Node{}
}
func (n versionRelative) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n YewPackageSource) Describe() (string, []api.Node) {
	return n.Type().String(), []api.Node{n.name, n.lang, n.body, n.mods}
}

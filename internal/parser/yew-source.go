package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/data"
)

func (ys YewSource) Header() data.Maybe[header] {
	return ys.header
}

func (ys YewSource) Imports() data.Maybe[data.List[importStatement]] {
	return data.MaybeMap((header).Snd)(ys.header)
}

func (ys YewSource) Module() data.Maybe[module] {
	return data.Bind(ys.header, (header).Fst)
}

func (ys YewSource) Body() data.Maybe[body] {
	return ys.body
}

func (ys YewSource) Footer() data.Maybe[annotations] {
	return ys.footer
}

func (ys YewSource) GetPos() api.Position {
	return ys.pos
}

func (ys YewSource) Pos() (int, int) {
	return ys.pos.Pos()
}
package parsing

import (
	"strings"
	err "yew/error"
	symbol "yew/symbol"
)



type Name struct {
	sym symbol.Symbolic
	left *Name
	right *Name
}

type searchDirection byte
const (
	searchLeft searchDirection = iota
	searchRight
	foundName
	nameNotFound 
)

type NameStack struct {
	names []*Name
	searchRoot *Name
}

func (names *NameStack) Push(name *Name) {
	names.names = append(names.names, name)
}

func (names *NameStack) Pop() *Name {
	if len(names.names) == 0 {
		err.PrintBug()
		panic("")
	}
	out := names.names[len(names.names)-1]
	names.names = names.names[:len(names.names)-1]
	return out
}

func (names *NameStack) Peek() *Name {
	if len(names.names) == 0 {
		err.PrintBug()
		panic("")
	}
	return names.names[len(names.names)-1]
}

func compare(root *Name, name *Name) searchDirection {
	res := strings.Compare(root.sym.GetIdToken().ToString(), name.sym.GetIdToken().ToString())
	if res == 0 {
		return foundName
	} else if res < 1 {
		return searchRight
	}
	return searchLeft
}

// expects n to be non-nil
func parentSearcher_(root *Name, n *Name) (*Name, bool) {
	if nil == root {
		return root, false
	}
	var out *Name = nil
	var swap bool = false
	switch compare(root, n) {
	case foundName:
		return root, true
	case searchRight:
		out, swap = parentSearcher_(root.right, n)
	case searchLeft:
		out, swap = parentSearcher_(root.left, n)
	}
	if swap {
		out = root
		swap = false
	}
	return out, swap
}

func parentSearcher(root *Name, n *Name) *Name {
	parent, swap := parentSearcher_(root, n)
	if swap { 
		// root node matched n, no parent
		return nil
	}
	return parent
}

func searcher(root *Name, n *Name) *Name {
	if nil == n {
		return nil
	}

	if compare(root, n) == foundName {
		return root
	}

	parent := parentSearcher(root, n)
	if nil == parent {
		return parent
	}

	if compare(parent.left, n) == foundName {
		return parent.left
	}
	return parent.right
}

func remover(root *Name, n *Name) *Name {
	if nil == n {
		return nil
	}

	if compare(root, n) == foundName {
		return root
	}

	parent := parentSearcher(root, n)
	if nil == parent {
		return parent
	}

	if compare(parent.left, n) == foundName {
		return parent.left
	}
	return parent.right
}

func placer(root *Name, n *Name) {
	switch compare(root, n) {
	case foundName:
		return
	case searchRight:
		if nil == root.right {
			root.right = n
			return
		}
		placer(root.right, n)
	case searchLeft:
		if nil == root.left {
			root.left = n
			return
		}
		placer(root.left, n)
	}
}

func (names *NameStack) Find(name *Name) *Name {
	if nil == name {
		return nil 
	}
	return searcher(names.searchRoot, name)
}
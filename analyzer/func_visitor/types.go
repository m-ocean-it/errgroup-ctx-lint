package func_visitor

import (
	"go/types"
	"slices"
)

type CommentPosition struct {
	Filename string
	Line     int
}

type errgroupStack []errgroupStackElement

type errgroupStackElement struct {
	groupObj types.Object
	ctxObj   types.Object
	depth    int
}

func (s errgroupStack) Trim(depth int) errgroupStack {
	if len(s) == 0 {
		return s
	}

	for i, elem := range s {
		if elem.depth > depth {
			return s[:i]
		}
	}

	return s
}

func (s errgroupStack) LastCtx() types.Object {
	for _, elem := range slices.Backward(s) {
		if elem.ctxObj != nil {
			return elem.ctxObj
		}
	}

	return nil
}

type goStack []int

func (gs goStack) Trim(depth int) goStack {
	if len(gs) == 0 {
		return gs
	}

	for i, d := range gs {
		if d > depth {
			return gs[:i]
		}
	}

	return gs
}

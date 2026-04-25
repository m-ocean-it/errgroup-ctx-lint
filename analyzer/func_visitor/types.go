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
	ctxName  string
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

// FindByGroup returns the most recent stack element matching the given group
// variable object.
func (s errgroupStack) FindByGroup(groupObj types.Object) *errgroupStackElement {
	for _, frame := range slices.Backward(s) {
		if frame.groupObj == groupObj {
			return &frame
		}
	}

	return nil
}

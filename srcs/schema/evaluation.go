package schema

import "strings"

type EvalResult struct {
	Course        string    `json:"course"`
	Pass          bool      `json:"pass"`
}

type SortableEvalRes []EvalResult

func (ser SortableEvalRes) Len() int {
	return len(ser)
}

func (ser SortableEvalRes) Less(i, j int) bool {
	switch cmp := strings.Compare(ser[i].Course, ser[j].Course); {
	case cmp < 0:
		return true
	case cmp == 0:
		if !ser[i].Pass {
			return false
		}
	default:
		return false
	}
	return false
}

func (ser SortableEvalRes) Swap(i, j int) {
	ser[i].Pass, ser[j].Pass = ser[j].Pass, ser[i].Pass
	ser[i].Course, ser[j].Course = ser[j].Course, ser[i].Course
}

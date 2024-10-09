package utils

type ArrayUtil struct{}

func (a *ArrayUtil) Contains(elm any, dst []any) bool {
	for _, v := range dst {
		if elm == v {
			return true
		}
	}
	return false
}

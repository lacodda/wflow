package core

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func NotNil(items ...interface{}) any {
	for _, item := range items {
		if item != nil {
			return item
		}
	}
	return nil
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

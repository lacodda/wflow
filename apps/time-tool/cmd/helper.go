package cmd

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

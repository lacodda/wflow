package core

import (
	"fmt"
	"strings"
)

const padLimit = 8192

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

func MinutesToTimeStr(totalMinutes int) string {
	minutes := fmt.Sprintf("%v", totalMinutes%60)
	hours := fmt.Sprintf("%v", float64(totalMinutes/60))

	return fmt.Sprintf("%v:%v", PadLeft(hours, 2, "0"), PadLeft(minutes, 2, "0"))
}

func PadLeft(in string, size int, sep string) string {
	if sep == "" {
		sep = " "
	}
	sepLen := len(sep)
	inLen := len(in)
	pads := size - inLen
	if pads <= 0 {
		return in
	}
	if sepLen == 1 && pads <= padLimit {
		return strings.Repeat(sep, pads) + in
	}
	if pads == sepLen {
		return sep + in
	} else if pads < sepLen {
		return sep[0:pads] + in
	} else {
		var padding string
		for i := 0; i < pads; i++ {
			padding += string(sep[i%sepLen])
		}
		return padding + in
	}
}

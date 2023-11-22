package core

import (
	"github.com/fatih/color"
)

var (
	Danger  = color.New(color.FgRed).PrintfFunc()
	Info    = color.New(color.FgCyan).PrintfFunc()
	Warning = color.New(color.FgYellow).PrintfFunc()
	Success = color.New(color.FgGreen).PrintfFunc()
)

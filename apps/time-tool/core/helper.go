package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
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

func CreateIfNotExists(filename string) error {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) || info.IsDir() {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Close()
		return nil
	}
	return err
}

func ArrayToString(ints []int, sep string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(ints), " ", sep, -1), "[]")
}

func DayRange(date time.Time) (from time.Time, to time.Time) {
	currentYear, currentMonth, currentDay := date.Date()
	currentLocation := date.Location()

	from = time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, currentLocation)
	to = time.Date(currentYear, currentMonth, currentDay, 23, 59, 59, 999, currentLocation)
	return
}

func MonthRange(date time.Time) (from time.Time, to time.Time) {
	currentYear, currentMonth, _ := date.Date()
	currentLocation := date.Location()

	from = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	to = from.AddDate(0, 1, -1)
	return
}

func LastWeekRange() (from time.Time, to time.Time) {
	date := time.Now()
	currentYear, currentMonth, currentDay := date.Date()
	currentLocation := date.Location()

	from = time.Date(currentYear, currentMonth, currentDay-7, 0, 0, 0, 0, currentLocation)
	to = time.Date(currentYear, currentMonth, currentDay-1, 23, 59, 59, 999, currentLocation)
	return
}

func GetUrl[T comparable](urlTemplate string, obj T) string {
	t, _ := template.New("").Parse(urlTemplate)
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, obj); err != nil {
		log.Fatalln(err)
	}

	return tpl.String()
}

func Unique[T comparable](s []T) []T {
	result := make([]T, 0)
	for _, v := range s {
		if !Contains(result, v) {
			result = append(result, v)
		}
	}
	return result
}

func Filter[T comparable](s []T, f func([]T, T) bool) []T {
	result := make([]T, 0)
	for _, v := range s {
		if f(result, v) {
			result = append(result, v)
		}
	}
	return result
}

func PrettyPrint(b []byte) []byte {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")

	if err != nil {
		log.Fatalln(err)
	}

	return out.Bytes()
}

func IntToLetters(number int32) (letters string) {
	number--
	if firstLetter := number / 26; firstLetter > 0 {
		letters += IntToLetters(firstLetter)
		letters += string('A' + number%26)
	} else {
		letters += string('A' + number)
	}

	return
}

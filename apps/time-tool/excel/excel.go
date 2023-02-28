package excel

import (
	excelize "github.com/xuri/excelize/v2"
)

var (
	err   error
	sheet = "Sheet1"
	zoom  = 85.0
	// rows height
	heights = map[int]float64{
		1:  12.6,
		2:  17.4,
		3:  13.2,
		4:  12.6,
		5:  12.6,
		6:  13.2,
		7:  15,
		8:  14,
		9:  30,
		10: 30,
		11: 30,
		12: 18.8,
	}
	// columns width
	widths = [][]interface{}{
		{"A", "A", 5.89},
		{"B", "C", 13.56},
		{"D", "D", 121.22},
		{"E", "E", 14.78},
		{"F", "F", 62.56},
	}
	// merge cells
	mergeCells = [][]string{
		{"B2", "C2"},
		{"B3", "C3"},
		{"B7", "C7"},
		{"D7", "D8"},
		{"E7", "E8"},
		{"F7", "F8"},
		{"D9", "D11"},
		{"B12", "D12"},
	}
	// result
	resultCells = []string{"E9", "E10", "E11"}
	// time cells
	timeCells = []string{"B9", "C9", "B10", "C10", "B11", "C11"}
	// text
	dailyReportText  = "Отчет за день"
	monthNames       = []string{"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь", "Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"}
	weekdayNames     = []string{"воскресенье", "понедельник", "вторник", "среда", "четверг", "пятница", "суббота"}
	dayTypeNames     = []string{"рабочий", "выходной"}
	workingHoursText = "Продолжительность рабочего дня:"
	dayText          = "День"
	startText        = "Начало"
	endText          = "Конец"
	hoursText        = "Часы"
	resultText       = "Результат"
)

func cellStyle(size float64, bold bool, italic bool, alignment []string, fill bool, border []int) *excelize.Style {
	style := &excelize.Style{
		Alignment: &excelize.Alignment{Vertical: alignment[0], Horizontal: alignment[1], WrapText: true},
		Font:      &excelize.Font{Color: "000000", Family: "Verdana", Size: size, Italic: italic, Bold: bold},
	}
	if fill {
		style.Fill = excelize.Fill{Type: "pattern", Color: []string{"#B2B2B2"}, Pattern: 1}
	}
	if bold {
		style.NumFmt = 20
	}
	if len(border) > 0 {
		style.Border = []excelize.Border{
			{Type: "top", Style: border[0], Color: "000000"},
			{Type: "right", Style: border[1], Color: "000000"},
			{Type: "bottom", Style: border[2], Color: "000000"},
			{Type: "left", Style: border[3], Color: "000000"},
		}
	}

	return style
}

package excel

import (
	"finlab/apps/time-tool/core"
	"time"

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
		{"D7", "D8"},
		{"D9", "D11"},
		{"B12", "D12"},
	}
	// styles
	top    = excelize.Border{Type: "top", Style: 1, Color: "DADEE0"}
	left   = excelize.Border{Type: "left", Style: 1, Color: "DADEE0"}
	right  = excelize.Border{Type: "right", Style: 1, Color: "DADEE0"}
	bottom = excelize.Border{Type: "bottom", Style: 1, Color: "DADEE0"}
	fill   = excelize.Fill{Type: "pattern", Color: []string{"EFEFEF"}, Pattern: 1}
	// text
	dailyReportText  = "Отчет за день"
	monthNames       = []string{"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь", "Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"}
	weekdayNames     = []string{"воскресенье", "понедельник", "вторник", "среда", "четверг", "пятница", "суббота"}
	dayTypeNames     = []string{"рабочий", "выходной"}
	workingHoursText = "Продолжительность рабочего дня:"
)

func textStyle(size float64, bold bool, italic bool, alignment string, fill bool) *excelize.Style {
	style := &excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center", Horizontal: alignment},
		Font:      &excelize.Font{Color: "000000", Family: "Verdana", Size: size, Italic: italic, Bold: bold},
	}
	if fill {
		style.Fill = excelize.Fill{Type: "pattern", Color: []string{"#E0EBF5"}, Pattern: 1}
	}

	return style
}

func SeveXlsx(time time.Time) error {
	xlsx := excelize.NewFile()

	h1Style, _ := xlsx.NewStyle(textStyle(14, true, false, "center", false))
	h2Style, _ := xlsx.NewStyle(textStyle(14, false, false, "center", false))
	h3Style, _ := xlsx.NewStyle(textStyle(10, true, false, "center", false))
	h4Style, _ := xlsx.NewStyle(textStyle(8, false, false, "center", false))
	h4IStyle, _ := xlsx.NewStyle(textStyle(8, false, true, "right", false))
	// tableHeadStyle, _ := xlsx.NewStyle(textStyle(10, true, false, "center", true))
	// tableBodyStyle, _ := xlsx.NewStyle(textStyle(10, false, false, "left", false))

	date := time.Format("02.01.2006")
	fileName := "Report_" + time.Format("20060102") + "_.xlsx"
	weekday := int(time.Weekday())
	isWeekend := weekday < 1 || weekday > 5
	dayType := 0
	workingHours := 9
	if isWeekend {
		dayType = 1
		workingHours = 8
	}

	// set custom row height
	for row, height := range heights {
		if err = xlsx.SetRowHeight(sheet, row, height); err != nil {
			return err
		}
	}

	// set custom column width
	for _, width := range widths {
		if err = xlsx.SetColWidth(sheet, width[0].(string), width[1].(string), width[2].(float64)); err != nil {
			return err
		}
	}

	// merge cells
	for _, mergeCell := range mergeCells {
		if err = xlsx.MergeCell(sheet, mergeCell[0], mergeCell[1]); err != nil {
			return err
		}
	}

	// daily report
	xlsx.SetCellStyle(sheet, "B2", "C2", h1Style)
	xlsx.SetCellValue(sheet, "B2", dailyReportText)

	// month name
	xlsx.SetCellStyle(sheet, "D2", "D2", h2Style)
	xlsx.SetCellValue(sheet, "D2", monthNames[time.Month()-1])

	// date
	xlsx.SetCellStyle(sheet, "B3", "C3", h3Style)
	xlsx.SetCellValue(sheet, "B3", date)

	// date
	xlsx.SetCellStyle(sheet, "B5", "C5", h4Style)
	xlsx.SetCellValue(sheet, "B5", weekdayNames[weekday])
	xlsx.SetCellValue(sheet, "C5", dayTypeNames[dayType])

	// working hours
	xlsx.SetCellStyle(sheet, "D5", "D5", h4IStyle)
	xlsx.SetCellValue(sheet, "D5", workingHoursText)
	xlsx.SetCellStyle(sheet, "E5", "E5", h4Style)
	xlsx.SetCellValue(sheet, "E5", workingHours)

	// Zoom
	if err := xlsx.SetSheetView(sheet, -1, &excelize.ViewOptions{
		ZoomScale: &zoom,
	}); err != nil {
		return err
	}

	// rename worksheet
	xlsx.SetSheetName(sheet, date)

	if err := xlsx.SaveAs(fileName); err != nil {
		return err
	}

	if err := xlsx.Close(); err != nil {
		return err
	}

	core.Success("Excel document is saved!\n")
	return nil
}

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
		{"B7", "C7"},
		{"D7", "D8"},
		{"E7", "E8"},
		{"F7", "F8"},
		{"D9", "D11"},
		{"B12", "D12"},
	}
	// formulas
	formulaCells = [][]string{
		{"E9", "=C9-B9"},
		{"E10", "=C10-B10"},
		{"E11", "=C11-B11"},
		{"E12", "=SUM(E9:E11)"},
	}
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
		Alignment: &excelize.Alignment{Vertical: alignment[0], Horizontal: alignment[1]},
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

func SeveXlsx(date time.Time, timestampsRes core.TimestampsRes) error {
	xlsx := excelize.NewFile()

	h1Style, _ := xlsx.NewStyle(cellStyle(14, true, false, []string{"center", "center"}, false, []int{}))
	h2Style, _ := xlsx.NewStyle(cellStyle(14, false, false, []string{"center", "center"}, false, []int{}))
	h3Style, _ := xlsx.NewStyle(cellStyle(10, true, false, []string{"center", "center"}, false, []int{}))
	h4Style, _ := xlsx.NewStyle(cellStyle(8, false, false, []string{"center", "center"}, false, []int{}))
	h4IStyle, _ := xlsx.NewStyle(cellStyle(8, false, true, []string{"center", "right"}, false, []int{}))
	tableHeadStyle, _ := xlsx.NewStyle(cellStyle(10, true, false, []string{"center", "center"}, true, []int{2, 2, 2, 2}))
	tableBodyStyle, _ := xlsx.NewStyle(cellStyle(10, true, false, []string{"center", "center"}, false, []int{1, 2, 1, 2}))
	tableBodyContentStyle, _ := xlsx.NewStyle(cellStyle(10, false, false, []string{"top", "left"}, false, []int{2, 2, 2, 2}))

	dateString := date.Format(core.DateDotTpl)
	fileName := "Report_" + date.Format(core.DateFileTpl) + "_.xlsx"
	weekday := int(date.Weekday())
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
	xlsx.SetCellValue(sheet, "D2", monthNames[date.Month()-1])

	// date
	xlsx.SetCellStyle(sheet, "B3", "C3", h3Style)
	xlsx.SetCellValue(sheet, "B3", dateString)

	// date
	xlsx.SetCellStyle(sheet, "B5", "C5", h4Style)
	xlsx.SetCellValue(sheet, "B5", weekdayNames[weekday])
	xlsx.SetCellValue(sheet, "C5", dayTypeNames[dayType])

	// working hours
	xlsx.SetCellStyle(sheet, "D5", "D5", h4IStyle)
	xlsx.SetCellValue(sheet, "D5", workingHoursText)

	xlsx.SetCellStyle(sheet, "E5", "E5", h4Style)
	xlsx.SetCellValue(sheet, "E5", workingHours)

	// table
	xlsx.SetCellStyle(sheet, "B9", "F11", tableBodyStyle)

	// table head
	xlsx.SetCellStyle(sheet, "B7", "C7", tableHeadStyle)
	xlsx.SetCellValue(sheet, "B7", dayText)

	xlsx.SetCellStyle(sheet, "B8", "B8", tableHeadStyle)
	xlsx.SetCellValue(sheet, "B8", startText)

	xlsx.SetCellStyle(sheet, "C8", "C8", tableHeadStyle)
	xlsx.SetCellValue(sheet, "C8", endText)

	xlsx.SetCellStyle(sheet, "D7", "D8", tableHeadStyle)
	xlsx.SetCellStyle(sheet, "E7", "E8", tableHeadStyle)

	xlsx.SetCellValue(sheet, "E7", hoursText)

	xlsx.SetCellStyle(sheet, "F7", "F8", tableHeadStyle)
	xlsx.SetCellValue(sheet, "F7", resultText)

	xlsx.SetCellStyle(sheet, "B12", "D12", tableHeadStyle)
	xlsx.SetCellStyle(sheet, "E12", "E12", tableHeadStyle)
	xlsx.SetCellStyle(sheet, "F12", "F12", tableHeadStyle)

	// table content
	xlsx.SetCellStyle(sheet, "D9", "D11", tableBodyContentStyle)

	for key, timestamp := range timestampsRes.Data {
		t, _ := time.Parse(core.DateISOTpl, timestamp.Timestamp)
		if err = xlsx.SetCellValue(sheet, timeCells[key], t.Format(core.TimeTpl)); err != nil {
			return err
		}
	}

	// formulas
	for _, formulaCell := range formulaCells {
		xlsx.SetCellFormula(sheet, formulaCell[0], formulaCell[1])
		xlsx.CalcCellValue(sheet, formulaCell[0])
	}

	// Zoom
	if err := xlsx.SetSheetView(sheet, -1, &excelize.ViewOptions{
		ZoomScale: &zoom,
	}); err != nil {
		return err
	}

	// rename worksheet
	xlsx.SetSheetName(sheet, dateString)

	if err := xlsx.SaveAs(fileName); err != nil {
		return err
	}

	if err := xlsx.Close(); err != nil {
		return err
	}

	core.Success("Excel document is saved!\n")
	return nil
}

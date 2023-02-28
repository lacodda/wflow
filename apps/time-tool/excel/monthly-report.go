package excel

import (
	"finlab/apps/time-tool/core"
	"fmt"
	"strings"
	"time"

	excelize "github.com/xuri/excelize/v2"
)

type CellFontStyle struct {
	Size   float64
	Bold   bool
	Italic bool
	Color  string
	Family string
}

type CellStyle struct {
	Font      CellFontStyle
	Alignment []string
	Fill      string
	Border    []int
	NumFmt    int
}

func GetCellStyle(s CellStyle) *excelize.Style {
	style := &excelize.Style{
		Alignment: &excelize.Alignment{Vertical: s.Alignment[0], Horizontal: s.Alignment[1], WrapText: true},
		Font:      &excelize.Font{Color: s.Font.Color, Family: s.Font.Family, Size: s.Font.Size, Italic: s.Font.Italic, Bold: s.Font.Bold},
	}
	if len(s.Fill) > 0 {
		style.Fill = excelize.Fill{Type: "pattern", Color: []string{s.Fill}, Pattern: 1}
	}
	if len(s.Border) > 0 {
		style.Border = []excelize.Border{
			{Type: "top", Style: s.Border[0], Color: "000000"},
			{Type: "right", Style: s.Border[1], Color: "000000"},
			{Type: "bottom", Style: s.Border[2], Color: "000000"},
			{Type: "left", Style: s.Border[3], Color: "000000"},
		}
	}

	return style
}

func getCell(colNum int32) (letter string, header string, body string) {
	letter = core.IntToLetters(colNum)
	header = letter + "1"
	body = letter + "2"
	return
}

func SeveXLSXMonthlyReport(from time.Time, to time.Time, summaryRes core.SummaryRes) (fileName string, err error) {
	zoom = 70.0
	monthName := monthNames[from.Month()-1]
	monthYearText := fmt.Sprintf("%v %v", strings.ToLower(monthName), from.Year())
	totalTimeText := "Итого:"
	totalWorkDaysText := "Раб:"

	fileName = "Report_" + monthName + ".xlsx"
	xlsx := excelize.NewFile()
	xlsx.SetDocProps(&excelize.DocProperties{Creator: "Kirill Lahtachev"})
	var basicCellStyle = CellStyle{Font: CellFontStyle{Color: "000000", Family: "Arial", Size: 10}, Alignment: []string{"center", "center"}, Border: []int{1, 1, 1, 1}}
	var headerCellStyle = basicCellStyle
	headerCellStyle.Fill = "#FF9900"
	headerCellStyle.Font.Bold = true

	var headerTotalCellStyle = headerCellStyle
	headerTotalCellStyle.Font.Color = "FF0000"
	var headerTotalDaysCellStyle = headerCellStyle
	headerTotalDaysCellStyle.Font.Color = "0000FF"

	var bodyCellStyle = basicCellStyle
	bodyCellStyle.NumFmt = 20
	var bodyTotalCellStyle = bodyCellStyle
	bodyTotalCellStyle.Font.Color = "FF0000"
	bodyTotalCellStyle.Font.Bold = true
	var bodyTotal2CellStyle = bodyTotalCellStyle
	bodyTotal2CellStyle.NumFmt = 3
	var bodyTotalDaysCellStyle = bodyTotal2CellStyle
	bodyTotalDaysCellStyle.Font.Color = "0000FF"
	var weekendCellStyle = bodyCellStyle
	weekendCellStyle.Fill = "FF99CC"

	basicStyle, _ := xlsx.NewStyle(GetCellStyle(basicCellStyle))
	headerStyle, _ := xlsx.NewStyle(GetCellStyle(headerCellStyle))
	headerTotalStyle, _ := xlsx.NewStyle(GetCellStyle(headerTotalCellStyle))
	headerTotalDaysStyle, _ := xlsx.NewStyle(GetCellStyle(headerTotalDaysCellStyle))
	bodyTotalStyle, _ := xlsx.NewStyle(GetCellStyle(bodyTotalCellStyle))
	bodyTotal2Style, _ := xlsx.NewStyle(GetCellStyle(bodyTotal2CellStyle))
	bodyTotalDaysStyle, _ := xlsx.NewStyle(GetCellStyle(bodyTotalDaysCellStyle))
	weekendStyle, _ := xlsx.NewStyle(GetCellStyle(weekendCellStyle))

	cells := [][]interface{}{
		{22.11, headerStyle, monthYearText, basicStyle, nil},
		{6.11, headerStyle, nil, []int{basicStyle, weekendStyle}, nil},
		{12.5, headerTotalStyle, totalTimeText, bodyTotalStyle, core.MinutesToTimeStr(summaryRes.TotalTime)},
		{8.89, headerTotalStyle, totalTimeText, bodyTotal2Style, summaryRes.TotalTime * 24},
		{8.89, headerTotalDaysStyle, totalWorkDaysText, bodyTotalDaysStyle, summaryRes.TotalTime * 24 / 8},
	}

	var (
		colNum int32  = 1
		letter string = ""
		header string = ""
		body   string = ""
	)
	for _, cell := range cells {
		if cell[2] != nil {
			letter, header, body = getCell(colNum)
			xlsx.SetColWidth(sheet, letter, letter, cell[0].(float64))
			xlsx.SetCellStyle(sheet, header, header, cell[1].(int))
			xlsx.SetCellStyle(sheet, body, body, cell[3].(int))
			xlsx.SetCellValue(sheet, header, cell[2].(string))
			if cell[4] != nil {
				xlsx.SetCellValue(sheet, body, cell[4])
			}
			colNum++
		} else {
			for day := from; !day.After(to); day = day.AddDate(0, 0, 1) {
				letter, header, body = getCell(colNum)
				xlsx.SetColWidth(sheet, letter, letter, cell[0].(float64))
				xlsx.SetCellStyle(sheet, header, header, cell[1].(int))
				xlsx.SetCellStyle(sheet, body, body, cell[3].([]int)[0])
				xlsx.SetCellValue(sheet, header, day.Day())
				if day.Weekday() == time.Saturday || day.Weekday() == time.Sunday {
					xlsx.SetCellStyle(sheet, body, body, cell[3].([]int)[1])
				}
				for _, summary := range summaryRes.Data {
					date, _ := time.Parse(core.DateISOTpl, summary.Date)
					if day.YearDay() == date.YearDay() {
						xlsx.SetCellValue(sheet, body, core.MinutesToTimeStr(summary.Time))
					}
				}
				colNum++
			}
		}
	}

	// Zoom
	if err = xlsx.SetSheetView(sheet, -1, &excelize.ViewOptions{
		ZoomScale: &zoom,
	}); err != nil {
		return
	}

	// rename worksheet
	xlsx.SetSheetName(sheet, monthYearText)

	if err = xlsx.SaveAs(fileName); err != nil {
		return
	}

	if err = xlsx.Close(); err != nil {
		return
	}

	core.Success("Excel document is saved!\n")
	return fileName, nil
}

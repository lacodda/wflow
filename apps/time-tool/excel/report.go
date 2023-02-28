package excel

import (
	"finlab/apps/time-tool/core"
	"fmt"
	"strings"
	"time"

	excelize "github.com/xuri/excelize/v2"
)

func SeveXLSXReport(date time.Time, timestampsRes core.TimestampsRes, tasksRes core.TasksRes) (fileName string, err error) {
	fileName = "Report_" + date.Format(core.DateFileTpl) + ".xlsx"
	xlsx := excelize.NewFile()
	xlsx.SetDocProps(&excelize.DocProperties{Creator: "Kirill Lahtachev"})

	h1Style, _ := xlsx.NewStyle(cellStyle(14, true, false, []string{"center", "center"}, false, []int{}))
	h2Style, _ := xlsx.NewStyle(cellStyle(14, false, false, []string{"center", "center"}, false, []int{}))
	h3Style, _ := xlsx.NewStyle(cellStyle(10, true, false, []string{"center", "center"}, false, []int{}))
	h4Style, _ := xlsx.NewStyle(cellStyle(8, false, false, []string{"center", "center"}, false, []int{}))
	h4IStyle, _ := xlsx.NewStyle(cellStyle(8, false, true, []string{"center", "right"}, false, []int{}))
	tableHeadStyle, _ := xlsx.NewStyle(cellStyle(10, true, false, []string{"center", "center"}, true, []int{2, 2, 2, 2}))
	tableBodyStyle, _ := xlsx.NewStyle(cellStyle(10, true, false, []string{"center", "center"}, false, []int{1, 2, 1, 2}))
	tableBodyContentStyle, _ := xlsx.NewStyle(cellStyle(10, false, false, []string{"top", "left"}, false, []int{2, 2, 2, 2}))

	dateString := date.Format(core.DateDotTpl)
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
			return
		}
	}

	// set custom column width
	for _, width := range widths {
		if err = xlsx.SetColWidth(sheet, width[0].(string), width[1].(string), width[2].(float64)); err != nil {
			return
		}
	}

	// merge cells
	for _, mergeCell := range mergeCells {
		if err = xlsx.MergeCell(sheet, mergeCell[0], mergeCell[1]); err != nil {
			return
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
			return
		}
	}

	var tasks []string
	for _, task := range tasksRes.Data {
		tasks = append(tasks, fmt.Sprintf("* %s (%v%%)", task.Name, task.Completeness))
	}
	tasksString := strings.Join(tasks, "\n")
	xlsx.SetCellValue(sheet, "D9", tasksString)

	// result
	for key, minutes := range timestampsRes.WorkTime {
		if err = xlsx.SetCellValue(sheet, resultCells[key], core.MinutesToTimeStr(minutes)); err != nil {
			return
		}
	}
	xlsx.SetCellValue(sheet, "E12", core.MinutesToTimeStr(timestampsRes.TotalTime))

	// Zoom
	if err = xlsx.SetSheetView(sheet, -1, &excelize.ViewOptions{
		ZoomScale: &zoom,
	}); err != nil {
		return
	}

	// rename worksheet
	xlsx.SetSheetName(sheet, dateString)

	if err = xlsx.SaveAs(fileName); err != nil {
		return
	}

	if err = xlsx.Close(); err != nil {
		return
	}

	core.Success("Excel document is saved!\n")
	return fileName, nil
}

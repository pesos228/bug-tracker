package excel

import (
	"bytes"
	"fmt"
	"time"

	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/service"
	"github.com/xuri/excelize/v2"
)

var statusTranslations = map[domain.CheckStatus]string{
	domain.NotChecked:       "не проверено",
	domain.Checked:          "проверено",
	domain.PartiallyChecked: "проверено частично",
	domain.Failed:           "провалено",
}

var resultTranslations = map[domain.CheckResult]string{
	domain.Success: "успешно",
	domain.Failure: "неуспешно",
	domain.Warning: "есть замечания",
}

type reportGenerator struct {
}

func (r *reportGenerator) Generate(tasks []*service.TaskReportRow) (*bytes.Buffer, error) {
	file := excelize.NewFile()
	defer file.Close()

	sheetName := "Отчет по таскам"
	file.SetSheetName("Sheet1", sheetName)

	styles, err := r.createStyles(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create styles: %w", err)
	}

	r.setupSheet(file, sheetName, styles)

	r.fillData(file, sheetName, tasks, styles)

	buffer := &bytes.Buffer{}
	if err := file.Write(buffer); err != nil {
		return nil, fmt.Errorf("failed to write in buffer: %w", err)
	}

	return buffer, nil
}

func (r *reportGenerator) createStyles(file *excelize.File) (map[string]int, error) {
	styles := make(map[string]int)

	headerStyle, err := file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#D3D3D3"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold: true,
			Size: 10,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: r.createBorder(),
	})
	if err != nil {
		return nil, err
	}
	styles["header"] = headerStyle

	defaultStyle, err := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: r.createBorder(),
	})
	if err != nil {
		return nil, err
	}
	styles["default"] = defaultStyle

	dateFormat := "dd.mm.yyyy"
	dateStyle, err := file.NewStyle(&excelize.Style{
		CustomNumFmt: &dateFormat,
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: r.createBorder(),
	})
	if err != nil {
		return nil, err
	}
	styles["date"] = dateStyle

	resultStyles := map[string]string{
		"success": "#4CAF50",
		"failure": "#F44336",
		"warning": "#FF9800",
		"other":   "#E0E0E0",
	}

	for name, color := range resultStyles {
		fontColor := "FFFFFF"
		if name == "warning" || name == "other" {
			fontColor = "000000"
		}

		style, err := file.NewStyle(&excelize.Style{
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{color},
				Pattern: 1,
			},
			Font: &excelize.Font{
				Color: fontColor,
				Bold:  true,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
				WrapText:   true,
			},
			Border: r.createBorder(),
		})
		if err != nil {
			return nil, err
		}
		styles[name] = style
	}

	return styles, nil
}

func (r *reportGenerator) setupSheet(file *excelize.File, sheetName string, styles map[string]int) {
	headers := []string{
		"ПО",
		"Номер заявки\nРазработчик/\nММ",
		"Описание задачи",
		"Ответственный от КБР",
		"Дата\nобновления\nтестовой\nсреды",
		"Дата\nпроверки\nФакт",
		"Статус проверки",
		"Результат\nпроверки",
		"Комментарий к тестированию",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		file.SetCellValue(sheetName, cell, header)
		file.SetCellStyle(sheetName, cell, cell, styles["header"])
	}

	file.SetRowHeight(sheetName, 1, 60)

	columnWidths := []float64{25, 20, 40, 20, 20, 15, 15, 15, 35}
	for i, width := range columnWidths {
		col := string(rune('A' + i))
		file.SetColWidth(sheetName, col, col, width)
	}
}

func (r *reportGenerator) fillData(file *excelize.File, sheetName string, tasks []*service.TaskReportRow, styles map[string]int) {
	for i, task := range tasks {
		row := i + 2
		r.writeTaskRow(file, sheetName, row, task, styles)
		file.SetRowHeight(sheetName, row, 40)
	}
}

func (r *reportGenerator) writeTaskRow(file *excelize.File, sheetName string, row int, task *service.TaskReportRow, styles map[string]int) {
	data := []interface{}{
		task.SoftName,
		task.RequestID,
		task.Description,
		task.AssigneePerson,
		task.TestEnvDateUpdate,
		task.CheckDate,
		r.getStatusDisplay(task.CheckStatus),
		r.getResultDisplay(task.CheckResult),
		task.Comment,
	}

	for col, value := range data {
		cell, _ := excelize.CoordinatesToCellName(col+1, row)
		style := styles["default"]

		switch col {
		case 4, 5:
			style = styles["date"]
			if t, ok := value.(time.Time); ok && t.IsZero() {
				value = nil
			}
		case 7:
			style = r.getResultStyle(task.CheckResult, styles)
		}

		file.SetCellValue(sheetName, cell, value)
		file.SetCellStyle(sheetName, cell, cell, style)
	}
}

func (r *reportGenerator) createBorder() []excelize.Border {
	return []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}
}

func (r *reportGenerator) getStatusDisplay(status domain.CheckStatus) string {
	if translation, ok := statusTranslations[status]; ok {
		return translation
	}
	return string(status)
}

func (r *reportGenerator) getResultStyle(result domain.CheckResult, styles map[string]int) int {
	switch result {
	case domain.Success:
		return styles["success"]
	case domain.Failure:
		return styles["failure"]
	case domain.Warning:
		return styles["warning"]
	default:
		return styles["other"]
	}
}

func (r *reportGenerator) getResultDisplay(result domain.CheckResult) string {
	if translation, ok := resultTranslations[result]; ok {
		return translation
	}
	if result == "" {
		return "-"
	}
	return string(result)
}

func NewReportGenerator() service.ReportGenerator {
	return &reportGenerator{}
}

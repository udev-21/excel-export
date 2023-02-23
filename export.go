package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type Sheet struct {
	ID     int
	styles map[string]excelize.Style
}

type Session struct {
	e            *excelize.File
	sheets       map[string]Sheet
	LastTimeUsed time.Time
}

func NewSession() string {
	newSession := Session{
		sheets:       make(map[string]Sheet),
		LastTimeUsed: time.Now(),
	}
	newFile := excelize.NewFile()
	newSession.e = newFile

	newSessionUUID := uuid.New().String()
	if _, ok := currentSessions[newSessionUUID]; ok {
		newSessionUUID = uuid.New().String()
	} else {
		currentSessions[newSessionUUID] = newSession
	}
	return newSessionUUID
}

func (s *Session) NewSheet(name string) (int, error) {
	if len(s.sheets) == 0 {
		idx := s.e.GetActiveSheetIndex()
		s.e.SetSheetName("Sheet1", name)
		s.sheets[name] = Sheet{
			ID:     idx,
			styles: make(map[string]excelize.Style),
		}
		return idx, nil
	}
	if _, ok := s.sheets[name]; ok {
		return 0, errors.New("sheet already exists")
	}
	s.sheets[name] = Sheet{
		ID:     s.e.NewSheet(name),
		styles: make(map[string]excelize.Style),
	}
	return s.sheets[name].ID, nil
}

func (s *Session) Save(fileName string) error {
	//commit cell styles before saving
	for sheetName, sheet := range s.sheets {
		for cell, style := range sheet.styles {
			// create new style inside file
			styleID, err := s.e.NewStyle(&style)
			if err != nil {
				return err
			}
			// log.Println(sheetName, cell, styleID, style.Font.Bold)
			err = s.e.SetCellStyle(sheetName, cell, cell, styleID)
			if err != nil {
				return err
			}
		}
	}

	return s.e.SaveAs(fileName)
}

func (s *Session) SetCellValue(sheetName, cell string, val interface{}) error {
	return s.e.SetCellValue(sheetName, cell, val)
}

func (s *Session) MergeCell(sheetName, hCell, vCell string) error {
	return s.e.MergeCell(sheetName, hCell, vCell)
}

func (s *Session) BoldCell(sheetName, cell string, bold bool) error {
	//check if sheet exists
	if sheet, ok := s.sheets[sheetName]; !ok {
		return errors.New("sheet doesn't exists")
	} else {
		if cellStyle, ok := sheet.styles[cell]; ok {
			if cellStyle.Font == nil {
				cellStyle.Font = &excelize.Font{
					Bold: bold,
				}
			} else {
				cellStyle.Font.Bold = bold
			}
			sheet.styles[cell] = cellStyle
		} else {
			cellStyle = excelize.Style{
				Font: &excelize.Font{
					Bold: bold,
				},
			}

			sheet.styles[cell] = cellStyle
		}
		s.sheets[sheetName] = sheet
	}
	return nil
}

func (s *Session) CenterCell(sheetName, cell string, h, v bool) error {
	horizontal := ""
	if h {
		horizontal = "center"
	}
	vertical := ""
	if v {
		vertical = "center"
	}

	//check if sheet exists
	if sheet, ok := s.sheets[sheetName]; !ok {
		return errors.New("sheet doesn't exists")
	} else {
		if cellStyle, ok := sheet.styles[cell]; ok {
			if cellStyle.Alignment == nil {
				cellStyle.Alignment = &excelize.Alignment{
					Horizontal: horizontal,
					Vertical:   vertical,
				}
			} else {
				cellStyle.Alignment.Horizontal = horizontal
				cellStyle.Alignment.Vertical = "center"
			}
			sheet.styles[cell] = cellStyle
		} else {
			cellStyle = excelize.Style{
				Alignment: &excelize.Alignment{
					Horizontal: horizontal,
					Vertical:   vertical,
				},
			}
			sheet.styles[cell] = cellStyle
		}
		s.sheets[sheetName] = sheet
	}
	return nil
}

func (s *Session) ItalicCell(sheetName, cell string, italic bool) error {
	//check if sheet exists
	if sheet, ok := s.sheets[sheetName]; !ok {
		return errors.New("sheet doesn't exists")
	} else {
		if cellStyle, ok := sheet.styles[cell]; ok {
			if cellStyle.Font == nil {
				cellStyle.Font = &excelize.Font{
					Italic: italic,
				}
			} else {
				cellStyle.Font.Italic = italic
			}
			sheet.styles[cell] = cellStyle
		} else {
			cellStyle = excelize.Style{
				Font: &excelize.Font{
					Italic: italic,
				},
			}

			sheet.styles[cell] = cellStyle
		}
		s.sheets[sheetName] = sheet
	}
	return nil
}

func (s *Session) SetCellColor(sheetName, cell, color string) error {
	//check if sheet exists
	if sheet, ok := s.sheets[sheetName]; !ok {
		return errors.New("sheet doesn't exists")
	} else {
		if cellStyle, ok := sheet.styles[cell]; ok {
			if cellStyle.Font == nil {
				cellStyle.Font = &excelize.Font{
					Color: color,
				}
			} else {
				cellStyle.Font.Color = color
			}
			sheet.styles[cell] = cellStyle
		} else {
			cellStyle = excelize.Style{
				Font: &excelize.Font{
					Color: color,
				},
			}
			sheet.styles[cell] = cellStyle
		}
		s.sheets[sheetName] = sheet
	}
	return nil
}

func (s *Session) SetCellFontSize(sheetName, cell string, size float64) error {
	//check if sheet exists
	if sheet, ok := s.sheets[sheetName]; !ok {
		return errors.New("sheet doesn't exists")
	} else {
		if cellStyle, ok := sheet.styles[cell]; ok {
			if cellStyle.Font == nil {
				cellStyle.Font = &excelize.Font{
					Size: size,
				}
			} else {
				cellStyle.Font.Size = size
			}
			sheet.styles[cell] = cellStyle
		} else {
			cellStyle = excelize.Style{
				Font: &excelize.Font{
					Size: size,
				},
			}
			sheet.styles[cell] = cellStyle
		}
		s.sheets[sheetName] = sheet
	}
	return nil
}

func (s *Session) SetCellBorder(sheetName, cell, color string, style int) error {
	//check if sheet exists
	if sheet, ok := s.sheets[sheetName]; !ok {
		return errors.New("sheet doesn't exists")
	} else {
		if cellStyle, ok := sheet.styles[cell]; ok {
			cellStyle.Border = []excelize.Border{
				{
					Type:  "left",
					Color: color,
					Style: style,
				},
				{
					Type:  "top",
					Color: color,
					Style: style,
				},
				{
					Type:  "bottom",
					Color: color,
					Style: style,
				},
				{
					Type:  "right",
					Color: color,
					Style: style,
				},
			}
			sheet.styles[cell] = cellStyle
		} else {
			cellStyle = excelize.Style{
				Border: []excelize.Border{
					{
						Type:  "left",
						Color: color,
						Style: style,
					},
					{
						Type:  "top",
						Color: color,
						Style: style,
					},
					{
						Type:  "bottom",
						Color: color,
						Style: style,
					},
					{
						Type:  "right",
						Color: color,
						Style: style,
					},
				},
			}
			sheet.styles[cell] = cellStyle
		}
		s.sheets[sheetName] = sheet
	}
	return nil
}

func (s *Session) SetColWidth(sheetName, column string, width float64) error {
	//check if sheet exists
	if _, ok := s.sheets[sheetName]; !ok {
		return errors.New("sheet doesn't exists")
	} else {
		return s.e.SetColWidth(sheetName, column, column, width)
	}
}

func (s *Session) SetRowHeight(sheetName string, row int, height float64) error {
	//check if sheet exists
	if _, ok := s.sheets[sheetName]; !ok {
		return errors.New("sheet doesn't exists")
	} else {
		return s.e.SetRowHeight(sheetName, row, height)
	}
}

func (s *Session) NewDefinedName(sheetName, name, vCell, hCell, scopeSheetName string) error {
	return s.e.SetDefinedName(&excelize.DefinedName{
		Name:     name,
		RefersTo: fmt.Sprintf("%s!%s:%s", sheetName, vCell, hCell),
		Scope:    scopeSheetName,
	})
}

func (s *Session) NewDataValidation(sheetName, definedName, hCell, vCell string) error {
	dvRange := excelize.NewDataValidation(true)
	dvRange.Sqref = fmt.Sprintf("%s:%s", hCell, vCell)
	dvRange.SetSqrefDropList(definedName)
	return s.e.AddDataValidation(sheetName, dvRange)
}

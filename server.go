package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
)

func RunServer(port, filesPort string, debug bool) {
	var router = new(bunrouter.Router)
	if debug {
		router = bunrouter.New(
			bunrouter.Use(reqlog.NewMiddleware()),
		)
	} else {
		router = bunrouter.New()
	}

	router.POST("/file", CreateFile)
	router.WithMiddleware(FileExistsMiddleware).WithGroup("/:fileID", func(g *bunrouter.Group) {
		g.POST("/sheet", CreateSheet)
		g.POST("/save", SaveFile)
		g.WithMiddleware(SheetExistsMiddleware).WithGroup("/:sheetName", func(g1 *bunrouter.Group) {
			g1.POST("/setCellValue", SetCellValue)
			g1.POST("/bulkSetCellValue", BulkSetCellValue)
			g1.POST("/mergeCell", MergeCell)
			g1.POST("/boldCell", BoldCell)
			g1.POST("/italicCell", ItalicCell)
			g1.POST("/centerCell", CenterCell)
			g1.POST("/setColWidth", SetColWidth)
			g1.POST("/setRowHeight", SetRowHeight)
			g1.POST("/setCellColor", SetCellColor)
			g1.POST("/definedName", SetDefinedName)
			g1.POST("/dataValidation", SetDataValidation)
			g1.POST("/setCellFontSize", SetCellFontSize)
			g1.POST("/setCellBorder", SetCellBorder)
			g1.POST("/setCellCenter", SetCellCenter)

		})
	})

	fs := http.FileServer(http.Dir("./outputs"))
	http.Handle("/", fs)
	go http.ListenAndServe(":"+port, router)
	http.ListenAndServe(":"+filesPort, nil)

}

type CreateSheetInput struct {
	Name string `json:"name"`
}

func CreateFile(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()

	fileID := NewSession()
	w.WriteHeader(http.StatusCreated)

	return json.NewEncoder(w).Encode(map[string]string{
		"fileID": fileID,
	})
}

func CreateSheet(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()
	input := CreateSheetInput{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	fileID := req.Params().ByName("fileID")

	file := currentSessions[fileID]
	if _, err := file.NewSheet(input.Name); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	currentSessions[fileID] = file

	return nil
}

func FileExistsMiddleware(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		fileID := req.Params().ByName("fileID")
		if _, ok := currentSessions[fileID]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		} else {
			session := currentSessions[fileID]
			session.LastTimeUsed = time.Now()
			currentSessions[fileID] = session
			return next(w, req)
		}
	}
}

func SheetExistsMiddleware(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		fileID := req.Params().ByName("fileID")
		if file, ok := currentSessions[fileID]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		} else {
			sheetName := req.Params().ByName("sheetName")
			if _, ok := file.sheets[sheetName]; !ok {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			} else {
				session := currentSessions[fileID]
				session.LastTimeUsed = time.Now()
				currentSessions[fileID] = session
				return next(w, req)
			}
		}
	}
}

type SetCellColorInput struct {
	Cell  string `json:"cell"`
	Color string `json:"color"`
}

func SetCellColor(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()
	input := SetCellColorInput{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	fileID := req.Params().ByName("fileID")
	file := currentSessions[fileID]
	sheetName := req.Params().ByName("sheetName")
	if err := file.SetCellColor(sheetName, input.Cell, input.Color); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return err
	}
	currentSessions[fileID] = file
	return nil
}

type SetCellFontSizeInput struct {
	Cell string  `json:"cell"`
	Size float64 `json:"size"`
}

func SetCellFontSize(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()
	input := SetCellFontSizeInput{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	fileID := req.Params().ByName("fileID")
	file := currentSessions[fileID]
	sheetName := req.Params().ByName("sheetName")
	if err := file.SetCellFontSize(sheetName, input.Cell, input.Size); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return err
	}
	currentSessions[fileID] = file
	return nil
}

type SetCellBorderInput struct {
	Cell    string `json:"cell"`
	Color   string `json:"color"`
	StyleID int    `json:"styleID"`
}

func SetCellBorder(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()
	input := SetCellBorderInput{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	fileID := req.Params().ByName("fileID")
	file := currentSessions[fileID]
	sheetName := req.Params().ByName("sheetName")
	if err := file.SetCellBorder(sheetName, input.Cell, input.Color, input.StyleID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return err
	}
	currentSessions[fileID] = file
	return nil
}

type SetCellValueInput struct {
	Cell  string      `json:"cell"`
	Value interface{} `json:"value"`
}

func SetCellValue(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()
	input := SetCellValueInput{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	fileID := req.Params().ByName("fileID")
	file := currentSessions[fileID]
	sheetName := req.Params().ByName("sheetName")
	if err := file.SetCellValue(sheetName, input.Cell, input.Value); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return err
	}
	currentSessions[fileID] = file
	return nil
}

func BulkSetCellValue(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()
	input := map[string]interface{}{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	fileID := req.Params().ByName("fileID")
	file := currentSessions[fileID]
	sheetName := req.Params().ByName("sheetName")
	for cell, v := range input {
		if err := file.SetCellValue(sheetName, cell, v); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return err
		}
	}
	currentSessions[fileID] = file
	return nil
}

func SaveFile(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()
	fileID := req.Params().ByName("fileID")
	file := currentSessions[fileID]
	if err := file.Save(fmt.Sprintf("outputs/%s.xlsx", fileID)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't save file:" + err.Error()))
		return err
	}
	// purge the file from memory
	delete(currentSessions, fileID)

	return nil
}

type MergeCellInput struct {
	HCell string `json:"hCell"`
	VCell string `json:"vCell"`
}

func MergeCell(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()
	input := MergeCellInput{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	fileID := req.Params().ByName("fileID")
	file := currentSessions[fileID]
	sheetName := req.Params().ByName("sheetName")
	if err := file.MergeCell(sheetName, input.HCell, input.VCell); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't merge cells:" + err.Error()))
		return err
	}
	currentSessions[fileID] = file
	return nil
}

type BoldCellInput struct {
	Cell string `json:"cell"`
	Bold bool   `json:"bold"`
}

func BoldCell(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()
	input := BoldCellInput{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	fileID := req.Params().ByName("fileID")
	file := currentSessions[fileID]
	sheetName := req.Params().ByName("sheetName")
	if err := file.BoldCell(sheetName, input.Cell, input.Bold); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't bold cell:" + err.Error()))
		return err
	}
	currentSessions[fileID] = file
	return nil
}

type ItalicCellInput struct {
	Cell   string `json:"cell"`
	Italic bool   `json:"italic"`
}

func ItalicCell(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()
	input := ItalicCellInput{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	fileID := req.Params().ByName("fileID")
	file := currentSessions[fileID]
	sheetName := req.Params().ByName("sheetName")
	if err := file.ItalicCell(sheetName, input.Cell, input.Italic); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't italize cell:" + err.Error()))
		return err
	}
	currentSessions[fileID] = file
	return nil
}

type CenterCellInput struct {
	Cell       string `json:"cell"`
	Vertical   bool   `json:"vertical"`
	Horizontal bool   `json:"horizontal"`
}

func CenterCell(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()
	input := CenterCellInput{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	fileID := req.Params().ByName("fileID")
	file := currentSessions[fileID]
	sheetName := req.Params().ByName("sheetName")
	if err := file.CenterCell(sheetName, input.Cell, input.Horizontal, input.Vertical); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't center cells:" + err.Error()))
		return err
	}
	currentSessions[fileID] = file
	return nil
}

type SetColWidthInput struct {
	Column string  `json:"column"`
	Width  float64 `json:"width"`
}

func SetColWidth(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()
	input := SetColWidthInput{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	fileID := req.Params().ByName("fileID")
	file := currentSessions[fileID]
	sheetName := req.Params().ByName("sheetName")
	if err := file.SetColWidth(sheetName, input.Column, input.Width); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't center cells:" + err.Error()))
		return err
	}
	currentSessions[fileID] = file
	return nil
}

type SetRowHeightInput struct {
	Row    int     `json:"row"`
	Height float64 `json:"height"`
}

func SetRowHeight(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()
	input := SetRowHeightInput{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	fileID := req.Params().ByName("fileID")
	file := currentSessions[fileID]
	sheetName := req.Params().ByName("sheetName")
	if err := file.SetRowHeight(sheetName, input.Row, input.Height); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't center cells:" + err.Error()))
		return err
	}
	currentSessions[fileID] = file
	return nil
}

type SetCellCenterinput struct {
	Cell       string `json:"cell"`
	Horizontal bool   `json:"horizontal"`
	Vertical   bool   `json:"vertical"`
}

func SetCellCenter(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()
	input := SetCellCenterinput{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	fileID := req.Params().ByName("fileID")
	file := currentSessions[fileID]
	sheetName := req.Params().ByName("sheetName")
	if err := file.CenterCell(sheetName, input.Cell, input.Horizontal, input.Vertical); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't center cells:" + err.Error()))
		return err
	}
	currentSessions[fileID] = file
	return nil
}

type SetDefinedNameInput struct {
	Name           string `json:"name"`
	VCell          string `json:"vCell"`
	HCell          string `json:"hCell"`
	ScopeSheetName string `json:"scopeSheetName"`
}

func SetDefinedName(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()
	input := SetDefinedNameInput{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	fileID := req.Params().ByName("fileID")
	file := currentSessions[fileID]
	sheetName := req.Params().ByName("sheetName")
	if err := file.NewDefinedName(sheetName, input.Name, input.VCell, input.HCell, input.ScopeSheetName); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't center cells:" + err.Error()))
		return err
	}
	currentSessions[fileID] = file
	return nil
}

type SetDataValidationInput struct {
	HCell       string `json:"hCell"`
	VCell       string `json:"vCell"`
	DefinedName string `json:"definedName"`
}

func SetDataValidation(w http.ResponseWriter, req bunrouter.Request) error {
	defer req.Body.Close()
	input := SetDataValidationInput{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	fileID := req.Params().ByName("fileID")
	file := currentSessions[fileID]
	sheetName := req.Params().ByName("sheetName")
	if err := file.NewDataValidation(sheetName, input.DefinedName, input.HCell, input.VCell); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't center cells:" + err.Error()))
		return err
	}
	currentSessions[fileID] = file
	return nil
}

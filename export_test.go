package main_test

import (
	excelexport "productbank/excel-export"
	"testing"
)

func TestNewSession(t *testing.T) {
	_ = excelexport.NewSession()
}

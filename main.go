package main

var currentSessions = make(map[string]Session)

func main() {
	// session1uuid := NewSession()

	// session1 := currentSessions[session1uuid]

	// session1.NewSheet("Sharof")
	// session1.NewSheet("Sarvar")

	// session1.SetCellValue("Sharof", "A1", "something new")
	// session1.SetCellValue("Sharof", "A2", "something new2")
	// session1.SetCellValue("Sharof", "A3", "something new3")
	// session1.SetCellValue("Sharof", "A4", "something new4")
	// session1.SetCellValue("Sharof", "A5", "something new5")
	// session1.SetCellValue("Sharof", "A6", "something new6")

	// session1.SetCellValue("Sharof", "b1", 123.1)
	// session1.SetCellValue("Sarvar", "b1", 123.1)
	// session1.MergeCell("Sharof", "C2", "G10")

	// session1.SetCellValue("Sharof", "C3", "Something new inside merged cells")
	// session1.BoldCell("Sharof", "A1")
	// session1.CenterCell("Sharof", "C2")
	// session1.BoldCell("Sharof", "C2")
	// session1.ItalicCell("Sharof", "C2")
	// session1.SetColWidth("Sharof", "A", 25)
	// session1.SetRowHeight("Sharof", 2, 10)
	// session1.NewDefinedName("Sharof", "somethings", "$A$1", "$A$6", "Sharof")
	// // session1.BoldCell2("Sharof", "C2", "C2")

	// if err := session1.NewDataValidation("Sharof", "somethings", "H3"); err != nil {
	// 	panic(err)
	// }

	// if err := session1.Save("excel.xlsx"); err != nil {
	// 	panic(err)
	// }

	// log.Println("done")
	RunServer("8080")
}

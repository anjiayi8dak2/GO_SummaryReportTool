package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/sqltocsv"
)

func main() {
	//// Create the database handle, confirm driver is present
	//db, err := sql.Open("mysql", "moves:moves@/")
	//defer db.Close()
	//if err != nil {
	//	fmt.Println("Failed to connect MariaDB", err)
	//	return
	//}
	db := initDb()

	// Connect and check the server version
	getDBVersion(db)

	//Initialize
	a := app.New()
	w := a.NewWindow("Summary Report Tool")
	w.Resize(fyne.NewSize(400, 400))
	//Global var
	var dbSelection string
	var tableSelection string
	//var sqlResult []Movesoutput

	//menu bar
	menuitemMaria := fyne.NewMenuItem("Open Maria Data Folder", func() {
		openMariaFolder(db)
	}) // ignore functions
	//menuitemRefresh := fyne.NewMenuItem("Refresh (F5)", buttonSubmit(db, dbSelection, tableSelection)) // ignore functions
	menuitemRefresh := fyne.NewMenuItem("Refresh (F5)", func() {
		whiteList := getWhiteList(db, dbSelection, tableSelection)

		fmt.Println("printing white list")
		fmt.Printf("%v", &whiteList)

		dummy, _ := getQueryResult(db, dbSelection, tableSelection, whiteList)
		fmt.Println("printing dummy")
		fmt.Printf("%v", &dummy)
	})
	menuitemOpenlog := fyne.NewMenuItem("Open Log", nil)   // ignore functions
	menuitemClearlog := fyne.NewMenuItem("Clear Log", nil) // ignore functions
	menuitemAbout := fyne.NewMenuItem("About", nil)        // ignore functions
	menuitemManual := fyne.NewMenuItem("Manual", nil)      // ignore functions
	// New Menu
	newMenu1 := fyne.NewMenu("File", menuitemMaria)
	newMenu2 := fyne.NewMenu("Edit", menuitemRefresh)
	newMenu3 := fyne.NewMenu("Logs", menuitemOpenlog, menuitemClearlog)
	newMenu4 := fyne.NewMenu("Help", menuitemAbout, menuitemManual)
	// New main menu
	menu := fyne.NewMainMenu(newMenu1, newMenu2, newMenu3, newMenu4)
	w.SetMainMenu(menu)

	// Get DB List
	var dbList []string
	dbList = getDBList(db)

	//Create dropdown for database selection
	dbDropdownResult := widget.NewLabel("Select a Database")
	//Use dbList to update dropdown box option
	dbDropdown := widget.NewSelect(
		dbList,
		func(selection string) {
			fmt.Printf("I selected %selection as my input DB..", selection)
			dbDropdownResult.Text = selection
			dbSelection = selection
			dbDropdownResult.Refresh()
		})

	//Create dropdown for table selection
	tableDropdownResult := widget.NewLabel("Select a Table")
	tableList := []string{"movesoutput", "rateperdistance", "rateperhour", "rateperprofile", "rateperstart", "ratepervehicle", "startspervehicle"}
	//Use dbList to update dropdown box option
	tableDropdown := widget.NewSelect(
		tableList,
		func(selection string) {
			fmt.Printf("I selected %selection as my input DB..", selection)
			tableDropdownResult.Text = selection
			tableSelection = selection
			tableDropdownResult.Refresh()
		})

	//// more than one widget. so use container
	//c := container.NewBorder(nil, dbDropdownResult, dbDropdown, tableDropdown)
	//w.SetContent(c)

	dropdownGrid := container.New(layout.NewGridLayout(2), dbDropdown, tableDropdown)
	w.SetContent(dropdownGrid)
	//show and run
	w.ShowAndRun()
}

//func test(m []Movesoutput) *fyne.Container {
//	table := widget.NewTable(
//		func() (int, int) {
//			return 10, 25
//		},
//		func() fyne.CanvasObject {
//			return widget.NewLabel("table window")
//		},
//		func(i widget.TableCellID, o fyne.CanvasObject) {
//			switch i.Col {
//			case 0:
//				o.(*widget.Label).SetText(m[i.Row].MOVESRunID)
//			case 1:
//				o.(*widget.Label).SetText(Convert.ToString(m[i.Row].MOVESRunID))
//			case 2:
//				o.(*widget.Label).SetText(m[i.Row].Name)
//			case 3:
//				o.(*widget.Label).SetText(m[i.Row].Memo)
//			}
//		},
//	)
//	table.SetColumnWidth(0, 200)
//	table.SetColumnWidth(1, 100)
//	table.SetColumnWidth(2, 100)
//	table.SetColumnWidth(3, 300)
//	split := container.NewHSplit(makeLeftSidebar(), table)
//	split.Offset = 0.2
//	return container.NewMax(split)
//}

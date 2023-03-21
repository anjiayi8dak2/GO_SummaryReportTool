package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	_ "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/sqltocsv"
)

func main() {
	// Create the database handle, confirm driver is present
	db := initDb()

	// Connect and check the server version
	getDBVersion(db)

	//Initialize app main window
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	window1 := a.NewWindow("Summary Report Tool")
	window1.Resize(fyne.NewSize(400, 400))

	//Global var
	var dbSelection string
	var tableSelection string

	//Top menu bar
	menuitemMaria := fyne.NewMenuItem("Open Maria Data Folder", func() {
		openMariaFolder(db)
	}) // ignore functions
	//menuitemRefresh := fyne.NewMenuItem("Refresh (F5)", buttonSubmit(db, dbSelection, tableSelection)) // ignore functions
	//TODO: all the test goes refresh button here
	menuitemRefresh := fyne.NewMenuItem("Refresh (F5)", func() {
		//whiteList []string only contains column names, whiteListIndex [] bool contains all columns from movesoutput
		whiteList, whiteListIndex := getWhiteList(db, dbSelection, tableSelection)
		//fmt.Println("printing white list index in bool")
		//fmt.Printf("%v", whiteListIndex)
		queryResult, _ := getQueryResult(db, dbSelection, tableSelection, whiteList, "", "")

		makeWindowTwo(a, queryResult, db, dbSelection, tableSelection, whiteListIndex, whiteList)

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
	window1.SetMainMenu(menu)

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

	dropdownGrid := container.New(layout.NewGridLayout(2), dbDropdown, tableDropdown)
	window1.SetContent(dropdownGrid)
	//show and run
	window1.ShowAndRun()
}

package main

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	a := app.New()
	w := a.NewWindow("Summary Report Tool")
	w.Resize(fyne.NewSize(400, 400))

	//menu bar
	menuitemMaria := fyne.NewMenuItem("Open Maria Data Folder", nil) // ignore functions
	menuitemRefresh := fyne.NewMenuItem("Refresh (F5)", nil)         // ignore functions
	menuitemOpenlog := fyne.NewMenuItem("Open Log", nil)             // ignore functions
	menuitemClearlog := fyne.NewMenuItem("Clear Log", nil)           // ignore functions
	menuitemAbout := fyne.NewMenuItem("About", nil)                  // ignore functions
	menuitemManual := fyne.NewMenuItem("Manual", nil)                // ignore functions
	// New Menu
	newMenu1 := fyne.NewMenu("File", menuitemMaria)
	newMenu2 := fyne.NewMenu("Edit", menuitemRefresh)
	newMenu3 := fyne.NewMenu("Logs", menuitemOpenlog, menuitemClearlog)
	newMenu4 := fyne.NewMenu("Help", menuitemAbout, menuitemManual)
	// New main menu
	menu := fyne.NewMainMenu(newMenu1, newMenu2, newMenu3, newMenu4)
	w.SetMainMenu(menu)

	// Create the database handle, confirm driver is present
	db, _ := sql.Open("mysql", "moves:moves@/")
	defer db.Close()

	// Connect and check the server version
	getDBVersion(db)

	// Get DB List
	var dbList []string
	dbList = getDBList(db)

	//create dropdown for database selection
	dbDropdownResult := widget.NewLabel("Select a Database")
	// use dbList to update dropdown box option
	dbDropdown := widget.NewSelect(
		dbList,
		func(selection string) {
			fmt.Printf("I selected %selection as my input DB..", selection)
			dbDropdownResult.Text = selection
			dbDropdownResult.Refresh()
		})

	//create dropdown for table selection
	tableDropdownResult := widget.NewLabel("Select a Table")
	tableList := []string{"movesoutput", "rateperdistance", "rateperhour", "rateperprofile", "rateperstart", "ratepervehicle", "startspervehicle"}
	// use dbList to update dropdown box option
	tableDropdown := widget.NewSelect(
		tableList,
		func(selection string) {
			fmt.Printf("I selected %selection as my input DB..", selection)
			tableDropdownResult.Text = selection
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

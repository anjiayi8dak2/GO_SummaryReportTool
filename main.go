package main

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)

	// Get DB List
	var dbList []string
	dbList = getDBList(db)
	// use dbList to update dropdown box
	// lets show our selected entry in label
	dbDropdownResult := widget.NewLabel("...")
	// dropdown/ select entry
	//[]string{} all our option goes in slice
	// s is the variable to get the selected value
	dbDropdown := widget.NewSelect(
		dbList,
		func(s string) {
			fmt.Printf("I selected %s as my input DB..", s)
			dbDropdownResult.Text = s
			dbDropdownResult.Refresh()
		})
	// more than one widget. so use container
	c := container.NewBorder(nil, dbDropdownResult, dbDropdown, nil)
	w.SetContent(c)
	//show and run
	w.ShowAndRun()
}

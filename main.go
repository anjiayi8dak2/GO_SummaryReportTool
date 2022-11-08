package main

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	a := app.New()
	w := a.NewWindow("Select entry widget, drop down")
	w.Resize(fyne.NewSize(400, 400))

	//toolbar
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			log.Println("New document")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentCutIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)
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
	c := container.NewBorder(toolbar, dbDropdownResult, dbDropdown, nil)
	w.SetContent(c)
	//show and run
	w.ShowAndRun()
}

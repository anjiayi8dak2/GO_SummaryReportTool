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
	// Create the database handle, confirm driver is present
	db := initDb()

	// Connect and check the server version
	getDBVersion(db)

	//Initialize app main window
	a := app.New()
	w := a.NewWindow("Summary Report Tool")
	w.Resize(fyne.NewSize(400, 400))
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
		whiteList, whiteListIndex := getWhiteList(db, dbSelection, tableSelection)
		fmt.Println("printing white list index in bool")
		fmt.Printf("%v", whiteListIndex)
		dummy, _ := getQueryResult(db, dbSelection, tableSelection, whiteList)
		fmt.Println("opening window #2")
		w2 := a.NewWindow("window #2")
		w2.SetContent(widget.NewLabel("window #2 label"))
		w2.Resize(fyne.NewSize(1000, 1000))

		tableData := widget.NewTable(
			func() (int, int) {
				return 1000, len(dummy[0]) // row size max
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("wide content")
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(dummy[i.Row][i.Col])
			})

		//TODO; make a function to create window or containers not in the main.go
		//create  buttons
		MOVESRunID := widget.NewButton("MOVESRunID", func() {
		})

		iterationID := widget.NewButton("iterationID", func() {
		})

		yearID := widget.NewButton("yearID", func() {
		})

		monthID := widget.NewButton("monthID", func() {
		})

		dayID := widget.NewButton("dayID", func() {
		})

		hourID := widget.NewButton("hourID", func() {
		})

		stateID := widget.NewButton("stateID", func() {
		})

		countyID := widget.NewButton("countyID", func() {
		})

		zoneID := widget.NewButton("zoneID", func() {
		})

		linkID := widget.NewButton("linkID", func() {
		})

		//pollutant
		pollutantidButton := widget.NewButton("pollutantID", func() {
		})
		distinctPollutant := getDistinct(db, dbSelection, tableSelection, "pollutantID")
		pollutantidCheckGroup := widget.NewCheckGroup(distinctPollutant, func(s []string) { fmt.Println("selected", s) })
		pollutantContainer := container.NewVBox(pollutantidButton, pollutantidCheckGroup)

		processID := widget.NewButton("processID", func() {
		})

		sourceTypeID := widget.NewButton("sourceTypeID", func() {
		})

		regClassID := widget.NewButton("regClassID", func() {
		})

		fuelTypeID := widget.NewButton("fuelTypeID", func() {
		})

		fuelSubTypeID := widget.NewButton("fuelSubTypeID", func() {
		})

		//model year
		modelYearIDButton := widget.NewButton("modelYearID", func() {
		})
		distinctModelYear := getDistinct(db, dbSelection, tableSelection, "modelYearID")
		modelYearCheckGroup := widget.NewCheckGroup(distinctModelYear, func(s []string) { fmt.Println("selected", s) })
		modelYearContainer := container.NewVBox(modelYearIDButton, modelYearCheckGroup)

		roadTypeID := widget.NewButton("roadTypeID", func() {
		})

		SCC := widget.NewButton("SCC", func() {
		})

		engTechID := widget.NewButton("engTechID", func() {
		})

		sectorID := widget.NewButton("sectorID", func() {
		})

		hpID := widget.NewButton("hpID", func() {
		})

		innerContainer := container.NewVBox(
			MOVESRunID,
			iterationID,
			yearID,
			monthID,
			dayID,
			hourID,
			stateID,
			countyID,
			zoneID,
			linkID,
			pollutantContainer,
			processID,
			sourceTypeID,
			regClassID,
			fuelTypeID,
			fuelSubTypeID,
			modelYearContainer,
			roadTypeID,
			SCC,
			engTechID,
			sectorID,
			hpID)

		//dynamic filter buttons, Use the record of whitelist, delete corresponding filter buttons above.
		for index, boo := range whiteListIndex {
			if boo {
				innerContainer.Objects[index].Visible()
			} else {
				innerContainer.Objects[index].Hide()
			}
		}

		//TODO: the filter button section has no scroll bar, contents are out of screen
		scrollContainer := container.NewVScroll(innerContainer)
		outerContainer := container.NewHSplit(
			scrollContainer,
			tableData,
		)

		w2.SetContent(outerContainer)
		w2.Show()

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

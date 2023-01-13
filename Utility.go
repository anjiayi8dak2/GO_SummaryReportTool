package main

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"os/exec"
	"runtime"
	"strings"
)

func openMariaFolder(db *sql.DB) {
	dataDir := getDataDir(db)
	dataDir = "\"" + dataDir + "\""
	fmt.Println(dataDir)
	cmd := "open"
	if runtime.GOOS == "windows" {
		cmd = "explorer"
	}
	fmt.Println(cmd)
	exec.Command(cmd, dataDir).Start()
	fmt.Println("finishing utility/openMariaFolder func")

}

// take a list of columns name, return one piece string with commas like roadTypeID, sourceTypeID, emissionQuant
func convertColumnsComma(columns []string) string {

	// prepend single quote, perform joins, append single quote
	ColumnsComma := strings.Join(columns, `,`)

	fmt.Println("printing comma seperated columns::: ", ColumnsComma)
	return ColumnsComma
}

func makeWindowTwo(a fyne.App, queryResult [][]string, db *sql.DB, dbSelection string, tableSelection string, whiteListIndex []bool) {
	fmt.Println("opening window #2")
	w2 := a.NewWindow("window #2")
	w2.SetContent(widget.NewLabel("window #2 label"))
	w2.Resize(fyne.NewSize(1000, 1000))

	tableData := widget.NewTable(
		func() (int, int) {
			return 1000, len(queryResult[0]) // row size max
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(queryResult[i.Row][i.Col])
		})

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
}

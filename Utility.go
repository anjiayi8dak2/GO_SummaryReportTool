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

	//create buttons
	//To these filters suppose to have group of checkbox
	//CheckGroup
	//= pollutantContainer
	//= title button + check group in vertical
	//For example
	//= pollutantidButton + pollutantContainer

	MOVESRunIDContainer := createNewButton(db, "MOVESRunID", dbSelection, tableSelection)

	iterationIDContainer := createNewButton(db, "iterationID", dbSelection, tableSelection)

	yearIDContainer := createNewButton(db, "yearID", dbSelection, tableSelection)

	monthIDContainer := createNewButton(db, "monthID", dbSelection, tableSelection)

	dayIDContainer := createNewButton(db, "dayID", dbSelection, tableSelection)

	hourIDContainer := createNewButton(db, "hourID", dbSelection, tableSelection)

	stateIDContainer := createNewButton(db, "stateID", dbSelection, tableSelection)

	countyIDContainer := createNewButton(db, "countyID", dbSelection, tableSelection)

	zoneIDContainer := createNewButton(db, "zoneID", dbSelection, tableSelection)

	linkIDContainer := createNewButton(db, "linkID", dbSelection, tableSelection)

	pollutantContainer := createNewButton(db, "pollutantID", dbSelection, tableSelection)

	processIDContainer := createNewButton(db, "processID", dbSelection, tableSelection)

	sourceTypeIDContainer := createNewButton(db, "sourceTypeID", dbSelection, tableSelection)

	regClassIDContainer := createNewButton(db, "regClassID", dbSelection, tableSelection)

	fuelTypeIDContainer := createNewButton(db, "fuelTypeID", dbSelection, tableSelection)

	fuelSubTypeIDContainer := createNewButton(db, "fuelSubTypeID", dbSelection, tableSelection)

	modelYearContainerContainer := createNewButton(db, "modelYearID", dbSelection, tableSelection)

	roadTypeIDContainer := createNewButton(db, "roadTypeID", dbSelection, tableSelection)

	SCCContainer := createNewButton(db, "SCC", dbSelection, tableSelection)

	engTechIDContainer := createNewButton(db, "engTechID", dbSelection, tableSelection)

	sectorIDContainer := createNewButton(db, "sectorID", dbSelection, tableSelection)

	hpIDContainer := createNewButton(db, "hpID", dbSelection, tableSelection)

	innerContainer := container.NewVBox(
		MOVESRunIDContainer,
		iterationIDContainer,
		yearIDContainer,
		monthIDContainer,
		dayIDContainer,
		hourIDContainer,
		stateIDContainer,
		countyIDContainer,
		zoneIDContainer,
		linkIDContainer,
		pollutantContainer,
		processIDContainer,
		sourceTypeIDContainer,
		regClassIDContainer,
		fuelTypeIDContainer,
		fuelSubTypeIDContainer,
		modelYearContainerContainer,
		roadTypeIDContainer,
		SCCContainer,
		engTechIDContainer,
		sectorIDContainer,
		hpIDContainer)

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

func createNewButton(db *sql.DB, columnsName string, dbSelection string, tableSelection string) *fyne.Container {
	xButton := widget.NewButton(columnsName, func() {
	})
	distinctX := getDistinct(db, dbSelection, tableSelection, columnsName)
	xCheckGroup := widget.NewCheckGroup(distinctX, func(s []string) { fmt.Println("selected", s) })
	xContainer := container.NewVBox(xButton, xCheckGroup)

	return xContainer
}

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"log"
	"os"
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

func makeWindowTwo(a fyne.App, queryResult [][]string, db *sql.DB, dbSelection string, tableSelection string, whiteListIndex []bool, whiteList []string) {
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
			if len(queryResult) >= 2 {
				o.(*widget.Label).SetText(queryResult[i.Row][i.Col])
			} else {
				o.(*widget.Label).SetText("no data")
				//runPopUp(w2)
				//dialog.ShowCustom("title", "Ok", nil, w2)
			}
		})

	//map to hold movesoutput filters
	moFilter := make(map[string][]string)

	//create buttons
	MOVESRunIDContainer := createNewCheckBoxGroup(db, "MOVESRunID", dbSelection, tableSelection, moFilter)

	iterationIDContainer := createNewCheckBoxGroup(db, "iterationID", dbSelection, tableSelection, moFilter)

	yearIDContainer := createNewCheckBoxGroup(db, "yearID", dbSelection, tableSelection, moFilter)

	monthIDContainer := createNewCheckBoxGroup(db, "monthID", dbSelection, tableSelection, moFilter)

	dayIDContainer := createNewCheckBoxGroup(db, "dayID", dbSelection, tableSelection, moFilter)

	hourIDContainer := createNewCheckBoxGroup(db, "hourID", dbSelection, tableSelection, moFilter)

	stateIDContainer := createNewCheckBoxGroup(db, "stateID", dbSelection, tableSelection, moFilter)

	countyIDContainer := createNewCheckBoxGroup(db, "countyID", dbSelection, tableSelection, moFilter)

	zoneIDContainer := createNewCheckBoxGroup(db, "zoneID", dbSelection, tableSelection, moFilter)

	linkIDContainer := createNewCheckBoxGroup(db, "linkID", dbSelection, tableSelection, moFilter)

	pollutantContainer := createNewCheckBoxGroup(db, "pollutantID", dbSelection, tableSelection, moFilter)

	processIDContainer := createNewCheckBoxGroup(db, "processID", dbSelection, tableSelection, moFilter)

	sourceTypeIDContainer := createNewCheckBoxGroup(db, "sourceTypeID", dbSelection, tableSelection, moFilter)

	regClassIDContainer := createNewCheckBoxGroup(db, "regClassID", dbSelection, tableSelection, moFilter)

	fuelTypeIDContainer := createNewCheckBoxGroup(db, "fuelTypeID", dbSelection, tableSelection, moFilter)

	fuelSubTypeIDContainer := createNewCheckBoxGroup(db, "fuelSubTypeID", dbSelection, tableSelection, moFilter)

	modelYearContainerContainer := createNewCheckBoxGroup(db, "modelYearID", dbSelection, tableSelection, moFilter)

	roadTypeIDContainer := createNewCheckBoxGroup(db, "roadTypeID", dbSelection, tableSelection, moFilter)

	SCCContainer := createNewCheckBoxGroup(db, "SCC", dbSelection, tableSelection, moFilter)

	engTechIDContainer := createNewCheckBoxGroup(db, "engTechID", dbSelection, tableSelection, moFilter)

	sectorIDContainer := createNewCheckBoxGroup(db, "sectorID", dbSelection, tableSelection, moFilter)

	hpIDContainer := createNewCheckBoxGroup(db, "hpID", dbSelection, tableSelection, moFilter)

	//TODO: update button with icon, icon way too small
	imgUpdate, err := os.Open("update.jpg")
	r := bufio.NewReader(imgUpdate)

	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	updateButton := widget.NewButtonWithIcon("UPDATE", fyne.NewStaticResource("icon", b), func() {
		whereClause := " WHERE "
		fmt.Println("pressed UPDATE button")
		//loop through all the keys in mo map and generate a where clause
		fmt.Println("loop throuth the whiteList keys and print")
		for index := 0; index < len(whiteList); index++ {
			var partialWhere string
			value, ok := moFilter[whiteList[index]]
			if ok { //none 0 value
				fmt.Println(whiteList[index], " key found: ", value)
				//TODO: detect if need AND
				if len(whereClause) > 7 { //predefined whereClause with string " where ", so that the size should be 7, if the size over 7, that means we need put AND in the beginning
					inValue := convertColumnsComma(value)
					fmt.Println("print in values ", inValue)
					partialWhere = " AND " + whiteList[index] + " IN ( " + inValue + " ) "
					fmt.Println("print dummy clause ", partialWhere)
				} else { // otherwise do not put AND
					inValue := convertColumnsComma(value)
					fmt.Println("print in values ", inValue)
					partialWhere = whiteList[index] + " IN ( " + inValue + " ) "
					fmt.Println("print dummy clause ", partialWhere)
				}
			} else {
				fmt.Println(whiteList[index], " key not found")
				//delete map that has empty value, they looks like hpID:[], they will eventually cause an empty IN() in the where clause
				for key, value := range moFilter {
					if len(value) == 0 {
						delete(moFilter, key)
					}
				}
				fmt.Println("print the map at then end of button function")
				fmt.Println(moFilter)
			}
			//append inner string to outer string
			whereClause += partialWhere
		}
		//catch if no checkbox were selected, then remove the default WHERE, since there is no filter
		if whereClause == " WHERE " {
			whereClause = ""
		}

		fmt.Println("printing the WHERE clause")
		fmt.Println(whereClause)

		//update the matrix with the new where clause we just made
		queryResult, err = getQueryResult(db, dbSelection, tableSelection, whiteList, whereClause)
		fmt.Println("printing error query result WHERE clause")
		fmt.Println(err)

		//TODO: make some sort of dialog pop out warning for no result query
		if len(queryResult) < 2 {
			runPopUp(w2, "Filter combination return no data, please try different filter")
		}

	})

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
		hpIDContainer,
		updateButton,
	)
	//dynamic filter buttons, Use the record of whiteListIndex [] bool, show and hide base on 1 or 0.
	for index, boo := range whiteListIndex {
		if boo {
			innerContainer.Objects[index].Visible()
		} else {
			innerContainer.Objects[index].Hide()
		}
	}

	// the filter button section scroll bar
	scrollContainer := container.NewVScroll(innerContainer)
	outerContainer := container.NewHSplit(
		scrollContainer,
		tableData,
	)
	//screen width distribution of filter panel VS data table panel
	outerContainer.Offset = 0.08
	w2.SetContent(outerContainer)
	w2.Show()
}

func createNewCheckBoxGroup(db *sql.DB, columnsName string, dbSelection string, tableSelection string, mo map[string][]string) *fyne.Container {
	//To these filters suppose to have group of checkbox
	//CheckGroup
	//= pollutantContainer
	//= title button + checkbox group in vertical
	//For example
	//pollutantidButton + pollutantContainer
	xButton := widget.NewButton(columnsName, func() {
	})
	distinctX := getDistinct(db, dbSelection, tableSelection, columnsName)
	xCheckGroup := widget.NewCheckGroup(distinctX, func(value []string) {
		fmt.Println("selected", value)
		//update map  from checked boxes statues
		mo[columnsName] = value
		fmt.Println("print entire filter map for  ", columnsName, " inside func createNewCheckBoxGroup")
		fmt.Println(mo)
	})

	xContainer := container.NewVBox(xButton, xCheckGroup)
	return xContainer
}

func RemoveElementFromSlice[T comparable](l []T, item T) []T {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func runPopUp(w fyne.Window, msg string) (modal *widget.PopUp) {
	modal = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel(msg),
			widget.NewButton("ok", func() { modal.Hide() }),
		),
		w.Canvas(),
	)
	modal.Show()
	return modal
}

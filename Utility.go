package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	grob "github.com/MetalBlueberry/go-plotly/graph_objects"
	"github.com/MetalBlueberry/go-plotly/offline"
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
	window2 := a.NewWindow("window #2")
	window2.SetContent(widget.NewLabel("window #2 label"))
	window2.Resize(fyne.NewSize(1000, 800))

	//TODO: This is the message tab on top of screen, should update this text on the fly

	dummyToolBarString_Where := fmt.Sprint(dbSelection)      //test some random variable print in the label, maybe update later somewhere
	dummyToolBarString_GroupBy := fmt.Sprint(tableSelection) //test some random variable print in the label, maybe update later somewhere
	dummyToolBarString := "Filters: " + dummyToolBarString_Where + "Aggregated by: " + dummyToolBarString_GroupBy
	dummyLable := widget.NewLabel(dummyToolBarString)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() { //update
			fmt.Println("I pressed update button")
		}),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() { //plot
			fmt.Println("I pressed plot button")
			//TODO: make X and Y here
			// bar names should be elements on the first row, except last element(emissionQuant)
			// how many bar, count of X = total row-1, because the first row is header
			// value of Y should be at last column(emissionQuant) except for first row

			temp_queryResult := queryResult
			plot_row_count := len(temp_queryResult) - 1 //total row - 1 because the first row is title
			fmt.Println("Printing row count without header", plot_row_count)
			//plot_column_count := len(temp_queryResult[0])
			//fmt.Println("Printing column count", plot_column_count)

			//TODO: add the checker in the beginning for size of queryResult, if it is empty then skip all
			//get X title, bar title
			var X_title []string
			for row := 1; row < len(temp_queryResult); row++ {
				X_title = append(X_title, temp_queryResult[row][0])
			}

			//X_title := temp_queryResult[0]
			//delete element emissionQuant from X title
			//RemoveElementFromSlice(X_title, "sum(emissionQuant)")
			fmt.Printf("Printing X title %v\n", X_title)
			//get Y title
			Y_title := "sum(emissionQuant)"
			fmt.Printf("Printing Y title %v\n\n", Y_title)
			//get Y value
			var Y_value []string
			for row := 1; row < len(temp_queryResult); row++ { //skip first row, start from second row
				Y_value = append(Y_value, temp_queryResult[row][len(temp_queryResult[0])-1]) //the index for last column = total column count -1
			}
			fmt.Printf("Printing Y value %v\n\n", Y_value)

			fig := &grob.Fig{
				Data: grob.Traces{
					&grob.Bar{
						Type: grob.TraceTypeBar,
						X:    X_title,
						Y:    Y_value,
					},
				},

				Layout: &grob.Layout{
					Title: &grob.LayoutTitle{
						Text: "Aggregation Result Plot",
					},
				},
			}

			offline.Show(fig)
		}),
		widget.NewToolbarAction(theme.DownloadIcon(), func() { //download CSV
			fmt.Println("I pressed download csv button")
			csvExport(queryResult)

		}),
		widget.NewToolbarSpacer(),
		//TODO: the text space on the toolbar, to show what filtered has applied and/or aggregated by
		//dummyLable:= widget.NewLabel(dummyToolBarString),
		dummyLable,
	)

	tableData := widget.NewTable(
		func() (int, int) {
			return len(queryResult), len(queryResult[0]) // row size, columns size
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if len(queryResult) >= 2 {
				o.(*widget.Label).SetText(queryResult[i.Row][i.Col])
			} else {
				o.(*widget.Label).SetText("no data")
			}
		})

	//map to hold movesoutput filters selection in checkbox group
	moFilter := make(map[string][]string)

	//fyne containers, create buttons for filters with checkbox selection saved in the map moFilter
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

	//TODO: aggregation selection
	//map to hold group by check boxes selection
	moGroupBy := make(map[string][]string)
	//pass whitelist?
	aggregationContainer := createNewAggregationGroup(whiteList, moGroupBy)

	updateButton := widget.NewButtonWithIcon("Update", theme.MediaReplayIcon(), func() {
		whereClause := " WHERE "
		fmt.Println("pressed UPDATE button")
		//loop through all the keys in mo map and generate a where clause
		fmt.Println("loop throuth the whiteList keys and print")
		// there is no easy loop solution for select and then unselect operation on the run time, because it will generate a key with empty value such as hpid {}
		// these empty values will cause empty IN() statement in the where clause that make future problems.
		// hence, we should detect empty value and delete that key before disaster happen
		for index := 0; index < len(whiteList); index++ {
			var partialWhere string
			value, ok := moFilter[whiteList[index]]
			if ok { //none 0 value
				fmt.Println(whiteList[index], " key found: ", value)
				//detect if need AND
				if len(whereClause) > 7 { //predefined whereClause with string " where ", size = 7, if  size > 7, that means we need put AND in the beginning
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
				//dummyToolBarString_Where = fmt.Sprint(moFilter)
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
		//test update dummy toolbar label on the fly
		dummyToolBarString_Where = whereClause
		dummyLable.SetText("this is updated dummy lable")

		//TODO: enter GROUP BY claus
		// if there is 1 item is checked, have GROUP BY xxx
		// else if, >1 items are checked, call func convertColumnsComma, get comma seperated field name, then have GROUP BY xxx,xxx, .....
		// else there is 0 item is checked, remove the "GROUP BY"

		// TODO: if 1 or more items are checked, update the select columns same as the GROUP BY, Plus the sum of emissionQuant, maybe override func getQueryResult with direct sql statement

		groupbyClause := " GROUP BY "
		var columnSelection []string
		var groupbySelection []string

		if len(moGroupBy["Aggregation"]) == 0 { //if nothing in the group by map
			//TODO: when select 0 checkbox, do not need group clause AND the select column should same as whiteList
			groupbyClause = ""
			columnSelection = whiteList
		} else if len(moGroupBy["Aggregation"]) >= 1 { //if there is anything in the group by map
			//TODO: update SELECT clause PLUS sum of emissionQuant
			//TODO: update GROUP BY
			//loop through the group by map, and copy selected box into a temp slice
			for _, value := range moGroupBy["Aggregation"] {
				groupbySelection = append(groupbySelection, value)
				fmt.Println("copy group by map value into slice", value)
				fmt.Println("printing updated group by slice", groupbySelection)
			}

			//pass the selected box name, PLUS sum(emissionQuant) AS emissionQuant) in the end of select statement
			//since columnSelection is []string, just throw the "sum(emissionQuant) AS emissionQuant)" in, as an element
			//getQuery will handle comma separate format later
			columnSelection = groupbySelection
			columnSelection = append(columnSelection, "sum(emissionQuant) ")
			//pass the selected box name to GROUP BY clause, convert list of name into comma seperated format
			groupbyClause += convertColumnsComma(groupbySelection)
		} else {
			panic("detect length of moGroupBy map size <0, WHY")
		}

		//update the matrix with the new where clause and group by we just made
		var err error
		queryResult, err = getQueryResult(db, dbSelection, tableSelection, columnSelection, whereClause, groupbyClause)
		fmt.Println("printing error query result WHERE clause")
		fmt.Println(err)

		//dialog box pop out warning for no result query
		if len(queryResult) < 2 {
			runPopUp(window2, "Filter combination returns no data, please try different filter")
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
		aggregationContainer,
		updateButton,
	)
	//dynamic filter buttons, Use the record of whiteListIndex [] bool, show and hide base on 1 or 0.
	//we initialized all 25 columns when the window #2 started
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
	//screen width horizontal distribution of filter panel VS data table panel
	outerContainer.Offset = 0.08
	//window2.SetContent(outerContainer)
	window2.SetContent(container.NewBorder(toolbar, nil, nil, nil, outerContainer))
	window2.Show()
}

func updateButton() {

}

func createNewAggregationGroup(whitelist []string, groupBy map[string][]string) *fyne.Container {
	xButton := widget.NewButton("Aggregation", func() {
	})
	whitelist2 := whitelist
	//whitelist has emssionQuant, delete the last element
	if len(whitelist2) > 0 {
		whitelist2 = whitelist2[:len(whitelist2)-1]
	}
	fmt.Println("printing whitelist2 slice inside createNewAggregationGroup", whitelist2)
	xCheckGroup := widget.NewCheckGroup(whitelist2, func(value []string) {
		fmt.Println("selected", value)
		//update map from checked boxes statues
		groupBy["Aggregation"] = value
		fmt.Println("print entire group by map for  ", "Aggregation", " inside func createNewAggregationGroup")
		fmt.Println(groupBy)

		//check if map is empty
		if len(groupBy["Aggregation"]) == 0 {
			fmt.Println("Aggregation map has no value", groupBy["Aggregation"])
			fmt.Println("BEFORE", groupBy["Aggregation"])
			delete(groupBy, "Aggregation")
			fmt.Println("AFTER", groupBy["Aggregation"])
		}
	})

	xContainer := container.NewVBox(xButton, xCheckGroup)
	return xContainer
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
		//TODO: put check empty value key here??
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

func csvExport(data [][]string) error {
	file, err := os.Create("result.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		if err := writer.Write(value); err != nil {
			return err // let's return errors if necessary, rather than having a one-size-fits-all error handler
		}
	}
	return nil
}

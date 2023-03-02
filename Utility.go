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
	//default display for DB and Table selection
	ToolbarLabel := widget.NewLabel("DB Selection: " + dbSelection + "Table Selection: " + tableSelection)

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
		ToolbarLabel,
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

	//TODO: switch here to split 7 table selections
	//map to hold filters selection in checkbox group for where clause
	filter := make(map[string][]string)
	//map to hold group by check boxes selection for group by clause
	groupBy := make(map[string][]string)
	//container sections
	innerContainer := container.NewVBox()

	switch tableSelection {
	case "movesoutput":
		//fyne containers, create buttons for filters with checkbox selection saved in the map filter
		MOVESRunIDContainer := createNewCheckBoxGroup(db, "MOVESRunID", dbSelection, tableSelection, filter)
		iterationIDContainer := createNewCheckBoxGroup(db, "iterationID", dbSelection, tableSelection, filter)
		yearIDContainer := createNewCheckBoxGroup(db, "yearID", dbSelection, tableSelection, filter)
		monthIDContainer := createNewCheckBoxGroup(db, "monthID", dbSelection, tableSelection, filter)
		dayIDContainer := createNewCheckBoxGroup(db, "dayID", dbSelection, tableSelection, filter)
		hourIDContainer := createNewCheckBoxGroup(db, "hourID", dbSelection, tableSelection, filter)
		stateIDContainer := createNewCheckBoxGroup(db, "stateID", dbSelection, tableSelection, filter)
		countyIDContainer := createNewCheckBoxGroup(db, "countyID", dbSelection, tableSelection, filter)
		zoneIDContainer := createNewCheckBoxGroup(db, "zoneID", dbSelection, tableSelection, filter)
		linkIDContainer := createNewCheckBoxGroup(db, "linkID", dbSelection, tableSelection, filter)
		pollutantContainer := createNewCheckBoxGroup(db, "pollutantID", dbSelection, tableSelection, filter)
		processIDContainer := createNewCheckBoxGroup(db, "processID", dbSelection, tableSelection, filter)
		sourceTypeIDContainer := createNewCheckBoxGroup(db, "sourceTypeID", dbSelection, tableSelection, filter)
		regClassIDContainer := createNewCheckBoxGroup(db, "regClassID", dbSelection, tableSelection, filter)
		fuelTypeIDContainer := createNewCheckBoxGroup(db, "fuelTypeID", dbSelection, tableSelection, filter)
		fuelSubTypeIDContainer := createNewCheckBoxGroup(db, "fuelSubTypeID", dbSelection, tableSelection, filter)
		modelYearContainer := createNewCheckBoxGroup(db, "modelYearID", dbSelection, tableSelection, filter)
		roadTypeIDContainer := createNewCheckBoxGroup(db, "roadTypeID", dbSelection, tableSelection, filter)
		SCCContainer := createNewCheckBoxGroup(db, "SCC", dbSelection, tableSelection, filter)
		engTechIDContainer := createNewCheckBoxGroup(db, "engTechID", dbSelection, tableSelection, filter)
		sectorIDContainer := createNewCheckBoxGroup(db, "sectorID", dbSelection, tableSelection, filter)
		hpIDContainer := createNewCheckBoxGroup(db, "hpID", dbSelection, tableSelection, filter)

		innerContainer.Add(MOVESRunIDContainer)
		innerContainer.Add(iterationIDContainer)
		innerContainer.Add(yearIDContainer)
		innerContainer.Add(monthIDContainer)
		innerContainer.Add(dayIDContainer)
		innerContainer.Add(hourIDContainer)
		innerContainer.Add(stateIDContainer)
		innerContainer.Add(countyIDContainer)
		innerContainer.Add(zoneIDContainer)
		innerContainer.Add(linkIDContainer)
		innerContainer.Add(pollutantContainer)
		innerContainer.Add(processIDContainer)
		innerContainer.Add(sourceTypeIDContainer)
		innerContainer.Add(regClassIDContainer)
		innerContainer.Add(fuelTypeIDContainer)
		innerContainer.Add(fuelSubTypeIDContainer)
		innerContainer.Add(modelYearContainer)
		innerContainer.Add(roadTypeIDContainer)
		innerContainer.Add(SCCContainer)
		innerContainer.Add(engTechIDContainer)
		innerContainer.Add(sectorIDContainer)
		innerContainer.Add(hpIDContainer)
	case "rateperdistance":
		//MOVESScenarioID, MOVESRunID, yearID, monthID, dayID, hourID, linkID, pollutantID, processID, sourceTypeID,
		//regClassID, SCC, fuelTypeID, modelYearID, roadTypeID, avgSpeedBinID, temperature, relHumidity, ratePerDistance
		MOVESScenarioIDContainer := createNewCheckBoxGroup(db, "MOVESScenarioID", dbSelection, tableSelection, filter)
		MOVESRunIDContainer := createNewCheckBoxGroup(db, "MOVESRunID", dbSelection, tableSelection, filter)
		yearIDContainer := createNewCheckBoxGroup(db, "yearID", dbSelection, tableSelection, filter)
		monthIDContainer := createNewCheckBoxGroup(db, "monthID", dbSelection, tableSelection, filter)
		dayIDContainer := createNewCheckBoxGroup(db, "dayID", dbSelection, tableSelection, filter)
		hourIDContainer := createNewCheckBoxGroup(db, "hourID", dbSelection, tableSelection, filter)
		linkIDContainer := createNewCheckBoxGroup(db, "linkID", dbSelection, tableSelection, filter)
		pollutantIDContainer := createNewCheckBoxGroup(db, "pollutantID", dbSelection, tableSelection, filter)
		processIDContainer := createNewCheckBoxGroup(db, "processID", dbSelection, tableSelection, filter)
		sourceTypeIDContainer := createNewCheckBoxGroup(db, "sourceTypeID", dbSelection, tableSelection, filter)
		regClassIDContainer := createNewCheckBoxGroup(db, "regClassID", dbSelection, tableSelection, filter)
		SCCContainer := createNewCheckBoxGroup(db, "SCC", dbSelection, tableSelection, filter)
		fuelTypeIDContainer := createNewCheckBoxGroup(db, "fuelTypeID", dbSelection, tableSelection, filter)
		modelYearIDContainer := createNewCheckBoxGroup(db, "modelYearID", dbSelection, tableSelection, filter)
		roadTypeIDContainer := createNewCheckBoxGroup(db, "roadTypeID", dbSelection, tableSelection, filter)
		avgSpeedBinIDContainer := createNewCheckBoxGroup(db, "avgSpeedBinID", dbSelection, tableSelection, filter)

		innerContainer.Add(MOVESScenarioIDContainer)
		innerContainer.Add(MOVESRunIDContainer)
		innerContainer.Add(yearIDContainer)
		innerContainer.Add(monthIDContainer)
		innerContainer.Add(dayIDContainer)
		innerContainer.Add(hourIDContainer)
		innerContainer.Add(linkIDContainer)
		innerContainer.Add(pollutantIDContainer)
		innerContainer.Add(processIDContainer)
		innerContainer.Add(sourceTypeIDContainer)
		innerContainer.Add(regClassIDContainer)
		innerContainer.Add(SCCContainer)
		innerContainer.Add(fuelTypeIDContainer)
		innerContainer.Add(modelYearIDContainer)
		innerContainer.Add(roadTypeIDContainer)
		innerContainer.Add(avgSpeedBinIDContainer)

	case "rateperhour":
		//123
	case "rateperprofile":
		//123
	case "rateperstart":
	case "ratepervehicle":
		//
	case "startspervehicle":
	default:
		//unknown table selection
	}

	//aggregation container

	aggregationContainer := container.NewVBox()
	if tableSelection == "movesoutput" {
		aggregationContainer = createNewAggregationGroup(whiteList, groupBy, 1)
	} else {
		aggregationContainer = createNewAggregationGroup(whiteList, groupBy, 3)
	}

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
			value, ok := filter[whiteList[index]]
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
				for key, value := range filter {
					if len(value) == 0 {
						delete(filter, key)
					}
				}
				fmt.Println("print the map at then end of button function")
				fmt.Println(filter)
				//dummyToolBarString_Where = fmt.Sprint(filter)
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

		// enter GROUP BY claus
		// if there is 1 item is checked, have GROUP BY xxx
		// else if, >1 items are checked, call func convertColumnsComma, get comma seperated field name, then have GROUP BY xxx,xxx, .....
		// else there is 0 item is checked, remove the "GROUP BY"

		//  if 1 or more items are checked, update the select columns same as the GROUP BY, Plus the sum of emissionQuant, maybe override func getQueryResult with direct sql statement

		groupbyClause := " GROUP BY "
		var columnSelection []string
		var groupbySelection []string

		if len(groupBy["Aggregation"]) == 0 { //if nothing in the group by map
			//when select 0 checkbox, do not need group clause AND the select column should same as whiteList
			groupbyClause = ""
			columnSelection = whiteList
		} else if len(groupBy["Aggregation"]) >= 1 { //if there is anything in the group by map
			//update SELECT clause PLUS sum of emissionQuant
			//update GROUP BY
			//loop through the group by map, and copy selected box into a temp slice
			for _, value := range groupBy["Aggregation"] {
				groupbySelection = append(groupbySelection, value)
				fmt.Println("copy group by map value into slice", value)
				fmt.Println("printing updated group by slice", groupbySelection)
			}

			//pass the selected box name, PLUS sum(emissionQuant) AS emissionQuant) in the end of select statement
			//since columnSelection is []string, just throw the "sum(emissionQuant) AS emissionQuant)" in, as an element
			//getQuery will handle comma separate format later
			columnSelection = groupbySelection
			//TODO: switch here, depends on different table, sum activity? emissionQuant? or rate, how to sum rate, does it make sense? maybe average
			//TODO: add activity
			if tableSelection == "movesoutput" {
				columnSelection = append(columnSelection, "sum(emissionQuant) ")
			} else if tableSelection == "startspervehicle" {
				columnSelection = append(columnSelection, "ROUND(avg(startsPerVehicle),2) AS average_startsPerVehicle ")
			} else {
				columnSelection = append(columnSelection, "ROUND( avg(temperature),2) AS average_temperature ")
				columnSelection = append(columnSelection, "ROUND( avg(relHumidity) , 2)AS average_relHumidity ")
				columnSelection = append(columnSelection, "ROUND( avg(rateperdistance) , 2)AS average_rateperdistance ")
			}

			//pass the selected box name to GROUP BY clause, convert list of name into comma seperated format
			groupbyClause += convertColumnsComma(groupbySelection)
		} else {
			panic("detect length of groupBy map size <0, WHY")
		}

		//update the matrix with the new where clause and group by we just made
		var err error
		queryResult, err = getQueryResult(db, dbSelection, tableSelection, columnSelection, whereClause, groupbyClause)
		fmt.Println("printing error query result WHERE clause")
		fmt.Println(err)
		updateToolbarMessage(ToolbarLabel, whereClause, groupbyClause)

		//dialog box pop out warning for no result query
		if len(queryResult) < 2 {
			runPopUp(window2, "Filter combination returns no data, please try different filter")
		}

	})

	innerContainer.Add(aggregationContainer)
	innerContainer.Add(updateButton)

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

func createNewAggregationGroup(whitelist []string, groupBy map[string][]string, numericColumnsInTheEnd int) *fyne.Container {
	xButton := widget.NewButton("Aggregation", func() {
	})
	whitelist2 := whitelist
	//following table has different count of columns in the very end that always have non-filter value,
	// for example raterpervehicle has startsPerVehicle, movesoutput has emissionQuant, meanwhile rateperdistance has temperature, relHumidity and ratePerDistance
	// This is based on how we defined the struct, NOT always look at MOVES DB schema, check it in the dataType.go
	// count 1: Movesoutput, TODO: add activity
	// count 3: Rateperdistance, Rateperhour, Rateperprofile, Rateperstart, Ratepervehicle Startspervehicle
	if len(whitelist2) > 0 {
		whitelist2 = whitelist2[:len(whitelist2)-numericColumnsInTheEnd]
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

func updateToolbarMessage(l *widget.Label, where string, group string) {
	var message string
	message = "Filters: " + where + "Aggregated by : " + group
	l.SetText(message)
}

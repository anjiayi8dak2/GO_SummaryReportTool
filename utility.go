package main

import (
	"fyne.io/fyne/v2/layout"
	_ "github.com/pkg/browser"
)
import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	distanceUnits string
	massUnits     string
	energyUnits   string
	tableList     []string
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

// take a slice of string names, return one piece string with commas like roadTypeID, sourceTypeID, emissionQuant
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

	//TODO: switch here to split 7 table selections
	//map to hold filters selection in checkbox group for where clause
	filter := make(map[string][]string)
	//map to hold group by check boxes selection for group by clause
	groupBy := make(map[string][]string)
	//container sections
	innerContainer := container.NewVBox()

	//the message tab on top of screen, should update this text on the fly
	//default display for DB and Table selection
	distanceUnits = getMOVESrun(db, dbSelection, "distanceUnits")
	massUnits = getMOVESrun(db, dbSelection, "massUnits")
	energyUnits = getMOVESrun(db, dbSelection, "energyUnits")
	//var field1, field2 string
	ToolbarLabel := widget.NewLabel("DB Selection: " + dbSelection + "Table Selection: " + tableSelection + " Energy Unit: " + energyUnits + " Distance Unit: " + distanceUnits + " Mass Unit: " + massUnits)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() { //update button
			fmt.Println("I pressed update button")
			updateButtonToolbar(db, window2, tableSelection, dbSelection, whiteList, filter, groupBy, &queryResult, ToolbarLabel)
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() { //plot button TODO
			fmt.Println("I pressed plot button")
			selectAggregationField(a, queryResult)

			//runPlot(distanceUnits, massUnits, energyUnits, queryResult) //TODO: need to select two field or use first two column?

		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.DownloadIcon(), func() { //download CSV
			fmt.Println("I pressed download csv button")
			csvExport(queryResult)

		}),
		widget.NewToolbarSpacer(),
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

	case "rateperprofile":
		MOVESScenarioIDContainer := createNewCheckBoxGroup(db, "MOVESScenarioID", dbSelection, tableSelection, filter)
		MOVESRunIDContainer := createNewCheckBoxGroup(db, "MOVESRunID", dbSelection, tableSelection, filter)
		temperatureProfileIDContainer := createNewCheckBoxGroup(db, "temperatureProfileID", dbSelection, tableSelection, filter)
		yearIDContainer := createNewCheckBoxGroup(db, "yearID", dbSelection, tableSelection, filter)
		dayIDContainer := createNewCheckBoxGroup(db, "dayID", dbSelection, tableSelection, filter)
		hourIDContainer := createNewCheckBoxGroup(db, "hourID", dbSelection, tableSelection, filter)
		pollutantIDContainer := createNewCheckBoxGroup(db, "pollutantID", dbSelection, tableSelection, filter)
		processIDContainer := createNewCheckBoxGroup(db, "processID", dbSelection, tableSelection, filter)
		sourceTypeIDContainer := createNewCheckBoxGroup(db, "sourceTypeID", dbSelection, tableSelection, filter)
		regClassIDContainer := createNewCheckBoxGroup(db, "regClassID", dbSelection, tableSelection, filter)
		SCCContainer := createNewCheckBoxGroup(db, "SCC", dbSelection, tableSelection, filter)
		fuelTypeIDContainer := createNewCheckBoxGroup(db, "fuelTypeID", dbSelection, tableSelection, filter)
		modelYearIDContainer := createNewCheckBoxGroup(db, "modelYearID", dbSelection, tableSelection, filter)

		innerContainer.Add(MOVESScenarioIDContainer)
		innerContainer.Add(MOVESRunIDContainer)
		innerContainer.Add(temperatureProfileIDContainer)
		innerContainer.Add(yearIDContainer)
		innerContainer.Add(dayIDContainer)
		innerContainer.Add(hourIDContainer)
		innerContainer.Add(pollutantIDContainer)
		innerContainer.Add(processIDContainer)
		innerContainer.Add(sourceTypeIDContainer)
		innerContainer.Add(regClassIDContainer)
		innerContainer.Add(SCCContainer)
		innerContainer.Add(fuelTypeIDContainer)
		innerContainer.Add(modelYearIDContainer)

	case "rateperstart":
		MOVESScenarioIDContainer := createNewCheckBoxGroup(db, "MOVESScenarioID", dbSelection, tableSelection, filter)
		MOVESRunIDContainer := createNewCheckBoxGroup(db, "MOVESRunID", dbSelection, tableSelection, filter)
		yearIDContainer := createNewCheckBoxGroup(db, "yearID", dbSelection, tableSelection, filter)
		monthIDContainer := createNewCheckBoxGroup(db, "monthID", dbSelection, tableSelection, filter)
		dayIDContainer := createNewCheckBoxGroup(db, "dayID", dbSelection, tableSelection, filter)
		hourIDContainer := createNewCheckBoxGroup(db, "hourID", dbSelection, tableSelection, filter)
		zoneIDContainer := createNewCheckBoxGroup(db, "zoneID", dbSelection, tableSelection, filter)
		sourceTypeIDContainer := createNewCheckBoxGroup(db, "sourceTypeID", dbSelection, tableSelection, filter)
		regClassIDContainer := createNewCheckBoxGroup(db, "regClassID", dbSelection, tableSelection, filter)
		SCCContainer := createNewCheckBoxGroup(db, "SCC", dbSelection, tableSelection, filter)
		fuelTypeIDContainer := createNewCheckBoxGroup(db, "fuelTypeID", dbSelection, tableSelection, filter)
		modelYearIDContainer := createNewCheckBoxGroup(db, "modelYearID", dbSelection, tableSelection, filter)
		pollutantIDContainer := createNewCheckBoxGroup(db, "pollutantID", dbSelection, tableSelection, filter)
		processIDContainer := createNewCheckBoxGroup(db, "processID", dbSelection, tableSelection, filter)

		innerContainer.Add(MOVESScenarioIDContainer)
		innerContainer.Add(MOVESRunIDContainer)
		innerContainer.Add(yearIDContainer)
		innerContainer.Add(monthIDContainer)
		innerContainer.Add(dayIDContainer)
		innerContainer.Add(hourIDContainer)
		innerContainer.Add(zoneIDContainer)
		innerContainer.Add(sourceTypeIDContainer)
		innerContainer.Add(regClassIDContainer)
		innerContainer.Add(SCCContainer)
		innerContainer.Add(fuelTypeIDContainer)
		innerContainer.Add(modelYearIDContainer)
		innerContainer.Add(pollutantIDContainer)
		innerContainer.Add(processIDContainer)

	case "ratepervehicle":
		//, , , , , , , , , , , , , , temperature, relHumidity, ratePerVehicle
		MOVESScenarioIDContainer := createNewCheckBoxGroup(db, "MOVESScenarioID", dbSelection, tableSelection, filter)
		MOVESRunIDContainer := createNewCheckBoxGroup(db, "MOVESRunID", dbSelection, tableSelection, filter)
		yearIDContainer := createNewCheckBoxGroup(db, "yearID", dbSelection, tableSelection, filter)
		monthIDContainer := createNewCheckBoxGroup(db, "monthID", dbSelection, tableSelection, filter)
		dayIDContainer := createNewCheckBoxGroup(db, "dayID", dbSelection, tableSelection, filter)
		hourIDContainer := createNewCheckBoxGroup(db, "hourID", dbSelection, tableSelection, filter)
		zoneIDContainer := createNewCheckBoxGroup(db, "zoneID", dbSelection, tableSelection, filter)
		pollutantIDContainer := createNewCheckBoxGroup(db, "pollutantID", dbSelection, tableSelection, filter)
		processIDContainer := createNewCheckBoxGroup(db, "processID", dbSelection, tableSelection, filter)
		sourceTypeIDContainer := createNewCheckBoxGroup(db, "sourceTypeID", dbSelection, tableSelection, filter)
		regClassIDContainer := createNewCheckBoxGroup(db, "regClassID", dbSelection, tableSelection, filter)
		SCCContainer := createNewCheckBoxGroup(db, "SCC", dbSelection, tableSelection, filter)
		fuelTypeIDContainer := createNewCheckBoxGroup(db, "fuelTypeID", dbSelection, tableSelection, filter)
		modelYearIDContainer := createNewCheckBoxGroup(db, "modelYearID", dbSelection, tableSelection, filter)

		innerContainer.Add(MOVESScenarioIDContainer)
		innerContainer.Add(MOVESRunIDContainer)
		innerContainer.Add(yearIDContainer)
		innerContainer.Add(monthIDContainer)
		innerContainer.Add(dayIDContainer)
		innerContainer.Add(hourIDContainer)
		innerContainer.Add(zoneIDContainer)
		innerContainer.Add(pollutantIDContainer)
		innerContainer.Add(processIDContainer)
		innerContainer.Add(sourceTypeIDContainer)
		innerContainer.Add(regClassIDContainer)
		innerContainer.Add(SCCContainer)
		innerContainer.Add(fuelTypeIDContainer)
		innerContainer.Add(modelYearIDContainer)

	case "startspervehicle":
		MOVESScenarioIDContainer := createNewCheckBoxGroup(db, "MOVESScenarioID", dbSelection, tableSelection, filter)
		MOVESRunIDContainer := createNewCheckBoxGroup(db, "MOVESRunID", dbSelection, tableSelection, filter)
		yearIDContainer := createNewCheckBoxGroup(db, "yearID", dbSelection, tableSelection, filter)
		monthIDContainer := createNewCheckBoxGroup(db, "monthID", dbSelection, tableSelection, filter)
		dayIDContainer := createNewCheckBoxGroup(db, "dayID", dbSelection, tableSelection, filter)
		hourIDContainer := createNewCheckBoxGroup(db, "hourID", dbSelection, tableSelection, filter)
		zoneIDContainer := createNewCheckBoxGroup(db, "zoneID", dbSelection, tableSelection, filter)
		sourceTypeIDContainer := createNewCheckBoxGroup(db, "sourceTypeID", dbSelection, tableSelection, filter)
		regClassIDContainer := createNewCheckBoxGroup(db, "regClassID", dbSelection, tableSelection, filter)
		SCCContainer := createNewCheckBoxGroup(db, "SCC", dbSelection, tableSelection, filter)
		fuelTypeIDContainer := createNewCheckBoxGroup(db, "fuelTypeID", dbSelection, tableSelection, filter)
		modelYearIDContainer := createNewCheckBoxGroup(db, "modelYearID", dbSelection, tableSelection, filter)

		innerContainer.Add(MOVESScenarioIDContainer)
		innerContainer.Add(MOVESRunIDContainer)
		innerContainer.Add(yearIDContainer)
		innerContainer.Add(monthIDContainer)
		innerContainer.Add(dayIDContainer)
		innerContainer.Add(hourIDContainer)
		innerContainer.Add(zoneIDContainer)
		innerContainer.Add(sourceTypeIDContainer)
		innerContainer.Add(regClassIDContainer)
		innerContainer.Add(SCCContainer)
		innerContainer.Add(fuelTypeIDContainer)
		innerContainer.Add(modelYearIDContainer)
	default:
		//unknown table selection
	}

	//aggregation container
	aggregationContainer := container.NewVBox()
	if tableSelection == "movesoutput" || tableSelection == "startspervehicle" { //these two table have 1 numeric column in the end that shows result
		aggregationContainer = createNewAggregationGroup(whiteList, groupBy, 1)
	} else {
		aggregationContainer = createNewAggregationGroup(whiteList, groupBy, 3)
	}

	// TODO: temporary disable aggregation for rate
	if tableSelection != "movesoutput" {
		aggregationContainer.Hide()
	}
	innerContainer.Add(aggregationContainer)

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

// open new window when hit the plot button, user should select 1 or 2 field for plotting
// then this function will pass all the parameter to the plotting library
func selectAggregationField(a fyne.App, queryResult [][]string) {
	var fieldSelection1, fieldSelection2, pollutantSelection string

	selectAggregationFieldWindow := a.NewWindow("Select 1 or 2 field that for X axis")
	selectAggregationFieldWindow.Resize(fyne.NewSize(400, 400))

	//get the list of field columns and the result column
	//field columns serve for dropdown box
	//iteration for the first row of the data grid, first row must be header.
	// TODO: #1 selection shall be full list
	// #2 list should be full list - the first selection

	headersList := queryResult[0]
	var headerList2 []string
	var resultColumn string

	if len(headersList) > 0 {
		//assign last element in the header into resultColumn before delete
		resultColumn = headersList[len(headersList)-1]
		//remove the last element in the header, this should be usually be result column such as activity or emissionQuant
		headersList = headersList[:len(headersList)-1]

		headerList2 = headersList
	}

	//TODO: Select pollutant

	tableList := []string{"pollutant1", "pollutant2", "pollutant3", "pollutant4"}
	//Create dropdown for pollutant
	pollutantSelectionResult := widget.NewLabel("Select A Pollutant")
	//pollutant dropdown box option
	pollutantSelectionDropdown := widget.NewSelect(
		tableList,
		func(selection string) {
			fmt.Printf("I selected %selection as pollutant..", selection)
			pollutantSelectionResult.Text = selection
			pollutantSelection = selection
			pollutantSelectionResult.Refresh()
		})

	//Create dropdown for field selection #1
	fieldSelectionResult1 := widget.NewLabel("Select field 1")
	//Use headersList to update dropdown box option
	fieldSelectionDropdown1 := widget.NewSelect(
		headersList,
		func(selection string) {
			fmt.Printf("I selected %selection as field 1..", selection)
			fieldSelectionResult1.Text = selection
			fieldSelection1 = selection
			fieldSelectionResult1.Refresh()
		})

	//Create dropdown for field selection #2
	fieldSelectionResult2 := widget.NewLabel("Select field 2")
	//Use header list EXCLUDE first selection to update dropdown box option
	headerList2 = RemoveElementFromSlice(headerList2, fieldSelection1)
	fieldSelectionDropdown2 := widget.NewSelect(
		headerList2,
		func(selection string) {
			fieldSelectionResult2.Refresh()
			fmt.Printf("I selected %selection as field 2..", selection)
			fieldSelectionResult2.Text = selection
			fieldSelection2 = selection
			fieldSelectionResult2.Refresh()
		})

	//TODO: put ok button
	submitButton := widget.NewButton("Submit", func() {
		fmt.Println("Submit button pressed")
		fmt.Println("Printing field selection #1  " + fieldSelection1 + " field #2 " + fieldSelection2 + " value column " + resultColumn + " pollutant selection " + pollutantSelection)
		runPlot(distanceUnits, massUnits, energyUnits, pollutantSelection, fieldSelection1, fieldSelection2, resultColumn, queryResult)
	})

	cancelButton := widget.NewButton("Cancel", func() {
		fmt.Println("Cancel button pressed")
		selectAggregationFieldWindow.Close()
	})
	buttonContainer := container.New(layout.NewGridLayout(2), submitButton, cancelButton)

	dropdownGrid := container.New(layout.NewGridLayout(3), pollutantSelectionDropdown, fieldSelectionDropdown1, fieldSelectionDropdown2)
	outerContainer := container.NewVSplit(dropdownGrid, buttonContainer)
	outerContainer.Offset = 0.8
	selectAggregationFieldWindow.SetContent(outerContainer)

	selectAggregationFieldWindow.Show()

}

// TODO: move update button on top tool bar, not all the way in the bottom of scrollbar
func updateButtonToolbar(db *sql.DB, window2 fyne.Window, tableSelection string, dbSelection string, whiteList []string, filter map[string][]string,
	groupBy map[string][]string, queryResult *[][]string, ToolbarLabel *widget.Label) {
	fmt.Println("print from update button function")

	whereClause := " WHERE "
	fmt.Println("pressed UPDATE button")
	//loop through all the keys in mo map and generate a where clause
	fmt.Println("loop throuth the whiteList keys and print")
	// there is no easy loop solution for select and then unselect operation on the run time, because it will generate a key with empty value such as hpid {}
	// these empty values will cause empty IN() statement in the where clause that make future problems.
	// hence, we should detect empty value in map and delete that key before disaster happen
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
		//update GROUP BY
		//loop through the group by map, and copy selected box into a temp slice
		for _, value := range groupBy["Aggregation"] {
			groupbySelection = append(groupbySelection, value)
			fmt.Println("copy group by map value into slice", value)
			fmt.Println("printing updated group by slice", groupbySelection)
		}

		//update SELECT clause PLUS sum of emissionQuant or acivity or rates
		columnSelection = groupbySelection
		//TODO: switch here, depends on different table, sum activity? emissionQuant?
		//TODO: disable rate aggregation for now, averaging/summing rate can be misleading
		// because it is not considered many factors such as population distribution, and all kinds of adjustments.
		if tableSelection == "movesoutput" {
			columnSelection = append(columnSelection, "sum(emissionQuant) ")
		} else if tableSelection == "startspervehicle" {
			columnSelection = append(columnSelection, "ROUND(avg(startsPerVehicle),2) AS average_startsPerVehicle ")
		} else { //this should include "rateperdistance", "rateperhour", "rateperprofile", "rateperstart", and "ratepervehicle", because they all have temperature and relHumidity columns
			columnSelection = append(columnSelection, "ROUND( avg(temperature),2) AS average_temperature ")
			columnSelection = append(columnSelection, "ROUND( avg(relHumidity) , 2)AS average_relHumidity ")
			//then add the last column rateperxxx to the end
			switch tableSelection {
			case "rateperdistance":
				columnSelection = append(columnSelection, "ROUND( avg(rateperdistance) , 2)AS average_rateperdistance ")
			case "rateperhour":
				columnSelection = append(columnSelection, "ROUND( avg(rateperhour) , 2)AS average_rateperhour ")
			case "rateperprofile":
				columnSelection = append(columnSelection, "ROUND( avg(rateperprofile) , 2)AS average_rateperprofile ")
			case "rateperstart":
				columnSelection = append(columnSelection, "ROUND( avg(rateperstart) , 2)AS average_rateperstart ")
			case "ratepervehicle":
				columnSelection = append(columnSelection, "ROUND( avg(ratepervehicle) , 2)AS average_ratepervehicle ")
			}
		}

		//pass the selected box name to GROUP BY clause, convert list of name into comma seperated format
		groupbyClause += convertColumnsComma(groupbySelection)
	} else {
		panic("detect length of groupBy map size <0, WHY")
	}

	//update the matrix with the new where clause and group by we just made
	var err error
	*queryResult, err = getQueryResult(db, dbSelection, tableSelection, columnSelection, whereClause, groupbyClause)
	fmt.Println("printing error query result WHERE clause")
	fmt.Println(err)
	updateToolbarMessage(ToolbarLabel, whereClause, groupbyClause, db, dbSelection)

	//dialog box pop out warning for no result query
	if len(*queryResult) < 2 {
		runPopUp(window2, "Filter combination returns no data, please try different filter")
	}

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

func updateToolbarMessage(l *widget.Label, where string, group string, db *sql.DB, dbSelection string) {
	var message string
	distanceUnits := getMOVESrun(db, dbSelection, "distanceUnits")
	massUnits := getMOVESrun(db, dbSelection, "massUnits")
	energyUnits := getMOVESrun(db, dbSelection, "energyUnits")

	message = "Filters: " + where + "Aggregated by : " + group + " Energy Unit: " + energyUnits + " Distance Unit: " + distanceUnits + " Mass Unit: " + massUnits
	l.SetText(message)
}

// pass distanceUnits/massUnits/energyUnits in string, return unit name in string such as "kg" or "mile"
func getMOVESrun(db *sql.DB, dbSelection string, columnName string) string {
	var query string
	var unit string
	query = "SELECT " + columnName + " FROM " + dbSelection + ".movesrun LIMIT 1;"
	db.QueryRow(query).Scan(&unit)

	return unit
}

func mapCopy[M1, M2 ~map[K]V, K comparable, V any](dst M1, src M2) {
	for k, v := range src {
		dst[k] = v
	}
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func removeDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// Generic solution
func removeDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

package main

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	_ "github.com/pkg/browser"
	"github.com/skratchdot/open-golang/open" // Open the folder in the file explorer, also handle different OS
	"strings"
	"time"
)

var (
	distanceUnits   string
	massUnits       string
	energyUnits     string
	dbSelection     string
	tableSelection  string
	whiteListIndex  []bool
	whiteList       []string
	fieldSelection1 string
	fieldSelection2 string
)

// TODO : not working at all, why open the document folder by default?
func openMariaFolder(db *sql.DB) {
	//dataDir := getDataDir(db)
	//dataDir = "\"" + dataDir + "\""
	//fmt.Println("print data dir", dataDir)

	// Open the folder in the file explorer
	err := open.Start("C:\\ProgramData\\MariaDB\\MariaDB 10.4\\data\\")
	fmt.Println("print err", err)
	if err != nil {
		fmt.Println("print err", err)
	}

}

// take a slice of string names, return one piece string with comma seperated format for example, dummy = {"roadTypeID", "sourceTypeID", "emissionQuant"} ==> "roadTypeID, sourceTypeID, emissionQuant"
func convertColumnsComma(columns []string) string {
	// prepend single quote, perform joins, append single quote
	ColumnsComma := strings.Join(columns, `,`)
	return ColumnsComma
}

// data browsing main window, include the data table and filters
func makeWindowTwo(a fyne.App, queryResult [][]string, db *sql.DB) {
	fmt.Println("opening Summary Report Tool Data Viewing Window")
	window2 := a.NewWindow("Data Viewing Window")
	window2.SetContent(widget.NewLabel("Data Viewing Window label"))
	window2.Resize(fyne.NewSize(1000, 800))

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

	tableData := widget.NewTable(
		func() (int, int) {
			return len(queryResult), len(queryResult[0]) // row size, columns size
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if len(queryResult) >= 2 { //if there is any data other than header
				o.(*widget.Label).SetText(queryResult[i.Row][i.Col])
			} else { //otherwise fill cells with "no data"
				o.(*widget.Label).SetText("no data")
			}
		})

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() { //update button
			fmt.Println("I pressed update button")
			updateButtonToolbar(db, window2, filter, groupBy, &queryResult, ToolbarLabel)
			tableAutoSize(queryResult, tableData)
			decodeDataTable(queryResult)
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() { //plot button
			fmt.Println("I pressed plot button")
			//check before plot if there are more than 2 columns are selected
			if len(queryResult[0]) > 3 {
				//popup a dialog box
				dialog.ShowInformation(
					"Hold Up",
					"Please make sure only check 1 or 2 checkbox in the aggregation section for plotting",
					window2,
				)
			} else {
				selectPlotColumn(a, queryResult)
			}

		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.DownloadIcon(), func() { //download CSV
			fmt.Println("I pressed download csv button")
			csvExport(queryResult)
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarSpacer(),
		ToolbarLabel,
	)

	createFilterButtons(db, filter, innerContainer)

	//aggregation container
	aggregationContainer := container.NewVBox()
	if tableSelection == "movesoutput" || tableSelection == "startspervehicle" || tableSelection == "movesactivityoutput" { //these 3 table have 1 numeric column in the end that shows result
		aggregationContainer = createNewAggregationGroup(whiteList, groupBy, 1)
	} else {
		aggregationContainer = createNewAggregationGroup(whiteList, groupBy, 3)
	}

	// TODO: temporary disable aggregation for rate, sum rate together does not make sense
	if tableSelection == "movesoutput" || tableSelection == "movesactivityoutput" {
		aggregationContainer.Visible()
	} else {
		aggregationContainer.Hide()
	}
	innerContainer.Add(aggregationContainer)

	//dynamic filter buttons, Use the record of whiteListIndex [] bool, show and hide base on 1 or 0.
	//we initialized all columns when the Data Viewing Window started
	for index, ok := range whiteListIndex {
		if ok {
			innerContainer.Objects[index].Visible()
		} else {
			innerContainer.Objects[index].Hide()
		}
	}

	// the filter button section scroll bar on the left, this is different one than the data table scroll bar
	scrollContainer := container.NewVScroll(innerContainer)
	outerContainer := container.NewHSplit(
		scrollContainer,
		tableData,
	)
	//screen width horizontal distribution of filter panel VS data table panel
	outerContainer.Offset = 0.08
	window2.SetContent(container.NewBorder(toolbar, nil, nil, nil, outerContainer))
	window2.Show()
}

func createFilterButtons(db *sql.DB, filter map[string][]string, innerContainer *fyne.Container) {
	switch tableSelection {
	//TODO: loop smart, something loop through struct field? instead of hardcoded
	case "movesactivityoutput":
		MOVESRunIDContainer := createNewCheckBoxGroup(db, "MOVESRunID", filter)
		iterationIDContainer := createNewCheckBoxGroup(db, "iterationID", filter)
		yearIDContainer := createNewCheckBoxGroup(db, "yearID", filter)
		monthIDContainer := createNewCheckBoxGroup(db, "monthID", filter)
		dayIDContainer := createNewCheckBoxGroup(db, "dayID", filter)
		hourIDContainer := createNewCheckBoxGroup(db, "hourID", filter)
		stateIDContainer := createNewCheckBoxGroup(db, "stateID", filter)
		countyIDContainer := createNewCheckBoxGroup(db, "countyID", filter)
		zoneIDContainer := createNewCheckBoxGroup(db, "zoneID", filter)
		linkIDContainer := createNewCheckBoxGroup(db, "linkID", filter)
		sourceTypeIDContainer := createNewCheckBoxGroup(db, "sourceTypeID", filter)
		regClassIDContainer := createNewCheckBoxGroup(db, "regClassID", filter)
		fuelTypeIDContainer := createNewCheckBoxGroup(db, "fuelTypeID", filter)
		fuelSubTypeIDContainer := createNewCheckBoxGroup(db, "fuelSubTypeID", filter)
		modelYearContainer := createNewCheckBoxGroup(db, "modelYearID", filter)
		roadTypeIDContainer := createNewCheckBoxGroup(db, "roadTypeID", filter)
		SCCContainer := createNewCheckBoxGroup(db, "SCC", filter)
		engTechIDContainer := createNewCheckBoxGroup(db, "engTechID", filter)
		sectorIDContainer := createNewCheckBoxGroup(db, "sectorID", filter)
		hpIDContainer := createNewCheckBoxGroup(db, "hpID", filter)
		activityTypeID := createNewCheckBoxGroup(db, "activityTypeID", filter)

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
		innerContainer.Add(activityTypeID)

	case "movesoutput":
		//fyne containers, create buttons for filters with checkbox selection saved in the map filter
		MOVESRunIDContainer := createNewCheckBoxGroup(db, "MOVESRunID", filter)
		iterationIDContainer := createNewCheckBoxGroup(db, "iterationID", filter)
		yearIDContainer := createNewCheckBoxGroup(db, "yearID", filter)
		monthIDContainer := createNewCheckBoxGroup(db, "monthID", filter)
		dayIDContainer := createNewCheckBoxGroup(db, "dayID", filter)
		hourIDContainer := createNewCheckBoxGroup(db, "hourID", filter)
		stateIDContainer := createNewCheckBoxGroup(db, "stateID", filter)
		countyIDContainer := createNewCheckBoxGroup(db, "countyID", filter)
		zoneIDContainer := createNewCheckBoxGroup(db, "zoneID", filter)
		linkIDContainer := createNewCheckBoxGroup(db, "linkID", filter)
		pollutantContainer := createNewCheckBoxGroup(db, "pollutantID", filter)
		processIDContainer := createNewCheckBoxGroup(db, "processID", filter)
		sourceTypeIDContainer := createNewCheckBoxGroup(db, "sourceTypeID", filter)
		regClassIDContainer := createNewCheckBoxGroup(db, "regClassID", filter)
		fuelTypeIDContainer := createNewCheckBoxGroup(db, "fuelTypeID", filter)
		fuelSubTypeIDContainer := createNewCheckBoxGroup(db, "fuelSubTypeID", filter)
		modelYearContainer := createNewCheckBoxGroup(db, "modelYearID", filter)
		roadTypeIDContainer := createNewCheckBoxGroup(db, "roadTypeID", filter)
		SCCContainer := createNewCheckBoxGroup(db, "SCC", filter)
		engTechIDContainer := createNewCheckBoxGroup(db, "engTechID", filter)
		sectorIDContainer := createNewCheckBoxGroup(db, "sectorID", filter)
		hpIDContainer := createNewCheckBoxGroup(db, "hpID", filter)

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
		MOVESScenarioIDContainer := createNewCheckBoxGroup(db, "MOVESScenarioID", filter)
		MOVESRunIDContainer := createNewCheckBoxGroup(db, "MOVESRunID", filter)
		yearIDContainer := createNewCheckBoxGroup(db, "yearID", filter)
		monthIDContainer := createNewCheckBoxGroup(db, "monthID", filter)
		dayIDContainer := createNewCheckBoxGroup(db, "dayID", filter)
		hourIDContainer := createNewCheckBoxGroup(db, "hourID", filter)
		linkIDContainer := createNewCheckBoxGroup(db, "linkID", filter)
		pollutantIDContainer := createNewCheckBoxGroup(db, "pollutantID", filter)
		processIDContainer := createNewCheckBoxGroup(db, "processID", filter)
		sourceTypeIDContainer := createNewCheckBoxGroup(db, "sourceTypeID", filter)
		regClassIDContainer := createNewCheckBoxGroup(db, "regClassID", filter)
		SCCContainer := createNewCheckBoxGroup(db, "SCC", filter)
		fuelTypeIDContainer := createNewCheckBoxGroup(db, "fuelTypeID", filter)
		modelYearIDContainer := createNewCheckBoxGroup(db, "modelYearID", filter)
		roadTypeIDContainer := createNewCheckBoxGroup(db, "roadTypeID", filter)
		avgSpeedBinIDContainer := createNewCheckBoxGroup(db, "avgSpeedBinID", filter)

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
		MOVESScenarioIDContainer := createNewCheckBoxGroup(db, "MOVESScenarioID", filter)
		MOVESRunIDContainer := createNewCheckBoxGroup(db, "MOVESRunID", filter)
		yearIDContainer := createNewCheckBoxGroup(db, "yearID", filter)
		monthIDContainer := createNewCheckBoxGroup(db, "monthID", filter)
		dayIDContainer := createNewCheckBoxGroup(db, "dayID", filter)
		hourIDContainer := createNewCheckBoxGroup(db, "hourID", filter)
		linkIDContainer := createNewCheckBoxGroup(db, "linkID", filter)
		pollutantIDContainer := createNewCheckBoxGroup(db, "pollutantID", filter)
		processIDContainer := createNewCheckBoxGroup(db, "processID", filter)
		sourceTypeIDContainer := createNewCheckBoxGroup(db, "sourceTypeID", filter)
		regClassIDContainer := createNewCheckBoxGroup(db, "regClassID", filter)
		SCCContainer := createNewCheckBoxGroup(db, "SCC", filter)
		fuelTypeIDContainer := createNewCheckBoxGroup(db, "fuelTypeID", filter)
		modelYearIDContainer := createNewCheckBoxGroup(db, "modelYearID", filter)
		roadTypeIDContainer := createNewCheckBoxGroup(db, "roadTypeID", filter)

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
		MOVESScenarioIDContainer := createNewCheckBoxGroup(db, "MOVESScenarioID", filter)
		MOVESRunIDContainer := createNewCheckBoxGroup(db, "MOVESRunID", filter)
		temperatureProfileIDContainer := createNewCheckBoxGroup(db, "temperatureProfileID", filter)
		yearIDContainer := createNewCheckBoxGroup(db, "yearID", filter)
		dayIDContainer := createNewCheckBoxGroup(db, "dayID", filter)
		hourIDContainer := createNewCheckBoxGroup(db, "hourID", filter)
		pollutantIDContainer := createNewCheckBoxGroup(db, "pollutantID", filter)
		processIDContainer := createNewCheckBoxGroup(db, "processID", filter)
		sourceTypeIDContainer := createNewCheckBoxGroup(db, "sourceTypeID", filter)
		regClassIDContainer := createNewCheckBoxGroup(db, "regClassID", filter)
		SCCContainer := createNewCheckBoxGroup(db, "SCC", filter)
		fuelTypeIDContainer := createNewCheckBoxGroup(db, "fuelTypeID", filter)
		modelYearIDContainer := createNewCheckBoxGroup(db, "modelYearID", filter)

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
		MOVESScenarioIDContainer := createNewCheckBoxGroup(db, "MOVESScenarioID", filter)
		MOVESRunIDContainer := createNewCheckBoxGroup(db, "MOVESRunID", filter)
		yearIDContainer := createNewCheckBoxGroup(db, "yearID", filter)
		monthIDContainer := createNewCheckBoxGroup(db, "monthID", filter)
		dayIDContainer := createNewCheckBoxGroup(db, "dayID", filter)
		hourIDContainer := createNewCheckBoxGroup(db, "hourID", filter)
		zoneIDContainer := createNewCheckBoxGroup(db, "zoneID", filter)
		sourceTypeIDContainer := createNewCheckBoxGroup(db, "sourceTypeID", filter)
		regClassIDContainer := createNewCheckBoxGroup(db, "regClassID", filter)
		SCCContainer := createNewCheckBoxGroup(db, "SCC", filter)
		fuelTypeIDContainer := createNewCheckBoxGroup(db, "fuelTypeID", filter)
		modelYearIDContainer := createNewCheckBoxGroup(db, "modelYearID", filter)
		pollutantIDContainer := createNewCheckBoxGroup(db, "pollutantID", filter)
		processIDContainer := createNewCheckBoxGroup(db, "processID", filter)

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
		MOVESScenarioIDContainer := createNewCheckBoxGroup(db, "MOVESScenarioID", filter)
		MOVESRunIDContainer := createNewCheckBoxGroup(db, "MOVESRunID", filter)
		yearIDContainer := createNewCheckBoxGroup(db, "yearID", filter)
		monthIDContainer := createNewCheckBoxGroup(db, "monthID", filter)
		dayIDContainer := createNewCheckBoxGroup(db, "dayID", filter)
		hourIDContainer := createNewCheckBoxGroup(db, "hourID", filter)
		zoneIDContainer := createNewCheckBoxGroup(db, "zoneID", filter)
		pollutantIDContainer := createNewCheckBoxGroup(db, "pollutantID", filter)
		processIDContainer := createNewCheckBoxGroup(db, "processID", filter)
		sourceTypeIDContainer := createNewCheckBoxGroup(db, "sourceTypeID", filter)
		regClassIDContainer := createNewCheckBoxGroup(db, "regClassID", filter)
		SCCContainer := createNewCheckBoxGroup(db, "SCC", filter)
		fuelTypeIDContainer := createNewCheckBoxGroup(db, "fuelTypeID", filter)
		modelYearIDContainer := createNewCheckBoxGroup(db, "modelYearID", filter)

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
		MOVESScenarioIDContainer := createNewCheckBoxGroup(db, "MOVESScenarioID", filter)
		MOVESRunIDContainer := createNewCheckBoxGroup(db, "MOVESRunID", filter)
		yearIDContainer := createNewCheckBoxGroup(db, "yearID", filter)
		monthIDContainer := createNewCheckBoxGroup(db, "monthID", filter)
		dayIDContainer := createNewCheckBoxGroup(db, "dayID", filter)
		hourIDContainer := createNewCheckBoxGroup(db, "hourID", filter)
		zoneIDContainer := createNewCheckBoxGroup(db, "zoneID", filter)
		sourceTypeIDContainer := createNewCheckBoxGroup(db, "sourceTypeID", filter)
		regClassIDContainer := createNewCheckBoxGroup(db, "regClassID", filter)
		SCCContainer := createNewCheckBoxGroup(db, "SCC", filter)
		fuelTypeIDContainer := createNewCheckBoxGroup(db, "fuelTypeID", filter)
		modelYearIDContainer := createNewCheckBoxGroup(db, "modelYearID", filter)

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
}

// open new window when hit the plot button, user should select 1 or 2 field for plotting
// then this function will pass the selection to plot function
func selectPlotColumn(a fyne.App, queryResult [][]string) {

	selectAggregationFieldWindow := a.NewWindow("Plot Configuration Window")
	selectAggregationFieldWindow.Resize(fyne.NewSize(400, 400))

	// get the list of field columns and the result column, for example dayID, fuelID, emissionQuant
	// field columns serve for 2 dropdown boxes
	// user select what is for X1 and X2
	// copy for the first row of the data grid, first row must be header.
	// headersList [] string for dropdown box

	headersList := queryResult[0]
	var resultColumn string //the value column names such as emissionQuant , activity or gramPerDistance

	if len(headersList) > 0 {
		//assign last element in the header into resultColumn before delete
		resultColumn = headersList[len(headersList)-1]
		//remove the last element in the header, this should be result column such as activity or emissionQuant
		headersList = headersList[:len(headersList)-1]
	}

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
	//text show on the dropdown box before click
	fieldSelectionDropdown1.PlaceHolder = "Select field X1"

	//Create dropdown for field selection #2
	fieldSelectionResult2 := widget.NewLabel("Select field 2")
	fieldSelectionDropdown2 := widget.NewSelect(
		headersList,
		func(selection string) {
			fieldSelectionResult2.Refresh()
			fmt.Printf("I selected %selection as field 2..", selection)
			fieldSelectionResult2.Text = selection
			fieldSelection2 = selection
			fieldSelectionResult2.Refresh()
		})
	//text show on the dropdown box before click
	fieldSelectionDropdown2.PlaceHolder = "Select field X2"

	submitButton := widget.NewButton("Submit", func() {
		fmt.Println("Submit button pressed")
		fmt.Println("Printing field selection #1  " + fieldSelection1 + " field #2 " + fieldSelection2 + " value column " + resultColumn)
		runPlot(distanceUnits, massUnits, energyUnits, fieldSelection1, fieldSelection2, resultColumn, queryResult)
	})

	cancelButton := widget.NewButton("Cancel", func() {
		fmt.Println("Cancel button pressed")
		selectAggregationFieldWindow.Close()
	})
	// layout setting
	// 1 row, 2 columns for submit and cancel grid
	buttonContainer := container.New(layout.NewGridLayout(2), submitButton, cancelButton)
	// 1 row, 2 columns for 2 dropdown box grid
	dropdownGrid := container.New(layout.NewGridLayout(2), fieldSelectionDropdown1, fieldSelectionDropdown2)
	// vertical split container 2 rows, 1 column, and it split up and down
	outerContainer := container.NewVSplit(dropdownGrid, buttonContainer)
	// upstairs own 80% of the space, downstairs own 20% space
	outerContainer.Offset = 0.8
	//fyne standard, set content to window and show
	selectAggregationFieldWindow.SetContent(outerContainer)
	selectAggregationFieldWindow.Show()

}

func updateButtonToolbar(db *sql.DB, window2 fyne.Window, filter map[string][]string,
	groupBy map[string][]string, queryResult *[][]string, ToolbarLabel *widget.Label) {
	fmt.Println("print from update button function")

	whereClause := " WHERE "
	fmt.Println("pressed UPDATE button")
	//loop through all the keys in mo map and generate a where clause
	fmt.Println("loop throuth the whiteList keys and print")
	// there is no easy loop solution for select and then unselect operation on the run time, because it will generate a key with empty value such as hpid {}
	// these empty values will cause empty IN() statement in the where clause that make sql syntax problems.
	// hence, detect empty value in map and delete that key before disaster happen
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
				fmt.Println("print partialWhere clause ", partialWhere)
			} else { // otherwise do not put AND
				inValue := convertColumnsComma(value)
				fmt.Println("print in values ", inValue)
				partialWhere = whiteList[index] + " IN ( " + inValue + " ) "
				fmt.Println("print partialWhere clause ", partialWhere)
			}
		} else {
			fmt.Println(whiteList[index], " key not found")
			//delete map that has empty value, they looks like hpID:[], they will eventually cause an empty IN() in the where clause
			for key, filterValue := range filter {
				if len(filterValue) == 0 {
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

	// enter GROUP BY claus
	// if there is 1 aggregation box is checked, have GROUP BY xxx
	// else if, >1 aggregation boxes are checked, call func convertColumnsComma, get comma seperated field name, then have GROUP BY xxx,xxx, .....
	// else there is 0 aggregation box is checked, remove the "GROUP BY" from the sql statement
	//  if 1 or more items are checked, update the select columns same as the GROUP BY, Plus the sum of emissionQuant

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
		//switch here, depends on different table, sum activity? emissionQuant?
		//TODO: disable rate aggregation for now, averaging/summing rate can be misleading
		// because it is not considered many factors such as population distribution, and all kinds of adjustments.
		if tableSelection == "movesoutput" {
			columnSelection = append(columnSelection, "sum(emissionQuant) ")
		} else if tableSelection == "movesactivityoutput" {
			columnSelection = append(columnSelection, "sum(activity) ")
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
	*queryResult, err = getQueryResult(db, columnSelection, whereClause, groupbyClause)
	fmt.Println("printing error query result WHERE clause")
	fmt.Println(err)
	updateToolbarMessage(ToolbarLabel, whereClause, groupbyClause, db, dbSelection)
	fmt.Println("printing query result re plot")
	fmt.Println(queryResult)

	//dialog box pop out warning for no result query
	if len(*queryResult) < 2 {
		runPopUp(window2, "Filter combination returns no data, please try different filter")
	}

}

func createNewAggregationGroup(whitelist []string, groupBy map[string][]string, numericColumnsInTheEnd int) *fyne.Container {
	xButton := widget.NewButton("Aggregation", func() {
	})
	whitelist2 := whitelist
	//following table has different count of columns in the very end that always have a value column,
	// for example raterpervehicle has startsPerVehicle, movesoutput has emissionQuant, meanwhile rateperdistance has temperature, relHumidity and ratePerDistance
	// This is based on how we defined the struct, NOT always look at MOVES DB schema, check it in the dataType.go
	// count 1: Movesoutput, activity
	// count 3: Rateperdistance, Rateperhour, Rateperprofile, Rateperstart, Ratepervehicle Startspervehicle

	// This aggregation buttons are not for value column, remove these columns from end of the slice, base on how many value columns a table has
	if len(whitelist2) > 0 {
		whitelist2 = whitelist2[:len(whitelist2)-numericColumnsInTheEnd]
	}
	fmt.Println("printing whitelist2 slice inside createNewAggregationGroup", whitelist2)
	// create checkbox button group
	xCheckGroup := widget.NewCheckGroup(whitelist2, func(value []string) {
		fmt.Println("selected", value)
		//update map from checked boxes selection statues
		groupBy["Aggregation"] = value
		fmt.Println("print entire group by map for  ", "Aggregation", " inside func createNewAggregationGroup")
		fmt.Println(groupBy)

		//check if map key combination is empty, if it is empty, delete that key combo
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

// TODO performance?
// base on the column name passed, generate a group of checkbox, what checkbox and how many checkbox will depend on how many distinct value that column has
func createNewCheckBoxGroup(db *sql.DB, columnsName string, filter map[string][]string) *fyne.Container {
	//To these filters suppose to have group of checkbox
	//CheckGroup
	//= pollutantContainer
	//= title button + checkbox group in vertical
	//For example
	//pollutantidButton + pollutantContainer
	xButton := widget.NewButton(columnsName, func() {
		//TODO: expand & collapse on click at checkbox group title
	})

	// TODO: How distinct query speed performance on the null value?
	//  we already have white list, search the columnsName that passed in to this function
	//  search it in the whitelist[], if it is not exist that means it is null column and it can be skipped

	//get list of distinct values from a ID column
	distinctX := getDistinct(db, columnsName)
	// TODO: the value here = the checkbox name, how to show full name but when select value by ID? for example I want checkbox show as fuelType 'gas' but value stay as 1
	//  The fuelType =1 is the way to query the filter, but if I cange the checkbox name into 'gas' the value I got in the backend will became gas.

	//
	xCheckGroup := widget.NewCheckGroup(distinctX, func(value []string) {
		fmt.Println("selected", value)
		//update map  from checked boxes statues
		filter[columnsName] = value
		fmt.Println("print entire filter map for  ", columnsName, " inside func createNewCheckBoxGroup")
		fmt.Println(filter)
	})

	xContainer := container.NewVBox(xButton, xCheckGroup)
	return xContainer
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

func updateToolbarMessage(l *widget.Label, where string, group string, db *sql.DB, dbSelection string) {
	var message string
	distanceUnits = getMOVESrun(db, dbSelection, "distanceUnits")
	massUnits = getMOVESrun(db, dbSelection, "massUnits")
	energyUnits = getMOVESrun(db, dbSelection, "energyUnits")
	message = "|Filters: " + where + "|Aggregated by : " + group + "| Energy Unit: " + energyUnits + "| Distance Unit: " + distanceUnits + "| Mass Unit: " + massUnits
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

// pass table data matrix and table object, update the column width base on the longest cell
func tableAutoSize(queryResult [][]string, tableData *widget.Table) {
	go func() {
		time.Sleep(1 * time.Second) //DELETE ME
		wi := getColWidths(queryResult)
		for i, v := range wi {
			tableData.SetColumnWidth(i, v)
		}
		tableData.Refresh()
	}()
}

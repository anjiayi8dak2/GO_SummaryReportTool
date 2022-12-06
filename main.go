package main

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
)

func main() {
	// Create the database handle, confirm driver is present
	db, err := sql.Open("mysql", "moves:moves@/")
	defer db.Close()
	if err != nil {
		fmt.Println("Failed to connect MariaDB", err)
		return
	}

	// Connect and check the server version
	getDBVersion(db)

	//Initialize
	a := app.New()
	w := a.NewWindow("Summary Report Tool")
	w.Resize(fyne.NewSize(400, 400))
	//Global var
	var dbSelection string
	var tableSelection string
	//var sqlResult []Movesoutput
	ifNullSQLMovesoutput := `SELECT  ifnull(MOVESRunID, -1) AS MOVESRunID,
		ifnull(MOVESRunID, -1) AS MOVESRunID,
		ifnull(iterationID, -1) AS iterationID,
		ifnull(yearID, -1) AS yearID,
		ifnull(monthID, -1) AS monthID,
		ifnull(dayID, -1) AS dayID,
		ifnull(stateID, -1) AS stateID,
		ifnull(countyID, -1) AS countyID,
		ifnull(zoneID, -1) AS zoneID,
		ifnull(linkID, -1) AS linkID,
		ifnull(pollutantID, -1) AS pollutantID,
		ifnull(processID, -1) AS processID,
		ifnull(sourceTypeID, -1) AS sourceTypeID,
		ifnull(regClassID, -1) AS regClassID,
		ifnull(fuelTypeID, -1) AS fuelTypeID,
		ifnull(fuelSubTypeID, -1) AS fuelSubTypeID,
		ifnull(modelYearID, -1) AS modelYearID,
		ifnull(roadTypeID, -1) AS roadTypeID,
		ifnull(SCC, -1) AS SCC,
		ifnull(engTechID, -1) AS engTechID,
		ifnull(sectorID, -1) AS sectorID,
		ifnull(hpID, -1) AS hpID,
		ifnull(emissionQuant, 0) AS emissionQuant
		FROM `

	//menu bar
	menuitemMaria := fyne.NewMenuItem("Open Maria Data Folder", func() {
		openMariaFolder(db)
	}) // ignore functions
	//menuitemRefresh := fyne.NewMenuItem("Refresh (F5)", buttonSubmit(db, dbSelection, tableSelection)) // ignore functions
	menuitemRefresh := fyne.NewMenuItem("Refresh (F5)", func() {
		var sqlStatement string
		sqlStatement = ifNullSQLMovesoutput + dbSelection + "." + tableSelection + " LIMIT 1"
		oneRowResult, _ := getOneRow(db, sqlStatement)
		fmt.Println("start Print default sql null")
		fmt.Printf("%v", &oneRowResult)
		fmt.Println("end Print default sql null")

		values := reflect.ValueOf(oneRowResult)
		types := values.Type()
		//fmt.Println("the field count for my struct ", values.NumField(), " my type is ", values.Type())

		var whiteList []string

		for i := 0; i < values.NumField(); i++ {
			// int to int
			if values.Field(i).Type() == reflect.TypeOf(1) {
				if values.Field(i).Int() != -1 {
					fmt.Println("found not null column, add it to white list ", types.Field(i).Name, values.Field(i))
					whiteList = append(whiteList, types.Field(i).Name)
					// != -1, add to whitelist
				}
				// float to float
			} else if values.Field(i).Type() == reflect.TypeOf(3.14) {
				if values.Field(i).Float() != -1 {
					fmt.Println("found not null column, add it to blacklist ", types.Field(i).Name, values.Field(i))
					whiteList = append(whiteList, types.Field(i).Name)
				}

			}

			//fmt.Println(types.Field(i).Index[0], types.Field(i).Name, values.Field(i))
		}

		fmt.Println("printing white list")
		fmt.Printf("%v", &whiteList)
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

//func test(m []Movesoutput) *fyne.Container {
//	table := widget.NewTable(
//		func() (int, int) {
//			return 10, 25
//		},
//		func() fyne.CanvasObject {
//			return widget.NewLabel("table window")
//		},
//		func(i widget.TableCellID, o fyne.CanvasObject) {
//			switch i.Col {
//			case 0:
//				o.(*widget.Label).SetText(m[i.Row].MOVESRunID)
//			case 1:
//				o.(*widget.Label).SetText(Convert.ToString(m[i.Row].MOVESRunID))
//			case 2:
//				o.(*widget.Label).SetText(m[i.Row].Name)
//			case 3:
//				o.(*widget.Label).SetText(m[i.Row].Memo)
//			}
//		},
//	)
//	table.SetColumnWidth(0, 200)
//	table.SetColumnWidth(1, 100)
//	table.SetColumnWidth(2, 100)
//	table.SetColumnWidth(3, 300)
//	split := container.NewHSplit(makeLeftSidebar(), table)
//	split.Offset = 0.2
//	return container.NewMax(split)
//}

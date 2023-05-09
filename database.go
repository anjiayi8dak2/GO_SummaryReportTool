package main

import (
	"database/sql"
	"fmt"
	_ "github.com/blockloop/scan"
	_ "github.com/joho/sqltocsv"
	"log"
	"reflect"
	_ "strconv"
	"strings"
)

func getDataDir(db *sql.DB) string {
	var dataDir string
	db.QueryRow("select @@datadir as dataDir;").Scan(&dataDir)
	fmt.Println("MariaDB data folder is in :", dataDir)
	//TODO: maybe make this a pop up?
	return dataDir
}

// PASS dbConnector, RETURN DB searching result as [] string
func getDBList(db *sql.DB) []string {
	// Query the DB
	var row string
	rows, err := db.Query(`SHOW DATABASES;`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var listSlice []string
	for rows.Next() {
		err = rows.Scan(&row)
		if err != nil {
			log.Fatal(err)
		}
		listSlice = append(listSlice, row)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return listSlice
}

// TODO improve performance
func getDistinct(db *sql.DB, dbSelection string, tableSelection string, targetColumn string) []string {
	sqlStatement := "SELECT DISTINCT " + targetColumn + " AS distinctValues FROM " + dbSelection + "." + tableSelection + " ORDER BY distinctValues ASC ; "
	fmt.Println("Running sql statement  :", sqlStatement)
	rows, err := db.Query(sqlStatement)
	var distinctResults []string
	if err != nil {
		panic(err)
		return distinctResults
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var distinctResult string
		rows.Scan(&distinctResult)
		if err != nil {
			panic(err) // Error related to the scan
		}
		if err = rows.Err(); err != nil {
			panic(err) // Error related to the iteration of rows
		}
		distinctResults = append(distinctResults, distinctResult)
	}
	fmt.Println("distinct result is %v :", distinctResults)
	return distinctResults
}

// PASS dbConnector, print DB version that is connected to
func getDBVersion(db *sql.DB) {
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)
}

func getQueryResult(db *sql.DB, dbSelection string, tableSelection string, whiteList []string, whereClause string, groupClause string) ([][]string, error) {
	columns := convertColumnsComma(whiteList)
	//check if there is an empty IN statement made by check/uncheck checkbox filters rapidly, and delete it if needed
	noFaultWhereClause := strings.ReplaceAll(whereClause, "IN (  )", "")
	//TODO user define Limit?
	sqlStatement := "SELECT " + columns + " FROM " + dbSelection + "." + tableSelection + " " + noFaultWhereClause + " " + groupClause + " LIMIT 1000 ; "

	fmt.Println("printing sql statement: " + sqlStatement)
	// A 2D array string represent the data table
	var outFlat [][]string
	// add the column names in first row from the remaining columns headers after dynamically filtered
	outFlat = append(outFlat, whiteList)

	// exe sql statement
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}

	count := len(whiteList)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		// string array to hold 1 row of query result
		var innerFlat []string
		for i := range whiteList {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		for i, _ := range whiteList {
			val := values[i]
			b, ok := val.([]byte)
			var v interface{}
			if ok {
				v = string(b)
			} else {
				v = val
			}
			innerFlat = append(innerFlat, v.(string))
		}
		// stick all 1D array into 2D for data table
		outFlat = append(outFlat, innerFlat)
	}
	//TODO: adding data, add the mass/distance unit for the (first column, last row)
	// also, depending what table is selected, the unit can be different. for example, grams or grams per mil
	// outFlat[0][len(whiteList)-1] = getUnit(db, dbSelection, tableSelection)
	return outFlat, err
}

func getUnit(db *sql.DB, dbSelection string, tableSelection string) string {
	var massUnit string
	var distanceUnit string
	var columnName string
	var query string

	query = "SELECT massUnits FROM " + dbSelection + ".movesrun LIMIT 1;"
	db.QueryRow(query).Scan(&massUnit)
	fmt.Println("Mass unit is  :", massUnit)

	query = "SELECT distanceUnits FROM " + dbSelection + ".movesrun LIMIT 1;"
	db.QueryRow(query).Scan(&distanceUnit)
	fmt.Println("Distance unit is  :", distanceUnit)

	switch tableSelection {
	case "movesactivityoutput":
		columnName = "activityQuant"
	case "movesoutput":
		columnName = "emissionQuant_" + massUnit
	case "rateperdistance":
		columnName = massUnit + "_per_" + distanceUnit
	case "rateperhour":
		columnName = massUnit + "_perHour"
	case "rateperprofile":
		columnName = massUnit + "_perVehicle"
	case "rateperstart":
		columnName = massUnit + "_perStarts"
	case "ratepervehicle":
		columnName = massUnit + "_perVehicle"
	case "startspervehicle":
		columnName = "starts_perVehicle"
	default:
		columnName = "unit error"

	}
	return columnName
}

// go-sql driver does not read null value that returned from query, therefore we use ifnull() and assign -1 as an indicator for the null value
// so far the best solution that does not rely on the random third party repo
func getOneRow(db *sql.DB, dbSelection string, tableSelection string) (interface{}, error) {
	var ifNullSQL string
	// there should be a smart way to do it, but I could not find any. stupid but works :(
	// editing SELECT clause sql statement depends on which table got selected
	switch tableSelection {
	case "movesactivityoutput":
		ifNullSQL = `SELECT
					ifnull(MOVESRunID, -1) AS MOVESRunID,
					ifnull(iterationID, -1) AS iterationID,
					ifnull(yearID, -1) AS yearID,
					ifnull(monthID, -1) AS monthID,
					ifnull(dayID, -1) AS dayID,
					ifnull(hourID, -1) AS hourID,
					ifnull(stateID, -1) AS stateID,
					ifnull(countyID, -1) AS countyID,
					ifnull(zoneID, -1) AS zoneID,
					ifnull(linkID, -1) AS linkID,
					ifnull(sourceTypeID, -1) AS sourceTypeID,
					ifnull(regClassID, -1) AS regClassID,
					ifnull(fuelTypeID, -1) AS fuelTypeID,
					ifnull(fuelSubTypeID, -1) AS fuelSubTypeID,
					ifnull(modelYearID, -1) AS modelYearID,
					ifnull(roadTypeID, -1) AS roadTypeID,
					IF(SCC IS NULL or SCC = '', -1, SCC) as SCC,
					ifnull(engTechID, -1) AS engTechID,
					ifnull(sectorID, -1) AS sectorID,
					ifnull(hpID, -1) AS hpID,
                    ifnull(activityTypeID, -1) AS activityTypeID,
					activity 
					FROM `

	case "movesoutput":
		ifNullSQL = `SELECT
					ifnull(MOVESRunID, -1) AS MOVESRunID,
					ifnull(iterationID, -1) AS iterationID,
					ifnull(yearID, -1) AS yearID,
					ifnull(monthID, -1) AS monthID,
					ifnull(dayID, -1) AS dayID,
					ifnull(hourID, -1) AS hourID,
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
					IF(SCC IS NULL or SCC = '', -1, SCC) as SCC,
					ifnull(engTechID, -1) AS engTechID,
					ifnull(sectorID, -1) AS sectorID,
					ifnull(hpID, -1) AS hpID,
					emissionQuant 
					FROM `

	case "rateperdistance":
		fmt.Println("I AM INSIDE RATE PER DISTANCE")
		ifNullSQL = `SELECT  
					ifnull(MOVESScenarioID, -1) AS MOVESScenarioID,
					ifnull(MOVESRunID, -1) AS MOVESRunID,
					ifnull(yearID, -1) AS yearID,
					ifnull(monthID, -1) AS monthID,
					ifnull(dayID, -1) AS dayID,
					ifnull(hourID, -1) AS hourID,
					ifnull(linkID, -1) AS linkID,
					ifnull(pollutantID, -1) AS pollutantID,
					ifnull(processID, -1) AS processID,
					ifnull(sourceTypeID, -1) AS sourceTypeID,
					ifnull(regClassID, -1) AS regClassID,
					IF(SCC IS NULL or SCC = '', -1, SCC) as SCC,
					ifnull(fuelTypeID, -1) AS fuelTypeID,
					ifnull(modelYearID, -1) AS modelYearID,
					ifnull(roadTypeID, -1) AS roadTypeID,
					ifnull(avgSpeedBinID, -1) AS avgSpeedBinID,
					ifnull(temperature, -1) AS temperature,
					ifnull(relHumidity, -1) AS relHumidity,
					ifnull(ratePerDistance, -1) AS ratePerDistance
					FROM `

	case "rateperhour":
		ifNullSQL = `SELECT  
					ifnull(MOVESScenarioID, -1) AS MOVESScenarioID,
					ifnull(MOVESRunID, -1) AS MOVESRunID,
					ifnull(yearID, -1) AS yearID,
					ifnull(monthID, -1) AS monthID,
					ifnull(dayID, -1) AS dayID,
					ifnull(hourID, -1) AS hourID,
					ifnull(linkID, -1) AS linkID,
					ifnull(pollutantID, -1) AS pollutantID,
					ifnull(processID, -1) AS processID,
					ifnull(sourceTypeID, -1) AS sourceTypeID,
					ifnull(regClassID, -1) AS regClassID,
					IF(SCC IS NULL or SCC = '', -1, SCC) as SCC,
					ifnull(fuelTypeID, -1) AS fuelTypeID,
					ifnull(modelYearID, -1) AS modelYearID,
					ifnull(roadTypeID, -1) AS roadTypeID,
					ifnull(temperature, -1) AS temperature,
					ifnull(relHumidity, -1) AS relHumidity,
					ifnull(ratePerHour, -1) AS ratePerHour
					FROM `

	case "rateperprofile":
		ifNullSQL = `SELECT  
					ifnull(MOVESScenarioID, -1) AS MOVESScenarioID,
					ifnull(MOVESRunID, -1) AS MOVESRunID,
					ifnull(temperatureProfileID, -1) AS temperatureProfileID,
					ifnull(yearID, -1) AS yearID,
					ifnull(dayID, -1) AS dayID,
					ifnull(hourID, -1) AS hourID,
					ifnull(pollutantID, -1) AS pollutantID,
					ifnull(processID, -1) AS processID,
					ifnull(sourceTypeID, -1) AS sourceTypeID,
					ifnull(regClassID, -1) AS regClassID,
					IF(SCC IS NULL or SCC = '', -1, SCC) as SCC,
					ifnull(fuelTypeID, -1) AS fuelTypeID,
					ifnull(modelYearID, -1) AS modelYearID,
					ifnull(temperature, -1) AS temperature,
					ifnull(relHumidity, -1) AS relHumidity,
					ifnull(ratePerVehicle, -1) AS ratePerVehicle
					FROM  `

	case "rateperstart":
		ifNullSQL = `SELECT  
					ifnull(MOVESScenarioID, -1) AS MOVESScenarioID,
					ifnull(MOVESRunID, -1) AS MOVESRunID,
					ifnull(yearID, -1) AS yearID,
					ifnull(monthID, -1) AS monthID,
					ifnull(dayID, -1) AS dayID,
					ifnull(hourID, -1) AS hourID,
					ifnull(zoneID, -1) AS zoneID,
					ifnull(sourceTypeID, -1) AS sourceTypeID,
					ifnull(regClassID, -1) AS regClassID,
					IF(SCC IS NULL or SCC = '', -1, SCC) as SCC,
					ifnull(fuelTypeID, -1) AS fuelTypeID,
					ifnull(modelYearID, -1) AS modelYearID,
					ifnull(pollutantID, -1) AS pollutantID,
					ifnull(processID, -1) AS processID,
					ifnull(temperature, -1) AS temperature,
					ifnull(relHumidity, -1) AS relHumidity,
					ifnull(ratePerStart, -1) AS ratePerStart
					FROM  `

	case "ratepervehicle":
		ifNullSQL = `SELECT  
					ifnull(MOVESScenarioID, -1) AS MOVESScenarioID,
					ifnull(MOVESRunID, -1) AS MOVESRunID,
					ifnull(yearID, -1) AS yearID,
					ifnull(monthID, -1) AS monthID,
					ifnull(dayID, -1) AS dayID,
					ifnull(hourID, -1) AS hourID,
					ifnull(zoneID, -1) AS zoneID,
					ifnull(pollutantID, -1) AS pollutantID,
					ifnull(processID, -1) AS processID,
					ifnull(sourceTypeID, -1) AS sourceTypeID,
					ifnull(regClassID, -1) AS regClassID,
					IF(SCC IS NULL or SCC = '', -1, SCC) as SCC,
					ifnull(fuelTypeID, -1) AS fuelTypeID,
					ifnull(modelYearID, -1) AS modelYearID,
					ifnull(temperature, -1) AS temperature,
					ifnull(relHumidity, -1) AS relHumidity,
					ifnull(ratePerVehicle, -1) AS ratePerVehicle
					FROM  `

	case "startspervehicle":
		ifNullSQL = `SELECT  
					ifnull(MOVESScenarioID, -1) AS MOVESScenarioID,
					ifnull(MOVESRunID, -1) AS MOVESRunID,
					ifnull(yearID, -1) AS yearID,
					ifnull(monthID, -1) AS monthID,
					ifnull(dayID, -1) AS dayID,
					ifnull(hourID, -1) AS hourID,
					ifnull(zoneID, -1) AS zoneID,
					ifnull(sourceTypeID, -1) AS sourceTypeID,
					ifnull(regClassID, -1) AS regClassID,
					IF(SCC IS NULL or SCC = '', -1, SCC) as SCC,
					ifnull(fuelTypeID, -1) AS fuelTypeID,
					ifnull(modelYearID, -1) AS modelYearID,
					ifnull(startsPerVehicle, -1) AS startsPerVehicle
					FROM  `

	default:
		fmt.Println("unknown table selection inside getOneRow")
		panic("unknown table selection inside getOneRow")
	}

	// put sql statement together and select one row
	sql := ifNullSQL + dbSelection + "." + tableSelection + " LIMIT 1;"
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
		return nil, err
	}
	defer rows.Close()

	// depends on the table selection string value, create different instance of the struct,
	// then scan the query result into struct specific field for next steps
	// the Scan function from go-sql driver has no way to select all columns into target struct in one command like "SELECT * ..." into struct_name
	// I could not find a better way than hard coded column names, this is fine as long as we don't change schemas of movesoutput
	switch tableSelection {
	case "movesactivityoutput":
		var output Movesactivityoutput
		// Loop through each column, using Scan to assign column data to struct fields.
		for rows.Next() {
			rows.Scan(&output.MOVESRunID, &output.iterationID, &output.yearID,
				&output.monthID, &output.dayID, &output.hourID, &output.stateID,
				&output.countyID, &output.zoneID, &output.linkID,
				&output.sourceTypeID, &output.regClassID, &output.fuelTypeID,
				&output.fuelSubTypeID, &output.modelYearID, &output.roadTypeID, &output.SCC,
				&output.engTechID, &output.sectorID, &output.hpID, &output.activityTypeID, &output.activity)
			if err != nil {
				panic(err) // Error related to the scan
			}
			if err = rows.Err(); err != nil {
				panic(err) // Error related to the iteration of rows
			}
		}
		return output, nil

	case "movesoutput":
		var output Movesoutput
		for rows.Next() {
			rows.Scan(&output.MOVESRunID, &output.iterationID, &output.yearID,
				&output.monthID, &output.dayID, &output.hourID, &output.stateID,
				&output.countyID, &output.zoneID, &output.linkID, &output.pollutantID,
				&output.processID, &output.sourceTypeID, &output.regClassID, &output.fuelTypeID,
				&output.fuelSubTypeID, &output.modelYearID, &output.roadTypeID, &output.SCC,
				&output.engTechID, &output.sectorID, &output.hpID, &output.emissionQuant)
			if err != nil {
				panic(err) // Error related to the scan
			}
			if err = rows.Err(); err != nil {
				panic(err) // Error related to the iteration of rows
			}
		}
		return output, nil
	case "rateperdistance":
		var output Rateperdistance
		for rows.Next() {
			rows.Scan(&output.MOVESScenarioID, &output.MOVESRunID, &output.yearID, &output.monthID,
				&output.dayID, &output.hourID, &output.linkID, &output.pollutantID, &output.processID,
				&output.sourceTypeID, &output.regClassID, &output.SCC, &output.fuelTypeID, &output.modelYearID,
				&output.roadTypeID, &output.avgSpeedBinID, &output.temperature, &output.relHumidity, &output.ratePerDistance)
			if err != nil {
				panic(err) // Error related to the scan
			}
			if err = rows.Err(); err != nil {
				panic(err) // Error related to the iteration of rows
			}
		}
		return output, nil
	case "rateperhour":
		var output Rateperhour
		for rows.Next() {
			rows.Scan(&output.MOVESScenarioID, &output.MOVESRunID, &output.yearID, &output.monthID,
				&output.dayID, &output.hourID, &output.linkID, &output.pollutantID,
				&output.processID, &output.sourceTypeID, &output.regClassID, &output.SCC,
				&output.fuelTypeID, &output.modelYearID, &output.roadTypeID, &output.temperature,
				&output.relHumidity, &output.ratePerHour)
			if err != nil {
				panic(err) // Error related to the scan
			}
			if err = rows.Err(); err != nil {
				panic(err) // Error related to the iteration of rows
			}
		}
		return output, nil
	case "rateperprofile":
		var output Rateperprofile
		for rows.Next() {
			rows.Scan(&output.MOVESScenarioID, &output.MOVESRunID, &output.temperatureProfileID,
				&output.yearID, &output.dayID, &output.hourID, &output.pollutantID,
				&output.processID, &output.sourceTypeID, &output.regClassID, &output.SCC,
				&output.fuelTypeID, &output.modelYearID, &output.temperature, &output.relHumidity,
				&output.ratePerVehicle)
			if err != nil {
				panic(err) // Error related to the scan
			}
			if err = rows.Err(); err != nil {
				panic(err) // Error related to the iteration of rows
			}
		}
		return output, nil
	case "rateperstart":
		var output Rateperstart
		for rows.Next() {
			rows.Scan(&output.MOVESScenarioID, &output.MOVESRunID, &output.yearID, &output.monthID, &output.dayID,
				&output.hourID, &output.zoneID, &output.sourceTypeID, &output.regClassID, &output.SCC,
				&output.fuelTypeID, &output.modelYearID, &output.pollutantID, &output.processID, &output.temperature,
				&output.relHumidity, &output.ratePerStart)
			if err != nil {
				panic(err) // Error related to the scan
			}
			if err = rows.Err(); err != nil {
				panic(err) // Error related to the iteration of rows
			}
		}
		return output, nil
	case "ratepervehicle":
		var output Ratepervehicle
		for rows.Next() {
			rows.Scan(&output.MOVESRunID, &output.yearID, &output.monthID, &output.dayID, &output.hourID, &output.zoneID,
				&output.pollutantID, &output.processID, &output.sourceTypeID, &output.regClassID, &output.SCC, &output.fuelTypeID,
				&output.modelYearID, &output.temperature, &output.relHumidity, &output.ratePerVehicle)
			if err != nil {
				panic(err) // Error related to the scan
			}
			if err = rows.Err(); err != nil {
				panic(err) // Error related to the iteration of rows
			}
		}
		return output, nil
	case "startspervehicle":
		var output Startspervehicle
		for rows.Next() {
			rows.Scan(&output.MOVESScenarioID, &output.MOVESRunID, &output.yearID, &output.monthID, &output.dayID,
				&output.hourID, &output.zoneID, &output.sourceTypeID, &output.regClassID, &output.SCC,
				&output.fuelTypeID, &output.modelYearID, &output.startsPerVehicle)
			if err != nil {
				panic(err) // Error related to the scan
			}
			if err = rows.Err(); err != nil {
				panic(err) // Error related to the iteration of rows
			}
		}
		return output, nil
	default:
		//unknow selection
		panic("unknown selection found in the table selection drop down box")
		break
	}

	//should not run to here
	return nil, nil

}

// return []string contains whitelist columns names that have meaningful values, size depends on how many column are survived
// return []bool that indicate valid/invalid values for each column, []bool size = total column that a table has in struct
// for example: a 4 column table has 4 attributes {"yearID", "pollutantID", "roadTypeID", "emissionQuant"}
// then the filter found yearID and roadTypeID are null
// returned whitelist will be {"pollutantID","emissionQuant"}
// returned []bool will be {0,1,0,1}
func getWhiteList(con *sql.DB, dbSelection string, tableSelection string) ([]string, []bool) {
	//var fieldNames []string
	//
	//switch tableSelection {
	//case "movesactivityoutput":
	//	fieldNames = []string{"MOVESRunID", "iterationID", "yearID", "monthID", "dayID", "hourID", "stateID", "countyID",
	//		"zoneID", "linkID", "sourceTypeID", "regClassID", "fuelTypeID", "fuelSubTypeID",
	//		"modelYearID", "roadTypeID", "SCC", "engTechID", "sectorID", "hpID", "activityTypeID", "activity"}
	//case "movesoutput":
	//	fieldNames = []string{"MOVESRunID", "iterationID", "yearID", "monthID", "dayID", "hourID", "stateID", "countyID",
	//		"zoneID", "linkID", "pollutantID", "processID", "sourceTypeID", "regClassID", "fuelTypeID", "fuelSubTypeID",
	//		"modelYearID", "roadTypeID", "SCC", "engTechID", "sectorID", "hpID", "emissionQuant"}
	//case "rateperdistance":
	//	fieldNames = []string{"MOVESScenarioID", "MOVESRunID", "yearID", "monthID", "dayID", "hourID", "linkID", "pollutantID",
	//		"processID", "sourceTypeID", "regClassID", "SCC", "fuelTypeID", "modelYearID", "roadTypeID", "avgSpeedBinID",
	//		"temperature", "relHumidity", "ratePerDistance"}
	//case "rateperhour":
	//	fieldNames = []string{"MOVESScenarioID", "MOVESRunID", "yearID", "monthID", "dayID", "hourID", "linkID", "pollutantID",
	//		"processID", "sourceTypeID", "regClassID", "SCC", "fuelTypeID", "modelYearID", "roadTypeID", "temperature", "relHumidity", "ratePerHour"}
	//case "rateperprofile":
	//	fieldNames = []string{"MOVESScenarioID", "MOVESRunID", "temperatureProfileID", "yearID", "dayID", "hourID", "pollutantID",
	//		"processID", "sourceTypeID", "regClassID", "SCC", "fuelTypeID", "modelYearID", "temperature", "relHumidity", "ratePerVehicle"}
	//case "rateperstart":
	//	fieldNames = []string{"MOVESScenarioID", "MOVESRunID", "yearID", "monthID", "dayID", "hourID", "zoneID", "sourceTypeID",
	//		"regClassID", "SCC", "fuelTypeID", "modelYearID", "pollutantID", "processID", "temperature", "relHumidity", "ratePerStart"}
	//case "ratepervehicle":
	//	fieldNames = []string{"MOVESScenarioID", "MOVESRunID", "yearID", "monthID", "dayID", "hourID", "zoneID", "pollutantID",
	//		"processID", "sourceTypeID", "regClassID", "SCC", "fuelTypeID", "modelYearID", "temperature", "relHumidity", "ratePerVehicle"}
	//case "startspervehicle":
	//	fieldNames = []string{"MOVESScenarioID", "MOVESRunID", "yearID", "monthID", "dayID", "hourID", "zoneID",
	//		"sourceTypeID", "regClassID", "SCC", "fuelTypeID", "modelYearID", "startsPerVehicle"}
	//default:
	//	panic("unknown table selection ")
	//
	//}

	oneRowResult, _ := getOneRow(con, dbSelection, tableSelection)
	values := reflect.ValueOf(oneRowResult)
	types := values.Type()

	var whiteList []string
	var whiteListIndex []bool

	// loop through all the columns from OneRow, if a column has valid number(not equal -1), append it into whitelist [] string
	for i := 0; i < values.NumField(); i++ {
		if values.Field(i).Type() == reflect.TypeOf(1) { //type int
			if values.Field(i).Int() != -1 { // non -1 integer
				fmt.Println("found column with valid integer value, add it to whitelist \n", types.Field(i).Name, values.Field(i))
				whiteList = append(whiteList, types.Field(i).Name)
			}
		} else if values.Field(i).Type() == reflect.TypeOf(3.14) { // type float
			fmt.Println("found column with valid float value, add it to whitelist  \n", types.Field(i).Name, values.Field(i))
			whiteList = append(whiteList, types.Field(i).Name)
			// string to string, the MOVESScenarioID unfortunately can be a string :(, Go type restriction will panic on compare different types
		} else if values.Field(i).Type() == reflect.TypeOf("word") { //type string
			fmt.Println("found column with valid string value, add it to whitelist \n", types.Field(i).Name, values.Field(i))
			whiteList = append(whiteList, types.Field(i).Name)
		}
	}

	//loop through values and update its boolean value when detect -1
	//whiteListIndex [] bool size = corresponding target table struct, the order of the column is also same
	//whiteListIndex flags is being used to determinate which column/filter_check_box to show/hide containers in the future
	//TODO merge this for loop into the above one
	for i := 0; i < values.NumField(); i++ {
		if values.Field(i).Type() == reflect.TypeOf(1) {
			if values.Field(i).Int() != -1 {
				whiteListIndex = append(whiteListIndex, true)
			} else if values.Field(i).Int() == -1 {
				whiteListIndex = append(whiteListIndex, false)
			}
		} else if values.Field(i).Type() == reflect.TypeOf("word") {
			if values.Field(i).String() != "-1" {
				whiteListIndex = append(whiteListIndex, true)
			} else if values.Field(i).String() == "1" {
				whiteListIndex = append(whiteListIndex, false)
			}
		}

	}

	//// numeric columns count that a table has, for example movesoutput has "emissionQuant", but rateperhour has "temperature", "relHumidity","ratePerHour"
	//var numericColumnsInTheEnd int
	//if tableSelection == "movesoutput" || tableSelection == "startspervehicle" || tableSelection == "movesactivityoutput" {
	//	numericColumnsInTheEnd = 1
	//} else {
	//	numericColumnsInTheEnd = 3
	//}
	//
	////find distinct values for the survived column, then also ignore these columns that only has 1 distinct value
	////loop through whiteListIndex, for these columns survived for the null value check, check the count of distinct value = 1,
	////for example if the MOVESRUNID only has 1, ignore it, there is no point to show them as both column or filter
	//for i := 0; i < len(whiteListIndex)-numericColumnsInTheEnd; i++ { //loop to the position before numeric column such as emissionQuant/activity
	//	if whiteListIndex[i] && fieldNames[i] != "pollutantID" { //if the column value is not null, skip pollutantID
	//		//get distinct query, and see the count or len(returned slice)
	//		distinctResult := getDistinct(con, dbSelection, tableSelection, fieldNames[i])
	//		if len(distinctResult) <= 1 { // if the returned slice only has <= 1 distinct value, mark the index to false
	//			whiteListIndex[i] = false
	//			fmt.Print(" found column that only has 1 distinct value ", fieldNames[i])
	//			fmt.Print(" printing updated whiteList v% v%", fieldNames[i], whiteListIndex[i])
	//			whiteList = RemoveElementFromSlice(whiteList, fieldNames[i]) // call func that remove #the column that only have 1 distinct value as well
	//		}
	//	}
	//}

	return whiteList, whiteListIndex
}

// TODO not to hard code, if default credential connection failed let user to enter username and password
func initDb() *sql.DB {
	// Create the database handle, confirm driver is present
	db, err := sql.Open("mysql", "moves:moves@/")
	//defer db.Close()
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

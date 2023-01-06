package main

import (
	"database/sql"
	"fmt"
	_ "github.com/blockloop/scan"
	_ "github.com/joho/sqltocsv"
	"log"
	"reflect"
	"strconv"
	_ "strconv"
)

func getDataDir(db *sql.DB) string {
	var dataDir string
	db.QueryRow("select @@datadir as dataDir;").Scan(&dataDir)
	fmt.Println("MariaDB data folder is in :", dataDir)
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
		err := rows.Scan(&row)
		if err != nil {
			log.Fatal(err)
		}
		listSlice = append(listSlice, row)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	//fmt.Println("DB list is here :", listSlice)
	return listSlice
}

// PASS dbConnector, print DB version that is connected to
func getDBVersion(db *sql.DB) {
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)
}

// // PASS dbConnector, sql statement in string, and table selection RETURN query result
//
//	func getQueryResult(db *sql.DB, sql string) {
//		rows, _ := db.Query(sql) // Note: Ignoring errors for brevity
//		cols, _ := rows.Columns()
//		//mapResult := make(map[string]string)
//		for rows.Next() {
//			// Create a slice of interface{}'s to represent each column,
//			// and a second slice to contain pointers to each item in the columns slice.
//			columns := make([]interface{}, len(cols))
//			columnPointers := make([]interface{}, len(cols))
//			for i, _ := range columns {
//				columnPointers[i] = &columns[i]
//			}
//
//			// Scan the result into the column pointers...
//			if err := rows.Scan(columnPointers...); err != nil {
//				fmt.Println("scanning pointer failed")
//				//return err
//			}
//
//			// Create our map, and retrieve the value for each column from the pointers slice,
//			// storing it in the map with the name of the column as the key.
//			m := make(map[string]interface{})
//			for i, colName := range cols {
//				val := columnPointers[i].(*interface{})
//				m[colName] = *val
//			}
//
//			//Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
//			fmt.Print(m)
//		}
//
// }
// PASS dbConnector, sql statement in string, and table selection RETURN query result
//func getQueryResult(db *sql.DB, dbSelection string, tableSelection string, whiteList []string) ([]map[string]interface{}, error) {
//	//TODO: build sql statement here
//	columns := convertColumnsComma(whiteList)
//	sqlStatement := "SELECT " + columns + " FROM " + dbSelection + "." + tableSelection + " LIMIT 10 ; "
//
//	fmt.Println("sql statement is :", sqlStatement)
//	rows, err := db.Query(sqlStatement)
//	//cols, _ := rows.Columns()
//	//data := make(map[string]string)
//	//TODO:dumping a query to a file
//	sqlErr := sqltocsv.WriteFile("temp.csv", rows)
//	if sqlErr != nil {
//		panic(sqlErr)
//	}
//	//TODO:import csv into datatable, testing testing testing
//	readCSV()
//
//	if err != nil {
//		panic(err)
//		return nil, err
//	}
//	defer rows.Close()
//
//	//queryResultMap := scanMap(rows)
//	//
//	//TODO: scan lite query result into Movesoutput struct
//
//	return queryResultMap, err
//
//}

// TODO: function that have switch assign field value to a query row
func (m *Movesoutput) myUpdater(fieldName, fieldValue string) {
	switch fieldName {
	case "MOVESRunID":
		m.MOVESRunID, _ = strconv.Atoi(fieldValue)
	case "iterationID":
		m.iterationID, _ = strconv.Atoi(fieldValue)
	case "yearID":
		m.yearID, _ = strconv.Atoi(fieldValue)
	case "monthID":
		m.monthID, _ = strconv.Atoi(fieldValue)
	case "dayID":
		m.dayID, _ = strconv.Atoi(fieldValue)
	case "hourID":
		m.hourID, _ = strconv.Atoi(fieldValue)
	case "stateID":
		m.stateID, _ = strconv.Atoi(fieldValue)
	case "countyID":
		m.countyID, _ = strconv.Atoi(fieldValue)
	case "zoneID":
		m.zoneID, _ = strconv.Atoi(fieldValue)
	case "linkID":
		m.linkID, _ = strconv.Atoi(fieldValue)
	case "pollutantID":
		m.pollutantID, _ = strconv.Atoi(fieldValue)
	case "processID":
		m.processID, _ = strconv.Atoi(fieldValue)
	case "sourceTypeID":
		m.sourceTypeID, _ = strconv.Atoi(fieldValue)
	case "regClassID":
		m.regClassID, _ = strconv.Atoi(fieldValue)
	case "fuelTypeID":
		m.fuelTypeID, _ = strconv.Atoi(fieldValue)
	case "fuelSubTypeID":
		m.fuelSubTypeID, _ = strconv.Atoi(fieldValue)
	case "modelYearID":
		m.modelYearID, _ = strconv.Atoi(fieldValue)
	case "roadTypeID":
		m.roadTypeID, _ = strconv.Atoi(fieldValue)
	case "SCC":
		m.SCC, _ = strconv.Atoi(fieldValue)
	case "engTechID":
		m.engTechID, _ = strconv.Atoi(fieldValue)
	case "sectorID":
		m.hpID, _ = strconv.Atoi(fieldValue)
	case "hpID":
		m.processID, _ = strconv.Atoi(fieldValue)
	case "emissionQuant":
		m.emissionQuant, _ = strconv.ParseFloat(fieldValue, 64)
	default:
		panic("I don't know that field name: " + fieldName + "!")
	}
}

func getQueryResult(db *sql.DB, dbSelection string, tableSelection string, whiteList []string) {
	//columns := convertColumnsComma(whiteList)
	//sqlStatement := ifNullSQLMovesoutput + dbSelection + "." + tableSelection + " LIMIT 10 ; "
	//sqlStatement := "SELECT " + columns + " FROM " + dbSelection + "." + tableSelection + " LIMIT 10 ; "
	//sqlStatement := "SELECT * FROM " + dbSelection + "." + tableSelection + " LIMIT 3 ; "

	ifNullSQLMovesoutput := `SELECT
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
		ifnull(SCC, -1) AS SCC,
		ifnull(engTechID, -1) AS engTechID,
		ifnull(sectorID, -1) AS sectorID,
		ifnull(hpID, -1) AS hpID,
		emissionQuant
		FROM `
	// A movesoutput slice to hold data from returned rows.
	sql := ifNullSQLMovesoutput + dbSelection + "." + tableSelection + " LIMIT 3;"
	var outArr []Movesoutput
	fmt.Println("sql statement is :", sql)
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var movesout Movesoutput
		rows.Scan(&movesout.MOVESRunID, &movesout.iterationID, &movesout.yearID,
			&movesout.monthID, &movesout.dayID, &movesout.hourID, &movesout.stateID,
			&movesout.countyID, &movesout.zoneID, &movesout.linkID, &movesout.pollutantID,
			&movesout.processID, &movesout.sourceTypeID, &movesout.regClassID, &movesout.fuelTypeID,
			&movesout.fuelSubTypeID, &movesout.modelYearID, &movesout.roadTypeID, &movesout.SCC,
			&movesout.engTechID, &movesout.sectorID, &movesout.hpID, &movesout.emissionQuant)
		if err != nil {
			panic(err) // Error related to the scan
		}
		if err = rows.Err(); err != nil {
			panic(err) // Error related to the iteration of rows
		}
		fmt.Printf("inside getQuery current row is %v\\n", movesout)
		outArr = append(outArr, movesout)
	}

	fmt.Println("printing error")
	fmt.Printf("%#v", err)
	fmt.Println("printing outArr")
	fmt.Printf("%#v", outArr)

	//var result []Movesoutput
	//err = scan.RowsStrict(&result, rows)
	//
	//if err != nil {
	//	panic(err) // Error related to the scan
	//}
	//fmt.Println("printing error")
	//fmt.Printf("%#v", err)
	//if err = rows.Err(); err != nil {
	//	panic(err) // Error related to the iteration of rows
	//}
	//fmt.Println("printing rows.err")
	//fmt.Printf("%#v", err)
	//
	//fmt.Println("printing query result START")
	//fmt.Printf("%#v", result)
	//fmt.Println("printing query result DONE")
	//
	//fmt.Println("printing error")
	//fmt.Printf("%#v", err)

	//// create a field binding object.
	//var fArr []string
	//fb := NewFieldBinding()
	//if fArr, err = rows.Columns(); err != nil {
	//	return nil, err
	//}
	//fb.PutFields(fArr)
	//outArr := []interface{}{}
	//for rows.Next() {
	//	if err := rows.Scan(fb.GetFieldPtrArr()...); err != nil {
	//		return nil, err
	//	}
	//	fmt.Printf("Row: %v, %v, %v, %s\n", fb.Get("MOVESRunID"), fb.Get("yearID"), fb.Get("pollutantID"), fb.Get("emissionQuant"))
	//	outArr = append(outArr, fb.GetFieldArr())
	//}
	//if err := rows.Err(); err != nil {
	//	return nil, err
	//}
	//return outArr, nil
}

// PASS dbConnector, sql statement in string, and table selection RETURN query result
func getOneRow(db *sql.DB, dbSelection string, tableSelection string) (Movesoutput, error) {
	ifNullSQLMovesoutput := `SELECT
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
		ifnull(SCC, -1) AS SCC,
		ifnull(engTechID, -1) AS engTechID,
		ifnull(sectorID, -1) AS sectorID,
		ifnull(hpID, -1) AS hpID,
		emissionQuant
		FROM `
	// A movesoutput slice to hold data from returned rows.
	sql := ifNullSQLMovesoutput + dbSelection + "." + tableSelection + " LIMIT 1;"
	var movesout Movesoutput
	fmt.Println("sql statement is :", sql)
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
		return movesout, err
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		rows.Scan(&movesout.MOVESRunID, &movesout.iterationID, &movesout.yearID,
			&movesout.monthID, &movesout.dayID, &movesout.hourID, &movesout.stateID,
			&movesout.countyID, &movesout.zoneID, &movesout.linkID, &movesout.pollutantID,
			&movesout.processID, &movesout.sourceTypeID, &movesout.regClassID, &movesout.fuelTypeID,
			&movesout.fuelSubTypeID, &movesout.modelYearID, &movesout.roadTypeID, &movesout.SCC,
			&movesout.engTechID, &movesout.sectorID, &movesout.hpID, &movesout.emissionQuant)
		if err != nil {
			panic(err) // Error related to the scan
		}
		if err = rows.Err(); err != nil {
			panic(err) // Error related to the iteration of rows
		}
		fmt.Printf("current row is %v\\n", movesout)
	}

	return movesout, nil

}

func getWhiteList(con *sql.DB, dbSelection string, tableSelection string) []string {
	oneRowResult, _ := getOneRow(con, dbSelection, tableSelection)
	fmt.Println("Print one row")
	fmt.Printf("%v", &oneRowResult)

	values := reflect.ValueOf(oneRowResult)
	types := values.Type()

	var whiteList []string

	for i := 0; i < values.NumField(); i++ {
		// int to int
		if values.Field(i).Type() == reflect.TypeOf(1) {
			if values.Field(i).Int() != -1 {
				fmt.Println("found column with valid integer value, add it to whitelist \n", types.Field(i).Name, values.Field(i))
				whiteList = append(whiteList, types.Field(i).Name)
			}
			// float to float, this is only for emissionQuant column
		} else if values.Field(i).Type() == reflect.TypeOf(3.14) {
			fmt.Println("this is only for emissionQuant column, add it to whitelist \n", types.Field(i).Name, values.Field(i))
			whiteList = append(whiteList, types.Field(i).Name)
		}
	}
	return whiteList
}

func initDb() *sql.DB {
	// Create the database handle, confirm driver is present
	db, err := sql.Open("mysql", "moves:moves@/")
	//defer db.Close()

	if err != nil {
		log.Fatalln(err)
	}
	return db
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
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
func getQueryResult(db *sql.DB, dbSelection string, tableSelection string, whiteList []string) (map[string]float64, error) {
	//TODO: build sql statement here
	columns := convertColumnsComma(whiteList)
	sqlStatement := "SELECT " + columns + " FROM " + dbSelection + "." + tableSelection + " LIMIT 10 ; "

	fmt.Println("sql statement is :", sqlStatement)
	rows, err := db.Query(sqlStatement)
	cols, _ := rows.Columns()
	data := make(map[string]string)
	if err != nil {
		panic(err)
		return nil, err
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		for i, colName := range cols {
			data[colName] = columns[i]
		}
	}

	fmt.Print(data)
	return nil, err

}

// PASS dbConnector, sql statement in string, and table selection RETURN query result
func getOneRow(db *sql.DB, dbSelection string, tableSelection string) (Movesoutput, error) {
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
	return whiteList
}

//func buildSql(db string) {
//
//}

package main

import (
	"database/sql"
	"fmt"
	"log"
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
func getQueryResult(db *sql.DB, sql string) ([]Movesoutput, error) {
	fmt.Println("sql statement is :", sql)
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
		return nil, err
	}
	defer rows.Close()
	// A movesoutput slice to hold data from returned rows.
	var movesoutputs []Movesoutput
	var movesout Movesoutput
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
		movesoutputs = append(movesoutputs, movesout)
	}

	return movesoutputs, nil

}

// PASS dbConnector, sql statement in string, and table selection RETURN query result
func getOneRow(db *sql.DB, sql string) (Movesoutput, error) {
	// A movesoutput slice to hold data from returned rows.
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

//func buttonSubmit(con *sql.DB, db string, table string) {
//	var sql string
//	sql = "SELECT * FROM " + db + "." + table + " LIMIT 10"
//	getQueryResult(con, sql)
//}

//func buildSql(db string) {
//
//}

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/blockloop/scan"
	_ "github.com/joho/sqltocsv"
	"log"
	"reflect"
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

func getDistinct(db *sql.DB, dbSelection string, tableSelection string, targetColumn string) []string {
	sqlStatement := "SELECT DISTINCT " + targetColumn + " AS dummy FROM " + dbSelection + "." + tableSelection + " ORDER BY dummy ASC ; "
	fmt.Println("sql statement is :", sqlStatement)
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
	sqlStatement := "SELECT " + columns + " FROM " + dbSelection + "." + tableSelection + " " + whereClause + " " + groupClause + " LIMIT 1000 ; "
	fmt.Println("printing sql statement: " + sqlStatement)
	// A 2D array string to hold the table
	var outFlat [][]string
	// add the column names in first row
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
	//fmt.Println("printing column name before return ")
	//fmt.Println(whiteList)
	//fmt.Println("printing query result before return ")
	//for i := 0; i < len(outFlat); i++ {
	//	fmt.Println(outFlat[i])
	//}
	return outFlat, err
}

// go-sql driver does not read null value, therefore we use -1 as an indicator for the null value
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
	// select one row
	sql := ifNullSQLMovesoutput + dbSelection + "." + tableSelection + " LIMIT 1;"
	var movesout Movesoutput
	//fmt.Println("sql statement is :", sql)
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
		return movesout, err
	}
	defer rows.Close()

	// Loop through each column, using Scan to assign column data to struct fields.
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
		//fmt.Printf("current row is %v\\n", movesout)
	}
	return movesout, nil

}

func getWhiteList(con *sql.DB, dbSelection string, tableSelection string) ([]string, []bool) {
	moFieldNames := []string{"MOVESRunID", "iterationID", "yearID", "monthID", "dayID", "hourID", "stateID", "countyID",
		"zoneID", "linkID", "pollutantID", "processID", "sourceTypeID", "regClassID", "fuelTypeID", "fuelSubTypeID",
		"modelYearID", "roadTypeID", "SCC", "engTechID", "sectorID", "hpID", "emissionQuant"}
	oneRowResult, _ := getOneRow(con, dbSelection, tableSelection)

	values := reflect.ValueOf(oneRowResult)
	types := values.Type()

	var whiteList []string
	var whiteListIndex []bool

	// get whitelist in [] string
	for i := 0; i < values.NumField(); i++ {
		// int to int
		if values.Field(i).Type() == reflect.TypeOf(1) {
			if values.Field(i).Int() != -1 {
				fmt.Println("found column with valid integer value, add it to whitelist \n", types.Field(i).Name, values.Field(i))
				whiteList = append(whiteList, types.Field(i).Name)
			}
			// float to float, this is only for emissionQuant column, only emissionQuant column has float
		} else if values.Field(i).Type() == reflect.TypeOf(3.14) {
			fmt.Println("this is only for emissionQuant column, add it to whitelist \n", types.Field(i).Name, values.Field(i))
			whiteList = append(whiteList, types.Field(i).Name)
		}
	}

	//else if values.Field(i).Type() == reflect.TypeOf(-1.0)
	//get whitelist index in bool, length -1 because we don't take emissionQuant as filter button
	//assign true for not null , false for null value
	for i := 0; i < values.NumField()-1; i++ {
		if values.Field(i).Int() != -1 {
			whiteListIndex = append(whiteListIndex, true)
		} else if values.Field(i).Int() == -1 {
			whiteListIndex = append(whiteListIndex, false)
		} else {
			//catch? what else could be?
			panic("what else could be? after =-1, and != -1")
			whiteListIndex = append(whiteListIndex, false)
		}
	}
	// hard code the last column emissionQuant as true, always show
	//whiteListIndex = append(whiteListIndex, true)

	//loop through whiteListIndex, for these columns are not -1, check the count of distinct value = 1,
	//for example if the MOVESRUNID only has 1, ignore it, there is no point to show them as both column or filter
	for i := 0; i < len(whiteListIndex)-1; i++ { // -1 because the last column emissionQuant should always show
		if whiteListIndex[i] { //if the column value is not null
			//get distinct query, and see the count or len(returned slice)
			distinctResult := getDistinct(con, dbSelection, tableSelection, moFieldNames[i])
			if len(distinctResult) <= 1 { // if the returned slice only has 1 distinct value, mark the index to false
				whiteListIndex[i] = false
				fmt.Print(" found column that only has 1 distinct value ", moFieldNames[i])
				fmt.Print(" printing updated whiteList v% v%", moFieldNames[i], whiteListIndex[i])
				whiteList = RemoveElementFromSlice(whiteList, moFieldNames[i]) // call func that remove #the column that only have 1 distinct value as well
			}
		}
	}

	return whiteList, whiteListIndex
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

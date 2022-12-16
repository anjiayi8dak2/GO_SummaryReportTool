package main

import (
	"database/sql"
	"fmt"
	"github.com/datasweet/datatable/import/csv"
	"log"
	"os/exec"
	"reflect"
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

//func scanMap(list *sql.Rows) (rows []map[string]interface{}) {
//	fields, _ := list.Columns()
//	for list.Next() {
//		scans := make([]interface{}, len(fields))
//		row := make(map[string]interface{})
//
//		for i := range scans {
//			scans[i] = &scans[i]
//		}
//		list.Scan(scans...)
//		for i, v := range scans {
//			var value = ""
//			if v != nil {
//				value = fmt.Sprintf("%s", v)
//			}
//			row[fields[i]] = value
//		}
//		rows = append(rows, row)
//	}
//	return
//}

func readCSV() {
	dt, err := csv.Import("csv", "temp.csv",
		csv.HasHeader(true),
		csv.AcceptDate("02/01/06 15:04"),
		csv.AcceptDate("2006-01"),
	)
	if err != nil {
		log.Fatalf("reading csv: %v", err)
	}
	//fmt.Println("================== print dt with stdout ===========================================")
	//dt.Print(os.Stdout, datatable.PrintMaxRows(24))
	fmt.Println("================== print dt type===========================================")
	fmt.Println(reflect.TypeOf(dt.Row(0)))
	//
	//fmt.Println("================== test dt2 ===========================================")
	//dt2, err := dt.Aggregate(datatable.AggregateBy{Type: datatable.Count, Field: "pollutantID"})
	//if err != nil {
	//	log.Fatalf("aggregate COUNT('pollutantID'): %v", err)
	//}
	//fmt.Println("================== print dt2 ===========================================")
	//fmt.Println(dt2)
	//fmt.Println(reflect.TypeOf(dt2))
	//
	//groups, err := dt.GroupBy(datatable.GroupBy{
	//	Name: "year",
	//	Type: datatable.Int,
	//	Keyer: func(row datatable.Row) (interface{}, bool) {
	//		if d, ok := row["date"]; ok {
	//			if tm, ok := d.(time.Time); ok {
	//				return tm.Year(), true
	//			}
	//		}
	//		return nil, false
	//	},
	//})
	//if err != nil {
	//	log.Fatalf("GROUP BY 'year': %v", err)
	//}
	//dt3, err := groups.Aggregate(
	//	datatable.AggregateBy{Type: datatable.Sum, Field: "duration"},
	//	datatable.AggregateBy{Type: datatable.CountDistinct, Field: "network"},
	//)
	//if err != nil {
	//	log.Fatalf("Aggregate SUM('duration'), COUNT_DISTINCT('network') GROUP BY 'year': %v", err)
	//}
	//fmt.Println(dt3)
	//TODO: testing some data table functionality here
	//fmt.Println("================== print columns ===========================================")
	//fmt.Println(dt.Columns())
	//fmt.Println(reflect.TypeOf(dt.Columns()))
	//fmt.Println("================== print rows ===========================================")
	//fmt.Println(dt.Rows())
	//fmt.Println(reflect.TypeOf(dt.Rows()))

}

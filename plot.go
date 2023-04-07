package main

import (
	_ "flag"
	"fmt"
	_ "fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	_ "github.com/pkg/browser"
	_ "log"
	"net/http"
	_ "os"
	"strconv"
	"strings"
)

var (

	//Most part of the variables should pass from fyne window
	distanceUnits, massUnits, energyUnits string
	//pollutant name should be here, it should pass from the caller
	titleName = "Defualt Title"
	// This count is useful to generate plot, user can select 1 or 2 field for aggregation
	// if count =1, make regular bar, if count =2, then make stack bar
	xAxisCount = 2
	//Y name could also be activity, it must located at (first row, last element)
	YAxisName = "EmissionQuant"
	// X name is selected by user, they must located at first row, except last element
	XAxisName = "dummyRegClass"

	// first field, it should pass from the caller

	field1         []string
	field2         []string
	valueColumnRow []string
	field2Count    []int
	//field1 = []string{"reg20", "reg41"}
	//// second field, it should pass from the caller
	//field2 = []string{"gas", "diesel", "EV"}
	//// this is the data table passing in, possible have header??

	dummyEmissionQuant = [][]string{}
	dummyGlobalTable   = map[wideTableShapeStruct]float64{}
	dummyGlobalCounter = 0
)

// transpose a 2d matrix
func transpose(slice [][]string) [][]string {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]string, xl)
	for i := range result {
		result[i] = make([]string, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}

func generateBarItems(stackBarCount int, field1_key string, field2_key string) []opts.BarData {

	//eCharts only take BarData type to generate plot
	items := make([]opts.BarData, 0)

	for i := 0; i < stackBarCount; i++ {
		f := dummyGlobalTable[wideTableShapeStruct{field1_key, field2_key}]
		items = append(items, opts.BarData{Value: f})
		fmt.Println("inside generate bar item print key, ", field1_key, " ", field2_key+" ", f, " small stack bar count ", stackBarCount)
		fmt.Println("print global counter ", dummyGlobalCounter, " and field1 with counter index ", field1[dummyGlobalCounter])
		dummyGlobalCounter++
	}
	dummyGlobalCounter = 0
	//fmt.Println("generate bar item ", f)
	return items
}

func barStack() *charts.Bar {
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{
			Name: XAxisName,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: YAxisName,
		}),
		charts.WithTitleOpts(opts.Title{
			Title: titleName,
		}),

		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Bottom: "50%"}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "20%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Title: "Download PNG",
				},
				DataView: &opts.ToolBoxFeatureDataView{
					Show:  true,
					Title: "DataView",
					Lang:  []string{"data view", "turn off", "refresh"},
				},
			}},
		),
	)

	distinctField1 := removeDuplicateStr(field1) //2 ,5
	fmt.Println("print distinct ", distinctField1)
	distinctField2 := removeDuplicateStr(field2) //1, 100
	fmt.Println("print distinct ", distinctField2)

	//iterate small bars on big bar
	for i := 0; i < len(distinctField2); i++ {
		bar.SetXAxis(distinctField1). // big bars horizontal
			//small vertical bars
			AddSeries(distinctField2[i], generateBarItems(len(distinctField2), distinctField1[dummyGlobalCounter], distinctField2[i])). //pass small stack bar count, field1 and field2 for key combination
			SetSeriesOptions(charts.WithBarChartOpts(opts.BarChart{
				Stack: "dummy",
			}))
	}
	dummyGlobalCounter = 0
	return bar
}

func httpserver(w http.ResponseWriter, _ *http.Request) {
	// create a new stack bar instance

	stackBar := barStack()
	fmt.Println("after call bar")
	stackBar.Render(w)
	fmt.Println("after call render")
}
func runPlot(distanceUnits string, massUnits string, energyUnits string, pollutant string, X1 string, X2 string, Y string, queryResult [][]string) {

	//pollutant name should be here, it should pass from the caller
	titleName = "Plot for " + pollutant + " VS " + X1 + " with " + X2
	// This count is useful to generate plot, user can select 1 or 2 field for aggregation
	// if count =1, make regular bar, if count =2, then make stack bar
	xAxisCount = 2
	//Y name could also be activity, it must locate at (first row, last element)
	YAxisName = Y
	// X name is selected by user, they must locate at first row, except last element
	XAxisName = X1

	field1_position := 0
	field2_position := 0
	valueColumn_position := 0

	//loop the existing table, and reshape long to wide end up in map
	dummyMap := map[wideTableShapeStruct]float64{}

	//find column position for field1, field2, and value in the queryResult matrix by searching the first row and do string comparing
	for col := 0; col < len(queryResult[0]); col++ {
		if queryResult[0][col] == X1 { //find field1
			field1_position = col
		}
		if queryResult[0][col] == X2 {
			field2_position = col
		}
		if queryResult[0][col] == Y {
			valueColumn_position = col
		}
	}
	fmt.Println("print field1 field2 and value columns position ", field1_position, field2_position, valueColumn_position)

	//loop through the query result and assign key and value into map
	//skip first row, that is header
	for row := 1; row < len(queryResult); row++ {
		// {firstKey,secondKey} = value
		stringToFloat, _ := strconv.ParseFloat(strings.TrimSpace(queryResult[row][valueColumn_position]), 64)
		fmt.Println("print string convert to float", stringToFloat)
		dummyMap[wideTableShapeStruct{queryResult[row][field1_position], queryResult[row][field2_position]}] = stringToFloat
		//update field1 []string
		field1 = append(field1, queryResult[row][field1_position])
		//and field2 []string
		field2 = append(field2, queryResult[row][field2_position])
	}

	mapCopy(dummyGlobalTable, dummyMap)

	fmt.Println("print map")
	for key, value := range dummyMap {
		fmt.Printf("%s  value is %v\n", key, value)
	}

	fmt.Println("print field1", field1)

	//
	////transposed QueryResult
	//transposedQueryResult := transpose(queryResult)
	//// #1 field, search X1 in the transposed matrix, if detect then copy the rest of the row
	//
	//for i := 0; i < len(transposedQueryResult); i++ {
	//	//only loop the first element of each row for header
	//	if transposedQueryResult[i][0] == X1 { //find field1
	//		columnSize := len(transposedQueryResult[0])
	//		field1 = transposedQueryResult[i] //copy row
	//		field1 = field1[1:columnSize]     //delete header that located at first position of the row
	//	}
	//	if transposedQueryResult[i][0] == X2 { //find field2
	//		columnSize := len(transposedQueryResult[0])
	//		field2 = transposedQueryResult[i]
	//		field2 = field2[1:columnSize]
	//	}
	//	if transposedQueryResult[i][0] == Y { //find value column like emissionQuant or activity
	//		columnSize := len(transposedQueryResult[0])
	//		valueColumnRow = transposedQueryResult[i]
	//		valueColumnRow = valueColumnRow[1:columnSize]
	//	}
	//}
	//fmt.Println("inside runPlot before remove duplicate field 1  ", field1, " and field 2 ", field2)

	//dummyMatrixCounter := 0
	//field1 = removeDuplicate(field1)
	//field2 = removeDuplicate(field2)

	//dummy matrix size field1 x field2
	//dummyMatrix := make([][]string, len(field1))
	//for i := 0; i < len(field1); i++ {
	//	dummyMatrix[i] = make([]string, len(field2))
	//}

	//var dummyMatrix = [][]string{}
	//dummyMatrix = append(dummyMatrix, field1)
	//dummyMatrix = append(dummyMatrix, field2)
	//dummyMatrix = append(dummyMatrix, valueColumnRow)
	//
	////fmt.Println("inside runPlot after remove duplicate field 1  ", field1, " and field 2 ", field2)
	//fmt.Println("inside runPlot dummyMatrix = ", dummyMatrix)

	//for i := 0; i < len(field1); i++ {
	//	for j := 0; j < len(field2); j++ {
	//		dummyMatrix[i] = append(dummyMatrix[i], valueColumnRow[dummyMatrixCounter])
	//		dummyMatrixCounter++
	//	}
	//}
	//dummyMatrixCounter = 0

	//fmt.Println("inside runPlot printing field 1  ", field1, " and field 2 ", field2)
	//fmt.Println("inside runPlot printing dummy value column")
	//fmt.Println(valueColumnRow)

	//copy it to global var
	//dummyEmissionQuant = dummyMatrix
	//copy(dummyEmissionQuant, dummyMatrix)
	//dummyEmissionQuant = dummyMatrix
	//fmt.Println("copy dummyEmissionQuant")
	fmt.Println("before call http server")

	http.HandleFunc("/", httpserver)
	http.ListenAndServe(":8081", nil)
	fmt.Println("after call http server")
}

func setMatrix(dataTable [][]string) {

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

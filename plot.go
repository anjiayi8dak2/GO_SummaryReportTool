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

func generateBarItems(columnIndex int) []opts.BarData {

	//eCharts only take BarData type to generate plot
	items := make([]opts.BarData, 0)
	// transpose data matrix
	// the iteration order of eCharts for stack bar is loop through all series(a color on a bar) then the SetXAxis (a bar)
	// for example first field selection = regClass(20,41), second = fuelType(1,2,9)
	// It will iterate in order of reg20+fuel1, reg20+fuel2, reg20+fuel9 then do the same for 41+1 41+2 41+9 combination

	// what we need return here is a slice of emissionQuant that represent field#1 + field#2 combination,
	// but the data matrix format we have has no easy way to return the value of a column.
	// hence, we transpose the matrix first, then return the row to achieve same goal
	//dummyEmissionQuant2 := transpose(dummyEmissionQuant)
	//columnCount := len(dummyEmissionQuant2[0])
	//save all element that associated with row # selection

	//for j := 0; j < len(field2); j++ {
	//	println("inside generate bar printing dummyEmissionQuant[][] string", dummyEmissionQuant[rowNumber][j])
	//	f, _ := strconv.ParseFloat(dummyEmissionQuant[rowNumber][j], 16)
	//	items = append(items, opts.BarData{Value: f})
	//	println("inside generate bar printing dummyEmissionQuant[][] that assigned into BarData", f)
	//}

	//println("inside generate bar printing dummyEmissionQuant[][] string", valueColumnRow[columnIndex])
	f, _ := strconv.ParseFloat(valueColumnRow[columnIndex], 16)
	items = append(items, opts.BarData{Value: f})
	columnIndex++
	f, _ = strconv.ParseFloat(valueColumnRow[columnIndex], 16)
	items = append(items, opts.BarData{Value: f})

	//println("inside generate bar printing dummyEmissionQuant[][] that assigned into BarData", f)

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
			Title: "stack bar",
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

	fmt.Println("printing matrix dummyEmissionQuant before generating bar items")
	fmt.Println("this is debugging flags that see before")
	fmt.Println(dummyEmissionQuant)

	//
	//for i := 0; i < len(field2); i++ {
	//	bar.SetXAxis(field1). // []string for name of the bars
	//		AddSeries(field2[i], generateBarItems(i)).
	//		SetSeriesOptions(charts.WithBarChartOpts(opts.BarChart{
	//			Stack: "Stack",
	//		}))
	//}

	dummyValueColumnCounter := 0
	distinctField1 := removeDuplicateStr(field1)
	distinctField2 := removeDuplicateStr(field2)
	for i := 0; i < len(distinctField2); i++ {
		//fmt.Println("value column value in for loop", valueColumnRow[i])
		bar.SetXAxis(distinctField1). // []string for name of the bars
						AddSeries(distinctField2[i], generateBarItems(dummyValueColumnCounter)).
						SetSeriesOptions(charts.WithBarChartOpts(opts.BarChart{
				Stack: "Stack",
			}))
		dummyValueColumnCounter = dummyValueColumnCounter + len(distinctField2)
	}
	dummyValueColumnCounter = 0

	return bar
}

func httpserver(w http.ResponseWriter, _ *http.Request) {
	// create a new stack bar instance
	stackBar := barStack()

	stackBar.Render(w)
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
	//transposed QueryResult
	transposedQueryResult := transpose(queryResult)
	// #1 field, search X1 in the transposed matrix, if detect then copy the rest of the row

	for i := 0; i < len(transposedQueryResult); i++ {
		//only loop the first element of each row for header
		if transposedQueryResult[i][0] == X1 { //find field1
			columnSize := len(transposedQueryResult[0])
			field1 = transposedQueryResult[i] //copy row
			field1 = field1[1:columnSize]     //delete header that located at first position of the row
		}
		if transposedQueryResult[i][0] == X2 { //find field2
			columnSize := len(transposedQueryResult[0])
			field2 = transposedQueryResult[i]
			field2 = field2[1:columnSize]
		}
		if transposedQueryResult[i][0] == Y { //find value column like emissionQuant or activity
			columnSize := len(transposedQueryResult[0])
			valueColumnRow = transposedQueryResult[i]
			valueColumnRow = valueColumnRow[1:columnSize]
		}
	}
	fmt.Println("inside runPlot before remove duplicate field 1  ", field1, " and field 2 ", field2)

	//dummyMatrixCounter := 0
	//field1 = removeDuplicate(field1)
	//field2 = removeDuplicate(field2)

	//dummy matrix size field1 x field2
	//dummyMatrix := make([][]string, len(field1))
	//for i := 0; i < len(field1); i++ {
	//	dummyMatrix[i] = make([]string, len(field2))
	//}

	var dummyMatrix = [][]string{}
	dummyMatrix = append(dummyMatrix, field1)
	dummyMatrix = append(dummyMatrix, field2)
	dummyMatrix = append(dummyMatrix, valueColumnRow)

	//fmt.Println("inside runPlot after remove duplicate field 1  ", field1, " and field 2 ", field2)
	fmt.Println("inside runPlot dummyMatrix = ", dummyMatrix)

	//for i := 0; i < len(field1); i++ {
	//	for j := 0; j < len(field2); j++ {
	//		dummyMatrix[i] = append(dummyMatrix[i], valueColumnRow[dummyMatrixCounter])
	//		dummyMatrixCounter++
	//	}
	//}
	//dummyMatrixCounter = 0

	fmt.Println("inside runPlot printing field 1  ", field1, " and field 2 ", field2)
	fmt.Println("inside runPlot printing dummy value column")
	fmt.Println(valueColumnRow)

	//copy it to global var
	//dummyEmissionQuant = dummyMatrix
	//copy(dummyEmissionQuant, dummyMatrix)
	dummyEmissionQuant = dummyMatrix
	fmt.Println("copy dummyEmissionQuant")
	fmt.Println(dummyEmissionQuant)

	http.HandleFunc("/", httpserver)
	http.ListenAndServe(":8081", nil)
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

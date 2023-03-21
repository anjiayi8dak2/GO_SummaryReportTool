package main

import (
	_ "flag"
	_ "fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	_ "github.com/pkg/browser"
	_ "log"
	"net/http"
	_ "os"
)

var (

	//Most part of the variables should pass from fyne window

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
	dummyRegClassItems = []string{"reg20", "reg41"}
	// second field, it should pass from the caller
	dummyFuelTypeItems = []string{"gas", "diesel", "EV"}
	// this is the data table passing in, possible have header??
	dummyEmissionQuant = [][]float64{
		{201, 202, 209}, {411, 412, 419},
	}
)

// transpose a 2d matrix
func transpose(slice [][]float64) [][]float64 {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]float64, xl)
	for i := range result {
		result[i] = make([]float64, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}

func generateBarItems(rowSelection int) []opts.BarData {

	//eCharts only take BarData type to generate plot
	items := make([]opts.BarData, 0)
	// transpose data matrix
	// the iteration order of eCharts for stack bar is loop through all series then the SetXAxis
	// for example first field selection = regClass(20,41), second = fuelType(1,2,9)
	// It will read reg20+fuel1, reg20+fuel2, reg20+fuel9 then do the same for reg41

	// the data table passing into this class is in transposed format, hence we convert here
	dummyEmissionQuant2 := transpose(dummyEmissionQuant)
	columnCount := len(dummyEmissionQuant2[0])
	//save all element that associated with row # selection
	for j := 0; j < columnCount; j++ {
		items = append(items, opts.BarData{Value: dummyEmissionQuant2[rowSelection][j]})
	}

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

	for i := 0; i < len(dummyFuelTypeItems); i++ {
		bar.SetXAxis(dummyRegClassItems). // []string for name of the bars
							AddSeries(dummyFuelTypeItems[i], generateBarItems(i)).
							SetSeriesOptions(charts.WithBarChartOpts(opts.BarChart{
				Stack: "dummyStack",
			}))
	}

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
	//Y name could also be activity, it must located at (first row, last element)
	YAxisName = Y
	// X name is selected by user, they must located at first row, except last element
	XAxisName = X1
	// #1 field, it should pass from the caller
	field1 = []string{"reg20", "reg41"}
	// #2 field, it should pass from the caller
	dummyFuelTypeItems = []string{"gas", "diesel", "EV"}
	// this is the data table passing in, possible have header??
	dummyEmissionQuant = [][]float64{
		{201, 202, 209}, {411, 412, 419},
	}

	http.HandleFunc("/", httpserver)
	http.ListenAndServe(":8081", nil)
}

package main

import (
	_ "flag"
	"fmt"
	_ "fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	_ "github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	_ "github.com/pkg/browser"
	_ "log"
	"os"
	_ "os"
	"os/exec"
	_ "os/exec"
	"runtime"
	"strconv"
	"strings"
)

var (

	//Most part of the variables should pass from fyne window
	//distanceUnits, massUnits, energyUnits string
	//pollutant name should be here, it should pass from the caller
	titleName    = "dummyTitleName"
	subTitleName = "dummySubTitleName"
	// This count is useful to generate plot, user can select 1 or 2 field for aggregation
	// if count =1, make regular bar, if count =2, then make stack bar
	xAxisCount = 2
	// YAxisName Y name could also be activity, it must locate at (first row, last element)
	YAxisName = "dummyYAxisName"
	// XAxisName X name is selected by user, they must locate at first row, except last element
	XAxisName = "dummyXAxisName"

	// field in slice
	field1 []string
	field2 []string

	contingencyGlobalTable = map[wideTableShapeStruct]float64{}

	distinctField1 []string
	distinctField2 []string
)

func generateBarItems(field2_key string) []opts.BarData {
	//eCharts only take BarData type to generate plot
	items := make([]opts.BarData, 0)

	//use map with 2 keys to find corresponding value and assign to items
	for _, field1_key := range distinctField1 {
		f := contingencyGlobalTable[wideTableShapeStruct{field1_key, field2_key}]
		items = append(items, opts.BarData{Value: f})
		fmt.Println("inside generate bar item print key, ", field1_key, " ", field2_key+" ", f)
	}

	return items
}

func barStack() *charts.Bar {
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			PageTitle: "default PageTitle",
			Width:     "1000px",
			Height:    "500px",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
			Name: XAxisName,
			AxisLabel: &opts.AxisLabel{
				Show:         true,
				Rotate:       45,
				Interval:     "0",
				ShowMaxLabel: true,
				ShowMinLabel: true,
				// You may set it to be 0 to display all labels compulsively.
				// If it is set to be 1, it means that labels are shown once after one label.
				// And if it is set to be 2, it means labels are shown once after two labels, and so on
			},
			Show: true,
		}),

		charts.WithYAxisOpts(opts.YAxis{
			Type: "value",
			Name: YAxisName,
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    titleName,
			Subtitle: subTitleName,
			Top:      "5%",
			Left:     "5%",
		}),

		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true, Right: "right", Bottom: "center", Orient: "vertical"}),

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
					Lang:  []string{"Data View", "Turn Off", "Refresh"},
				},
			}},
		),
		//charts.WithDataZoomOpts(opts.DataZoom{
		//	Type:  "slider",
		//	Start: 0,
		//	End:   100,
		//}),
		charts.WithGridOpts(opts.Grid{
			Bottom: "200px",
			Left:   "100px",
			Right:  "300px",
			Top:    "100px",
			//ContainLabel: false,
		}),
	)

	distinctField1 = removeDuplicateStr(field1)
	distinctField2 = removeDuplicateStr(field2)

	//iterate small bars on big bar
	for i := 0; i < len(distinctField2); i++ {
		bar.SetXAxis(distinctField1). // big bars horizontal
			//small vertical bars
			AddSeries(distinctField2[i], generateBarItems(distinctField2[i])). //pass small stack bar count, field1 and field2 for key combination
			SetSeriesOptions(charts.WithBarChartOpts(opts.BarChart{
				Stack: "stack",
				//BarCategoryGap: "50%",
			}))
	}
	return bar
}

func runPlot(distanceUnits string, massUnits string, energyUnits string, X1 string, X2 string, Y string, queryResult [][]string) {

	//pollutant name should be here, it should pass from the caller
	titleName = "Plot for " + X1 + " VS " + X2
	// This count is useful to generate plot, user can select 1 or 2 field for aggregation
	// if count =1, make regular bar, if count =2, then make stack bar
	xAxisCount = 2
	//Y name could also be activity, it must locate at (first row, last element)
	YAxisName = Y
	subTitleName = " in " + massUnits + "/" + distanceUnits + "/" + energyUnits
	// X name is selected by user, they must locate at first row, except last element
	XAxisName = X1

	field1_position := 0
	field2_position := 0
	valueColumn_position := 0

	//loop the existing table, and reshape long to wide end up in map, aka contingency table
	longToWideMap := map[wideTableShapeStruct]float64{}

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
		longToWideMap[wideTableShapeStruct{queryResult[row][field1_position], queryResult[row][field2_position]}] = stringToFloat
		//update field1 []string
		field1 = append(field1, queryResult[row][field1_position])
		//and field2 []string
		field2 = append(field2, queryResult[row][field2_position])
	}
	//clear and copy new
	clearMap(contingencyGlobalTable)
	mapCopy(contingencyGlobalTable, longToWideMap)
	fmt.Println("========print map contingencyGlobalTable========= ", contingencyGlobalTable)

	//save the plot into html file
	folderPath := getAbsolutePath()
	f, _ := os.Create(folderPath + "\\" + "plot.html")
	fmt.Println("printing path+file name", folderPath+"\\"+"plot.html")
	stackBar := barStack()
	stackBar.Render(f)
	f.Close()

	//open that html just saved
	_, ExistErr := os.Stat(folderPath + "\\" + "plot.html")
	if ExistErr != nil {
		if os.IsNotExist(ExistErr) {
			fmt.Println(folderPath + "\\" + "plot.html does not exist")
		} else {
			fmt.Println(ExistErr)
		}
	} else {
		var cmd string
		switch runtime.GOOS { //command for different OS?
		case "linux":
			cmd = "xdg-open"
		case "darwin":
			cmd = "open"
		case "windows":
			cmd = "cmd"
		default:
			fmt.Println("unsupported operating system")
		}

		cmdErr := exec.Command(cmd, "/C", "start", folderPath+"\\"+"plot.html").Run()
		if cmdErr != nil {
			fmt.Println("cmd error", cmdErr)
		}
	}
}

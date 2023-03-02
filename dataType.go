package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"math"
)

//TODO: add activity output

// Movesoutput all the field in movesoutput table
type Movesoutput struct {
	MOVESRunID    int
	iterationID   int
	yearID        int
	monthID       int
	dayID         int
	hourID        int
	stateID       int
	countyID      int
	zoneID        int
	linkID        int
	pollutantID   int
	processID     int
	sourceTypeID  int
	regClassID    int
	fuelTypeID    int
	fuelSubTypeID int
	modelYearID   int
	roadTypeID    int
	SCC           int
	engTechID     int
	sectorID      int
	hpID          int
	emissionQuant float64
}

type Rateperdistance struct {
	MOVESScenarioID string
	MOVESRunID      int
	yearID          int
	monthID         int
	dayID           int
	hourID          int
	linkID          int
	pollutantID     int
	processID       int
	sourceTypeID    int
	regClassID      int
	SCC             int
	fuelTypeID      int
	modelYearID     int
	roadTypeID      int
	avgSpeedBinID   int
	temperature     float64
	relHumidity     float64
	ratePerDistance float64
}

type Rateperhour struct {
	MOVESScenarioID string
	MOVESRunID      int
	yearID          int
	monthID         int
	dayID           int
	hourID          int
	linkID          int
	pollutantID     int
	processID       int
	sourceTypeID    int
	regClassID      int
	SCC             int
	fuelTypeID      int
	modelYearID     int
	roadTypeID      int
	temperature     float64
	relHumidity     float64
	ratePerHour     float64
}

type Rateperprofile struct {
	MOVESScenarioID      string
	MOVESRunID           int
	temperatureProfileID int
	yearID               int
	dayID                int
	hourID               int
	pollutantID          int
	processID            int
	sourceTypeID         int
	regClassID           int
	SCC                  int
	fuelTypeID           int
	modelYearID          int
	temperature          float64
	relHumidity          float64
	ratePerVehicle       float64
}

type Rateperstart struct {
	MOVESScenarioID string
	MOVESRunID      int
	yearID          int
	monthID         int
	dayID           int
	hourID          int
	zoneID          int
	sourceTypeID    int
	regClassID      int
	SCC             int
	fuelTypeID      int
	modelYearID     int
	pollutantID     int
	processID       int
	temperature     float64
	relHumidity     float64
	ratePerStart    float64
}

type Ratepervehicle struct {
	MOVESScenarioID string
	MOVESRunID      int
	yearID          int
	monthID         int
	dayID           int
	hourID          int
	zoneID          int
	pollutantID     int
	processID       int
	sourceTypeID    int
	regClassID      int
	SCC             int
	fuelTypeID      int
	modelYearID     int
	temperature     float64
	relHumidity     float64
	ratePerVehicle  float64
}

type Startspervehicle struct {
	MOVESScenarioID  string
	MOVESRunID       int
	yearID           int
	monthID          int
	dayID            int
	hourID           int
	zoneID           int
	sourceTypeID     int
	regClassID       int
	SCC              int
	fuelTypeID       int
	modelYearID      int
	startsPerVehicle float64
}

var _ fyne.Widget = &HeaderTable{}

type HeaderTable struct {
	widget.BaseWidget
	TableOpts *TableOpts
	Header    *widget.Table
	Data      *widget.Table
}

type BindingConverter func(interface{}) string

func NewHeaderTable(tableOpts *TableOpts) *HeaderTable {
	t := &HeaderTable{
		TableOpts: tableOpts,
		Header: widget.NewTable(
			// Dimensions (rows, cols)
			func() (int, int) { return 1, len(tableOpts.ColAttrs) },
			// Default value
			func() fyne.CanvasObject { return widget.NewLabel("the content") },
			// Cell values
			func(cellID widget.TableCellID, o fyne.CanvasObject) {
				l := o.(*widget.Label)
				opts := tableOpts.ColAttrs[cellID.Col]
				l.TextStyle = opts.HeaderStyle.TextStyle
				l.Alignment = opts.HeaderStyle.Alignment
				l.Wrapping = opts.HeaderStyle.Wrapping
				l.SetText(opts.Header)
			},
		),
		Data: widget.NewTable(dataTableLengthFunc(tableOpts), dataTableCreateFunc, dataTableUpdateFunc(tableOpts)),
	}
	t.ExtendBaseWidget(t)

	// Set Column widths
	refWidth := widget.NewLabel(t.TableOpts.RefWidth).MinSize().Width
	for i, colAttr := range t.TableOpts.ColAttrs {
		if t.Data != nil {
			t.Data.SetColumnWidth(i, float32(colAttr.WidthPercent)/100.0*refWidth)
		}
		if t.Header != nil {
			t.Header.SetColumnWidth(i, float32(colAttr.WidthPercent)/100.0*refWidth)
		}
	}

	return t
}

// ****************** Renderer *******************************

var _ fyne.WidgetRenderer = headerTableRenderer{}

type headerTableRenderer struct {
	headerTable *HeaderTable
	container   *fyne.Container
}

func (h *HeaderTable) CreateRenderer() fyne.WidgetRenderer {
	return headerTableRenderer{
		headerTable: h,
		container:   container.NewBorder(h.Header, nil, nil, nil, h.Data),
		// container:   container.NewVBox(h.Header, h.Data),
	}
}

func (r headerTableRenderer) MinSize() fyne.Size {
	dataMinSize := fyne.NewSize(0, 0)
	if r.headerTable.Data != nil {
		dataMinSize = r.headerTable.Data.MinSize()
	}
	headerMinSize := fyne.NewSize(0, 0)
	if r.headerTable.Header != nil {
		dataMinSize = r.headerTable.Header.MinSize()
	}
	return fyne.NewSize(
		float32(math.Max(float64(dataMinSize.Width), float64(headerMinSize.Width))),
		dataMinSize.Height+headerMinSize.Height)
}

func (r headerTableRenderer) Layout(s fyne.Size) {
	r.container.Resize(s)
}

func (r headerTableRenderer) Destroy() {
}

func (r headerTableRenderer) Refresh() {
	r.container.Refresh()
}

func (r headerTableRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container}
}

type CellStyle struct {
	Alignment fyne.TextAlign
	TextStyle fyne.TextStyle
	Wrapping  fyne.TextWrap
}

type ColAttr struct {
	Converter    BindingConverter
	DataStyle    CellStyle
	Header       string
	HeaderStyle  CellStyle
	Name         string
	WidthPercent int
}

type TableOpts struct {
	Bindings         []binding.Struct
	ColAttrs         []ColAttr
	OnDataCellSelect func(cellID widget.TableCellID)
	RefWidth         string
}

type Header struct {
	widget.Table
}

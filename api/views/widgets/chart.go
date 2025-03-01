package widgets

import (
	"github.com/coroot/coroot/timeseries"
	"sort"
)

type ChartType string

type Annotation struct {
	Name string          `json:"name"`
	X1   timeseries.Time `json:"x1"`
	X2   timeseries.Time `json:"x2"`
	Icon string          `json:"icon"`
}

type Chart struct {
	Ctx timeseries.Context `json:"ctx"`

	Title       string       `json:"title"`
	Series      []*Series    `json:"series"`
	Threshold   *Series      `json:"threshold"`
	Featured    bool         `json:"featured"`
	IsStacked   bool         `json:"stacked"`
	IsSorted    bool         `json:"sorted"`
	IsColumn    bool         `json:"column"`
	ColorShift  int          `json:"color_shift"`
	Annotations []Annotation `json:"annotations"`
}

func NewChart(ctx timeseries.Context, title string) *Chart {
	return &Chart{Ctx: ctx, Title: title}
}

func (chart *Chart) Stacked() *Chart {
	chart.IsStacked = true
	return chart
}

func (chart *Chart) Sorted() *Chart {
	chart.IsSorted = true
	return chart
}

func (chart *Chart) Column() *Chart {
	chart.IsColumn = true
	chart.IsStacked = true
	return chart
}

func (chart *Chart) ShiftColors() *Chart {
	chart.ColorShift = 1
	return chart
}

func (chart *Chart) AddAnnotation(name string, start, end timeseries.Time, icon string) *Chart {
	chart.Annotations = append(chart.Annotations, Annotation{X1: start, X2: end, Name: name, Icon: icon})
	return chart
}

func (chart *Chart) AddMany(series []timeseries.Named) *Chart {
	for _, v := range series {
		chart.AddSeries(v.Name, v.Series)
	}
	return chart
}

func (chart *Chart) AddSeries(name string, data timeseries.TimeSeries, color ...string) *Chart {
	if data == nil || data.IsEmpty() {
		return chart
	}
	s := &Series{Name: name, Data: data}
	if len(color) > 0 {
		s.Color = color[0]
	}
	chart.Series = append(chart.Series, s)
	return chart
}

func (chart *Chart) SetThreshold(name string, data timeseries.TimeSeries, aggFunc timeseries.F) *Chart {
	if data == nil {
		return chart
	}
	if chart.Threshold == nil {
		chart.Threshold = &Series{Name: name, Data: timeseries.Aggregate(aggFunc), Color: "black"}
	}
	chart.Threshold.Data.(*timeseries.AggregatedTimeseries).AddInput(data)
	return chart
}

func (chart *Chart) Feature() *Chart {
	chart.Featured = true
	return chart
}

type Series struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Fill  bool   `json:"fill"`

	Data timeseries.TimeSeries `json:"data"`
}

type ChartGroup struct {
	Title  string   `json:"title"`
	Charts []*Chart `json:"charts"`
}

func (cg *ChartGroup) GetOrCreateChart(ctx timeseries.Context, title string) *Chart {
	for _, ch := range cg.Charts {
		if ch.Title == title {
			return ch
		}
	}
	ch := NewChart(ctx, title)
	cg.Charts = append(cg.Charts, ch)
	return ch
}

func (cg *ChartGroup) AutoFeatureChart() {
	if len(cg.Charts) < 2 {
		return
	}
	type weightedChart struct {
		ch *Chart
		w  float64
	}
	for _, ch := range cg.Charts {
		if ch.Featured {
			return
		}
	}
	charts := make([]weightedChart, 0, len(cg.Charts))
	for _, ch := range cg.Charts {
		var w float64
		for _, s := range ch.Series {
			w += timeseries.Reduce(timeseries.NanSum, s.Data)
		}
		charts = append(charts, weightedChart{ch: ch, w: w})
	}
	sort.Slice(charts, func(i, j int) bool {
		return charts[i].w > charts[j].w
	})
	if charts[0].w/charts[1].w > 1.2 {
		charts[0].ch.Featured = true
	}
}

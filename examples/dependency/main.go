package main

import (
	"fmt"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

func init() {
	fmt.Println("from init")
}

func main() {
	latencyView := &view.View{
		Name:        "myapp/latency",
		Measure:     stats.Int64("example.com/measure/openconns", "open connections", stats.UnitDimensionless),
		Description: "The distribution of the latencies",
		Aggregation: view.Distribution(0, 25, 100, 200, 400, 800, 10000),
	}

	view.Register(latencyView)

	fmt.Println("from dependency")
}

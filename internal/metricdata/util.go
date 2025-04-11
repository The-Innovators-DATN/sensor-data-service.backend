package metricdata

// import "internal/metricdata/model"

// GenerateForecast generates forecasted data points.
func GenerateForecast(series []MetricPoint, horizon, timeStep int32) []MetricPoint {
	var forecast []MetricPoint
	// Simulate forecast logic here
	for i := 0; i < int(horizon); i++ {
		forecast = append(forecast, MetricPoint{
			Datetime: "forecasted-time",
			Value:    float32(i * int(timeStep)), // placeholder
		})
	}
	return forecast
}

// ProcessAnomalies applies anomaly detection logic based on local error threshold.
// func ProcessAnomalies(seriesData *[]SeriesData, threshold float32) {
// 	for i := range *seriesData {
// 		for j := range (*seriesData)[i].Series {
// 			if (*seriesData)[i].Series[j].Value > threshold {
// 				(*seriesData)[i].Series[j].Anomaly = true
// 			}
// 		}
// 	}
// }

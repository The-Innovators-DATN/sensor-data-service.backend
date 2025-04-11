package metricdata

// TimeRange represents the time range for querying data.
type TimeRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// ForecastConfig represents the configuration for the forecast.
type ForecastConfig struct {
	Enabled  bool  `json:"enabled"`
	TimeStep int32 `json:"time_step"`
	Horizon  int32 `json:"horizon"`
}

// AnomalyDetectionConfig represents anomaly detection configuration.
type AnomalyDetectionConfig struct {
	Enabled             bool    `json:"enabled"`
	LocalErrorThreshold float32 `json:"local_error_threshold"`
}

// MetricPoint represents a point of data in the metric series.
type MetricPoint struct {
	Datetime   string  `json:"datetime"`
	Value      float32 `json:"value"`
	LocalError float32 `json:"local_error"`
	// Anomaly  bool    `json:"anomaly"`
}

// SeriesSelector represents a series that the client can select for the metric.
type SeriesSelector struct {
	RefID      string `json:"ref_id"`
	TargetType int32  `json:"target_type"`
	TargetID   int32  `json:"target_id"`
	MetricID   int32  `json:"metric_id"`
}

// SeriesData represents the series of data for a specific metric and target.
type SeriesData struct {
	RefID      string        `json:"ref_id"`
	TargetType int32         `json:"target_type"`
	TargetID   int32         `json:"target_id"`
	MetricID   int32         `json:"metric_id"`
	Series     []MetricPoint `json:"series"`
	Forecast   []MetricPoint `json:"forecast"`
}

// MetricSeriesRequest represents the request for fetching metric data.
type MetricSeriesRequest struct {
	ChartType        string                 `json:"chart_type"`
	TimeRange        TimeRange              `json:"time_range"`
	Forecast         ForecastConfig         `json:"forecast"`
	AnomalyDetection AnomalyDetectionConfig `json:"anomaly_detection"`
	Series           []SeriesSelector       `json:"series"`
}

// MetricSeriesResponse represents the response to the metric data request.
type MetricSeriesResponse struct {
	Results []SeriesData `json:"results"`
}

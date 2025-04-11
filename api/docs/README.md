# /api/v1/locations

LocationType	Mô tả
station	Trạm    quan trắc (có sensor cụ thể)
river_basin	    Lưu vực sông (tập hợp catchment)
catchment	    Diện tích thoát nước (bao gồm nhiều trạm)
water_body	    Đơn vị thủy văn tổng thể (hồ, sông, biển...)

# /api/v1/locations/{id}

# /api/v1/dashboard/data

GET /api/v0/metrics/data
{
  "chart_type": "line",
  "from": "2024-12-01T00:00:00Z",
  "to": "2024-12-31T23:59:59Z",
  "forecast": {
    "enabled": true,
    "time_step": 3600,
    "horizon": 5
  },
  "anomaly_detection": {
    "enabled": true,
    "local_error_threshold": 0.8
  },
  "series": [
    {
      "refId": "A",
      "target_type": "station",
      "target_id": "1",
      "metric_id": "1"
    },
    {
      "refId": "B",
      "target_type": "water_body",
      "target_id": "2",
      "metric_id": "1"
    }
  ]
}

# Response
 [
    {
      "refId": "A",
      "target_type": "station",
      "target_id": "1",
      "metric_id": "1",
      "series": [
        { "datetime": "2024-12-01T02:00:00Z", "value": 25.1, "anomaly": false },
        { "datetime": "2024-12-01T03:00:00Z", "value": 25.2, "anomaly": false },
        { "datetime": "2024-12-01T04:00:00Z", "value": 25.3, "anomaly": true },
        { "datetime": "2024-12-01T05:00:00Z", "value": 25.4, "anomaly": false },
        { "datetime": "2024-12-01T06:00:00Z", "value": 25.5, "anomaly": false }
      ],
      "forecast": [
        { "datetime": "2024-12-01T07:00:00Z", "value": 25.6, "anomaly": false },
        { "datetime": "2024-12-01T08:00:00Z", "value": 25.7, "anomaly": false },
        { "datetime": "2024-12-01T09:00:00Z", "value": 25.8, "anomaly": false },
        { "datetime": "2024-12-01T10:00:00Z", "value": 25.9, "anomaly": true },
        { "datetime": "2024-12-01T11:00:00Z", "value": 26.0, "anomaly": false }
      ]
    },
    {
      "refId": "B",
      "target_type": "water_body",
      "target_id": "2",
      "metric_id": "1",
      "series": [
        { "datetime": "2024-12-01T02:00:00Z", "value": 26.1, "anomaly": false },
        { "datetime": "2024-12-01T03:00:00Z", "value": 26.2, "anomaly": false },
        { "datetime": "2024-12-01T04:00:00Z", "value": 26.3, "anomaly": true },
        { "datetime": "2024-12-01T05:00:00Z", "value": 26.4, "anomaly": false },
        { "datetime": "2024-12-01T06:00:00Z", "value": 26.5, "anomaly": false }
      ],
      "forecast": [
        { "datetime": "2024-12-01T07:00:00Z", "value": 26.6, "anomaly": false },
        { "datetime": "2024-12-01T08:00:00Z", "value": 26.7, "anomaly": false },
        { "datetime": "2024-12-01T09:00:00Z", "value": 26.8, "anomaly": false },
        { "datetime": "2024-12-01T10:00:00Z", "value": 26.9, "anomaly": true },
        { "datetime": "2024-12-01T11:00:00Z", "value": 27.0, "anomaly": false }
      ]
    }
  ]
  

curl -X POST http://localhost:5000/predict_future -H "Content-Type: application/json" -d '{ "last_raw_data": { "metric": "pH", "value": 8.12, "station_id": 287, "datetime": 1505725800, "unit": "phunits" }, "horizon": 5, "time_step": 10}'
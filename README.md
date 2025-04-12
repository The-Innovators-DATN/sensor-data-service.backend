
# Sensor Data Service Backend

Backend service for environmental station monitoring, supporting gRPC and RESTful API for managing parameters (e.g. pH, TSS, DO, etc).

> 📚 [Danh sách sông của Việt Nam](https://vi.wikipedia.org/wiki/Th%E1%BB%83_lo%E1%BA%A1i:S%C3%B4ng_c%E1%BB%A7a_Vi%E1%BB%87t_Nam)

---

## 🛠 Prerequisites

### System Dependencies

```bash
sudo apt update
sudo apt install -y libpq-dev gdal-bin libgdal-dev golang
```

> GDAL is needed if dealing with spatial data.

---

### Python Env (optional for scraping tools)

```bash
pip install virtualenv
python3 -m virtualenv scrapper-env
source scrapper-env/bin/activate
pip install -r requirements.txt
```

---

## ⚙️ Golang Setup

### Install Go (if not yet)

```bash
wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
source ~/.bashrc   # or ~/.zshrc
```

---

### Install gRPC and Protobuf tools

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

---

### Project Init & Dependencies

```bash
go mod init github.com/The-Innovators-DATN/sensor-data-service.backend

go get github.com/joho/godotenv
go get github.com/spf13/viper
go get github.com/ClickHouse/clickhouse-go/v2
go get github.com/redis/go-redis/v9
```

---

### Optional: Golang Migrate

```bash
curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" | sudo tee /etc/apt/sources.list.d/migrate.list
sudo apt-get update
sudo apt-get install -y migrate
```

---

## 🧱 Build & Run

```bash
go mod tidy
go mod vendor
go build -mod=vendor
go run cmd/station/main.go
```

---

## 🧬 Protobuf Compilation

```bash
protoc -I api/proto   --go_out=paths=import:api   --go-grpc_out=paths=import:api   --grpc-gateway_out=import:api --go-grpc_out=paths=source_relative:api/pb \  api/proto/parameter.proto
```
```bash
protoc \
  -I=api/proto \
  -I=third_party \
  --go_out=paths=source_relative:api/pb/metricdatapb/ \
  --go-grpc_out=paths=source_relative:api/pb/metricdatapb/\
  --grpc-gateway_out=paths=source_relative:api/pb/metricdatapb/ \
  api/proto/metricdata.proto
```
---

## 🔍 API Testing

### 📦 Install grpcurl

```bash
wget https://github.com/fullstorydev/grpcurl/releases/download/v1.8.9/grpcurl_1.8.9_linux_x86_64.tar.gz
tar -xzf grpcurl_1.8.9_linux_x86_64.tar.gz
chmod +x grpcurl
sudo mv grpcurl /usr/local/bin/
```

---

### 🧪 gRPC Examples

```bash
grpcurl -plaintext localhost:8080 list
grpcurl -plaintext localhost:8080 parameter.ParameterService/ListParameters

grpcurl -plaintext -d '{"id": 1}' localhost:8080 parameter.ParameterService/GetParameter

grpcurl -plaintext -d '{
  "name": "pH",
  "unit": "mg/L",
  "parameter_group": "water",
  "description": "Acidity level"
}' localhost:8080 parameter.ParameterService/CreateParameter

grpcurl -plaintext -d '{"id": 1}' localhost:8080 parameter.ParameterService/DeleteParameter
```

---

### 🌐 RESTful Examples (via grpc-gateway)

```bash
curl http://localhost:8081/v0/parameters

curl -X POST http://localhost:8081/v0/parameters   -H "Content-Type: application/json"   -d '{
    "name": "TSS",
    "unit": "mg/L",
    "parameter_group": "solid",
    "description": "Total Suspended Solids"
}'

curl -X PUT http://localhost:8081/v0/parameters/1   -H "Content-Type: application/json"   -d '{
    "id": 1,
    "name": "TSS updated",
    "unit": "mg/L",
    "parameter_group": "solid",
    "description": "Updated TSS desc"
}'

curl -X DELETE http://localhost:8081/v0/parameters/1
```

---

## ✅ Status

- [x] gRPC API ready
- [x] REST gateway
- [x] Redis caching
- [x] PostgreSQL with pgx
- [ ] Authentication & RBAC
- [ ] Metrics data
- [ ] Dashboard Configuration load
- [ ] Refresh panel
---

> Need help? Ping the dev team or drop a GitHub issue. Happy hacking 🚀

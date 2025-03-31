# sensor-data-service.backend
https://vi.wikipedia.org/wiki/Th%E1%BB%83_lo%E1%BA%A1i:S%C3%B4ng_c%E1%BB%A7a_Vi%E1%BB%87t_Nam

sudo apt install libpq-dev gdal-bin libgdal-dev

pip install virtualenv
python -m virtualenv scrapper-env
python3 -m virtualenv scrapper-env
pip install -r requirements.txt

apt update 
apt install golang

wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
export PATH=/usr/local/go/bin:$PATH
source ~/.bashrc   # or ~/.zshrc

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

go mod init github.com/The-Innovators-DATN/sensor-data-service.backend
go get github.com/joho/godotenv

go get github.com/spf13/viper
go get github.com/ClickHouse/clickhouse-go/v2

curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey| apt-key add -
echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
apt-get update
apt-get install -y migrate

go get github.com/redis/go-redis/v9

go run cmd/station/main.go

go mod tidy

go mod vendor

go build -mod=vendor


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
go mod tidy
go get github.com/joho/godotenv
go run cmd/station/main.go

# sensor-data-service.backend
https://vi.wikipedia.org/wiki/Th%E1%BB%83_lo%E1%BA%A1i:S%C3%B4ng_c%E1%BB%A7a_Vi%E1%BB%87t_Nam

sudo apt install libpq-dev gdal-bin libgdal-dev

pip install virtualenv
python -m virtualenv scrapper-env
python3 -m virtualenv scrapper-env
pip install -r requirements.txt

docker run -d \
  --name emqx-enterprise \
  -p 1883:1883 \
  -p 18083:18083 \
  -e EMQX_LOADED_PLUGINS="emqx_bridge_kafka,emqx_management,emqx_dashboard" \
  emqx/emqx-enterprise:5.3.1

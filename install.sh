#!/bin/bash

set -e

echo "installing grafana …"
sudo apt-get install grafana

echo "download prometheus …"
wget "https://github.com/prometheus/prometheus/releases/download/v1.7.0/prometheus-1.7.0.linux-armv6.tar.gz"

echo "installing prometheus …"
tar -C /usr/local/prometheus -zxf prometheus-1.7.0.linux-armv6.tar.gz

echo "installing prometheus config"
sudo mkdir /etc/prometheus
sudo cp prometheus.yml /etc/prometheus/

echo "downloading Go …"
wget https://storage.googleapis.com/golang/go1.8.linux-armv6l.tar.gz

echo "extracting Go …"
sudo tar -C /usr/local -xzf go1.8.linux-armv6l.tar.gz

echo "setting up Go bin environment"
echo 'export PATH=$PATH:/usr/local/go/bin' >> $HOME/.profile
source $HOME/.profile

echo "getting dependencies for the climate.go server …"
go get -d .

echo "compiling C driver for the sensor …"
make -C $GOPATH/src/github.com/pakohan/dht/

echo "compiling climate.go server …"
go build climate.go

echo "installing climate server"
sudo mv climate /opt/

echo "installing systemd services"
sudo cp prometheus.service /etc/systemd/system/
sudo cp climate.service /etc/systemd/system/

echo "enabling systemd services …"
sudo systemctl enable climate.service
sudo systemctl enable prometheus.service
sudo systemctl enable grafana.service

IP=$(ifconfig eth0 | grep 'inet addr' | cut -d: -f2 | awk '{print $1}')

echo "finished!"
echo "Grafana    : http://$IP:3000"
echo "Prometheus : http://$IP:9090"
echo "Climate    : http://$IP:8080"

# rpi-climate-monitor

This repository contains everything to get a Grafana dashboard running for the
DHT11 temparature and humidity sensor. Changing it to the DHT22 should be really
easy.

Run ./install to install Grafana, Go, Prometheus and the climate.go server.
The server also contains code for letting and LED blink whenever the sensor is
being read. Remove it when the there is no LED attached.

It will systemd services for the Grafana, the Prometheus and the climate.go
server so it can. When everthing is done, you should access the RPi Prometheus
dashboard and check whether it is set up correctly. Then add a new Prometheus
data source to Grafana and add the grafana_dashboard.json as your dashboard and
you now should be seing what is going on around your sensor.

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/pakohan/dht"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

var (
	climate = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "climate",
		Help: "humidty and temparature measured",
	}, []string{"type"})
	errors = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "errors",
		Help: "shows whether the last measuring was an error",
	})
)

func init() {
	prometheus.MustRegister(climate)
	prometheus.MustRegister(errors)
}

func main() {
	log.SetFlags(log.LUTC | log.Lshortfile)
	http.Handle("/metrics", promhttp.Handler())
	go readSensor()
	log.Println("starting server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func readSensor() {
	r := raspi.NewAdaptor()
	defer r.Finalize()
	led := gpio.NewLedDriver(r, "11")

	measure(led)
	for range time.Tick(1 * time.Minute) {
		measure(led)
	}
}

func measure(led *gpio.LedDriver) {
	err := led.On()
	if err != nil {
		log.Printf("err turning led on: %s", err.Error())
	}

	h, t, err := dht.GetSensorData(dht.SensorDHT11, 4)
	if err != nil {
		log.Printf("err retrieving data from sensor: %s", err.Error())
		errors.Set(1)
	} else {
		climate.WithLabelValues("humidity").Set(float64(h))
		climate.WithLabelValues("temparature").Set(float64(t))
		errors.Set(0)
	}

	err = led.Off()
	if err != nil {
		log.Printf("err turning led off: %s", err.Error())
	}
}

package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/<%= organization %>/<%= repo %>/models"
	"github.com/<%= organization %>/<%= repo %>/services"
)

func main() {
	service := flag.String("service", "http", "run a service")
	flag.Parse()

	// database := durable.NewDatabase()
	models.InitDB()
	log.Println(*service)

	switch *service {
	case "http":
		go func() {
			runtime.SetBlockProfileRate(1) // 开启对阻塞操作的跟踪
			// models.StartWithHttpServiceJob()
			_ = http.ListenAndServe("0.0.0.0:6060", nil)
		}()
		err := StartHTTP()
		if err != nil {
			log.Println(err)
		}
	default:
		hub := services.NewHub()
		err := hub.StartService(*service)
		if err != nil {
			log.Println(err)
		}
	}
}

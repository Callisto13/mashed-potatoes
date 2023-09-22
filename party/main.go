package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/callisto13/mashed-potatoes/party/handler"
	httphandler "github.com/callisto13/mashed-potatoes/party/handler/http"
	"github.com/callisto13/mashed-potatoes/party/handler/nats"
)

const (
	DEFAULT_NATS_SERVER = "http://localhost:4222"
)

func main() {
	var (
		natsAddress string
		protocol    string
	)

	flag.StringVar(&natsAddress, "nats-address", DEFAULT_NATS_SERVER, "address + port of running NATs service")
	flag.StringVar(&protocol, "protocol", "", "cloudevent protocol to use")

	flag.Parse()

	// this is deliberately the wrong abstraction
	var h handler.ProtocolHandler

	switch protocol {
	case "nats":
		h = nats.Handler{
			NatsAddress: natsAddress,
		}
	case "http":
		h = httphandler.Handler{
			RegisteredProviders: httphandler.Providers,
		}
	default:
		log.Fatal("unrecognised protocol, choose 'nats' or 'http'")
	}

	http.HandleFunc("/", ping)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/enrol", h.Enrol)

	log.Println("starting on :8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal("failed to get this party started")
	}
}

func ping(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	message := map[string][]string{"available-endpoints": {"/ping", "/enrol"}}
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("could not marshal response")
	}

	if _, err := w.Write([]byte(data)); err != nil {
		log.Println("could not write response")
	}
	log.Println("ping hit")
}

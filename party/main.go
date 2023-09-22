package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/callisto13/mashed-potatoes/party/handler"
	grpchandler "github.com/callisto13/mashed-potatoes/party/handler/grpc"
	httphandler "github.com/callisto13/mashed-potatoes/party/handler/http"
	"github.com/callisto13/mashed-potatoes/party/handler/nats"
)

func main() {
	var (
		protocol string
	)

	flag.StringVar(&protocol, "protocol", "", "cloudevent protocol to use")

	flag.Parse()

	// this is deliberately the wrong abstraction
	var h handler.ProtocolHandler

	switch protocol {
	case "nats":
		h = nats.Handler{
			NatsAddress: nats.DEFAULT_NATS_SERVER,
		}
	case "http":
		h = httphandler.Handler{
			RegisteredProviders: httphandler.Providers,
		}
	case "grpc":
		h = grpchandler.Handler{
			RegisteredProviders: grpchandler.Providers,
		}
	default:
		log.Fatal("unrecognised protocol, choose 'nats', 'http' or 'grpc'")
	}

	http.HandleFunc("/", ping)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/enrol", h.Enrol)

	log.Println("starting on :8090 with protocol: " + protocol)
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

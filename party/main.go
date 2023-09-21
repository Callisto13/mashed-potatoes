package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/callisto13/mashed-potatoes/natsemitter"
)

const (
	DEFAULT_NATS_SERVER = "http://localhost:4222"
)

func main() {
	var (
		natsAddress string
	)

	flag.StringVar(&natsAddress, "nats-address", DEFAULT_NATS_SERVER, "address + port of running NATs service")

	http.HandleFunc("/", ping)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/enrol", enrol(natsAddress))

	log.Println("starting on :8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal("failed to get this party started")
	}
}

func enrol(natsAddress string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		natsSubject := req.URL.Query().Get("target")

		if natsSubject == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(errors.New("missing target param"))

			return
		}

		if err := natsemitter.SendEvent(natsAddress, natsSubject, "enrol"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)

			return
		}

		w.WriteHeader(http.StatusOK)

		message := map[string]string{
			"action": "enrol",
			"type":   natsSubject,
		}
		data, err := json.Marshal(message)
		if err != nil {
			log.Println("could not marshal response")
		}

		if _, err := w.Write([]byte(data)); err != nil {
			log.Println("could not write response")
		}

		log.Printf("enroling a %s", natsSubject)
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

package nats

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/callisto13/mashed-potatoes/party/natsemitter"
)

const DEFAULT_NATS_SERVER = "http://localhost:4222"

type Handler struct {
	NatsAddress string
}

func (h Handler) Enrol(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	natsSubject := req.URL.Query().Get("target")

	if natsSubject == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(errors.New("missing target param"))

		return
	}

	log.Printf("enroling a thing of type %s via nats protocol", natsSubject)

	if err := natsemitter.SendEvent(h.NatsAddress, natsSubject, "enrol"); err != nil {
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
}

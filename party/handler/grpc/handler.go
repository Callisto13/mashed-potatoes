package grpcsender

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/callisto13/mashed-potatoes/party/grpcsender"
)

type Handler struct {
	RegisteredProviders map[string]string
}

var Providers = map[string]string{"jokeprovider": "localhost:1430"}

func (h Handler) Enrol(w http.ResponseWriter, req *http.Request) {
	targetProvider := req.URL.Query().Get("target")

	if targetProvider == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(errors.New("missing target param"))

		return
	}

	target, ok := h.RegisteredProviders[targetProvider]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(fmt.Errorf("provider not registered: %s", targetProvider))

		return
	}

	log.Printf("enroling a thing of type %s via grpc protocol", targetProvider)

	if err := grpcsender.SendEvent(target, "enrol"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)

		return
	}

	w.WriteHeader(http.StatusOK)

	message := map[string]string{
		"action": "enrol",
		"type":   target,
	}
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("could not marshal response")
	}

	if _, err := w.Write([]byte(data)); err != nil {
		log.Println("could not write response")
	}
}

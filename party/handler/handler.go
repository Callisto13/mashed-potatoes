package handler

import "net/http"

type ProtocolHandler interface {
	Enrol(http.ResponseWriter, *http.Request)
}

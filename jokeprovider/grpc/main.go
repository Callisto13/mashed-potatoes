package main

import (
	"encoding/json"
	"log"
	"net"

	// _ "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	"github.com/callisto13/mashed-potatoes/jokeprovider/joke"
	"github.com/callisto13/mashed-potatoes/party/grpcsender"
	"github.com/callisto13/mashed-potatoes/proto"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedProviderServiceServer
}

func main() {
	lis, err := net.Listen("tcp", ":1430")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterProviderServiceServer(s, &server{})

	log.Printf("server starting at %v, waiting for events", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Event(req *proto.EventRequest, idk proto.ProviderService_EventServer) error {
	data := &grpcsender.EventData{}

	if err := json.Unmarshal(req.Event.GetBinaryData(), data); err != nil {
		return err
	}

	log.Printf("received action: %s", data.Action)

	switch data.Action {
	case "enrol":
		if err := joke.Enrol(); err != nil {
			return err
		}
	default:
		log.Printf("no response found for action: %s", data.Action)
	}

	log.Printf("action %s completed successfully", data.Action)

	return nil
}

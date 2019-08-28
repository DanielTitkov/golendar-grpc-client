package main

import (
	"context"
	"log"
	"time"

	pb "github.com/DanielTitkov/golendar/api/grpc/golendarpb"
	"google.golang.org/grpc"
)

func main() {
	log.Println("START")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	client := pb.NewEventServiceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()

	// create some new events
	uuids := []string{}
	for _, e := range []struct {
		t string
		u string
	}{
		{t: "Super protobuf event", u: "Dick Johnson"},
		{t: "Go hunt some cherry blossom", u: "QuuoQuuo"},
		{t: "Plow my own farrow", u: "Tobi bida"},
	} {
		createEventRes, _ := client.CreateEvent(ctx, &pb.CreateEventRequest{
			Event: &pb.Event{
				Title: e.t,
				User:  e.u,
			},
		})
		uuids = append(uuids, createEventRes.GetEvent().EventUUID)
	}

	// get events with known id
	log.Println("GET")
	getEventRes, _ := client.GetEvent(ctx, &pb.GetEventRequest{EventUUID: uuids})
	log.Println(getEventRes)

	// update one event
	log.Println("UPDATE")
	updateEventRes, _ := client.UpdateEvent(ctx, &pb.UpdateEventRequest{EventUUID: uuids[0], Event: &pb.Event{
		User:  "Dick Johnson the Great",
		Title: "Super protobuf event",
	}})
	log.Println(updateEventRes)

	// get events with known id
	log.Println("GET")
	getEventRes, _ = client.GetEvent(ctx, &pb.GetEventRequest{EventUUID: uuids})
	log.Println(getEventRes)

	// delete one event
	log.Println("DELETE")
	deleteEventRes, _ := client.DeleteEvent(ctx, &pb.DeleteEventRequest{EventUUID: uuids[0]})
	uuids = uuids[1:]
	log.Println(deleteEventRes)

	// get events with known id
	log.Println("GET")
	getEventRes, _ = client.GetEvent(ctx, &pb.GetEventRequest{EventUUID: uuids})
	log.Println(getEventRes)

	log.Println("EXIT")
}

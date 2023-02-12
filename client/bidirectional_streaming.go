package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/IRSHIT033/go-grpc/proto"
)

func callSayHelloBidirectionalStream(client pb.GreetServiceClient, names *pb.NameList) {
	log.Printf("Bidirectional streaming started")
	stream, err := client.SayHelloBidirectionalStreaming(context.Background())
	if err != nil {
		log.Fatalf("could not send names : %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while running %v", err)
			}
			log.Println(message)
		}
		close(waitc)
	}()

	for _, name := range names.Message {
		req := &pb.HelloRequest{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending %v", err)
		}
		time.Sleep(2 * time.Second)
	}
	stream.CloseSend()
	<-waitc
	log.Printf("Bidirectional streaming finished")
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"app.ir/internal/job"
	app "app.ir/internal/transport/grpc"
	"app.ir/pkg/logHandler"
	pb "app.ir/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 3000, "The server port")
)

func main() {

	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		logHandler.LogError(err.Error())
	}

	// db := db.NewDatabase()

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	pb.RegisterEstimateServer(s, &app.Server{})

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve %v", err)
		}
	}()
	logHandler.Log(fmt.Sprintf("server listening at %v", listener.Addr()))

	job.RunJobs()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("\nStopping the server... ")
	s.Stop()
	listener.Close()
	fmt.Println("Closing MongoDB conncetion")
	fmt.Println("Done")
}

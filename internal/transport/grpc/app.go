package grpc

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"app.ir/internal/handler"
	"app.ir/pkg/logHandler"
	pb "app.ir/proto"
	"google.golang.org/grpc"
)

func RunServer(listener net.Listener, segmentHandler handler.SegmentHandler) {
	// func RunServer(listener net.Listener, segmentHandler pb.EstimateServer) {
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	pb.RegisterEstimateServer(s, segmentHandler)

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve %v", err)
		}
	}()
	logHandler.Log(fmt.Sprintf("server listening at %v", listener.Addr()))

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("\nStopping the server... ")
	s.Stop()
	listener.Close()
	fmt.Println("Closing MongoDB conncetion")
	fmt.Println("Done")
}

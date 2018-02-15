package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	microclient "github.com/micro/go-micro/client"

	"golang.org/x/net/context"

	"github.com/micro/go-micro/cmd"
	pb "github.com/mooncaker816/shipper/consignment-service/proto/consignment"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment pb.Consignment
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &consignment)
	if err != nil {
		return nil, err
	}
	return &consignment, nil
}

func main() {
	// conn, err := grpc.Dial(address, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("dial failed %v", err)
	// }
	// defer conn.Close()
	cmd.Init()
	client := pb.NewShippingServiceClient("go.micro.srv.consignment", microclient.DefaultClient)
	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("parse failed %v", err)
	}
	res, err := client.CreateConsignment(context.TODO(), consignment)
	if err != nil {
		log.Fatalf("can not create consignment: %v", err)
	}
	log.Printf("Created: %v", res.Consignment)
	res, err = client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("can not get consignments: %v", err)
	}
	for _, v := range res.Consignments {
		log.Println(v)
	}
}

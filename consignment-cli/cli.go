package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"

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
	service := micro.NewService(
		micro.Flags(
			cli.StringFlag{
				Name:  "file",
				Usage: "consignment file name",
			},
			cli.StringFlag{
				Name:  "token",
				Usage: "token string",
			},
		),
	)

	cmd.Init()
	client := pb.NewShippingServiceClient("go.micro.srv.consignment", microclient.DefaultClient)
	var file, token string
	service.Init(
		micro.Action(func(c *cli.Context) {

			file = c.String("file")
			token = c.String("token")
			log.Printf("input token: %s", token)
			if len(file) <= 0 {
				file = defaultFilename
			}
			consignment, err := parseFile(file)
			if err != nil {
				log.Fatalf("parse failed %v", err)
			}
			ctx := metadata.NewContext(context.Background(), map[string]string{
				"token": token,
			})

			res, err := client.CreateConsignment(ctx, consignment)
			if err != nil {
				log.Fatalf("can not create consignment: %v", err)
			}
			log.Printf("Created: %v", res.Consignment)
			res, err = client.GetConsignments(ctx, &pb.GetRequest{})
			if err != nil {
				log.Fatalf("can not get consignments: %v", err)
			}
			for _, v := range res.Consignments {
				log.Println(v)
			}
			// let's just exit because
			os.Exit(0)
		}),
	)
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

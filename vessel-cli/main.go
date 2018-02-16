package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	pb "github.com/mooncaker816/shipper/vessel-service/proto/vessel"
)

func main() {
	// session, err := CreateSession(host)
	// if err != nil {

	// 	// We're wrapping the error returned from our CreateSession
	// 	// here to add some context to the error.
	// 	log.Panicf("Could not connect to datastore with host %s - %v", host, err)
	// }
	// // Mgo creates a 'master' session, we need to end that session
	// // before the main function closes.
	// defer session.Close()
	cmd.Init()
	file := "vessel.json"
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic("read file error!")
	}
	var vessel pb.Vessel
	err = json.Unmarshal(data, &vessel)
	if err != nil {
		panic("unmarshal failed")
	}
	client := pb.NewVesselServiceClient("go.micro.srv.vessel", microclient.DefaultClient)

	res, err := client.Create(context.Background(), &vessel)
	if err != nil {
		log.Fatalf("create vessel failed: %v", err)
	}
	log.Printf("created: %v", res.Vessel)

	res, err = client.GetAll(context.Background(), &pb.Specification{0, 0})
	if err != nil {
		log.Fatalf("get vessels failed: %v", err)
	}
	for _, v := range res.Vessels {
		log.Println(v)
	}
}

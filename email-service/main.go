package main

import (
	"context"
	"log"

	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/nats"
	pb "github.com/mooncaker816/shipper/user-service/proto/user"
)

const topic = "user.created"

type Subscriber struct{}

func (sub *Subscriber) Process(ctx context.Context, user *pb.User) error {
	log.Println("Picked up a new message")
	log.Println("Sending email to:", user.Name)
	return nil
}

func main() {
	srv := micro.NewService(
		micro.Name("go.micro.srv.email"),
		micro.Version("latest"),
	)

	srv.Init()

	// Get the broker instance using our environment variables
	// pubsub := srv.Server().Options().Broker
	// if err := pubsub.Connect(); err != nil {
	// 	log.Fatal(err)
	// }

	// // Subscribe to messages on the broker
	// _, err := pubsub.Subscribe(topic, func(p broker.Publication) error {
	// 	var user *pb.User
	// 	if err := json.Unmarshal(p.Message().Body, &user); err != nil {
	// 		return err
	// 	}
	// 	log.Println(user)
	// 	go sendEmail(user)
	// 	return nil
	// })

	// if err != nil {
	// 	log.Println(err)
	// }

	micro.RegisterSubscriber(topic, srv.Server(), new(Subscriber))

	// Run the server
	if err := srv.Run(); err != nil {
		log.Println(err)
	}
}

func sendEmail(user *pb.User) error {
	log.Println("Sending email to:", user.Name)
	return nil
}

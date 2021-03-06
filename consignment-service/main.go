package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	userService "github.com/mooncaker816/shipper/user-service/proto/user"

	vessel "github.com/mooncaker816/shipper/vessel-service/proto/vessel"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	pb "github.com/mooncaker816/shipper/consignment-service/proto/consignment"
)

const (
	defaultHost = "localhost:27017"
)

// type Repository interface {
// 	Create(*pb.Consignment) (*pb.Consignment, error)
// 	GetAll() ([]*pb.Consignment, error)
// }

// type ConsignmentRepository struct {
// 	consignments []*pb.Consignment
// }

// func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
// 	repo.consignments = append(repo.consignments, consignment)
// 	return consignment, nil
// }
// func (repo *ConsignmentRepository) GetAll() ([]*pb.Consignment, error) {
// 	return repo.consignments, nil
// }

// type service struct {
// 	repo         Repository
// 	vesselClient vessel.VesselServiceClient
// }

// func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
// 	// Here we call a client instance of our vessel service with our consignment weight,
// 	// and the amount of containers as the capacity value
// 	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vessel.Specification{
// 		MaxWeight: req.Weight,
// 		Capacity:  int32(len(req.Containers)),
// 	})
// 	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
// 	if err != nil {
// 		return err
// 	}

// 	// We set the VesselId as the vessel we got back from our
// 	// vessel service
// 	req.VesselId = vesselResponse.Vessel.Id
// 	consi, err := s.repo.Create(req)
// 	if err != nil {
// 		return err
// 	}
// 	res.Created = true
// 	res.Consignment = consi
// 	return nil
// }

// func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
// 	consignments, err := s.repo.GetAll()
// 	if err != nil {
// 		return err
// 	}
// 	res.Consignments = consignments
// 	return nil
// }

func main() {
	//repo := &ConsignmentRepository{}
	// Database host from the environment variables
	host := os.Getenv("DB_HOST")

	if host == "" {
		host = defaultHost
	}

	session, err := CreateSession(host)

	// Mgo creates a 'master' session, we need to end that session
	// before the main function closes.
	defer session.Close()

	if err != nil {

		// We're wrapping the error returned from our CreateSession
		// here to add some context to the error.
		log.Panicf("Could not connect to datastore with host %s - %v", host, err)
	}

	srv := micro.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
		micro.WrapHandler(AuthWrapper),
	)
	vesselClient := vessel.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	// Init will parse the command line flags.
	srv.Init()

	// Register handler
	pb.RegisterShippingServiceHandler(srv.Server(), &service{session, vesselClient})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
	// pb.RegisterShippingServiceServer(server, &service{repo}) //注册service到grpc服务器，使得pb.go中的代码与服务器绑定到一起
	// reflection.Register(server)                              //注册反射到grpc服务器
	// if err := server.Serve(lis); err != nil {                //启动服务
	// 	log.Fatalf("faild to serve %v", err)
	// }
}

// AuthWrapper is a high-order function which takes a HandlerFunc
// and returns a function, which takes a context, request and response interface.
// The token is extracted from the context set in our consignment-cli, that
// token is then sent over to the user service to be validated.
// If valid, the call is passed along to the handler. If not,
// an error is returned.
func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		if os.Getenv("DISABLE_AUTH") == "true" {
			return fn(ctx, req, resp)
		}
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}
		log.Println(meta)
		// Note this is now uppercase (not entirely sure why this is...)
		token := strings.TrimSpace(meta["Token"])
		log.Println("Authenticating with token:", token)

		// Auth here
		authClient := userService.NewUserServiceClient("go.micro.srv.user", client.DefaultClient)
		_, err := authClient.ValidateToken(context.Background(), &userService.Token{
			Token: token,
		})
		if err != nil {
			return err
		}
		err = fn(ctx, req, resp)
		return err
	}
}

package main

import (
	"context"
	"fmt"
	"log"

	vessel "github.com/mooncaker816/shipper/vessel-service/proto/vessel"

	micro "github.com/micro/go-micro"
	pb "github.com/mooncaker816/shipper/consignment-service/proto/consignment"
)

type Repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() ([]*pb.Consignment, error)
}

type ConsignmentRepository struct {
	consignments []*pb.Consignment
}

func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}
func (repo *ConsignmentRepository) GetAll() ([]*pb.Consignment, error) {
	return repo.consignments, nil
}

type service struct {
	repo         Repository
	vesselClient vessel.VesselServiceClient
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	// Here we call a client instance of our vessel service with our consignment weight,
	// and the amount of containers as the capacity value
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vessel.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}

	// We set the VesselId as the vessel we got back from our
	// vessel service
	req.VesselId = vesselResponse.Vessel.Id
	consi, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	res.Created = true
	res.Consignment = consi
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &ConsignmentRepository{}
	// lis, err := net.Listen("tcp", port) //监听端口
	// if err != nil {
	// 	log.Fatalf("failed to listen port %v", port)
	// }
	// server := grpc.NewServer()                               // grpc server
	srv := micro.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)
	vesselClient := vessel.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	// Init will parse the command line flags.
	srv.Init()

	// Register handler
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})

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
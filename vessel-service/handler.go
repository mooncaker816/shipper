package main

import (
	"context"

	"gopkg.in/mgo.v2"

	pb "github.com/mooncaker816/shipper/vessel-service/proto/vessel"
)

type service struct {
	//repo Repository
	session *mgo.Session
}

func (s *service) GetRepo() Repository {
	return &VesselRepository{s.session.Clone()}
}

func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	repo := s.GetRepo()
	defer repo.Close()
	// Find the next available vessel
	vessel, err := repo.FindAvailable(req)
	if err != nil {
		return err
	}

	// Set the vessel as part of the response message type
	res.Vessel = vessel
	return nil
}

func (s *service) Create(ctx context.Context, req *pb.Vessel, res *pb.Response) error {
	repo := s.GetRepo()
	defer repo.Close()
	err := repo.Create(req)
	if err != nil {
		return err
	}
	res.Vessel = req
	res.Created = true
	return nil
}

func (s *service) GetAll(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	repo := s.GetRepo()
	defer repo.Close()
	vessels, err := repo.GetAll(req)
	if err != nil {
		return err
	}
	res.Vessels = vessels
	return nil
}

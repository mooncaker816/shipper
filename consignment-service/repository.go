package main

import (
	pb "github.com/mooncaker816/shipper/consignment-service/proto/consignment"
	"gopkg.in/mgo.v2"
)

const (
	dbName                = "shipper"
	consignmentCollection = "consignments"
)

type Repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() ([]*pb.Consignment, error)
	Close()
}

type ConsignmentRepository struct {
	//consignments []*pb.Consignment
	session *mgo.Session
}

func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	err := repo.session.DB(dbName).C(consignmentCollection).Insert(consignment)
	//repo.consignments = append(repo.consignments, consignment)
	if err != nil {
		return nil, err
	}
	return consignment, nil
}

func (repo *ConsignmentRepository) GetAll() (consignments []*pb.Consignment, err error) {
	err = repo.session.DB(dbName).C(consignmentCollection).Find(nil).All(&consignments)
	if err != nil {
		return nil, err
	}
	return consignments, nil
	//return repo.consignments, nil
}

func (repo *ConsignmentRepository) Close() {
	repo.session.Close()
}

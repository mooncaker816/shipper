package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	pb "github.com/mooncaker816/shipper/vessel-service/proto/vessel"
)

const (
	dbName           = "shipper"
	vesselCollection = "vessels"
)

type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
	Close()
	Create(*pb.Vessel) error
	GetAll(*pb.Specification) ([]*pb.Vessel, error)
}

type VesselRepository struct {
	//vessels []*pb.Vessel
	session *mgo.Session
}

// FindAvailable - checks a specification against a map of vessels,
// if capacity and max weight are below a vessels capacity and max weight,
// then return that vessel.
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (vessel *pb.Vessel, err error) {
	err = repo.session.DB(dbName).C(vesselCollection).Find(bson.M{
		"capacity":  bson.M{"$gte": spec.Capacity},
		"maxweight": bson.M{"$gte": spec.MaxWeight},
	}).One(&vessel)
	// for _, vessel := range repo.vessels {
	// 	if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
	// 		return vessel, nil
	// 	}
	// }
	if err != nil {
		return nil, err
	}
	return vessel, nil
}

func (repo *VesselRepository) Close() {
	repo.session.Close()
}

func (repo *VesselRepository) Create(vessel *pb.Vessel) error {
	return repo.session.DB(dbName).C(vesselCollection).Insert(vessel)
}

func (repo *VesselRepository) GetAll(spec *pb.Specification) ([]*pb.Vessel, error) {
	var vessels []*pb.Vessel
	err := repo.session.DB(dbName).C(vesselCollection).Find(bson.M{
		"capacity":  bson.M{"$gte": spec.Capacity},
		"maxweight": bson.M{"$gte": spec.MaxWeight},
	}).All(&vessels)
	if err != nil {
		return nil, err
	}
	return vessels, nil
}

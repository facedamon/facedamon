package main

import (
	"context"
	"errors"
	pb "facedamon/shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
	"log"
)

type Repository interface {
	//查找可用的货轮
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

//货轮集合
type VesselRepository struct {
	vessels []*pb.Vessel
}

//接口实现
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	//选择最近一条容量、载重都符合的货轮
	for _, v := range repo.vessels {
		if v.Capacity >= spec.Capacity && v.MaxWeight >= spec.MaxWeight {
			return v, nil
		}
	}
	return nil, errors.New("No vessel can`t be use")
}

//定义货轮服务
type service struct {
	repo Repository
}

//实现服务端
func (s *service) FindAvailable(ctx context.Context, spec *pb.Specification, resp *pb.Response) error {
	//调用北部方法查找货轮
	v, err := s.repo.FindAvailable(spec)
	if err != nil {
		return err
	}
	resp.Vessel = v
	return nil
}

func main() {
	//停留在港口的货轮，此处先写死，之后用数据库替代
	vessels := []*pb.Vessel{
		{
			Id:        "vessel001",
			Name:      "Boaty McBoatface",
			MaxWeight: 200000,
			Capacity:  500,
		},
	}
	repo := &VesselRepository{vessels}
	server := micro.NewService(micro.Name("go.micro.srv.vessel"), micro.Version("latest"))
	server.Init()

	//将实现服务端的API注册到服务端
	pb.RegisterVesselServiceHandler(server.Server(), &service{repo})

	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

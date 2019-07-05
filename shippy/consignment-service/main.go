package main

import (
	"context"
	pb "facedamon/shippy/consignment-service/proto/consignment"
	vesselPb "facedamon/shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
	"log"
)

const PORT = ":50051"

//仓库接口
type IRepository interface {
	//存放新货物
	Create(consignment *pb.Consignment) (*pb.Consignment, error)
	//获取仓库中所有的货物
	GetAll() []*pb.Consignment
}

//我们存放多批获取的仓库，实现了IRepository接口
type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

//定义微服务
type service struct {
	repo Repository

	vesselClient vesselPb.VesselService
}

//service实现consignment.pb.go中的ShippingServiceServer接口
//使service作为grpc的服务端 托运新的货物
// handler (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response,error) {
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {
	//检查是否有合适的货轮
	vReq := &vesselPb.Specification{
		Capacity:  int32(len(req.Containers)),
		MaxWeight: req.Weight,
	}
	vResp, err := s.vesselClient.FindAvailable(context.Background(), vReq)
	if err != nil {
		return err
	}

	log.Printf("found vessel: %s\n", vResp.Vessel.Name)
	req.VesselId = vResp.Vessel.Id

	//接收承运的货物
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	//resp = &pb.Response{Created: true, Consignment: consignment}
	//这里地方困扰了老子一万年
	//结构体赋值必须逐一赋值
	resp.Created = true
	resp.Consignment = consignment
	return nil
}

//handler (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response,error) {
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
	allConsignments := s.repo.GetAll()
	//resp = &pb.Response{Consignments: allConsignments}
	resp.Consignments = allConsignments
	return nil
}

func main() {
	//micro 版本
	ser := micro.NewService(micro.Name("go_micro_srv_consignment"), micro.Version("latest"))
	ser.Init()
	repo := Repository{}

	//ser 作为vessel-service客户端注册
	vClient := vesselPb.NewVesselService("go.micro.srv.vessel", ser.Client())

	pb.RegisterShippingServiceHandler(ser.Server(), &service{repo: repo, vesselClient: vClient})

	if err := ser.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
	"context"
	"encoding/json"
	"errors"
	pb "facedamon/shippy/consignment-service/proto/consignment"
	"github.com/micro/go-micro"
	"io/ioutil"
	"log"
	"os"
)

const (
	ADDRESS           = "localhost:50051"
	DEFAULT_INFO_FILE = "consignment.json"
)

// 读取consignment.json 中记录的货物信息
func parseFile(fileName string) (*pb.Consignment, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var consignment *pb.Consignment
	err = json.Unmarshal(data, &consignment)
	if err != nil {
		return nil, errors.New("consignment.json file content error")
	}
	return consignment, nil
}

func main() {
	service := micro.NewService(micro.Name("go_micro_srv_consignment.client"))
	service.Init()

	ship := pb.NewShippingService("go_micro_srv_consignment", service.Client())

	//在命令行中指定新的货物信息
	infoFile := DEFAULT_INFO_FILE
	if len(os.Args) > 1 {
		infoFile = os.Args[1]
	}

	//解析货物信息
	consignment, err := parseFile(infoFile)
	if err != nil {
		log.Fatalf("parse info file error: %v", err)
	}

	//调用rpc
	//将货物存储到我们自己的仓库里
	resp, err := ship.CreateConsignment(context.Background(), consignment)

	if err != nil {
		log.Fatalf("create consignment error: %v", err)
	}

	log.Printf("created: %t", resp.Created)

	//列出目前所有托运的货物
	resp, err = ship.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("failed to list consignments: %v", err)
	}
	for _, c := range resp.Consignments {
		log.Printf("%+v", c)
	}
}

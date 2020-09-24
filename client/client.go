package client

import (
	"context"
	"fmt"
	myprotos "grpc-client/day01/goprotos"
	"io"
	"os"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	//"google.golang.org/genproto/googleapis/iam/credentials/v1"
)

func StartClinet() {
	creds, err := credentials.NewClientTLSFromFile("./cert/server.crt", "")
	if err != nil {
		panic(err.Error())
	}
	options := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial("localhost:8080", options...)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	client := myprotos.NewStudentClient(conn)
	//getRepneMess(client)
	//
	sendPhoto(client)
}

func sendPhoto(client myprotos.StudentClient) {
	file, err := os.Open("/Users/a123456/lc/Readme.md")
	if err != nil {
		fmt.Println("open file failed ,err: ", err)
		return
	}
	defer file.Close()
	md := metadata.New(map[string]string{"id": "3"})
	context := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := client.SendPhoto(context)
	if err != nil {
		panic(err.Error())
	}
	for {
		bytes := make([]byte, 1024)
		n, err := file.Read(bytes)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err.Error())
		}
		stream.Send(&myprotos.PhotoMessage{Data: bytes[:n]})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res.StatusCode)
}

func getOneMess(client myprotos.StudentClient) {
	response, err := client.GetMess(context.Background(), &myprotos.Person{Id: 4, Name: "lc"})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v", response)
}

func getRepneMess(client myprotos.StudentClient) {
	stream, err := client.GetAll(context.Background(), &myprotos.Person{})
	if err != nil {
		panic(err.Error())
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("response over")
			break
		} else if err != nil {
			panic(err.Error())
		}
		fmt.Printf("%+v\n", res)
	}
}

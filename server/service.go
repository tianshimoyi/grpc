package server

import (
	"io"
	"os"
	"time"

	"google.golang.org/grpc/metadata"

	//"grpc/goprotos"
	"context"
	"fmt"
	myprotos "grpc/day01/goprotos"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type StuServer struct {
}

var (
	personMessages = []myprotos.PMes{{
		Class:       "1706",
		TeacherName: "正宏",
	},
		{
			Class:       "1705",
			TeacherName: "田继红",
		},
		{
			Class:       "1707",
			TeacherName: "邹国红",
		},
	}
	commMes = []myprotos.ResponseComm{{
		Price: 38.3,
		Stock: 500,
	},
		{
			Price: 40.3,
			Stock: 200,
		},
	}
)

//一元
func (s *StuServer) GetMess(ctx context.Context, p *myprotos.Person) (*myprotos.PMes, error) {
	if p.Name == "lc" && p.Id == 3 {
		return &myprotos.PMes{
			Class:       "1707",
			TeacherName: "邹国红",
		}, nil
	} else {
		err := fmt.Errorf("%s", "People Undfined")
		return nil, err
	}
}

//stream Response

func (s *StuServer) GetAll(p *myprotos.Person, stream myprotos.Student_GetAllServer) error {
	for _, v := range personMessages {
		stream.Send(&v)
		time.Sleep(time.Second * 2)
	}
	return nil
}

//stream Request

func (s *StuServer) SendPhoto(stream myprotos.Student_SendPhotoServer) error {
	file, err := os.OpenFile("Readme.md", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		panic(err.Error())
	}
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		fmt.Println("ID: ", md["id"][0])
	}
	img := []byte{}
	index := 1
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("image size: ", len(img))
			file.Write(img)
			return stream.SendAndClose(&myprotos.PhotoResponse{StatusCode: 200})
		} else if err != nil {
			return nil
		}
		fmt.Println("第", index, "传输的字节大小为：", len(data.Data))
		img = append(img, data.Data...)
		time.Sleep(time.Second * 1)
		index++
	}
	//	fmt.Println(len(img))
	//	file.Write(img)
	return nil
}

//stream both

func (s *StuServer) SearchComm(stream myprotos.Student_SearchCommServer) error {

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("客户端发送数据完毕")
			break
		} else if err != nil {
			panic(err.Error())
		}
		if data.Id == 1 {
			fmt.Println("服务端接收到数据1")
			stream.Send(&commMes[0])
		} else if data.Id == 2 {
			fmt.Println("服务端接受到数据2")
			stream.Send(&commMes[1])
		}
	}
	return nil
}

func StartServer() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err.Error())
	}
	defer listen.Close()
	creds, err := credentials.NewServerTLSFromFile("./cert/server.crt", "./cert/server.key")
	if err != nil {
		panic(err.Error())
	}
	options := []grpc.ServerOption{grpc.Creds(creds)}
	server := grpc.NewServer(options...)
	myprotos.RegisterStudentServer(server, new(StuServer))
	fmt.Println("server started!")
	server.Serve(listen)
}

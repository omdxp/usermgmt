package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"

	pb "github.com/Omar-Belghaouti/usermgmt/usermgmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	port = ":50051"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func (server *UserManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, server)
	log.Printf("Server listening at %v", lis.Addr())
	return s.Serve(lis)
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	readBytes, err := ioutil.ReadFile("users.json")
	var users *pb.Users = &pb.Users{}
	var user_id int32 = int32(rand.Intn(1000))
	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}

	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("File not found. Creating new file.")
			users.Users = append(users.Users, created_user)
			jsonBytes, err := protojson.Marshal(users)
			if err != nil {
				log.Fatalf("Failed to marshal users: %v", err)
			}
			err = ioutil.WriteFile("users.json", jsonBytes, 0644)
			if err != nil {
				log.Fatalf("Failed to write users: %v", err)
			}
			return created_user, nil

		} else {
			log.Fatalf("Failed to read users: %v", err)
		}
	}

	err = protojson.Unmarshal(readBytes, users)
	if err != nil {
		log.Fatalf("Failed to unmarshal users: %v", err)
	}
	users.Users = append(users.Users, created_user)
	jsonBytes, err := protojson.Marshal(users)
	if err != nil {
		log.Fatalf("Failed to marshal users: %v", err)
	}
	err = ioutil.WriteFile("users.json", jsonBytes, 0644)
	if err != nil {
		log.Fatalf("Failed to write users: %v", err)
	}
	return created_user, nil

}

func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.Empty) (*pb.Users, error) {
	readBytes, err := ioutil.ReadFile("users.json")
	var users *pb.Users = &pb.Users{}
	if err != nil {
		log.Fatalf("Failed to read users: %v", err)
	}

	err = protojson.Unmarshal(readBytes, users)
	if err != nil {
		log.Fatalf("Failed to unmarshal users: %v", err)
	}
	return users, nil
}

func main() {
	var usermgmt_server *UserManagementServer = NewUserManagementServer()
	if err := usermgmt_server.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

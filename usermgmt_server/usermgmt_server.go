package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "github.com/Omar-Belghaouti/usermgmt/usermgmt"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
		users: &pb.Users{},
	}
}

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
	users *pb.Users
}

func (server *UserManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, &UserManagementServer{})
	log.Printf("Server listening at %v", lis.Addr())
	return s.Serve(lis)
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	var user_id int32 = int32(rand.Intn(1000))
	createed_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}
	s.users.Users = append(s.users.Users, createed_user)
	return createed_user, nil
}

func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.Empty) (*pb.Users, error) {
	return s.users, nil
}

func main() {
	var usermgmt_server *UserManagementServer = NewUserManagementServer()
	if err := usermgmt_server.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

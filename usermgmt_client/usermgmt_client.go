package main

import (
	"context"
	"log"
	"time"

	pb "github.com/Omar-Belghaouti/usermgmt/usermgmt"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var new_users = make(map[string]int32)
	new_users["Omar"] = 23
	new_users["Yasser"] = 8

	for name, age := range new_users {
		r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("User created: %v, %v, %v", r.GetName(), r.GetAge(), r.GetId())
	}

	r, err := c.GetUsers(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Users: %v", r.GetUsers())

}

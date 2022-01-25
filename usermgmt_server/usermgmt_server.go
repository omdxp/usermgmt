package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "github.com/Omar-Belghaouti/usermgmt/usermgmt"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

type UserManagementServer struct {
	conn *pgx.Conn
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
	createSql := `create table if not exists users (
		id serial primary key,
		name text,
		age integer
	);
	`
	_, err := s.conn.Exec(context.Background(), createSql)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
		os.Exit(1)
	}
	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge()}
	tx, err := s.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin failed: %v", err)
	}
	_, err = tx.Exec(context.Background(), "insert into users(name, age) values ($1, $2)", created_user.Name, created_user.Age)
	if err != nil {
		log.Fatalf("tx.Exec failed: %v", err)
	}
	tx.Commit(context.Background())

	return created_user, nil

}

func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.Empty) (*pb.Users, error) {

	rows, err := s.conn.Query(context.Background(), "select * from users")
	if err != nil {
		log.Fatalf("conn.Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()
	var users []*pb.User
	for rows.Next() {
		var id int
		var name string
		var age int
		err := rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatalf("rows.Scan failed: %v", err)
			return nil, err
		}
		user := &pb.User{Id: int32(id), Name: name, Age: int32(age)}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("rows.Err failed: %v", err)
	}
	return &pb.Users{Users: users}, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database_url := os.Getenv("DATABASE_URL")
	if database_url == "" {
		log.Fatalf("DATABASE_URL is not set")
	}
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())
	server := NewUserManagementServer()
	server.conn = conn
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

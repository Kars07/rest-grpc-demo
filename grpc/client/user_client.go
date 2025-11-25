// grpc/client/user_client.go
package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/Kars07/rest-grpc-demo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

func NewUserClient(address string) (*UserClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	client := pb.NewUserServiceClient(conn)
	return &UserClient{client: client, conn: conn}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

func (c *UserClient) GetUser(id int64) (*pb.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.GetUserRequest{Id: id}
	return c.client.GetUser(ctx, req)
}

func (c *UserClient) GetAllUsers() (*pb.UserListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.Empty{}
	return c.client.GetAllUsers(ctx, req)
}

func (c *UserClient) CreateUser(name, email, phone string) (*pb.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.CreateUserRequest{
		Name:  name,
		Email: email,
		Phone: phone,
	}
	return c.client.CreateUser(ctx, req)
}

func (c *UserClient) UpdateUser(id int64, name, email, phone string) (*pb.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.UpdateUserRequest{
		Id:    id,
		Name:  name,
		Email: email,
		Phone: phone,
	}
	return c.client.UpdateUser(ctx, req)
}

func (c *UserClient) DeleteUser(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.DeleteUserRequest{Id: id}
	_, err := c.client.DeleteUser(ctx, req)
	return err
}

func (c *UserClient) StreamUsers() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.Empty{}
	stream, err := c.client.StreamUsers(ctx, req)
	if err != nil {
		return err
	}

	for {
		user, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("Received user: %s (ID: %d)", user.Name, user.Id)
	}

	return nil
}

// Example usage
func ExampleUsage() {
	// Connect to the gRPC server
	client, err := NewUserClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	// Create a user
	newUser, err := client.CreateUser("John Doe", "john@example.com", "123-456-7890")
	if err != nil {
		log.Fatalf("CreateUser failed: %v", err)
	}
	log.Printf("Created user: %+v", newUser)

	// Get user by ID
	user, err := client.GetUser(newUser.Id)
	if err != nil {
		log.Fatalf("GetUser failed: %v", err)
	}
	log.Printf("Retrieved user: %+v", user)

	// Get all users
	allUsers, err := client.GetAllUsers()
	if err != nil {
		log.Fatalf("GetAllUsers failed: %v", err)
	}
	log.Printf("Total users: %d", len(allUsers.Users))

	// Update user
	updatedUser, err := client.UpdateUser(newUser.Id, "John Smith", "john.smith@example.com", "098-765-4321")
	if err != nil {
		log.Fatalf("UpdateUser failed: %v", err)
	}
	log.Printf("Updated user: %+v", updatedUser)

	// Stream users
	log.Println("Streaming users:")
	if err := client.StreamUsers(); err != nil {
		log.Fatalf("StreamUsers failed: %v", err)
	}

	// Delete user
	if err := client.DeleteUser(newUser.Id); err != nil {
		log.Fatalf("DeleteUser failed: %v", err)
	}
	log.Println("User deleted successfully")
}

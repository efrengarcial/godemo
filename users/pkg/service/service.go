package service

import (
	"context"
	"github.com/efrengarcial/godemo/notificator/pkg/grpc/pb"
	"google.golang.org/grpc"
	"log"
)

// UsersService describes the service.
type UsersService interface {
	Create(ctx context.Context, email string) error
}

type basicUsersService struct {
	notificatorServiceClient pb.NotificatorClient
}

func (b *basicUsersService) Create(ctx context.Context, email string) error {
	reply, err := b.notificatorServiceClient.SendEmail(context.Background(), &pb.SendEmailRequest{
		Email:   email,
		Content: "Hi! Thank you for registration...",
	})

	if reply != nil {
		log.Printf("Email ID: %s", reply.Id)
	}

	return err
}

// NewBasicUsersService returns a naive, stateless implementation of UsersService.
func NewBasicUsersService() UsersService {

	conn, err :=  grpc.Dial("notificator-service:8085", grpc.WithInsecure())
	if err != nil {
		log.Printf("unable to connect to notificator-service: %s", err.Error())
		return new(basicUsersService)
	}

	log.Printf("connected to notificator-service")

	return &basicUsersService{
		notificatorServiceClient: pb.NewNotificatorClient(conn),
	}
}

// New returns a UsersService with all of the expected middleware wired in.
func New(middleware []Middleware) UsersService {
	var svc UsersService = NewBasicUsersService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}

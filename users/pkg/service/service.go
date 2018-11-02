package service

import (
	"context"
	"github.com/efrengarcial/godemo/notificator/pkg/grpc/pb"
	sdetcd "github.com/go-kit/kit/sd/etcd"
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
	var etcdServer = "http://etcd:2379" //"http://localhost:23791"

	client, err := sdetcd.NewClient(context.Background(), []string{etcdServer}, sdetcd.ClientOptions{})
	if err != nil {
		log.Printf("unable to connect to etcd: %s", err.Error())
		return new(basicUsersService)
	}

	entries, err := client.GetEntries("/services/notificator/")
	if err != nil || len(entries) == 0 {
		log.Printf("unable to get prefix entries: %s", err)
		return new(basicUsersService)
	}

	conn, err :=  grpc.Dial(entries[0], grpc.WithInsecure())
	if err != nil {
		log.Printf("unable to connect to notificator: %s", err.Error())
		return new(basicUsersService)
	}

	log.Printf("connected to notificator")

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

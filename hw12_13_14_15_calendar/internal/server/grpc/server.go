package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/logger"
	grpcpb "github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/storage"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) (int, error)
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, eventID int) error
	ListEvent(ctx context.Context) ([]storage.Event, error)
	GetEvent(ctx context.Context, eventID int) (*storage.Event, error)
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type Server struct {
	grpcpb.UnimplementedCalendarServiceServer
	Host    string
	Port    string
	Server  *grpc.Server
	Storage Storage
	Logger  Logger
}

func NewServer(host, port string, storage Storage, logger Logger) *Server {
	return &Server{
		Host:    host,
		Port:    port,
		Storage: storage,
		Logger:  logger,
	}
}

func (s *Server) CreateEvent(ctx context.Context, req *grpcpb.CreateEventRequest) (*grpcpb.CreateEventResponse, error) {
	pbEvent := req.GetEvent()
	event := storage.Event{
		UserID:      int(pbEvent.GetUserId()),
		Title:       pbEvent.GetTitle(),
		Description: stringPtr(pbEvent.GetDescription()),
		StartAt:     pbEvent.GetStartAt().AsTime(),
		EndAt:       pbEvent.GetEndAt().AsTime(),
		NotifyAt:    timePtr(pbEvent.GetNotifyAt().AsTime()),
	}
	id, err := s.Storage.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	pbEvent.Id = int32(id)
	return &grpcpb.CreateEventResponse{Event: pbEvent}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *grpcpb.UpdateEventRequest) (*grpcpb.UpdateEventResponse, error) {
	pbEvent := req.GetEvent()
	event := storage.Event{
		ID:          int(pbEvent.GetId()),
		UserID:      int(pbEvent.GetUserId()),
		Title:       pbEvent.GetTitle(),
		Description: stringPtr(pbEvent.GetDescription()),
		StartAt:     pbEvent.GetStartAt().AsTime(),
		EndAt:       pbEvent.GetEndAt().AsTime(),
		NotifyAt:    timePtr(pbEvent.GetNotifyAt().AsTime()),
	}
	if err := s.Storage.UpdateEvent(ctx, event); err != nil {
		return nil, err
	}
	return &grpcpb.UpdateEventResponse{Event: pbEvent}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *grpcpb.DeleteEventRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.Storage.DeleteEvent(ctx, int(req.GetId()))
}

func (s *Server) ListEvent(ctx context.Context, _ *grpcpb.ListEventRequest) (*grpcpb.ListEventResponse, error) {
	var pbEvents []*grpcpb.Event
	events, err := s.Storage.ListEvent(ctx)
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		pbEvents = append(pbEvents, &grpcpb.Event{
			Id:          int32(event.ID),
			UserId:      int32(event.UserID),
			Title:       event.Title,
			Description: stringFromPtr(event.Description),
			StartAt:     timestamppb.New(event.StartAt),
			EndAt:       timestamppb.New(event.EndAt),
			NotifyAt:    timestampFromPtr(event.NotifyAt),
		})
	}
	return &grpcpb.ListEventResponse{Events: pbEvents}, nil
}

func (s *Server) GetEvent(ctx context.Context, req *grpcpb.GetEventRequest) (*grpcpb.GetEventResponse, error) {
	event, err := s.Storage.GetEvent(ctx, int(req.GetId()))
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, errors.New("event not found")
	}
	pbEvent := grpcpb.Event{
		Id:          int32(event.ID),
		UserId:      int32(event.UserID),
		Title:       event.Title,
		Description: stringFromPtr(event.Description),
		StartAt:     timestamppb.New(event.StartAt),
		EndAt:       timestamppb.New(event.EndAt),
		NotifyAt:    timestampFromPtr(event.NotifyAt),
	}
	return &grpcpb.GetEventResponse{Event: &pbEvent}, nil
}

func (s *Server) Start(ctx context.Context) error {
	var ch = make(chan error)

	s.Logger.Info("Starting grpc server on " + s.getAddr())

	go func() {
		lis, err := net.Listen("tcp", s.getAddr())
		if err != nil {
			ch <- err
		}
		opts := make([]grpc.ServerOption, 0)
		if logg, ok := s.Logger.(*logger.Logger); ok {
			opts = append(opts, grpc.ChainUnaryInterceptor(grpc_zap.UnaryServerInterceptor(logg.Zap)))
		}
		s.Server = grpc.NewServer(opts...)
		grpcpb.RegisterCalendarServiceServer(s.Server, s)
		ch <- s.Server.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-ch:
		return err
	}
}

func (s *Server) Stop(ctx context.Context) error {
	var ch = make(chan struct{})

	go func() {
		s.Server.GracefulStop()
		ch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return errors.New("timeout stopping grpc server")
	case <-ch:
		return nil
	}
}

func (s *Server) getAddr() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func stringFromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func timestampFromPtr(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return timestamppb.New(time.Time{})
	}
	return timestamppb.New(*t)
}

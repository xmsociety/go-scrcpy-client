package service

import (
	"context"
	"go-scrcpy-client/scrcpy"
	"log"
)

var (
	CancelMap = make(map[string]context.CancelFunc)
)

type Service struct {
	Client      *scrcpy.Client
	Config      string
	ErrReceiver chan error
}

func NewService(client *scrcpy.Client, ErrReceiver chan error, config string) *Service {
	return &Service{Client: client, Config: config, ErrReceiver: ErrReceiver}
}

func (s *Service) Start() {
	defer func() {
		log.Printf("[run] %v: function service.client.start quit！", s.Client.Device.Serial)
	}()
	go s.Client.Start()
	go s.handler()
}
func (s *Service) handler() {
	defer func() {
		log.Printf("[run] %v: goroutine handler quit！", s.Client.Device.Serial)
	}()
	for {
		select {
		case frame := <-s.Client.VideoSender:
			if !s.Client.Alive {
				return
			}
			// TODO AI
			log.Printf("[run] %v: get Frame -> %v", s.Client.Device.Serial, frame.Bounds())
			// send error when logic
			// s.ErrReceiver <- errors.New("fuck you")
		case err := <-s.Client.ErrReceiver:
			// receive frame error
			log.Printf("[run] %v: receive Err -> %v", s.Client.Device.Serial, err.Error())
			s.ErrReceiver <- err
			return
		}
	}
}

func (s *Service) Stop() {
	cancelFunc := CancelMap[s.Client.Device.Serial]
	cancelFunc()
}

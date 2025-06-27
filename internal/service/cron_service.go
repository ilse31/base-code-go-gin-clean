package service

import (
	"time"

	"github.com/robfig/cron/v3"
)

type CronService interface {
	Start()
	Stop()
	AddJob(spec string, cmd func()) (cron.EntryID, error)
}

type cronService struct {
	cron *cron.Cron
}

func NewCronService() CronService {
	return &cronService{
		cron: cron.New(cron.WithSeconds()),
	}
}

func (s *cronService) Start() {
	s.cron.Start()
}

func (s *cronService) Stop() {
	ctx := s.cron.Stop()
	// Wait for all running jobs to complete or timeout after 5 seconds
	timer := time.NewTimer(5 * time.Second)
	select {
	case <-timer.C:
		// Force stop after timeout
		return
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return
	}
}

func (s *cronService) AddJob(spec string, cmd func()) (cron.EntryID, error) {
	return s.cron.AddFunc(spec, cmd)
}

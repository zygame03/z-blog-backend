package infra

import "github.com/robfig/cron/v3"

type Job interface {
	RegisterJob(*cron.Cron)
}

func NewCron(jobs ...Job) *cron.Cron {
	c := cron.New()
	for _, job := range jobs {
		job.RegisterJob(c)
	}
	return c
}

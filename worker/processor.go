package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/yashagw/event-management-api/db"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendEmailVerify(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server   *asynq.Server
	provider db.Provider
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, provider db.Provider) TaskProcessor {
	server := asynq.NewServer(redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
		},
	)
	return &RedisTaskProcessor{
		server:   server,
		provider: provider,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, p.ProcessTaskSendEmailVerify)

	return p.server.Start(mux)
}

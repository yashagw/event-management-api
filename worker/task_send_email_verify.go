package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendEmailVerify struct {
	Email string `json:"email"`
}

func (d *RedisTaskDistributor) DistributeTaskSendEmailVerify(context context.Context, payload *PayloadSendEmailVerify, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	_, err = d.client.EnqueueContext(context, task)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %w", err)
	}
	// TODO: log task enqueued
	return nil
}

func (p *RedisTaskProcessor) ProcessTaskSendEmailVerify(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendEmailVerify
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("could not unmarshal payload: %w", err)
	}

	user, err := p.provider.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user with email %s not found", payload.Email)
		}
		return fmt.Errorf("could not get user by email: %w", err)
	}

	fmt.Println("sending email to", user.Email)
	// TODO: send email

	return nil
}

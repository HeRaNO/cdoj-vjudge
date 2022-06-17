package dal

import (
	"github.com/HeRaNO/cdoj-execution-worker/model"
	"github.com/HeRaNO/cdoj-vjudge/config"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

func SendMessage(req model.ExecRequest) (string, error) {
	ch := config.MQCh
	bd, err := sonic.Marshal(req)
	if err != nil {
		return "", err
	}
	corId := uuid.NewString()
	return corId, ch.Publish("", config.QueueName, false, false, amqp091.Publishing{
		ContentType:   "application/json",
		CorrelationId: corId,
		ReplyTo:       config.ReplyQueueName,
		Body:          bd,
	})
}

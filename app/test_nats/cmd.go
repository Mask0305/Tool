package main

import (
	"context"
	"encoding/json"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/tyr-tech-team/hawk/broker"
	"github.com/tyr-tech-team/hawk/broker/natsstreaming"
	"github.com/tyr-tech-team/hawk/status"
	"go.uber.org/zap"
)

type SmsInput struct {
	Phone string

	Content string
}

var (
	client broker.Broker
	log    func(context.Context) *zap.SugaredLogger
)

func Start() {

	req := &SmsInput{
		Phone:   "925331305",
		Content: "測試文字",
	}

	jq, _ := json.Marshal(req)

	Publish(context.Background(), "TIM-SEND-SMS", jq)

}

// New -
func NewNats() {
	natsStreaming := natsstreaming.New()

	client = broker.NewBroker(natsStreaming)
	log = func(ctx context.Context) *zap.SugaredLogger {
		return ctxzap.Extract(ctx).With(zap.String("entry", "natsStreaming")).Sugar()
	}
}

// Publish - 推送
// Topic - 推送主題
// Message - 訊息內容
func Publish(ctx context.Context, topic string, message []byte) error {

	natsStreaming := natsstreaming.New()

	client = broker.NewBroker(natsStreaming)
	log = func(ctx context.Context) *zap.SugaredLogger {
		return ctxzap.Extract(ctx).With(zap.String("entry", "natsStreaming")).Sugar()
	}
	if topic == "" {
		st := status.PermissionDenied.SetServiceCode(status.ServiceOfficeGateway).Err()
		log(ctx).Errorf("broker topic is empty: %s", st.Error())
		return st
	}

	// 發送訊息
	err := client.Publish(topic, &broker.Message{
		Body: message,
	})

	if err != nil {
		return err
	}

	return nil
}

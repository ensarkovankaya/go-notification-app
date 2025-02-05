package clients

import (
	"context"
	"github.com/ensarkovankaya/go-notification-app/clients/webhook/client"
	"github.com/ensarkovankaya/go-notification-app/clients/webhook/client/messages"
	_models "github.com/ensarkovankaya/go-notification-app/clients/webhook/models"
	openAPIClient "github.com/go-openapi/runtime/client"
	"go.uber.org/zap"
	"net/http"
)

type Payload = _models.MessageRequest

type WebhookClient struct {
	webhookId string
	API       *client.WebhookAPI
}

func NewWebhookClient(webhookId string) *WebhookClient {
	transport := openAPIClient.NewWithClient(client.DefaultHost, client.DefaultBasePath, []string{"https"}, &http.Client{})
	return &WebhookClient{
		webhookId: webhookId,
		API:       client.New(transport, nil),
	}
}

func (c *WebhookClient) Send(ctx context.Context, payload Payload) (string, error) {
	zap.L().Debug("Sending message", zap.Any("payload", payload))
	params := &messages.SendMessageParams{
		ID:      c.webhookId,
		Request: &payload,
		Context: ctx,
	}
	resp, err := c.API.Messages.SendMessage(params)
	if err != nil {
		zap.L().Error("Failed to send message", zap.Error(err))
		return "", err
	}
	zap.L().Debug("Message sent", zap.Any("response", resp))
	return resp.Payload.MessageID, nil
}

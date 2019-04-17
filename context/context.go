package context

import (
	"context"
	logrus "github.com/sirupsen/logrus"

	"github.com/chaostreff-flensburg/moc/models"
)

// GetLogger get context based logger
func GetLogger(ctx context.Context) *logrus.Entry {
	return ctx.Value("log").(*logrus.Entry)
}

// WithRequestID set a request id to context
func WithRequestID(ctx context.Context, id *string) context.Context {
	ctx = context.WithValue(ctx, "requestid", id)

	return ctx
}

// GetRequestID get context based request id
func GetRequestID(ctx context.Context) *string {
	id := ctx.Value("requestid")
	if id == nil {
		return nil
	}

	return id.(*string)
}

// WithMessage set a message to context
func WithMessage(ctx context.Context, message *models.Message) context.Context {
	ctx = context.WithValue(ctx, "message", message)

	return ctx
}

// GetMessage get context based request id
func GetMessage(ctx context.Context) *models.Message {
	message := ctx.Value("message")
	if message == nil {
		return nil
	}

	return message.(*models.Message)
}

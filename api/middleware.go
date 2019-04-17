package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	logrus "github.com/sirupsen/logrus"

	session "github.com/chaostreff-flensburg/moc/context"
	"github.com/chaostreff-flensburg/moc/models"
	"github.com/chaostreff-flensburg/moc/router"
)

// withMessageID load message entity by request param
func (api *API) withMessageID(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	messageID := chi.URLParam(r, "messageID")

	if _, err := uuid.Parse(messageID); err != nil {
		return nil, router.BadRequestError("bad messageID").WithInternalError(err)
	}

	var message models.Message
	if res := api.db.First(&message, models.Message{ID: messageID}); res.Error != nil {
		if gorm.IsRecordNotFoundError(res.Error) {
			return nil, router.NotFoundError("message not found")
		}

		return nil, router.HandleSQLError(res.Error)
	}

	ctx := r.Context()
	ctx = session.WithMessage(ctx, &message)

	return ctx, nil
}

// withLogger add request details to log output
func (api *API) withLogger(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	ctx := r.Context()

	logFields := logrus.Fields{}

	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)
	if reqID := session.GetRequestID(ctx); reqID != nil {
		logFields["req_id"] = *reqID
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	requestLogger := api.log.WithFields(logFields)

	ctx = context.WithValue(ctx, "log", requestLogger)

	return ctx, nil
}

// withRequestID add request id to log
func withRequestID(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	id := uuid.New().String()

	oldID := r.Header.Get("X-Request-ID")
	if oldID != "" {
		id = oldID
	}

	ctx := r.Context()
	ctx = session.WithRequestID(ctx, &id)

	w.Header().Set("X-Request-ID", id)

	return ctx, nil
}

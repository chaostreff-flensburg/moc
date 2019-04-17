package api

import (
	"encoding/json"
	"net/http"

	session "github.com/chaostreff-flensburg/moc/context"
	"github.com/chaostreff-flensburg/moc/models"
	"github.com/chaostreff-flensburg/moc/router"
)

// getMessages delivers all messages
func (api *API) getMessages(w http.ResponseWriter, r *http.Request) error {
	var message []*models.Message

	if res := api.db.Find(&message); res.Error != nil {
		return router.HandleSQLError(res.Error)
	}

	return router.SendJSON(w, http.StatusOK, message)
}

// createMessage
func (api *API) createMessage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	log := session.GetLogger(ctx)
	log.Info("request createMessage")

	message := &models.Message{}

	if err := json.NewDecoder(r.Body).Decode(&message.MessageRequest); err != nil {
		return router.BadRequestError("bad payload").WithInternalError(err)
	}

	if err := message.MessageRequest.Validate(); err != nil {
		return router.BadRequestError("bad payload").WithJsonError(err)
	}

	if res := api.db.Create(message); res.Error != nil {
		return router.HandleSQLError(res.Error)
	}

	return router.SendJSON(w, http.StatusOK, message)
}

// delivers message
func (api *API) getMessage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	log := session.GetLogger(ctx)
	log.Info("request getMessage")

	message := session.GetMessage(ctx)

	return router.SendJSON(w, http.StatusOK, message)
}

// delete a messages
func (api *API) deleteMessage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	log := session.GetLogger(ctx)
	log.Info("request deleteMessage")

	message := session.GetMessage(ctx)

	if res := api.db.Delete(message); res.Error != nil {
		return router.HandleSQLError(res.Error)
	}

	return router.SendJSON(w, http.StatusOK, message)
}

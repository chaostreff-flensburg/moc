package api

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/netlify/netlify-commons/graceful"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	"github.com/chaostreff-flensburg/moc/config"
	"github.com/chaostreff-flensburg/moc/router"
)

// API contains all the routes for the userstorage service
type API struct {
	handler http.Handler
	db      *gorm.DB
	router  *router.Router
	config  *config.Config
	log     *logrus.Entry
}

// NewAPI creates a new API object according to the configuration
func NewAPI(db *gorm.DB, config *config.Config) *API {
	r := router.NewRouter()
	log := logrus.WithField("component", "moc")

	api := &API{
		db:     db,
		router: r,
		config: config,
		log:    log,
	}

	r.Use(withRequestID)
	r.Use(router.Recoverer)
	r.Use(api.withLogger)

	log.Info("initialize Routes...")

	r.Route("/messages", func(r *router.Router) {
		r.Get("/", api.getMessages)
		r.Post("/", api.createMessage)

		r.Route("/{messageID}", func(r *router.Router) {
			r.Use(api.withMessageID)

			r.Get("/", api.getMessage)
			r.Delete("/", api.deleteMessage)
		})
	})

	corsHandler := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Total-Count"},
		AllowCredentials: true,
	})
	api.handler = corsHandler.Handler(r)

	return api
}

// ListenAndServe will finally start the real http server
func (api *API) ListenAndServe(addr string) {
	api.log.Info("Start App...")

	server := graceful.NewGracefulServer(api.handler, api.log)
	if err := server.Bind(addr); err != nil {
		api.log.WithError(err).Fatal("http server bind failed")
	}

	if err := server.Listen(); err != nil {
		api.log.WithError(err).Fatal("http server listen failed")
	}
}

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	logrus "github.com/sirupsen/logrus"

	"github.com/chaostreff-flensburg/moc/config"
	"github.com/chaostreff-flensburg/moc/models"
)

type APITest struct {
	DB      *gorm.DB
	Config  *config.Config
	BaseURL string
	T       *testing.T
}

// NewAPITest create a api test with tmp database and config
func NewAPITest(t *testing.T, baseURL string) *APITest {
	// read config
	config := config.ReadConfig()

	// prepare db
	if _, err := os.Stat("/tmp/test.db"); !os.IsNotExist(err) {
		os.Remove("/tmp/test.db")
	}

	db, _ := gorm.Open("sqlite3", "/tmp/test.db")

	// migrate
	db.AutoMigrate(&models.Message{})

	return &APITest{
		BaseURL: baseURL,
		T:       t,
		Config:  config,
		DB:      db,
	}
}

// TestEndpoint create a httptest recorder and a request to call a local endpoint
func (t *APITest) TestEndpoint(method string, url string, body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	r := httptest.NewRequest(method, fmt.Sprintf("%s%s", t.BaseURL, url), body)
	w := httptest.NewRecorder()

	// r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.Config.OperatorToken))

	logrus.SetLevel(logrus.ErrorLevel)
	api := NewAPI(t.DB, t.Config)
	api.handler.ServeHTTP(w, r)

	return w, r
}

// Request create a request with optional body
func (t *APITest) Request(method string, url string, data interface{}) *httptest.ResponseRecorder {
	var body io.Reader

	if data != nil {
		jsonBody, _ := json.Marshal(data)
		body = bytes.NewReader(jsonBody)
	}

	w, _ := t.TestEndpoint(method, url, body)

	return w
}

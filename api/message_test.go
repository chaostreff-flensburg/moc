package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/chaostreff-flensburg/moc/models"
	"github.com/chaostreff-flensburg/moc/router"
)

func TestGetMessages(t *testing.T) {
	name := "TestGetMessages"
	apiTest := NewAPITest(t, "http://localhost")

	// seed
	message := models.Seed(apiTest.DB)

	testCases := []struct {
		name    string
		method  string
		url     string
		code    int
		exepted []models.Message
	}{{
		name:   "with entry",
		method: "GET",
		url:    "/messages",
		code:   http.StatusOK,
		exepted: []models.Message{
			*message,
		},
	}}

	for i, testCase := range testCases {
		r := apiTest.Request(testCase.method, testCase.url, nil)

		var response []models.Message
		json.NewDecoder(r.Body).Decode(&response)

		options := []cmp.Option{
			cmpopts.IgnoreTypes(time.Time{}),
		}

		if diff := cmp.Diff(testCase.exepted, response, options...); diff != "" {
			t.Errorf("%s > %s #%d mismatch (-want +got):\n%s", name, testCase.name, i, diff)
		}

		assert.Equal(t, testCase.code, r.Code, fmt.Sprintf("%s > %s", name, testCase.name))
	}
}

func TestCreateMessage(t *testing.T) {
	name := "TestCreateMessage"
	apiTest := NewAPITest(t, "http://localhost")

	testCases := []struct {
		name    string
		method  string
		url     string
		code    int
		data    models.MessageRequest
		exepted interface{}
	}{{
		name:    "correct",
		method:  "POST",
		url:     "/messages",
		code:    http.StatusOK,
		data:    models.NewMessage("Testmessage").MessageRequest,
		exepted: *models.NewMessage("Testmessage"),
	}, {
		name:    "with empty data",
		method:  "POST",
		url:     "/messages",
		code:    http.StatusBadRequest,
		data:    models.MessageRequest{},
		exepted: *router.BadRequestError("bad payload").WithJsonError(map[string]interface{}{"message": "required"}),
	}, {
		name:    "with wrong data",
		method:  "POST",
		url:     "/messages",
		code:    http.StatusBadRequest,
		data:    models.NewMessage("ab").MessageRequest,
		exepted: *router.BadRequestError("bad payload").WithJsonError(map[string]interface{}{"message": "min"}),
	}}

	options := []cmp.Option{
		cmpopts.IgnoreTypes(time.Time{}),
		cmpopts.IgnoreFields(models.Message{}, "ID"),
	}

	for i, testCase := range testCases {
		r := apiTest.Request(testCase.method, testCase.url, testCase.data)

		assert.Equal(t, testCase.code, r.Code, fmt.Sprintf("%s > %s", name, testCase.name))

		var response models.Message
		var err router.HTTPError

		if r.Code != http.StatusOK {
			json.NewDecoder(r.Body).Decode(&err)

			if diff := cmp.Diff(testCase.exepted, err, options...); diff != "" {
				t.Errorf("%s > %s #%d mismatch (-want +got):\n%s", name, testCase.name, i, diff)
			}
		} else {
			json.NewDecoder(r.Body).Decode(&response)

			if diff := cmp.Diff(testCase.exepted, response, options...); diff != "" {
				t.Errorf("%s > %s #%d mismatch (-want +got):\n%s", name, testCase.name, i, diff)
			}
		}
	}
}

func TestGetMessage(t *testing.T) {
	name := "TestGetMessage"
	apiTest := NewAPITest(t, "http://localhost")

	// seed
	message := models.Seed(apiTest.DB)

	testCases := []struct {
		name    string
		method  string
		url     string
		code    int
		exepted interface{}
	}{{
		name:    "correct",
		method:  "GET",
		url:     fmt.Sprintf("/messages/%s", message.ID),
		code:    http.StatusOK,
		exepted: *models.NewMessage(message.Text),
	}, {
		name:    "with no uuid",
		method:  "GET",
		url:     fmt.Sprintf("/messages/%s", "teeest"),
		code:    http.StatusBadRequest,
		exepted: *router.BadRequestError("bad messageID"),
	}, {
		name:    "with not existing id",
		method:  "GET",
		url:     fmt.Sprintf("/messages/%s", uuid.New().String()),
		code:    http.StatusNotFound,
		exepted: *router.NotFoundError("message not found"),
	}}

	options := []cmp.Option{
		cmpopts.IgnoreTypes(time.Time{}),
		cmpopts.IgnoreFields(models.Message{}, "ID"),
	}

	for i, testCase := range testCases {
		r := apiTest.Request(testCase.method, testCase.url, nil)

		assert.Equal(t, testCase.code, r.Code, fmt.Sprintf("%s > %s", name, testCase.name))

		var response models.Message
		var err router.HTTPError

		if r.Code != http.StatusOK {
			json.NewDecoder(r.Body).Decode(&err)

			if diff := cmp.Diff(testCase.exepted, err, options...); diff != "" {
				t.Errorf("%s > %s #%d mismatch (-want +got):\n%s", name, testCase.name, i, diff)
			}
		} else {
			json.NewDecoder(r.Body).Decode(&response)

			if diff := cmp.Diff(testCase.exepted, response, options...); diff != "" {
				t.Errorf("%s > %s #%d mismatch (-want +got):\n%s", name, testCase.name, i, diff)
			}
		}
	}
}

func TestDeleteMessage(t *testing.T) {
	name := "TestDeleteMessage"
	apiTest := NewAPITest(t, "http://localhost")

	// seed
	message := models.Seed(apiTest.DB)

	testCases := []struct {
		name    string
		method  string
		url     string
		code    int
		exepted interface{}
	}{{
		name:    "correct",
		method:  "DELETE",
		url:     fmt.Sprintf("/messages/%s", message.ID),
		code:    http.StatusOK,
		exepted: *models.NewMessage(message.Text),
	}, {
		name:    "with no uuid",
		method:  "DELETE",
		url:     fmt.Sprintf("/messages/%s", "teeest"),
		code:    http.StatusBadRequest,
		exepted: *router.BadRequestError("bad messageID"),
	}, {
		name:    "with not existing id",
		method:  "DELETE",
		url:     fmt.Sprintf("/messages/%s", uuid.New().String()),
		code:    http.StatusNotFound,
		exepted: *router.NotFoundError("message not found"),
	}}

	options := []cmp.Option{
		cmpopts.IgnoreTypes(time.Time{}),
		cmpopts.IgnoreFields(models.Message{}, "ID"),
	}

	for i, testCase := range testCases {
		r := apiTest.Request(testCase.method, testCase.url, nil)

		assert.Equal(t, testCase.code, r.Code, fmt.Sprintf("%s > %s", name, testCase.name))

		var response models.Message
		var err router.HTTPError

		if r.Code != http.StatusOK {
			json.NewDecoder(r.Body).Decode(&err)

			if diff := cmp.Diff(testCase.exepted, err, options...); diff != "" {
				t.Errorf("%s > %s #%d mismatch (-want +got):\n%s", name, testCase.name, i, diff)
			}
		} else {
			json.NewDecoder(r.Body).Decode(&response)

			if diff := cmp.Diff(testCase.exepted, response, options...); diff != "" {
				t.Errorf("%s > %s #%d mismatch (-want +got):\n%s", name, testCase.name, i, diff)
			}
		}
	}
}

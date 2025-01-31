package misc_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkstas/tnr-backend/internal/testutils"
)

func TestNotFound(t *testing.T) {
	serv, cleanup, _ := testutils.NewTestApplication(t)
	defer cleanup()

	{
		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/", nil)

		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	}
	{
		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/asdf", nil)

		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	}
	{
		response := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/", nil)

		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	}
}

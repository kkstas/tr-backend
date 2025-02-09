package misc_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkstas/tnr-backend/internal/testutils"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()
	serv, _ := testutils.NewTestApplication(t)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/health-check", nil)
	serv.ServeHTTP(response, request)

	testutils.AssertStatus(t, response.Code, http.StatusOK)
}

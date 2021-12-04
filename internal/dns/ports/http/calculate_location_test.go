package http

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"space2218/internal/dns"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type locationCalculatorMock struct {
	Location dns.Location
	Err      error
}

func (m locationCalculatorMock) CalculateLocation(context.Context,
	dns.Coordinates3D, dns.Velocity,
) (dns.Location, error) {
	return m.Location, m.Err
}

func TestCalculateLocationHandlerOK(t *testing.T) {

	reqFd, err := os.Open(filepath.Join("testdata", "request_ok.json"))
	require.NoError(t, err)
	defer reqFd.Close()

	req, err := http.NewRequest(http.MethodPost, "", reqFd)
	if err != nil {
		t.Fatal(err)
	}

	const loc = 1389.57
	handler := http.HandlerFunc(
		CalculateLocationHandler(zap.NewNop(), locationCalculatorMock{loc, nil}),
	)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)

	expectedBody, err := ioutil.ReadFile(filepath.Join("testdata", "response_ok.json"))
	require.NoError(t, err)

	assert.JSONEq(t, string(expectedBody), rr.Body.String())
}

func TestCalculateLocationHandlerBadRequest(t *testing.T) {

	testFiles := []string{
		"request_bad_coordinate_x.json",
		"request_bad_coordinate_y.json",
		"request_bad_coordinate_z.json",
		"request_bad_velocity.json",
	}

	for _, filename := range testFiles {
		reqFd, err := os.Open(filepath.Join("testdata", filename))
		require.NoError(t, err)
		defer reqFd.Close()

		req, err := http.NewRequest(http.MethodPost, "", reqFd)
		if err != nil {
			t.Fatal(err)
		}

		handler := http.HandlerFunc(
			CalculateLocationHandler(zap.NewNop(), locationCalculatorMock{0, nil}),
		)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	}
}

func TestCalculateLocationHandlerErrorInvalidValue(t *testing.T) {

	reqFd, err := os.Open(filepath.Join("testdata", "request_ok.json"))
	require.NoError(t, err)
	defer reqFd.Close()

	req, err := http.NewRequest(http.MethodPost, "", reqFd)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(
		CalculateLocationHandler(zap.NewNop(), locationCalculatorMock{0, dns.ErrInvalidValue}),
	)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	assert.Contains(t, rr.Body.String(), dns.ErrInvalidValue.Error())
}

func TestCalculateLocationHandlerEmptyBody(t *testing.T) {

	req, err := http.NewRequest(http.MethodPost, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(
		CalculateLocationHandler(zap.NewNop(), locationCalculatorMock{0, dns.ErrInvalidValue}),
	)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Result().StatusCode)
}

func TestCalculateLocationHandlerRequestNotJson(t *testing.T) {

	reqFd, err := os.Open(filepath.Join("testdata", "request_not_json"))
	require.NoError(t, err)
	defer reqFd.Close()

	req, err := http.NewRequest(http.MethodPost, "", reqFd)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(
		CalculateLocationHandler(zap.NewNop(), locationCalculatorMock{0, dns.ErrInvalidValue}),
	)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Result().StatusCode)
}

func TestCalculateLocationHandlerUnexpectedError(t *testing.T) {

	reqFd, err := os.Open(filepath.Join("testdata", "request_ok.json"))
	require.NoError(t, err)
	defer reqFd.Close()

	req, err := http.NewRequest(http.MethodPost, "", reqFd)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(
		CalculateLocationHandler(zap.NewNop(), locationCalculatorMock{0, errors.New("armageddon")}),
	)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
}

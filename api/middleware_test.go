package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"simple-bank/token"
	"simple-bank/util"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,

) {
	tokenString, err := tokenMaker.CreateToken(username, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, tokenString)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)

}

func TestAuthMiddleware(t *testing.T) {

	tests := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, response *httptest.ResponseRecorder)
	}{
		{
			name: "No Authorization Added",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

			},
			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, response.Code)
			},
		},
		{
			name: "Invalid Authorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				username := util.RandomOwner()
				addAuthorization(t, request, tokenMaker, "invalid auth type", username, time.Minute)
			},
			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, response.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				username := util.RandomOwner()
				addAuthorization(t, request, tokenMaker, "", username, time.Minute)
			},
			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, response.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				username := util.RandomOwner()
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, -time.Minute)
			},
			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, response.Code)
			},
		},

		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				username := util.RandomOwner()
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, time.Minute)

			},
			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, response.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := newTestServer(t, nil)
			authPath := "/auth"

			server.router.GET(
				authPath,
				authMiddleware(server.token),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
				},
			)

			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			assert.NoError(t, err)

			test.setupAuth(t, request, server.token)
			server.router.ServeHTTP(recorder, request)
			test.checkResponse(t, recorder)
		})
	}

}

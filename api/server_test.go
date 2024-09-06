package api

import (
	"net/http"
	mockdb "simple-bank/db/mock"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestServer_SetRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	server := NewServer(store)
	testCases := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodPost,
			path:   "/account",
		},
		{
			method: http.MethodGet,
			path:   "/account/:id",
		},
		{
			method: http.MethodGet,
			path:   "/accounts",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.path, func(t *testing.T) {
			routeInfo := findRouteInfo(server.router, tt.method, tt.path)
			require.NotNil(t, routeInfo, "route should exist")

		})
	}
}

func findRouteInfo(router *gin.Engine, method string, path string) *gin.RouteInfo {
	for _, route := range router.Routes() {
		if route.Method == method && route.Path == path {
			return &route
		}
	}
	return nil
}

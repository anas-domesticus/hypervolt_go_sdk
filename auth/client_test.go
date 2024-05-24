package auth

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetToken(t *testing.T) {
	tests := []struct {
		name             string
		mockResponse     string
		mockStatusCode   int
		user             string
		password         string
		expectedToken    *Token
		expectedHasError bool
	}{
		{
			name:           "Successful Login",
			mockResponse:   `{"access_token":"test_token","expires_in":3600,"token_type":"Bearer"}`,
			mockStatusCode: http.StatusOK,
			user:           "testuser",
			password:       "testpass",
			expectedToken: &Token{
				AccessToken: "test_token",
				ExpiresIn:   3600,
				TokenType:   "Bearer",
			},
			expectedHasError: false,
		},
		{
			name:             "Unsuccessful Login",
			mockResponse:     `{"error": "invalid_grant"}`,
			mockStatusCode:   http.StatusBadRequest,
			user:             "baduser",
			password:         "badpass",
			expectedToken:    nil,
			expectedHasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authEndpointMockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.mockStatusCode)
				io.WriteString(w, test.mockResponse)
			}))
			defer authEndpointMockServer.Close()

			authEndpoint = authEndpointMockServer.URL

			resultToken, err := GetToken(test.user, test.password)

			if test.expectedHasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectedToken, resultToken)
			}
		})
	}
}

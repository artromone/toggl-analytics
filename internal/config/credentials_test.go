package config

import (
	// "errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"os"
	"testing"
)

type MockRequester struct {
	mock.Mock
}

func (m *MockRequester) MakeRequest(method, url string) (*http.Response, error) {
	args := m.Called(method, url)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestGetUserCredentials(t *testing.T) {
	os.Setenv("USER_API_KEY", "api_key_value")
	os.Setenv("USER_WORKSPACE_ID", "workspace_id_value")
	os.Setenv("USER_USER_NAME", "username_value")
	os.Setenv("USER_PAY_PER_HOUR", "10")
	os.Setenv("USER_CLIENT_PAY", "20")
	defer func() {
		os.Unsetenv("USER_API_KEY")
		os.Unsetenv("USER_WORKSPACE_ID")
		os.Unsetenv("USER_USER_NAME")
		os.Unsetenv("USER_PAY_PER_HOUR")
		os.Unsetenv("USER_CLIENT_PAY")
	}()

	credentials := GetUserCredentials("USER")

	assert.Equal(t, "api_key_value", credentials.APIKey)
	assert.Equal(t, "workspace_id_value", credentials.WorkspaceID)
	assert.Equal(t, "username_value", credentials.FileName)
	assert.Equal(t, "10", credentials.PayPerHour)
}

func TestGetAllUserCredentials(t *testing.T) {
	os.Setenv("USER1_API_KEY", "api_key_value_1")
	os.Setenv("USER1_WORKSPACE_ID", "workspace_id_value_1")
	os.Setenv("USER2_API_KEY", "api_key_value_2")
	os.Setenv("USER2_WORKSPACE_ID", "workspace_id_value_2")
	defer func() {
		os.Unsetenv("USER1_API_KEY")
		os.Unsetenv("USER1_WORKSPACE_ID")
		os.Unsetenv("USER2_API_KEY")
		os.Unsetenv("USER2_WORKSPACE_ID")
	}()

	credentials := GetAllUserCredentials()

	assert.Len(t, credentials, 2)

	assert.Equal(t, "api_key_value_1", credentials["USER1"].APIKey)
	assert.Equal(t, "workspace_id_value_1", credentials["USER1"].WorkspaceID)

	assert.Equal(t, "api_key_value_2", credentials["USER2"].APIKey)
	assert.Equal(t, "workspace_id_value_2", credentials["USER2"].WorkspaceID)
}

// func TestCheckCredentials_Success(t *testing.T) {
// 	mockRequester := new(MockRequester)
// 	mockRequester.On("MakeRequest", http.MethodGet, "https://api.track.toggl.com/api/v9/me").Return(&http.Response{
// 		StatusCode: http.StatusOK,
// 	}, nil)
//
// 	os.Setenv("CLIENT_PAY", "50")
// 	defer os.Unsetenv("CLIENT_PAY")
//
// 	err := CheckCredentials(mockRequester, "api_key")
// 	assert.NoError(t, err)
//
// 	mockRequester.AssertExpectations(t)
// }

// func TestCheckCredentials_RequestError(t *testing.T) {
// 	mockRequester := new(MockRequester)
// 	mockRequester.On("MakeRequest", http.MethodGet, "https://api.track.toggl.com/api/v9/me").Return(nil, errors.New("request error"))
//
// 	err := CheckCredentials(mockRequester, "api_key")
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "Error sending request")
//
// 	mockRequester.AssertExpectations(t)
// }

// func TestCheckCredentials_UnsuccessfulResponse(t *testing.T) {
// 	mockRequester := new(MockRequester)
// 	mockRequester.On("MakeRequest", http.MethodGet, "https://api.track.toggl.com/api/v9/me").Return(&http.Response{
// 		StatusCode: http.StatusForbidden,
// 	}, nil)
//
// 	err := CheckCredentials(mockRequester, "api_key")
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "Unexpected status code")
//
// 	mockRequester.AssertExpectations(t)
// }
//
// func TestCheckCredentials_MissingClientPay(t *testing.T) {
// 	mockRequester := new(MockRequester)
// 	mockRequester.On("MakeRequest", http.MethodGet, "https://api.track.toggl.com/api/v9/me").Return(&http.Response{
// 		StatusCode: http.StatusOK,
// 	}, nil)
//
// 	os.Unsetenv("CLIENT_PAY")
//
// 	err := CheckCredentials(mockRequester, "api_key")
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "CLIENT_PAY environment variable is not set")
//
// 	mockRequester.AssertExpectations(t)
// }

package xcrop

import (
	"context"
	"fmt"
)

// AccountService handles X account connection management.
type AccountService struct {
	http *httpClient
}

func newAccountService(h *httpClient) *AccountService {
	return &AccountService{http: h}
}

// ConnectResponse wraps the account connect result.
type ConnectResponse struct {
	Data WriteResult `json:"data"`
	Meta Meta        `json:"meta"`
}

// Connect connects an X account to the API user for write operations.
// Provide username, password, and optionally a TOTP secret for 2FA accounts.
func (s *AccountService) Connect(ctx context.Context, params *ConnectParams) (*WriteResult, error) {
	if params == nil {
		return nil, fmt.Errorf("xcrop: ConnectParams must not be nil")
	}
	if params.Username == "" {
		return nil, fmt.Errorf("xcrop: ConnectParams.Username must not be empty")
	}
	if params.Password == "" {
		return nil, fmt.Errorf("xcrop: ConnectParams.Password must not be empty")
	}
	var resp ConnectResponse
	err := s.http.do(ctx, requestOptions{
		method: "POST",
		path:   "/account/connect",
		body:   params,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// StatusResponse wraps the account status result.
type StatusResponse struct {
	Data AccountStatus `json:"data"`
	Meta Meta          `json:"meta"`
}

// Status checks the connection status of the X account.
func (s *AccountService) Status(ctx context.Context) (*AccountStatus, error) {
	var resp StatusResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/account/status",
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// DisconnectResponse wraps the account disconnect result.
type DisconnectResponse struct {
	Data WriteResult `json:"data"`
	Meta Meta        `json:"meta"`
}

// Disconnect disconnects the X account and removes stored credentials.
func (s *AccountService) Disconnect(ctx context.Context) (*WriteResult, error) {
	var resp DisconnectResponse
	err := s.http.do(ctx, requestOptions{
		method: "DELETE",
		path:   "/account/disconnect",
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

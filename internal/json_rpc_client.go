package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	// ErrJSONRPCClientRequestError is an error that is raised JSON-RPC request was failed.
	ErrJSONRPCClientRequestError    = errors.New("failed a request of JSON RPC to zabbix")
	errUnexpectedJSONRPCClientError = errors.New("unexpected JSON RPC client error")
)

const defaultJSONRPCVersion = "2.0"
const zabbixInternalItemType = 5 // more info: https://www.zabbix.com/documentation/current/manual/api/reference/item/object
const defaultTimeoutDuration = 3 * time.Second

const extendOutput = "extend"

const itemGetMethod = "item.get"
const userLoginMethod = "user.login"

// JSONRPCClient is a JSON-RPC client for Zabbix.
type JSONRPCClient struct {
	ZabbixBaseURL string
}

type itemGetRequestParam struct {
	Output                 string                     `json:"output"`
	Type                   uint8                      `json:"type"`
	Search                 *itemGetRequestSearchParam `json:"search"`
	SearchWildcardsEnabled bool                       `json:"searchWildcardsEnabled"`
}

type itemGetRequestSearchParam struct {
	Key string `json:"key_"`
}

type itemGetRequest struct {
	JSONRPCVersion string               `json:"jsonrpc"`
	Method         string               `json:"method"`
	AuthToken      string               `json:"auth"`
	ID             uint64               `json:"id"`
	Params         *itemGetRequestParam `json:"params"`
}

// ItemGetResponse is a struct that represents a envelope of responses for JSON-RPC `item.get` request.
type ItemGetResponse struct {
	Result []*ItemGetResponseResult `json:"result"`
}

// ItemGetResponseResult is a struct that represents a response for JSON-RPC `item.get` request.
type ItemGetResponseResult struct {
	Name      string      `json:"name"`
	Key       string      `json:"key_"`
	Status    string      `json:"status"`
	LastValue interface{} `json:"lastvalue"`
	PrevValue interface{} `json:"prevvalue"`
}

func makeItemGetRequest(authToken string, itemKey string) *itemGetRequest {
	return &itemGetRequest{
		JSONRPCVersion: defaultJSONRPCVersion,
		Method:         itemGetMethod,
		AuthToken:      authToken,
		ID:             1,
		Params: &itemGetRequestParam{
			Output:                 extendOutput,
			Type:                   zabbixInternalItemType,
			SearchWildcardsEnabled: true,
			Search: &itemGetRequestSearchParam{
				Key: itemKey,
			},
		},
	}
}

// GetItem retrieves item according to the given key through `item.get` JSON-RPC request.
func (c *JSONRPCClient) GetItem(authToken string, internalChecksKey string) (*ItemGetResponse, error) {
	reqBody, err := json.Marshal(makeItemGetRequest(authToken, internalChecksKey))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, errUnexpectedJSONRPCClientError)
	}

	resp, err := c.doRequest(reqBody)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrJSONRPCClientRequestError)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code = %d: %w", resp.StatusCode, ErrJSONRPCClientRequestError)
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, errUnexpectedJSONRPCClientError)
	}

	var itemGetResponse ItemGetResponse
	err = json.Unmarshal(bs, &itemGetResponse)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, errUnexpectedJSONRPCClientError)
	}

	return &itemGetResponse, nil
}

type userLoginRequest struct {
	JSONRPCVersion string                 `json:"jsonrpc"`
	Method         string                 `json:"method"`
	ID             uint64                 `json:"id"`
	Params         *userLoginRequestParam `json:"params"`
}

type userLoginRequestParam struct {
	UserName string `json:"user"`
	Password string `json:"password"`
}

// UserLoginResponse is a struct that represents a response of JSON-RPC `user.login`.
type UserLoginResponse struct {
	AuthToken string `json:"result"`
}

func makeUserLoginRequest(userName string, password string) *userLoginRequest {
	return &userLoginRequest{
		JSONRPCVersion: defaultJSONRPCVersion,
		Method:         userLoginMethod,
		ID:             1,
		Params: &userLoginRequestParam{
			UserName: userName,
			Password: password,
		},
	}
}

// UserLogin attempts authentication with given user name and password through the JSON-RPC `user.login` request.
func (c *JSONRPCClient) UserLogin(userName string, password string) (*UserLoginResponse, error) {
	reqBody, err := json.Marshal(makeUserLoginRequest(userName, password))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, errUnexpectedJSONRPCClientError)
	}

	resp, err := c.doRequest(reqBody)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrJSONRPCClientRequestError)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code = %d: %w", resp.StatusCode, ErrJSONRPCClientRequestError)
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, errUnexpectedJSONRPCClientError)
	}

	var userLoginResponse UserLoginResponse
	err = json.Unmarshal(bs, &userLoginResponse)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, errUnexpectedJSONRPCClientError)
	}

	return &userLoginResponse, nil
}

func (c *JSONRPCClient) doRequest(reqBody []byte) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/zabbix/api_jsonrpc.php", c.ZabbixBaseURL), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"Content-Type": {"application/json-rpc"},
	}

	client := &http.Client{
		Timeout: defaultTimeoutDuration,
	}
	return client.Do(req)
}

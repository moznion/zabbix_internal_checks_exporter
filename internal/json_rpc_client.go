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
	ErrUnexpectedJsonRPCClientError = errors.New("unexpected JSON RPC client error")
	ErrJsonRPCClientRequestError    = errors.New("failed a request of JSON RPC to zabbix")
)

const defaultJsonRPCVersion = "2.0"
const zabbixInternalItemType = 5 // more info: https://www.zabbix.com/documentation/current/manual/api/reference/item/object
const defaultTimeoutDuration = 3 * time.Second

const extendOutput = "extend"

const itemGetMethod = "item.get"
const userLoginMethod = "user.login"

type JsonRPCClient struct {
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
	JsonRPCVersion string               `json:"jsonrpc"`
	Method         string               `json:"method"`
	AuthToken      string               `json:"auth"`
	ID             uint64               `json:"id"`
	Params         *itemGetRequestParam `json:"params"`
}

type itemGetResponse struct {
	Result []*itemGetResponseResult `json:"result"`
}

type itemGetResponseResult struct {
	Name      string      `json:"name"`
	Key       string      `json:"key_"`
	Status    string      `json:"status"`
	LastValue interface{} `json:"lastvalue"`
	PrevValue interface{} `json:"prevvalue"`
}

func makeItemGetRequest(authToken string, itemKey string) *itemGetRequest {
	return &itemGetRequest{
		JsonRPCVersion: defaultJsonRPCVersion,
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

func (c *JsonRPCClient) GetItem(authToken string, internalChecksKey string) (*itemGetResponse, error) {
	reqBody, err := json.Marshal(makeItemGetRequest(authToken, internalChecksKey))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrUnexpectedJsonRPCClientError)
	}

	resp, err := c.doRequest(reqBody)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrJsonRPCClientRequestError)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code = %d: %w", resp.StatusCode, ErrJsonRPCClientRequestError)
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrUnexpectedJsonRPCClientError)
	}

	var itemGetResponse itemGetResponse
	err = json.Unmarshal(bs, &itemGetResponse)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrUnexpectedJsonRPCClientError)
	}

	return &itemGetResponse, nil
}

type userLoginRequest struct {
	JsonRPCVersion string                 `json:"jsonrpc"`
	Method         string                 `json:"method"`
	ID             uint64                 `json:"id"`
	Params         *userLoginRequestParam `json:"params"`
}

type userLoginRequestParam struct {
	UserName string `json:"user"`
	Password string `json:"password"`
}

type userLoginResponse struct {
	AuthToken string `json:"result"`
}

func makeUserLoginRequest(userName string, password string) *userLoginRequest {
	return &userLoginRequest{
		JsonRPCVersion: defaultJsonRPCVersion,
		Method:         userLoginMethod,
		ID:             1,
		Params: &userLoginRequestParam{
			UserName: userName,
			Password: password,
		},
	}
}

func (c *JsonRPCClient) UserLogin(userName string, password string) (*userLoginResponse, error) {
	reqBody, err := json.Marshal(makeUserLoginRequest(userName, password))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrUnexpectedJsonRPCClientError)
	}

	resp, err := c.doRequest(reqBody)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrJsonRPCClientRequestError)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code = %d: %w", resp.StatusCode, ErrJsonRPCClientRequestError)
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrUnexpectedJsonRPCClientError)
	}

	var userLoginResponse userLoginResponse
	err = json.Unmarshal(bs, &userLoginResponse)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrUnexpectedJsonRPCClientError)
	}

	return &userLoginResponse, nil
}

func (c *JsonRPCClient) doRequest(reqBody []byte) (*http.Response, error) {
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

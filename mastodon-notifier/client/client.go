package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (c *Client) DoRequest(ctx context.Context, r MastodonClientRequest) ([]byte, error) {
	_, err := url.Parse(c.Config.Server)
	if err != nil {
		return nil, fmt.Errorf("parse server URL on %s failed: %v", r.Prefix, err)
	}

	requestUrl := fmt.Sprintf("%s%s", c.Config.Server, r.Path)

	req, err := http.NewRequest(r.Method, requestUrl, r.RequestBody)
	if err != nil {
		return nil, fmt.Errorf("wrap request to %s failed: %v", r.Prefix, err)
	}

	req = req.WithContext(ctx)
	req.Header = r.Headers

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s failed: %v", r.Prefix, err)
	}
	defer resp.Body.Close()

	if responseBody, err := handleResponseBody(r.Prefix, resp); err != nil {
		return nil, err
	} else {
		return responseBody, nil
	}
}

func handleResponseBody(origin string, resp *http.Response) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s failed, HTTP response code: %s Host: %s, URL: %s, responsebody: %s",
			origin, resp.Status, resp.Request.URL.Hostname(), resp.Request.URL.Path, body)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s failed, HTTP response code: %s, reading response body failed: %w",
			origin, resp.Status, err)
	}
	return body, nil
}

func NewClient(config *Config) *Client {
	return &Client{
		Client: http.Client{Timeout: 10 * time.Second},
		Config: config,
	}
}

func (c *Client) Authenticate(ctx context.Context) error {
	params := url.Values{
		"client_id":     {c.Config.ClientID},
		"client_secret": {c.Config.ClientSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {"urn:ietf:wg:oauth:2.0:oob"},
		"scope":         {"read write"},
		"code":          {c.Config.ClientCode},
	}

	request := MastodonClientRequest{
		RequestBody: strings.NewReader(params.Encode()),
		Headers: http.Header{
			"Content-Type": {"application/x-www-form-urlencoded"},
		},
		Prefix: "Authenticate to mastodon",
		Path:   "/oauth/token",
		Method: http.MethodPost,
	}

	body, err := c.DoRequest(ctx, request)
	if err != nil {
		return err
	}

	var res struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return fmt.Errorf("failed to unmarshal authentication response: %v", err)
	}

	c.Config.AccessToken = res.AccessToken

	return nil
}

func (c *Client) ValidateCredentials(ctx context.Context) error {
	request := MastodonClientRequest{
		Headers: http.Header{
			"Authorization": {"Bearer " + c.Config.AccessToken},
		},
		Prefix: "Validate Credentials",
		Path:   "/api/v1/apps/verify_credentials",
		Method: http.MethodGet,
	}

	r, err := c.DoRequest(ctx, request)
	if err != nil {
		return err
	}

	var v = ValidCreds{}
	if err := json.Unmarshal(r, &v); err != nil {
		return err
	}
	fmt.Println("name: " + v.Name + " website: " + v.WebsiteURL)
	return nil
}

func (c *Client) PostStatus(ctx context.Context, toot TootRequest) (TootResponse, error) {
	var statusResponse = TootResponse{}
	payload, err := json.Marshal(toot)
	if err != nil {
		return statusResponse, fmt.Errorf("postStatus json marshaling failed: %w", err)
	}

	request := MastodonClientRequest{
		RequestBody: bytes.NewBuffer(payload),
		Headers: http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer " + c.Config.AccessToken},
		},
		Prefix: "Post Toot",
		Path:   "/api/v1/statuses",
		Method: http.MethodPost,
	}

	body, err := c.DoRequest(ctx, request)
	if err != nil {
		return statusResponse, err
	}

	if err := json.Unmarshal(body, &statusResponse); err != nil {
		return statusResponse, fmt.Errorf("postStatus unmarshaling status response failed: %w", err)
	}

	log.Printf("Status Posted: %s, URL: %s", statusResponse.ID, statusResponse.URL)

	return statusResponse, nil
}

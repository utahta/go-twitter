package twitter

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/pkg/errors"
)

const (
	BaseURL       = "https://api.twitter.com/1.1"
	UploadBaseURL = "https://upload.twitter.com/1.1"
)

type Client struct {
	HTTPClient    *http.Client
	Credentials   *oauth.Credentials
	BaseURL       string
	UploadBaseURL string
}

type ClientOption func(*Client) error

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
}

func SetConsumerCredentials(consumerKey, consumerSecret string) {
	oauthClient.Credentials.Token = consumerKey
	oauthClient.Credentials.Secret = consumerSecret
}

func New(accessToken, accessTokenSecret string, options ...ClientOption) (*Client, error) {
	c := &Client{
		HTTPClient: &http.Client{},
		Credentials: &oauth.Credentials{
			Token:  accessToken,
			Secret: accessTokenSecret,
		},
		BaseURL:       BaseURL,
		UploadBaseURL: UploadBaseURL,
	}

	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *Client) error {
		client.HTTPClient = c
		return nil
	}
}

// get issues a GET request to the Twitter API and decodes the response JSON to data.
func (c *Client) get(urlStr string, form url.Values, data interface{}) error {
	resp, err := oauthClient.Get(c.HTTPClient, c.Credentials, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// post issues a POST request to the Twitter API and decodes the response JSON to data.
func (c *Client) post(urlStr string, form url.Values, data interface{}) error {
	resp, err := oauthClient.Post(c.HTTPClient, c.Credentials, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// decodeResponse decodes the JSON response from the Twitter API.
func decodeResponse(resp *http.Response, data interface{}) error {
	if resp.StatusCode != http.StatusOK {
		p, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.Errorf("get %s returned status %d, %s", resp.Request.URL, resp.StatusCode, p)
	}
	return json.NewDecoder(resp.Body).Decode(data)
}

func makeValues(v url.Values) url.Values {
	if v == nil {
		return url.Values{}
	}
	return v
}

package twitter

import (
	"net/url"

	"github.com/utahta/go-twitter/types"
)

func (c *Client) Tweet(message string, v url.Values) (*types.Tweets, error) {
	v = makeValues(v)
	v.Set("status", message)

	tweets := &types.Tweets{}
	return tweets, c.post(c.BaseURL+"/statuses/update.json", v, tweets)
}

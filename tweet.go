package twitter

import (
	"net/url"
)

func (c *Client) Tweet(message string, v url.Values) (*Tweets, error) {
	v = makeValues(v)
	v.Set("status", message)

	tweets := &Tweets{}
	return tweets, c.post(c.BaseURL+"/statuses/update.json", v, tweets)
}

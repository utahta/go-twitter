package twitter

import (
	"io"
	"net/url"
	"strings"

	"github.com/utahta/go-twitter/types"
)

func (c *Client) Tweet(msg string, v url.Values) (*types.Tweets, error) {
	v = makeValues(v)
	v.Set("status", msg)

	tweets := &types.Tweets{}
	return tweets, c.post(c.BaseURL+"/statuses/update.json", v, tweets)
}

func (c *Client) TweetImageURLs(msg string, urlsStr []string, v url.Values) (*types.Tweets, error) {
	v = makeValues(v)
	medias, err := c.UploadMediaImageURLs(urlsStr)
	if err != nil {
		return nil, err
	}

	mediaIDs := []string{}
	for _, media := range medias {
		mediaIDs = append(mediaIDs, media.MediaIDString)
	}

	if len(mediaIDs) > 0 {
		v.Set("media_ids", strings.Join(mediaIDs, ","))
	}
	return c.Tweet(msg, v)
}

func (c *Client) TweetImages(msg string, images []io.Reader, v url.Values) (*types.Tweets, error) {
	v = makeValues(v)
	medias, err := c.UploadMediaImages(images)
	if err != nil {
		return nil, err
	}

	var mediaIDs []string
	for _, media := range medias {
		mediaIDs = append(mediaIDs, media.MediaIDString)
	}

	if len(mediaIDs) > 0 {
		v.Set("media_ids", strings.Join(mediaIDs, ","))
	}
	return c.Tweet(msg, v)
}

func (c *Client) TweetVideoURL(msg, urlStr, mediaType string, v url.Values) (*types.Tweets, error) {
	v = makeValues(v)
	media, err := c.UploadMediaVideoURL(urlStr, mediaType)
	if err != nil {
		return nil, err
	}
	v.Set("media_ids", media.MediaIDString)
	return c.Tweet(msg, v)
}

func (c *Client) TweetVideo(msg string, video io.Reader, mediaType string, v url.Values) (*types.Tweets, error) {
	v = makeValues(v)
	media, err := c.UploadMediaVideo(video, mediaType)
	if err != nil {
		return nil, err
	}
	v.Set("media_ids", media.MediaIDString)
	return c.Tweet(msg, v)
}

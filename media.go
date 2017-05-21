package twitter

import (
	"encoding/base64"
	"io/ioutil"
	"sync"

	"github.com/utahta/go-twitter/types"
	"golang.org/x/sync/errgroup"
)

func (c *Client) DownloadFile(urlStr string) ([]byte, error) {
	resp, err := c.HTTPClient.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) UploadMediaImageURLs(urlsStr []string) ([]*types.Media, error) {
	var (
		medias = []*types.Media{}
		mux    = new(sync.Mutex)
		eg     = new(errgroup.Group)
	)

	for _, urlStr := range urlsStr {
		urlStr := urlStr
		eg.Go(func() error {
			media, err := c.UploadMediaImageURL(urlStr)
			if err != nil {
				return err
			}

			mux.Lock()
			defer mux.Unlock()
			medias = append(medias, media)
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return medias, nil
}

func (c *Client) UploadMediaImageURL(urlStr string) (*types.Media, error) {
	data, err := c.DownloadFile(urlStr)
	if err != nil {
		return nil, err
	}
	return c.UploadMediaImage(data)
}

func (c *Client) UploadMediaImages(images [][]byte) ([]*types.Media, error) {
	var (
		medias = []*types.Media{}
		mux    = new(sync.Mutex)
		eg     = new(errgroup.Group)
	)

	for _, image := range images {
		image := image
		eg.Go(func() error {
			media, err := c.UploadMediaImage(image)
			if err != nil {
				return err
			}

			mux.Lock()
			defer mux.Unlock()
			medias = append(medias, media)
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return medias, nil
}

func (c *Client) UploadMediaImage(image []byte) (*types.Media, error) {
	v := makeValues(nil)
	v.Set("media_data", base64.StdEncoding.EncodeToString(image))

	media := &types.Media{}
	return media, c.post(c.UploadBaseURL+"/media/upload.json", v, media)
}

func (c *Client) UploadMediaImageAsync(image []byte) (*types.Media, error) {
	//TODO implement this
	return nil, nil
}

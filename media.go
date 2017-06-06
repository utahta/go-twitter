package twitter

import (
	"bytes"
	"encoding/base64"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/go-twitter/types"
	"golang.org/x/sync/errgroup"
)

func (c *Client) DownloadFile(urlStr string) (io.ReadCloser, int, error) {
	resp, err := c.HTTPClient.Get(urlStr)
	if err != nil {
		return nil, 0, err
	}
	return resp.Body, int(resp.ContentLength), nil
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
	body, _, err := c.DownloadFile(urlStr)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return c.UploadMediaImage(body)
}

func (c *Client) UploadMediaImages(images []io.Reader) ([]*types.Media, error) {
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

func (c *Client) UploadMediaImage(image io.Reader) (*types.Media, error) {
	buff := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, buff)
	if _, err := io.Copy(encoder, image); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}

	v := makeValues(nil)
	v.Set("media_data", buff.String())

	media := &types.Media{}
	return media, c.post(c.UploadBaseURL+"/media/upload.json", v, media)
}

func (c *Client) UploadMediaVideoURL(urlStr, mediaType string) (*types.Media, error) {
	body, length, err := c.DownloadFile(urlStr)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return c.UploadMediaVideo(body, length, mediaType)
}

func (c *Client) UploadMediaVideo(video io.Reader, length int, mediaType string) (*types.Media, error) {
	return c.UploadMediaAsync(video, length, mediaType, MediaCategoryVideo)
}

const (
	MediaCategoryImage = "tweet_image"
	MediaCategoryGif   = "tweet_gif"
	MediaCategoryVideo = "tweet_video"
)

// Upload large media file asynchronously
func (c *Client) UploadMediaAsync(body io.Reader, length int, mediaType, mediaCategory string) (*types.Media, error) {
	if length <= 0 {
		tmpBuff := &bytes.Buffer{}
		if _, err := io.Copy(tmpBuff, body); err != nil {
			return nil, err
		}
		length = tmpBuff.Len()
		body = tmpBuff
	}
	totalBytes := length

	// base64 encode
	buff := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, buff)
	if _, err := io.Copy(encoder, body); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	data := buff.String()
	dataLen := buff.Len()

	// media/upload INIT
	v := makeValues(nil)
	v.Set("command", "INIT")
	v.Set("total_bytes", strconv.FormatInt(int64(totalBytes), 10))
	v.Set("media_type", mediaType)
	v.Set("media_category", mediaCategory)

	media := &types.Media{}
	if err := c.post(c.UploadBaseURL+"/media/upload.json", v, media); err != nil {
		return nil, errors.Wrap(err, "failed to media/upload INIT")
	}
	c.Logger.Printf("media/upload INIT total_bytes:%v media_type:%v media_category:%v", totalBytes, mediaType, mediaCategory)

	// media/upload APPEND
	const mediaMaxLen = 5 * 1024 * 1024 // 5MB
	segment := 0
	for i := 0; i < dataLen; i += mediaMaxLen {
		var mediaData string
		if i+mediaMaxLen < dataLen {
			mediaData = data[i : i+mediaMaxLen]
		} else {
			mediaData = data[i:]
		}

		v = makeValues(nil)
		v.Set("command", "APPEND")
		v.Set("media_id", media.MediaIDString)
		v.Set("media_data", mediaData)
		v.Set("segment_index", strconv.FormatInt(int64(segment), 10))

		if err := c.post(c.UploadBaseURL+"/media/upload.json", v, nil); err != nil {
			return nil, errors.Wrapf(err, "failed to media/upload APPEND segment:%v", segment)
		}
		c.Logger.Printf("media/upload APPEND media_id:%v segment_index:%v", media.MediaIDString, segment)

		segment += 1
	}

	// media/upload FINALIZE
	v = makeValues(nil)
	v.Set("command", "FINALIZE")
	v.Set("media_id", media.MediaIDString)

	media = &types.Media{}
	if err := c.post(c.UploadBaseURL+"/media/upload.json", v, media); err != nil {
		return nil, errors.Wrap(err, "failed to media/upload FINALIZE")
	}
	c.Logger.Printf("media/upload FINALIZE media_id:%v state:%v check_after_secs:%v", media.MediaIDString, media.ProcessingInfo.State, media.ProcessingInfo.CheckAfterSecs)

	// media/upload STATUS
	eg := new(errgroup.Group)
	eg.Go(func() error {
		for {
			if media.ProcessingInfo.State == "succeeded" {
				break
			} else if media.ProcessingInfo.State == "failed" {
				return errors.Errorf(
					"%v:%s:%s",
					media.ProcessingInfo.Error.Code,
					media.ProcessingInfo.Error.Name,
					media.ProcessingInfo.Error.Message,
				)
			}
			time.Sleep(time.Duration(media.ProcessingInfo.CheckAfterSecs) * time.Second)

			v = makeValues(nil)
			v.Set("command", "STATUS")
			v.Set("media_id", media.MediaIDString)

			media = &types.Media{}
			if err := c.get(c.UploadBaseURL+"/media/upload.json", v, media); err != nil {
				return err
			}
			c.Logger.Printf("media/upload STATUS media_id:%v state:%v check_after_secs:%v", media.MediaIDString, media.ProcessingInfo.State, media.ProcessingInfo.CheckAfterSecs)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, errors.Wrap(err, "failed to media/upload STATUS")
	}
	c.Logger.Print("media/upload done")
	return media, nil
}

package types

type Media struct {
	MediaID          int64          `json:"media_id"`
	MediaIDString    string         `json:"media_id_string"`
	Size             int            `json:"size"`
	ExpiresAfterSecs int            `json:"expires_after_secs"`
	Image            Image          `json:"image"`
	Video            Video          `json:"video"`
	ProcessingInfo   ProcessingInfo `json:"processing_info"`
}

type Image struct {
	W         int    `json:"w"`
	H         int    `json:"h"`
	ImageType string `json:"image_type"`
}

type Video struct {
	VideoType string `json:"video_type"`
}

type ProcessingInfo struct {
	State           string `json:"state"`
	CheckAfterSecs  int    `json:"check_after_secs"`
	ProgressPercent int    `json:"progress_percent"`
	Error           struct {
		Code    int    `json:"code"`
		Name    string `json:"name"`
		Message string `json:"message"`
	} `json:"error"`
}

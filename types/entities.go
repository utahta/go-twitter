package types

// https://dev.twitter.com/overview/api/entities#entities
type Entities struct {
	Hashtags     []EntitiesHashtags    `json:"hashtags"`
	Media        []EntitiesMedia       `json:"media"`
	URLs         []EntitiesURL         `json:"urls"`
	UserMentions []EntitiesUserMention `json:"user_mentions"`
}

type EntitiesHashtags struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type EntitiesURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type EntitiesMedia struct {
	DisplayURL        string             `json:"display_url"`
	ExpandedURL       string             `json:"expanded_url"`
	ID                int64              `json:"id"`
	IDStr             string             `json:"id_str"`
	Indices           []int              `json:"indices"`
	MediaURL          string             `json:"media_url"`
	MediaURLHTTPS     string             `json:"media_url_https"`
	Sizes             EntitiesMediaSizes `json:"sizes"`
	SourceStatusID    int64              `json:"source_status_id"`
	SourceStatusIDStr string             `json:"source_status_id_str"`
	Type              string             `json:"type"`
	URL               string             `json:"url"`
	VideoInfo         EntitiesVideoInfo  `json:"video_info"`
}

type EntitiesMediaSizes struct {
	Medium EntitiesMediaSize `json:"medium"`
	Thumb  EntitiesMediaSize `json:"thumb"`
	Small  EntitiesMediaSize `json:"small"`
	Large  EntitiesMediaSize `json:"large"`
}

type EntitiesMediaSize struct {
	W      int    `json:"w"`
	H      int    `json:"h"`
	Resize string `json:"resize"`
}

type EntitiesUserMention struct {
	ID         int64  `json:"id"`
	IDStr      string `json:"id_str"`
	Indices    []int  `json:"indices"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type EntitiesVideoInfo struct {
	AspectRatio    []int             `json:"aspect_ratio"`
	DurationMillis int64             `json:"duration_millis"`
	Variants       []EntitiesVariant `json:"variants"`
}

type EntitiesVariant struct {
	Bitrate     int    `json:"bitrate"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}

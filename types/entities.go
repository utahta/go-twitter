package types

// https://dev.twitter.com/overview/api/entities#entities
type Entities struct {
	Hashtags     []Hashtags    `json:"hashtags"`
	Media        []Media       `json:"media"`
	URLs         []URL         `json:"urls"`
	UserMentions []UserMention `json:"user_mentions"`
}

type Hashtags struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type URL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type Media struct {
	DisplayURL        string     `json:"display_url"`
	ExpandedURL       string     `json:"expanded_url"`
	ID                int64      `json:"id"`
	IDStr             string     `json:"id_str"`
	Indices           []int      `json:"indices"`
	MediaURL          string     `json:"media_url"`
	MediaURLHTTPS     string     `json:"media_url_https"`
	Sizes             MediaSizes `json:"sizes"`
	SourceStatusID    int64      `json:"source_status_id"`
	SourceStatusIDStr string     `json:"source_status_id_str"`
	Type              string     `json:"type"`
	URL               string     `json:"url"`
	VideoInfo         VideoInfo  `json:"video_info"`
}

type MediaSizes struct {
	Medium MediaSize `json:"medium"`
	Thumb  MediaSize `json:"thumb"`
	Small  MediaSize `json:"small"`
	Large  MediaSize `json:"large"`
}

type MediaSize struct {
	W      int    `json:"w"`
	H      int    `json:"h"`
	Resize string `json:"resize"`
}

type UserMention struct {
	ID         int64  `json:"id"`
	IDStr      string `json:"id_str"`
	Indices    []int  `json:"indices"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type VideoInfo struct {
	AspectRatio    []int     `json:"aspect_ratio"`
	DurationMillis int64     `json:"duration_millis"`
	Variants       []Variant `json:"variants"`
}

type Variant struct {
	Bitrate     int    `json:"bitrate"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}

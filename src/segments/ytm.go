package segments

import (
	"encoding/json"

	"github.com/jandedobbeleer/oh-my-posh/src/properties"
)

type Ytm struct {
	base

	MusicPlayer
}

const (
	// APIURL is the YTMDA Remote Control API URL property.
	APIURL properties.Property = "api_url"
)

func (y *Ytm) Template() string {
	return " {{ .Icon }}{{ if ne .Status \"stopped\" }}{{ .Artist }} - {{ .Track }}{{ end }} "
}

func (y *Ytm) Enabled() bool {
	err := y.setStatus()
	// If we don't get a response back (error), the user isn't running
	// YTMDA, or they don't have the RC API enabled.
	return err == nil
}

type ytmdaStatusResponse struct {
	track  `json:"track"`
	player `json:"player"`
}

type player struct {
	SeekbarCurrentPositionHuman string  `json:"seekbarCurrentPositionHuman"`
	LikeStatus                  string  `json:"likeStatus"`
	RepeatType                  string  `json:"repeatType"`
	VolumePercent               int     `json:"volumePercent"`
	SeekbarCurrentPosition      int     `json:"seekbarCurrentPosition"`
	StatePercent                float64 `json:"statePercent"`
	HasSong                     bool    `json:"hasSong"`
	IsPaused                    bool    `json:"isPaused"`
}

type track struct {
	Author          string `json:"author"`
	Title           string `json:"title"`
	Album           string `json:"album"`
	Cover           string `json:"cover"`
	DurationHuman   string `json:"durationHuman"`
	URL             string `json:"url"`
	ID              string `json:"id"`
	Duration        int    `json:"duration"`
	IsVideo         bool   `json:"isVideo"`
	IsAdvertisement bool   `json:"isAdvertisement"`
	InLibrary       bool   `json:"inLibrary"`
}

func (y *Ytm) setStatus() error {
	// https://github.com/ytmdesktop/ytmdesktop/wiki/Remote-Control-API
	url := y.props.GetString(APIURL, "http://127.0.0.1:9863")
	httpTimeout := y.props.GetInt(APIURL, properties.DefaultHTTPTimeout)
	body, err := y.env.HTTPRequest(url+"/query", nil, httpTimeout)
	if err != nil {
		return err
	}
	q := new(ytmdaStatusResponse)
	err = json.Unmarshal(body, &q)
	if err != nil {
		return err
	}
	y.Status = playing
	y.Icon = y.props.GetString(PlayingIcon, "\uE602 ")
	if !q.HasSong {
		y.Status = stopped
		y.Icon = y.props.GetString(StoppedIcon, "\uF04D ")
	} else if q.IsPaused {
		y.Status = paused
		y.Icon = y.props.GetString(PausedIcon, "\uF8E3 ")
	}
	y.Artist = q.Author
	y.Track = q.Title
	return nil
}

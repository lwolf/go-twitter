package twitter

import (
	"encoding/base64"
	"github.com/dghubble/sling"
	"log"
	"net/http"
	"strconv"
)

/*
  "media_id": 710511363345354753,
  "media_id_string": "710511363345354753",
  "size": 11065,
  "expires_after_secs": 86400,
  "image": {
    "image_type": "image/jpeg",
    "w": 800,
    "h": 320
  }
*/

type MediaResponse struct {
	MediaId          uint64 `json:"media_id"`
	MediaIdString    string `json:"media_id_string"`
	Size             uint64 `json:"size"`
	ExpiresAfterSecs uint64 `json:"expires_after_secs"`
}

type MediaService struct {
	baseSling *sling.Sling
	sling     *sling.Sling
}

func newMediaService(sling *sling.Sling) *MediaService {
	return &MediaService{
		baseSling: sling.New(),
		sling:     sling.Path("media/"),
	}
}

type InitFormData struct {
	Command    string `url:"command"`
	MediaType  string `url:"media_type"`
	TotalBytes int    `url:"total_bytes"`
}

type AppendFormData struct {
	Command      string `url:"command"`
	MediaID      string `url:"media_id"`
	SegmentIndex string `url:"segment_index"`
	MediaData    string `url:"media_data,omitempty"`
	Media        []byte `url:"media,omitempty"`
}

type FinalizeFormData struct {
	Command string `url:"command"`
	MediaID string `url:"media_id"`
}

func (s *MediaService) Init(payloadSize int) (*MediaResponse, *http.Response, error) {
	form := InitFormData{
		Command:    "INIT",
		MediaType:  "image/jpeg",
		TotalBytes: payloadSize,
	}
	initResponse := new(MediaResponse)
	apiError := new(APIError)
	resp, err := s.sling.New().Post("upload.json").BodyForm(form).Receive(initResponse, apiError)
	return initResponse, resp, relevantError(err, *apiError)
}

// Append sends binary data,
// TODO: make it up to the client to split the big blob into the pieces if it's too big
func (s *MediaService) Append(data []byte, mediaID string, index int, b64 bool) (*http.Response, error) {
	params := AppendFormData{
		Command:      "APPEND",
		MediaID:      mediaID,
		SegmentIndex: strconv.Itoa(index),
	}
	if b64 {
		params.MediaData = base64.StdEncoding.EncodeToString(data)
	} else {
		params.Media = data
	}
	apiError := new(APIError)
	resp, err := s.sling.New().Post("upload.json").BodyForm(params).ReceiveSuccess(apiError)
	log.Printf("status_code: %v\n", resp.StatusCode)
	return resp, relevantError(err, *apiError)

}

func (s *MediaService) Finalize(mediaID string) (*MediaResponse, *http.Response, error) {
	form := FinalizeFormData{
		Command: "FINALIZE",
		MediaID: mediaID,
	}
	finalizeResponse := new(MediaResponse)
	apiError := new(APIError)
	resp, err := s.sling.New().Post("upload.json").BodyForm(form).Receive(finalizeResponse, apiError)
	return finalizeResponse, resp, relevantError(err, *apiError)

}

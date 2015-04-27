package gmail

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	BASEURL = "https://www.googleapis.com/"
)

type Gmail struct {
	AccessToken string
}

type EmailsListResponse struct {
	Messages []struct {
		ID       string `json:"id"`
		ThreadId string `json:"threadId"`
	} `json:"messages"`
	NextPageToken      string  `json:"nextPageToken"`
	ResultSizeEstimate float64 `json:"resultSizeEstimate"`
}

func (g *Gmail) GetEmails(MaxResults int) (response EmailsListResponse, err error) {
	toReturn := EmailsListResponse{}
	if len(g.AccessToken) > 0 {
		URL := BASEURL + "gmail/v1/users/me/messages?access_token=" + g.AccessToken + "&maxResults=" + strconv.Itoa(MaxResults)
		response, err := http.Get(URL)
		if err != nil {
			fmt.Printf("%s", err)
			return toReturn, err
		} else {
			defer response.Body.Close()
			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return toReturn, err
			}
			err = json.Unmarshal(contents, &toReturn)
			if err != nil {
				return toReturn, err
			}
			return toReturn, nil
		}
	}
	return EmailsListResponse{}, errors.New("Missing API key")
}

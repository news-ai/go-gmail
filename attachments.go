package gmail

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine/urlfetch"
)

type AttachmentResponse struct {
	Data string  `json:"data"`
	Size float64 `json:"size"`
}

func (g *Gmail) GetAttachmentById(c context.Context, messageId string, attachmentId string) (response AttachmentResponse, err error) {
	toReturn := AttachmentResponse{}
	if len(g.AccessToken) > 0 {
		contextWithTimeout, _ := context.WithTimeout(c, time.Second*15)
		client := urlfetch.Client(contextWithTimeout)

		URL := BASEURL + "gmail/v1/users/me/messages/" + messageId + "/attachments/" + attachmentId + "?access_token=" + g.AccessToken
		req, _ := http.NewRequest("GET", URL, nil)

		response, err := client.Do(req)
		if err != nil {
			fmt.Printf("%s", err)
			return toReturn, err
		}

		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return toReturn, err
		}

		err = json.Unmarshal(contents, &toReturn)
		if err != nil {
			return toReturn, err
		}

		if toReturn.Size == 0 {
			return toReturn, errors.New("Incorrect API key, message Id, or attachment Id.")
		}

		return toReturn, nil

	}
	return toReturn, errors.New("Missing API key")
}

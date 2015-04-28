package gmail

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AttachmentResponse struct {
	Data string  `json:"data"`
	Size float64 `json:"size"`
}

func (g *Gmail) GetAttachmentById(messageId string, attachmentId string) (response AttachmentResponse, err error) {
	toReturn := AttachmentResponse{}
	if len(g.AccessToken) > 0 {
		URL := BASEURL + "gmail/v1/users/me/messages/" + messageId + "/attachments/" + attachmentId + "?access_token=" + g.AccessToken
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
			if toReturn.Size == 0 {
				return toReturn, errors.New("Incorrect API key, message Id, or attachment Id.")
			}
			return toReturn, nil
		}
	}
	return toReturn, errors.New("Missing API key")
}

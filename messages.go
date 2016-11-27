package gmail

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"github.com/news-ai/tabulae/attach"
	"github.com/news-ai/tabulae/models"

	"golang.org/x/net/context"
	gmailv1 "google.golang.org/api/gmail/v1"
)

type EmailListResponse struct {
	Messages []struct {
		ID       string `json:"id"`
		ThreadId string `json:"threadId"`
	} `json:"messages"`
	NextPageToken      string  `json:"nextPageToken"`
	ResultSizeEstimate float64 `json:"resultSizeEstimate"`
}

type EmailIdResponse struct {
	HistoryId string   `json:"historyId"`
	ID        string   `json:"id"`
	LabelIds  []string `json:"labelIds"`
	Payload   struct {
		Body struct {
			Size float64 `json:"size"`
		} `json:"body"`
		Filename string `json:"filename"`
		Headers  []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"headers"`
		MimeType string `json:"mimeType"`
		Parts    []struct {
			Body struct {
				Size float64 `json:"size"`
			} `json:"body"`
			Filename string `json:"filename"`
			Headers  []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"headers"`
			MimeType string `json:"mimeType"`
			Parts    []struct {
				Body struct {
					Data string  `json:"data"`
					Size float64 `json:"size"`
				} `json:"body"`
				Filename string `json:"filename"`
				Headers  []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"headers"`
				MimeType string `json:"mimeType"`
				PartId   string `json:"partId"`
			} `json:"parts"`
		} `json:"parts"`
	} `json:"payload"`
	SizeEstimate float64 `json:"sizeEstimate"`
	Snippet      string  `json:"snippet"`
	ThreadId     string  `json:"threadId"`
}

type Message struct {
	Raw string `json:"raw,omitempty"`
}

func (g *Gmail) SendEmailWithAttachments(r *http.Request, c context.Context, from string, to string, subject string, body string, email models.Email, files []models.File) (string, string, error) {

	if len(g.AccessToken) > 0 {
		contextWithTimeout, _ := context.WithTimeout(c, time.Second*15)
		client := urlfetch.Client(contextWithTimeout)

		nl := "\r\n" // newline
		boundary := "__newsai_tabulae__"

		var message Message
		temp := []byte("MIME-Version: 1.0" + nl +
			"To:  " + to + nl +
			"From: " + from + nl +
			"reply-to: " + from + nl +
			"Subject: " + subject + nl +

			"Content-Type: multipart/mixed; boundary=\"" + boundary + "\"" + nl + nl +

			// Boundary one is email itself
			"--" + boundary + nl +

			"Content-Type: text/html; charset=UTF-8" + nl +
			"MIME-Version: 1.0" + nl +
			"Content-Transfer-Encoding: base64" + nl + nl +

			// Body itself
			body + nl + nl)

		for i := 0; i < len(files); i++ {
			bytesArray, attachmentType, fileNames, err := attach.GetAttachmentsForEmail(r, email, files)
			if err == nil {
				for i := 0; i < len(bytesArray); i++ {
					log.Infof(c, "%v", attachmentType[i])
					log.Infof(c, "%v", fileNames[i])

					str := base64.StdEncoding.EncodeToString(bytesArray[i])

					attachment := []byte(
						"--" + boundary + nl +
							"Content-Type: " + attachmentType[i] + nl +
							"MIME-Version: 1.0" + nl +
							"Content-Disposition: attachment; filename=\"" + fileNames[i] + "\"" + nl +
							"Content-Transfer-Encoding: base64" + nl + nl +
							str + nl + nl,
					)

					temp = append(temp, attachment...)
				}
			}
		}

		finalBoundry := []byte(
			"--" + boundary + "--",
		)

		temp = append(temp, finalBoundry...)

		message.Raw = base64.StdEncoding.EncodeToString(temp)
		message.Raw = strings.Replace(message.Raw, "/", "_", -1)
		message.Raw = strings.Replace(message.Raw, "+", "-", -1)
		message.Raw = strings.Replace(message.Raw, "=", "", -1)

		messageJson, err := json.Marshal(message)
		if err != nil {
			log.Errorf(c, "%v", err)
			return "", "", err
		}

		messageQuery := bytes.NewReader(messageJson)

		log.Infof(c, "%v", messageQuery)

		URL := BASEURL + "gmail/v1/users/me/messages/send?uploadType=multipart"
		req, _ := http.NewRequest("POST", URL, messageQuery)

		req.Header.Add("Authorization", "Bearer "+g.AccessToken)
		req.Header.Add("Content-Type", "application/json")

		response, err := client.Do(req)
		if err != nil {
			log.Errorf(c, "%v", err)
			return "", "", err
		}

		// Decode JSON from Google
		decoder := json.NewDecoder(response.Body)
		var gmailMessage gmailv1.Message
		err = decoder.Decode(&gmailMessage)
		if err != nil {
			log.Errorf(c, "%v", err)
			return "", "", err
		}

		log.Infof(c, "%v", gmailMessage)

		return gmailMessage.Id, gmailMessage.ThreadId, nil
	}

	return "", "", errors.New("No access token supplied")
}

func (g *Gmail) SendEmail(c context.Context, from string, to string, subject string, body string) (string, string, error) {
	if len(g.AccessToken) > 0 {
		contextWithTimeout, _ := context.WithTimeout(c, time.Second*15)
		client := urlfetch.Client(contextWithTimeout)

		var message Message
		temp := []byte("From: " + from + "\r\n" +
			"reply-to: " + from + "\r\n" +
			"Content-type: text/html;charset=iso-8859-1\r\n" +
			"MIME-Version: 1.0\r\n" +
			"To:  " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" + body)

		message.Raw = base64.StdEncoding.EncodeToString(temp)
		message.Raw = strings.Replace(message.Raw, "/", "_", -1)
		message.Raw = strings.Replace(message.Raw, "+", "-", -1)
		message.Raw = strings.Replace(message.Raw, "=", "", -1)

		messageJson, err := json.Marshal(message)
		if err != nil {
			log.Errorf(c, "%v", err)
			return "", "", err
		}

		messageQuery := bytes.NewReader(messageJson)

		URL := BASEURL + "gmail/v1/users/me/messages/send"
		req, _ := http.NewRequest("POST", URL, messageQuery)

		req.Header.Add("Authorization", "Bearer "+g.AccessToken)
		req.Header.Add("Content-Type", "application/json")

		response, err := client.Do(req)
		if err != nil {
			log.Errorf(c, "%v", err)
			return "", "", err
		}

		// Decode JSON from Google
		decoder := json.NewDecoder(response.Body)
		var gmailMessage gmailv1.Message
		err = decoder.Decode(&gmailMessage)
		if err != nil {
			log.Errorf(c, "%v", err)
			return "", "", err
		}

		return gmailMessage.Id, gmailMessage.ThreadId, nil
	}

	return "", "", errors.New("No access token supplied")
}

func (g *Gmail) GetEmails(c context.Context, MaxResults int) (response EmailListResponse, err error) {
	toReturn := EmailListResponse{}
	if len(g.AccessToken) > 0 {
		contextWithTimeout, _ := context.WithTimeout(c, time.Second*15)
		client := urlfetch.Client(contextWithTimeout)

		URL := BASEURL + "gmail/v1/users/me/messages?access_token=" + g.AccessToken + "&maxResults=" + strconv.Itoa(MaxResults)
		req, _ := http.NewRequest("GET", URL, nil)

		response, err := client.Do(req)
		if err != nil {
			log.Errorf(c, "%v", err)
			return toReturn, err
		}

		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Errorf(c, "%v", err)
			return toReturn, err
		}

		err = json.Unmarshal(contents, &toReturn)
		if err != nil {
			log.Errorf(c, "%v", err)
			return toReturn, err
		}

		if toReturn.ResultSizeEstimate == 0 {
			return toReturn, errors.New("Incorrect API key")
		}

		return toReturn, nil
	}

	return toReturn, errors.New("Missing API key")
}

func (g *Gmail) GetEmailById(c context.Context, emailId string) (response EmailIdResponse, err error) {
	toReturn := EmailIdResponse{}
	if len(g.AccessToken) > 0 {
		contextWithTimeout, _ := context.WithTimeout(c, time.Second*15)
		client := urlfetch.Client(contextWithTimeout)

		URL := BASEURL + "gmail/v1/users/me/messages/" + emailId + "?access_token=" + g.AccessToken
		req, _ := http.NewRequest("GET", URL, nil)

		response, err := client.Do(req)
		if err != nil {
			log.Errorf(c, "%v", err)
			return toReturn, err
		}

		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Errorf(c, "%v", err)
			return toReturn, err
		}

		err = json.Unmarshal(contents, &toReturn)
		if err != nil {
			log.Errorf(c, "%v", err)
			return toReturn, err
		}

		if toReturn.SizeEstimate == 0 {
			return toReturn, errors.New("Incorrect API key")
		}

		return toReturn, nil
	}

	return toReturn, errors.New("Missing API key")
}

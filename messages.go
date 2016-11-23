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

	"golang.org/x/net/context"
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

func (g *Gmail) SendEmail(c context.Context, from string, to string, subject string, body string) error {
	if len(g.AccessToken) > 0 {
		contextWithTimeout, _ := context.WithTimeout(c, time.Second*15)
		client := urlfetch.Client(contextWithTimeout)

		var message Message
		temp := []byte("From: 'me'\r\n" +
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
			return err
		}

		messageQuery := bytes.NewReader(messageJson)

		URL := BASEURL + "gmail/v1/users/me/messages/send"
		req, _ := http.NewRequest("POST", URL, messageQuery)

		req.Header.Add("Authorization", "Bearer "+g.AccessToken)
		req.Header.Add("Content-Type", "application/json")

		response, err := client.Do(req)
		if err != nil {
			log.Errorf(c, "%v", err)
			return err
		}

		log.Infof(c, "%v", response)

		return nil
	}

	return errors.New("No access token supplied")
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

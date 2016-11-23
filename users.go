package gmail

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"google.golang.org/appengine/urlfetch"

	"golang.org/x/net/context"
)

type UserProfile struct {
	Email         string `json:"email"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	Hd            string `json:"hd"`
	ID            string `json:"id"`
	Locale        string `json:"locale"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

func (g *Gmail) GetProfile(c context.Context) (response UserProfile, err error) {
	toReturn := UserProfile{}
	if len(g.AccessToken) > 0 {
		contextWithTimeout, _ := context.WithTimeout(c, time.Second*15)
		client := urlfetch.Client(contextWithTimeout)

		URL := BASEURL + "oauth2/v1/userinfo?access_token=" + g.AccessToken
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

		if toReturn.Email == "" {
			return toReturn, errors.New("Incorrect API key")
		}

		return toReturn, nil
	}

	return toReturn, errors.New("Missing API key")
}

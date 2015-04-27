package gmail

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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

func (g *Gmail) GetProfile() (response UserProfile, err error) {
	toReturn := UserProfile{}
	if len(g.AccessToken) > 0 {
		URL := BASEURL + "oauth2/v1/userinfo?access_token=" + g.AccessToken
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
	return toReturn, errors.New("Missing API key")
}

package gmail

import (
	"fmt"
	"testing"
)

func TestGetProfile(t *testing.T) {
	var gmail Gmail
	gmail.AccessToken = AccessToken
	singleResult, err := gmail.GetProfile()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Sprintf("%s\n", singleResult.Email)
}

package gmail

import (
	"fmt"
	"testing"
)

func TestGetEmails(t *testing.T) {
	var gmail Gmail
	gmail.AccessToken = AccessToken
	allResults, err := gmail.GetEmails(100)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Sprintf("", allResults)
	return
}

func TestGetEmailById(t *testing.T) {
	var gmail Gmail
	gmail.AccessToken = AccessToken
	allResults, err := gmail.GetEmails(100)
	if err != nil {
		fmt.Println(err)
		return
	}
	singleResult, err := gmail.GetEmailById(allResults.Messages[0].ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Sprintf("", singleResult.Snippet)
	return
}

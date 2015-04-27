package gmail

import (
	"fmt"
	"testing"
)

func TestGetEmails(t *testing.T) {
	var gmail Gmail
	gmail.AccessToken = ""
	allResults, err := gmail.GetEmails(100)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(allResults)
}

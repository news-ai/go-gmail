package gmail

import (
	"fmt"
	"testing"
)

func TestGetAttachmentById(t *testing.T) {
	var gmail Gmail
	gmail.AccessToken = AccessToken
	singleAttachment, err := gmail.GetAttachmentById("", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(singleAttachment.Size)
}

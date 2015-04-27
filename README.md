# go-gmail

A simple wrapper ontop of the Gmail API. This library is useful only after you have an access token from the user. 

### Gmail API

The Gmail REST API is structured by:

1. Drafts
2. History
3. Labels
4. Messages
5. Message Attachments
6. Threads
7. Users

The rate limit will be set using:

1. Daily Usage at 1,000,000,000 quota units per day
2. Per User Rate Limit at 250 quota units per user per second

### Users

GetProfile of a single user

```go
var gmail Gmail
gmail.AccessToken = AccessToken
singleResult, err := gmail.GetProfile()
if err != nil {
	fmt.Println(err)
	return
}
fmt.Sprintf("%s\n", singleResult.Email)
```
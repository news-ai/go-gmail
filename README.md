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

### Documentation

Get the package, and install it:

`go get github.com/abhiagarwal/go-gmail`, `import Gmail "github.com/abhiagarwal/go-gmail"`

All of the methods below require an access token for a single user. 

#### 1. Drafts

#### 2. History

#### 3. Labels

#### 4. Messages

Get the emails of the user.

```go
var gmail Gmail.Gmail
gmail.AccessToken = AccessToken
allResults, err := gmail.GetEmails(100)
if err != nil {
	fmt.Println(err)
	return
}
fmt.Println(allResults)
return
```

Get a single email of the user given an Id.

```go
var gmail Gmail.Gmail
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
fmt.Println(singleResult)
return
```

#### 5. Message Attachments

Get attachment from a specific email that's received by a user.

```go
var gmail Gmail
gmail.AccessToken = AccessToken
singleAttachment, err := gmail.GetAttachmentById("", "")
if err != nil {
	fmt.Println(err)
	return
}
fmt.Println(singleAttachment)
```

#### 6. Threads

#### 7. Users

Get the profile of the user.

```go
var gmail Gmail.Gmail
gmail.AccessToken = AccessToken
singleResult, err := gmail.GetProfile()
if err != nil {
	fmt.Println(err)
	return
}
fmt.Println(singleResult)
```

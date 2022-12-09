package webfinger

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/carlmjohnson/requests"
)

type Webfinger struct {
	Subject string   `json:"subject"`
	Aliases []string `json:"aliases"`
	Links   []Link   `json:"links"`
}

func (wf *Webfinger) ActivityPub() (string, error) {
	for _, link := range wf.Links {
		if link.Type == "application/activity+json" {
			return link.Href, nil
		}
	}
	return "", fmt.Errorf("no ActivityPub link found")
}

type Link struct {
	Rel      string `json:"rel"`
	Type     string `json:"type"`
	Href     string `json:"href"`
	Template string `json:"template"`
}

type Acct struct {
	User string
	Host string
}

func (a *Acct) String() string {
	return "acct:" + a.User + "@" + a.Host
}

// Webfinger returns the URL for the webfinger resource for this Acct.
func (a *Acct) Webfinger() string {
	return "https://" + a.Host + "/.well-known/webfinger?resource=" + url.QueryEscape(a.String())
}

// ID returns the URL for the ID resource for this Acct.
func (a *Acct) ID() string {
	return "https://" + a.Host + "/users/" + a.User
}

// Followers returns the URL for the followers collection for this Acct.
func (a *Acct) Followers() string {
	return a.ID() + "/followers"
}

// Following returns the URL for the following collection for this Acct.
func (a *Acct) Following() string {
	return a.ID() + "/following"
}

// Inbox returns the URL for the inbox collection for this Acct.
func (a *Acct) Inbox() string {
	return a.ID() + "/inbox"
}

// Outbox returns the URL for the outbox collection for this Acct.
func (a *Acct) Outbox() string {
	return a.ID() + "/outbox"
}

// SharedInbox returns the URL for the shared inbox collection for this Acct.
func (a *Acct) SharedInbox() string {
	return "https://" + a.Host + "/inbox"
}

func (a *Acct) Fetch(ctx context.Context) (*Webfinger, error) {
	var webfinger Webfinger
	err := requests.URL(a.Webfinger()).ToJSON(&webfinger).Fetch(ctx)
	return &webfinger, err
}

func Parse(acct string) (*Acct, error) {
	u, err := url.Parse(acct)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "acct" {
		return nil, fmt.Errorf("invalid scheme: %s", u.Scheme)
	}
	parts := strings.SplitN(u.Opaque, "@", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid opaque: %s", u.Opaque)
	}

	return &Acct{
		User: parts[0],
		Host: parts[1],
	}, nil
}
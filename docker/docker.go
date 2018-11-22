package docker

// this package recieves the Docker Hub Automated Build webhook
// https://docs.docker.com/docker-hub/webhooks/
// NOT the Docker Trusted Registry webhook
// https://docs.docker.com/ee/dtr/user/create-and-manage-webhooks/

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// parse errors
var (
	ErrInvalidHTTPMethod = errors.New("invalid HTTP Method")
	ErrParsingPayload    = errors.New("error parsing payload")
)

// Event defines a GitHub hook event type
type Event string

// GitHub hook types
const (
	BuildEvent Event = "build"
)

// BuildPayload a docker hub build notice
// https://docs.docker.com/docker-hub/webhooks/
type BuildPayload struct {
	CallbackURL string `json:"callback_url"`
	PushData    struct {
		Images   []string `json:"images"`
		PushedAt float32  `json:"pushed_at"`
		Pusher   string   `json:"pusher"`
		Tag      string   `json:"tag"`
	} `json:"push_data"`
	Repository struct {
		CommentCount    int     `json:"comment_count"`
		DateCreated     float32 `json:"date_created"`
		Description     string  `json:"description"`
		Dockerfile      string  `json:"dockerfile"`
		FullDescription string  `json:"full_description"`
		IsOfficial      bool    `json:"is_official"`
		IsPrivate       bool    `json:"is_private"`
		IsTrusted       bool    `json:"is_trusted"`
		Name            string  `json:"name"`
		Namespace       string  `json:"namespace"`
		Owner           string  `json:"owner"`
		RepoName        string  `json:"repo_name"`
		RepoURL         string  `json:"repo_url"`
		StarCount       int     `json:"star_count"`
		Status          string  `json:"status"`
	} `json:"repository"`
}

// there are no options for docker webhooks
// however I'm leaving this here for now in anticipation of future support for Docker Trusted Registry

// // Option is a configuration option for the webhook
// type Option func(*Webhook) error

// // Options is a namespace var for configuration options
// var Options = WebhookOptions{}

// // WebhookOptions is a namespace for configuration option methods
// type WebhookOptions struct{}

// // Secret registers the GitHub secret
// func (WebhookOptions) Secret(secret string) Option {
// 	return func(hook *Webhook) error {
// 		hook.secret = secret
// 		return nil
// 	}
// }

// Webhook instance contains all methods needed to process events
type Webhook struct {
	secret string
}

// New creates and returns a WebHook instance denoted by the Provider type
func New() (*Webhook, error) {
	// func New(options ...Option) (*Webhook, error) {
	hook := new(Webhook)
	// for _, opt := range options {
	// 	if err := opt(hook); err != nil {
	// 		return nil, errors.New("Error applying Option")
	// 	}
	// }
	return hook, nil
}

// Parse verifies and parses the events specified and returns the payload object or an error
func (hook Webhook) Parse(r *http.Request, events ...Event) (interface{}, error) {
	defer func() {
		_, _ = io.Copy(ioutil.Discard, r.Body)
		_ = r.Body.Close()
	}()

	if r.Method != http.MethodPost {
		return nil, ErrInvalidHTTPMethod
	}

	// event := r.Header.Get("X-GitHub-Event")
	// if event == "" {
	// 	return nil, ErrMissingGithubEventHeader
	// }
	// gitHubEvent := Event(event)

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil || len(payload) == 0 {
		log.Error(ErrParsingPayload)
		log.Error(err)
		return nil, ErrParsingPayload
	}

	var pl BuildPayload
	err = json.Unmarshal([]byte(payload), &pl)
	if err != nil {
		log.Error(ErrParsingPayload)
		log.Error(err)
		return nil, ErrParsingPayload
	}
	return pl, err

}

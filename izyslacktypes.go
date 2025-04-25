package izyslackgo

import (
	"github.com/slack-go/slack"
)

type IzySlack struct {
	Credentials *IzySlackCredentials
	Client      *slack.Client
}

type IzySlackCredentials struct {
	SlackSignSecret string
	BotToken        string
}

type IzySlackEvent struct {
	Token          string                  `json:"token"`
	TeamID         string                  `json:"team_id"`
	APIAppID       string                  `json:"api_app_id"`
	Event          IzySlackSubEvent        `json:"event"`
	Type           string                  `json:"type"`
	EventID        string                  `json:"event_id"`
	EventTime      int64                   `json:"event_time"`
	Authorizations []IzySlackAuthorization `json:"authorizations"`
	IsExtShared    bool                    `json:"is_ext_shared_channel"`
	EventContext   string                  `json:"event_context"`
}

type IzySlackSubEvent struct {
	User        string         `json:"user"`
	Type        string         `json:"type"`
	TS          string         `json:"ts"`
	ClientMsgID string         `json:"client_msg_id"`
	ThreadTs    string         `json:"thread_ts"`
	Text        string         `json:"text"`
	Team        string         `json:"team"`
	Blocks      []IzSlackBlock `json:"blocks"`
	Channel     string         `json:"channel"`
	EventTS     string         `json:"event_ts"`
}

type IzSlackBlock struct {
	Type     string            `json:"type"`
	BlockID  string            `json:"block_id"`
	Elements []IzySlackElement `json:"elements"`
}

type IzySlackElement struct {
	Type     string               `json:"type"`
	Elements []IzySlackSubElement `json:"elements,omitempty"`
	UserID   string               `json:"user_id,omitempty"`
	Text     string               `json:"text,omitempty"`
}

type IzySlackSubElement struct {
	Type   string `json:"type"`
	UserID string `json:"user_id,omitempty"`
	Text   string `json:"text,omitempty"`
}

type IzySlackAuthorization struct {
	EnterpriseID        *string `json:"enterprise_id"`
	TeamID              string  `json:"team_id"`
	UserID              string  `json:"user_id"`
	IsBot               bool    `json:"is_bot"`
	IsEnterpriseInstall bool    `json:"is_enterprise_install"`
}

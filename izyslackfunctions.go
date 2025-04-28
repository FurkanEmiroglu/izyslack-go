package izyslackgo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
)

func Initialize(signSecret string, token string) *IzySlack {
	creds := IzySlackCredentials{
		SlackSignSecret: signSecret,
		BotToken:        token,
	}

	slackClient := slack.New(token)

	izySlack := IzySlack{
		Credentials: &creds,
		Client:      slackClient,
	}

	return &izySlack
}

func (receiver *IzySlack) IsChallengeRequest(body []byte) (bool, *string, error) {
	var data map[string]any
	err := json.Unmarshal(body, &data)

	if err != nil {
		return false, nil, err
	}

	if challenge, ok := data["challenge"].(string); ok {
		return true, &challenge, nil
	} else {
		return false, nil, nil
	}
}

func (receiver *IzySlack) HandleChallengeRequestLambda(challenge string) events.LambdaFunctionURLResponse {
	return events.LambdaFunctionURLResponse{
		StatusCode:      http.StatusOK,
		Headers:         map[string]string{"Content-Type": "text/plain"},
		Body:            challenge,
		IsBase64Encoded: false,
		Cookies:         nil,
	}
}

func (receiver *IzySlack) HandleChallengeRequest(w http.ResponseWriter, challenge string) {
	fmt.Fprintln(w, challenge, http.StatusOK)
}

func (receiver *IzySlack) UploadFile(channel string, filePath string) error {
	params := slack.UploadFileV2Parameters{
		Channel: channel,
		File:    filePath,
		Title:   "File Upload",
	}

	_, err := receiver.Client.UploadFileV2(params)

	if err != nil {
		return err
	}

	return nil
}

func (receiver *IzySlack) ReceiveEvent(body []byte) (*IzySlackEvent, error) {
	var event IzySlackEvent
	err := json.Unmarshal(body, &event)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (receiver *IzySlack) VerifySignatureLambda(header map[string]string, body []byte) bool {
	slackTimestamp := header["X-Slack-Request-Timestamp"]
	slackSignature := header["X-Slack-Signature"]

	ts, err := time.ParseDuration(slackTimestamp + "s")
	if err != nil || time.Since(time.Unix(int64(ts.Seconds()), 0)) > 5*time.Minute {
		fmt.Println("Timestamp is too old or invalid")
		return false
	}

	baseString := fmt.Sprintf("v0:%s:%s", slackTimestamp, string(body))
	mac := hmac.New(sha256.New, []byte(receiver.Credentials.SlackSignSecret))
	mac.Write([]byte(baseString))
	expectedSignature := "v0=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expectedSignature), []byte(slackSignature))
}

func (receiver *IzySlack) VerifySignature(header http.Header, body []byte) bool {
	slackTimestamp := header.Get("X-Slack-Request-Timestamp")
	slackSignature := header.Get("X-Slack-Signature")

	ts, err := time.ParseDuration(slackTimestamp + "s")
	if err != nil || time.Since(time.Unix(int64(ts.Seconds()), 0)) > 5*time.Minute {
		fmt.Println("Timestamp is too old or invalid")
		return false
	}

	baseString := fmt.Sprintf("v0:%s:%s", slackTimestamp, string(body))
	mac := hmac.New(sha256.New, []byte(receiver.Credentials.SlackSignSecret))
	mac.Write([]byte(baseString))
	expectedSignature := "v0=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expectedSignature), []byte(slackSignature))
}

func (receiver *IzySlack) SendMessage(channel string, text string) error {
	_, _, err := receiver.Client.PostMessage(channel, slack.MsgOptionText(text, false))
	if err != nil {
		fmt.Println("Error while replying in slack")
		return err
	}
	return nil
}

func (receiver *IzySlack) ReplyMessage(channel string, text string, ts string) error {
	_, _, err := receiver.Client.PostMessage(channel,
		slack.MsgOptionText(text, false),
		slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			ThreadTimestamp: ts,
		}),
	)

	if err != nil {
		fmt.Println("Error while replying in slack")
		return err
	}
	return nil
}

func (receiver *IzySlack) ReplyWithFile(channel string, filePath string, threadTs string) error {
	// Dosya yolunun geçerli olduğundan emin ol
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("dosya açılamadı: %v", err)
	}
	defer file.Close()

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("dosya boyutu alınamadı: %v", err)
	}

	fileSize := fileInfo.Size()

	params := slack.UploadFileV2Parameters{
		Channel:         channel,
		File:            filePath, // Dosya yolunu burada kullanıyoruz
		Filename:        "response.txt",
		ThreadTimestamp: threadTs,
		FileSize:        int(fileSize),
	}

	// Dosyayı Slack'e yükle
	fileUpload, err := receiver.Client.UploadFileV2(params)
	if err != nil {
		return fmt.Errorf("dosya yüklenemedi: %v", err)
	}

	// Dosya başarıyla yüklendi
	fmt.Println("Dosya başarıyla yüklendi:", fileUpload.Title)
	return nil
}

func (receiver *IzySlack) GetUserInfo(userId string) (*slack.User, error) {
	identity, err := receiver.Client.GetUserInfo(userId)

	if err != nil {
		return nil, err
	}

	return identity, nil
}

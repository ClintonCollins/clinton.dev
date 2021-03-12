package email

import (
	"log"

	"github.com/keighl/postmark"
)

type Instance struct {
	postmarkClient *postmark.Client
}

func New(postmarkServerToken, postmarkAccountToken string) *Instance {
	return &Instance{
		postmarkClient: postmark.NewClient(postmarkServerToken, postmarkAccountToken),
	}
}

func (inst *Instance) SendEmail(from, replyTo, to, subject, body string) error {
	email := postmark.Email{
		From:        from,
		To:          to,
		Cc:          "",
		Bcc:         "",
		Subject:     subject,
		Tag:         "",
		HtmlBody:    body,
		ReplyTo:     replyTo,
		TrackOpens:  false,
		Attachments: nil,
		Metadata:    nil,
	}
	_, err := inst.postmarkClient.SendEmail(email)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (inst *Instance) SendEmailTemplate(from, replyTo, to, contactName, contactMessage string, templateID int64) error {
	templateVariables := map[string]interface{}{
		"contactFormNameEl":    contactName,
		"contactFormMessageEl": contactMessage,
	}
	emailTemplate := postmark.TemplatedEmail{
		// TemplateId:    19927936,
		TemplateId:    templateID,
		TemplateAlias: "Contact Form Response",
		TemplateModel: templateVariables,
		InlineCss:     true,
		From:          from,
		To:            to,
		Cc:            "",
		Bcc:           "",
		Tag:           "",
		ReplyTo:       replyTo,
		Headers:       nil,
		TrackOpens:    false,
		Attachments:   nil,
	}
	_, err := inst.postmarkClient.SendTemplatedEmail(emailTemplate)
	if err != nil {
		log.Println(err)
	}
	return err
}

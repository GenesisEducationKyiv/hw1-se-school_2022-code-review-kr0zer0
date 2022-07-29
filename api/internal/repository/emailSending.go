package repository

import (
	"fmt"
	"github.com/mailjet/mailjet-apiv3-go"
	"os"
)

type EmailSendingRepository struct {
}

func NewEmailSendingRepository() *EmailSendingRepository {
	return &EmailSendingRepository{}
}

func (r EmailSendingRepository) SendToList(emails []string, message string) error {
	if len(emails) == 0 {
		return nil
	}

	var sendingList []mailjet.InfoSendMail
	for _, email := range emails {
		info := mailjet.InfoSendMail{
			FromEmail: "krokkozerro@gmail.com",
			FromName:  "BTC app",
			Recipients: []mailjet.Recipient{
				mailjet.Recipient{
					Email: email,
				},
			},
			Subject:  "BTC exchange rate",
			HTMLPart: fmt.Sprintf("<h3>%s</h3>", message),
		}
		sendingList = append(sendingList, info)
	}

	mailjetClient := mailjet.NewMailjetClient(os.Getenv("MAILJET_PUBLIC_KEY"), os.Getenv("MAILJET_PRIVATE_KEY"))
	email := &mailjet.InfoSendMail{
		Messages: sendingList,
	}

	_, err := mailjetClient.SendMail(email)
	if err != nil {
		return err
	}

	return nil
}

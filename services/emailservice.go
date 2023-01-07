package services

import (
	"crypto/tls"
	"fmt"
	"hypegym-backend/models"

	gomail "gopkg.in/gomail.v2"
)

func SendContactUsMail(contact models.Contact) {
	msgBody := fmt.Sprintf("<b>Dear %s %s, we have recieved your message. We are going to response as soon as possible. Keep in hype!<b>", contact.FirstName, contact.LastName)
	msg := gomail.NewMessage()
	msg.SetHeader("From", "hypeline.hypegym@gmail.com")
	msg.SetHeader("To", contact.Email)
	msg.SetHeader("Subject", "HYPEGYM")
	msg.SetBody("text/html", msgBody)

	n := gomail.NewDialer("smtp.gmail.com", 587, "hypeline.hypegym@gmail.com", "qmkkxtssqkanickh")
	n.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

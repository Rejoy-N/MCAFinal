package main

import (
	"log"
	"net/smtp"
)

func orderEmail(el string) {
	    from := "123redoc@gmail.com"
    password := "something@1"

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	to := el
	msg := "From: " + from + "\r\n" +
           "To: " + to + "\r\n" + 
           "Subject: Order Confirmation" + "\r\n\r\n" +
           "You have sucessfully placed your order with Mobile Store" + "\r\n" +
		   "Sincerely" + "\r\n" +	
		   "The Online Mobile Store" + "\r\n" +		
		   "IGNOU Enrollment ID: 137132696" + "\r\n" + "\r\n\r\n" + 
		   "This is a test email of a academic project" + "\r\n" 
    /* Ports 465 and 587 are intended for email client to email server communication (sending email). Port 465 is for smtps. SSL encryption is started automatically before any SMTP level communication. Port 587 is for msa. It is almost like standard SMTP port. */
           
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("Error: %s", err)
        return
	}
	log.Print("message sent")
}

/*

func main() {
	el := "rejoy_nair@yahoo.com"
	orderEmail(el)
}
*/

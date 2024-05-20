package utils

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

var monthMap = map[int]string{
	1:  "January",
	2:  "February",
	3:  "March",
	4:  "April",
	5:  "May",
	6:  "June",
	7:  "July",
	8:  "August",
	9:  "September",
	10: "October",
	11: "November",
	12: "December",
}

func SendVerifyAccountHealthSpecialistEmail(mailClient mailer.Mailer, toName string, toEmail string) error {
	token, err := GenerateVerificationToken(toEmail)
	if err != nil {
		return errors.New("Failed to generate verification token")
	}

	verificationLink := "http://localhost:3000/confirmation/" + token

	return mailClient.Send(&mailer.Message{
		From: mail.Address{
			Address: "hello@noreply.com",
		},
		To:      []mail.Address{{Name: toName, Address: toEmail}},
		Subject: "Email Verification",
		HTML: fmt.Sprintf(`
			<p>Hello,</p>
			<p>Thank you for joining us at WellnessWave.</p>
			<p>Click on the button below to verify your email address.</p>
			<p>
			<a class="btn" href="%s" target="_blank" rel="noopener">Verify</a>
			</p>
			<p>
			Thanks,<br/>
			WellnessWave team
			</p>
		`, verificationLink),
	})
}

func SendVerifyAccountPatientEmail(mailClient mailer.Mailer, toName string, toEmail string, oldPassword string, id string) error {
	token, err := GeneratePatientVerificationToken(toEmail, oldPassword, id)
	if err != nil {
		return errors.New("Failed to generate verification token")
	}

	verificationLink := "http://localhost:3000/patientConfirmation/" + token

	return mailClient.Send(&mailer.Message{
		From: mail.Address{
			Address: "hello@noreply.com",
		},
		To:      []mail.Address{{Name: toName, Address: toEmail}},
		Subject: "Email Verification",
		HTML: fmt.Sprintf(`
			<p>Hello,</p>
			<p>Thank you for joining us at WellnessWave.</p>
			<p>Click on the button below to verify your email address.</p>
			<p>This is your current randomly generated password </p>
			<p>%s</p> 
			<p>You will be asked to change it during the verification process</p>
			<p>
			<a class="btn" href="%s" target="_blank" rel="noopener">Verify</a>
			</p>
			<p>
			Thanks,<br/>
			WellnessWave team
			</p>
		`, oldPassword, verificationLink),
	})
}

func SendEventEmailToPatient(mailClient mailer.Mailer, toName string, toEmail string, eventRecord *models.Record) error {
	dateAndTime := eventRecord.GetDateTime("event_date").Time()
	day, month, year, hour, minute := dateAndTime.Day(), dateAndTime.Month(), dateAndTime.Year(), dateAndTime.Hour(), dateAndTime.Minute()

	return mailClient.Send(&mailer.Message{
		From: mail.Address{
			Address: "hello@noreply.com",
		},
		To:      []mail.Address{{Name: toName, Address: toEmail}},
		Subject: "Event Reminder",
		HTML: fmt.Sprintf(`
			<p>Hello,</p>
			<p>Your practitioner scheduled an event with you.</p>
			<p>See you on the %d of %s %d at %d:%d!</p>
			Thanks,<br/>
			WellnessWave team
			</p>
		`, day, monthMap[int(month)], year, hour, minute),
	})
}

func SendRescheduleEventEmailToPatient(mailClient mailer.Mailer, toName string, toEmail string, eventRecord *models.Record) error {
	dateAndTime := eventRecord.GetDateTime("event_date").Time()
	day, month, year, hour, minute := dateAndTime.Day(), dateAndTime.Month(), dateAndTime.Year(), dateAndTime.Hour(), dateAndTime.Minute()

	return mailClient.Send(&mailer.Message{
		From: mail.Address{
			Address: "hello@noreply.com",
		},
		To:      []mail.Address{{Name: toName, Address: toEmail}},
		Subject: "Event Reminder",
		HTML: fmt.Sprintf(`
			<p>Hello,</p>
			<p>Your practitioner has re-scheduled an event with you.</p>
			<p>See you on the %d of %s %d at %d:%d!</p>
			Thanks,<br/>
			WellnessWave team
			</p>
		`, day, monthMap[int(month)], year, hour, minute),
	})
}

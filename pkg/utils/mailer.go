package utils

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/pocketbase/pocketbase/tools/mailer"
)

func SendVerifyAccountHealthSpecialistEmail(mailClient mailer.Mailer, from string, toName string, toEmail string) error {
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

func SendVerifyAccountPatientEmail(mailClient mailer.Mailer, from string, toName string, toEmail string, oldPassword string) error {
	token, err := GeneratePatientVerificationToken(toEmail, oldPassword)
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
			<p>This is your current randomly generated password: %s. You will be asked to change it during the verification process</p>
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

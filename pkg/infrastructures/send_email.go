package infrastructures

import (
	"fmt"
	"loan-management/config"
	"net/smtp"
)

func sendEmail(subject, data string, to []string) error {
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}
	from := config.Email.Address
	key := config.Email.Key
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	message := []byte("Subject: " + subject + "\n" + "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" + data)
	auth := smtp.PlainAuth("", from, key, host)
	if err = smtp.SendMail(address, auth, from, to, message); err != nil {
		return err
	}
	return nil
}

func SendVerificationEmail(email, token string) error {
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}
	verificationLink := fmt.Sprintf(config.Server.Url+config.Server.Port+"/users/verify-email?email=%s&token=%s", email, token)
	subject := "Verify Your Email Address"
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Email Verification</title>
		</head>
		<body>
			<p>Hello,</p>
			<p>Thank you for registering with Loan Manager! Please click the link below to verify your email address:</p>
			<p><a href="%s">Verify Email</a></p>
			<p>If you did not register for this account, you can ignore this email.</p>
			<p>Best regards,<br>The Loan Manager Team</p>
		</body>
		</html>`, verificationLink)
	return sendEmail(subject, body, []string{email})
}

func SendPasswordResetEmail(email, token string) error {
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	resetLink := fmt.Sprintf("%s%s/users/password-update?email=%s&token=%s", config.Server.Url, config.Server.Port, email, token)
	subject := "Reset Your Password"
	body := fmt.Sprintf(`
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Password Reset</title>
        </head>
        <body>
            <p>Hello,</p>
            <p>We received a request to reset your password. Please click the link below to choose a new password:</p>
            <p><a href="%s">Reset Password</a></p>
            <p>If you did not request a password reset, you can ignore this email.</p>
            <p>Best regards,<br>The Loan Manager Team</p>
        </body>
        </html>`, resetLink)
	return sendEmail(subject, body, []string{email})
}

package email

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

// Sender wraps the Resend client and the verified "from" address - one per worker
type Sender struct {
	client		*resend.Client
	from 		string
}

// NewSender builds the Resend-Backed sender (apiKey + from config)
func NewSender(apiKey, from string) *Sender {
	return &Sender{
		client: resend.NewClient(apiKey),
		from: from,
	}
}

// Send picks content for the job's kind and delivers it through Resend
func (s *Sender) Send (job EmailJob) error {
	var subject, html string

	switch job.Kind {
	case KindVerification:
			subject = "Verify Your LocalVault email"
			html = fmt.Sprintf("Hi %s, confirm your email: <a href=%q>Verify</a>", job.Name, job.URL)
	case KindPasswordReset:
		subject = "Reset your LocalVault password"
		html = fmt.Sprintf("Hi %s, reset your password: <a href=%q>Reset</a>", job.Name, job.URL)
	default:
		return fmt.Errorf("unknown email kind: %s", job.Kind) // guards against malformed messages
	}

	// SendEmailRequest is Resend's payment; Html is the rendered body
	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From: 		s.from,
		To: 		[]string{string(job.To)},
		Subject: 	subject,
		Html: 		html,	
	})
	return err
}

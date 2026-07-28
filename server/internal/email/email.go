package email

import (
	"github.com/resend/resend-go/v2"
)

// Sender wraps the Resend client and the verified "from" address - one per worker
type Sender struct {
	client		*resend.Client
	from 		string
	renderer 	*Renderer
}

// NewSender builds the Resend-Backed sender (apiKey + from config)
func NewSender(apiKey, from string) *Sender {
	// templates are compliled into the binary; a parse failure is a build time bug
	renderer, err := NewRenderer()

	if err != nil {
		panic(err)
	}

	return &Sender{
		client: resend.NewClient(apiKey),
		from: from,
		renderer: renderer,
	}
}

// Send picks content for the job's kind and delivers it through Resend
func (s *Sender) Send(job EmailJob) error {
	subject, html, err := s.renderer.Render(job)
	if err != nil {
		return err
	}

	// SendEmailRequest is Resend's payment; Html is the rendered body
	_, err = s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from,
		To:      []string{job.To}, // job.To is a plain string after Task 1
		Subject: subject,
		Html:    html,
	})
	return err
}

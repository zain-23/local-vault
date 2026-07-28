package email

// EmailKind names which email to send - the worker switches on this value
type EmailKind string

const (
	KindVerification 	EmailKind = "verification"
	KindPasswordReset 	EmailKind = "password_reset"
	KindSendMagicLink	EmailKind = "send_magic_link"
)

// EmailJob is the message we put on the queue - json-encoded in the body
type EmailJob struct {
	Kind	EmailKind	`json:"kind"`
	To		string	`json:"to"`
	Name	string	`json:"name"`
	URL		string	`json:"url"`
}

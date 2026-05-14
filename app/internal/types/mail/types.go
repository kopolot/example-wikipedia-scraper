package mail

type Mail struct {
	To                  []string
	Bcc                 []string
	Cc                  []string
	Subject             string
	Body                string
	AttachmentsFileName []string
}

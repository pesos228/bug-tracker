package notification

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/templates"
	"github.com/wneessen/go-mail"
)

type Notifier interface {
	NotifyAboutNewTask(user *domain.User, task *domain.Task)
}

type emailNotifier struct {
	Host      string
	Port      int
	Username  string
	Password  string
	From      string
	publicURL string
}

type newTaskEmailData struct {
	FirstName string
	SoftName  string
	TaskURL   string
}

func (e *emailNotifier) NotifyAboutNewTask(user *domain.User, task *domain.Task) {
	m := mail.NewMsg()
	if err := m.From(e.From); err != nil {
		log.Printf("EMAIL_NOTIFICATION_ERROR: couldn't identify sender: %v", err)
		return
	}
	if err := m.To(user.Email); err != nil {
		log.Printf("EMAIL_NOTIFICATION_ERROR: couldn't identify the recipient: %v", err)
		return
	}

	tmpl, err := template.ParseFS(templates.Files, "new_task_email.html")
	if err != nil {
		log.Printf("EMAIL_NOTIFICATION_ERROR: couldn't parse the template: %v", err)
		return
	}

	data := newTaskEmailData{
		FirstName: user.FirstName,
		SoftName:  task.SoftName,
		TaskURL:   fmt.Sprintf("%s/tasks/%s", e.publicURL, task.ID),
	}

	var bodyBuffer bytes.Buffer
	if err := tmpl.Execute(&bodyBuffer, data); err != nil {
		log.Printf("EMAIL_NOTIFICATION_ERROR: couldn't execute the template: %v", err)
		return
	}

	m.Subject(fmt.Sprintf("Новая задача: %s", task.SoftName))
	m.SetBodyString(mail.TypeTextHTML, bodyBuffer.String())

	client, err := mail.NewClient(
		e.Host,
		mail.WithPort(e.Port),
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithUsername(e.Username),
		mail.WithPassword(e.Password),
	)
	if err != nil {
		log.Printf("EMAIL_NOTIFICATION_ERROR: couldn't create SMTP client: %v", err)
		return
	}

	if err := client.DialAndSend(m); err != nil {
		log.Printf("EMAIL_NOTIFICATION_ERROR: couldn't send email: %v", err)
		return
	}
}

func NewEmailNotifier(host string, port int, username, password, from, publicURL string) Notifier {
	return &emailNotifier{
		Host:      host,
		Port:      port,
		Username:  username,
		Password:  password,
		From:      from,
		publicURL: publicURL,
	}
}

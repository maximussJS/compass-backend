package services

import (
	fx_utils "compass-backend/pkg/common/fx"
	"compass-backend/pkg/sender/lib"
	"crypto/tls"
	"fmt"
	"go.uber.org/fx"
	gomail "gopkg.in/mail.v2"
)

type IMailService interface {
	Send(to string, subject string, body string) error
}

type mailServiceParams struct {
	fx.In

	Env lib.IEnv
}

type mailService struct {
	from   string
	dialer *gomail.Dialer
}

func FxMailService() fx.Option {
	return fx_utils.AsProvider(newMailService, new(IMailService))
}

func newMailService(params mailServiceParams) IMailService {
	dialer := gomail.NewDialer(
		params.Env.GetMailHost(),
		params.Env.GetMailPort(),
		params.Env.GetMailUsername(),
		params.Env.GetMailPassword(),
	)

	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &mailService{
		dialer: dialer,
		from:   params.Env.GetMailUsername(),
	}
}

func (s mailService) Send(to string, subject string, body string) error {
	message := gomail.NewMessage()

	message.SetHeader("From", s.from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	message.SetBody("text/html", body)

	if err := s.dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

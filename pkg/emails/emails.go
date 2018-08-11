package emails

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"sync"

	"github.com/yakovzaytsev/ysz/pkg/ysz"
)

func Send(email, fromEmail, authEmail, EMAIL_PASSWORD, subj, body string) error {
	from := mail.Address{"", fromEmail}
	to := mail.Address{"", email}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html"
	headers["charset"] = "UTF-8"

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := "smtp.gmail.com:465"

	host, _, _ := net.SplitHostPort(servername)

	// Set up authentication information.
	auth := smtp.PlainAuth("", authEmail, EMAIL_PASSWORD, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()
	return nil
}

type EmailVerificationOrder struct {
	Email string
	Hash  string
}

var emailVerificationOrders = map[string]EmailVerificationOrder{}
var emailVerificationOrdersLock = sync.RWMutex{}

func saveEmailVerificationOrder(order EmailVerificationOrder) string {
	token := ysz.RandSeq(4)
	emailVerificationOrdersLock.Lock()
	defer emailVerificationOrdersLock.Unlock()
	emailVerificationOrders[token] = order
	return token
}

func getAndRmEmailVerificationOrder(token string) EmailVerificationOrder {
	emailVerificationOrdersLock.RLock()
	defer emailVerificationOrdersLock.RUnlock()
	o := emailVerificationOrders[token]
	delete(emailVerificationOrders, token)
	return o
}

// returns hash
// send it
func VerifyEmail(email string) string {
	hash := ysz.RandSeq(4)

	order := EmailVerificationOrder{Email: email, Hash: hash}
	token := saveEmailVerificationOrder(order)

	return token
}

func CheckEmailCode(token, got_hash string) bool {
	o := getAndRmEmailVerificationOrder(token)

	if got_hash != o.Hash {
		log.Printf("checkEmailCode: expected %s got %s", o.Hash, got_hash)

		return false
	}

	return true
}

package emails

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"sync"

	"github.com/yakovzaytsev/ysz/pkg/ysz"
)

func Send(email, fromEmail, authEmail, EMAIL_PASSWORD, subj, body string) {
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
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()
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

func verifyEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("e")
	if len(email) == 0 {
		log.Print("verifyEmail: no email")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	log.Printf("verifyEmail: email: %s", email)

	hash := ysz.RandSeq(4)

	// TODO send hash in email
	Send(email, "yakov@freshGDPR.com", os.Getenv("GMAIL_AUTH_EMAIL"), os.Getenv("EMAIL_PASSWORD"), "Your freshGDPR verification code", fmt.Sprintf(`<p>Hello,</p>
<p>please enter this verification code %s</p>
`, hash))

	order := EmailVerificationOrder{Email: email, Hash: hash}
	token := saveEmailVerificationOrder(order)

	// w.WriteHeader(http.StatusForbidden)
	// return

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
{
	"t": "%s"
}`, token)))
}

func checkEmailCode(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("t")
	if len(token) == 0 {
		log.Print("checkEmailCode: no token")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	log.Printf("checkEmailCode: token: %s", token)

	got_hash := r.URL.Query().Get("h")
	if len(got_hash) == 0 {
		log.Print("checkEmailCode: no hash")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	log.Printf("checkEmailCode: hash: %s", got_hash)

	o := getAndRmEmailVerificationOrder(token)

	if got_hash != o.Hash {
		log.Printf("checkEmailCode: expected %s got %s", o.Hash, got_hash)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}

package emails

import (
	"sync"

	"github.com/yakovzaytsev/ysz"
)

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

	hash := freshGDPR.RandSeq(4)

	// TODO send hash in email
	myemail.Send(email, "yakov@freshGDPR.com", os.Getenv("GMAIL_AUTH_EMAIL"), os.Getenv("EMAIL_PASSWORD"), "Your freshGDPR verification code", fmt.Sprintf(`<p>Hello,</p>
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
	log.Print("checkEmailCode: token: %s", token)

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

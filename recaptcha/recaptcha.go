package recaptcha

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	SECRET     = os.Getenv("SECRET_KEY")
	URL_GOOGLE = "https://www.google.com/recaptcha/api/siteverify"
)

type Response struct {
	Success     bool
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string
	ErrorCodes  interface{} `json:"error-codes"`
}

func verify(token string) bool {
	// data := url.Values{
	// 	"secret":   {SECRET},
	// 	"response": {token},
	// }
	// body := strings.NewReader(data.Encode())

	var r http.Request
	r.ParseForm()
	r.Form.Add("secret", SECRET)
	r.Form.Add("response", token)
	body := strings.NewReader(r.Form.Encode())
	resp, err := http.Post(URL_GOOGLE, "application/x-www-form-urlencoded", body)
	if err != nil {
		log.Errorf("verify: %s", err)
		return false
	}
	var res Response
	err = unmarshal(resp.Body, &res)
	log.Println(res)
	return res.Success
}

func unmarshal(body io.Reader, v interface{}) error {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Errorf("ioutil.ReadAll: %s", err)
		return err
	}
	bodyBytes = bytes.TrimPrefix(bodyBytes, []byte("\xef\xbb\xbf"))
	err = json.Unmarshal(bodyBytes, &v)
	log.Println(v)

	if err != nil {
		log.Errorf("unmarshal: %s", err)
		return err
	}
	return nil
}

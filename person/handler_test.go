package person_test

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type testFormValue struct {
	data []string
	code int
}

func testus(urlStr string, usecases []testFormValue, client *http.Client, t *testing.T) {
	for _, uc := range usecases {
		form := url.Values{}
		peak := len(uc.data) - 1
		for i := 0; i < peak; i++ {
			form.Add(uc.data[i], uc.data[i+1])
		}

		req, err := http.NewRequest("POST", urlStr, strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		req.PostForm = form
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != uc.code {
			t.Errorf("Expected %d got %d", uc.code, res.StatusCode)
		}
	}
}

//
func TestHandlerPostInvalidData(t *testing.T) {

	t.Log("Posting invalid data expects response with 412 status code")
	usecases := []testFormValue{
		{[]string{"", ""}, 412},
		{[]string{"longname", "aaa", "phone", "072153", "email", "gb@mob.ro", "is_male", "1", "address", "", "is_client", "0", "is_contractor", "1"}, 412},
		{[]string{"is_contactor", "one"}, 412},
		{[]string{"is_client", "one"}, 412},
		{[]string{"is_client", ""}, 412},
		{[]string{"donotexists", "42", "nilkey", "24"}, 412},
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	urlStr := "https://localhost:3000/person"
	testus(urlStr, usecases, client, t)

}

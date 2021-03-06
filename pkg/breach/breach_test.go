package breach

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Sleep(t *testing.T) {
	cases := []struct {
		seconds      string
		expectingErr bool
	}{
		{"10E99999", true},
		{"1", false},
	}

	for i, c := range cases {
		fmt.Printf("Running case %d\n", i+1)
		err := Sleep(c.seconds)

		if c.expectingErr {
			if err == nil {
				t.Errorf("Expecting an error in case %d\n", i+1)
			} else {
				fmt.Println(err)
			}
			continue
		}
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_VerifyResponse(t *testing.T) {
	cases := []struct {
		status       int
		expectingErr bool
	}{
		{400, true},
		{200, false},
		{429, false},
	}

	for i, c := range cases {
		fmt.Printf("Running case %d\n", i+1)
		_, err := VerifyResponse(c.status)

		if c.expectingErr {
			if err == nil {
				t.Errorf("Expecting an error in case %d\n", i+1)
			} else {
				fmt.Println(err)
			}
			continue
		}
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_VerifyAndRetry(t *testing.T) {
	cases := []struct {
		res          http.Response
		expectingErr bool
	}{
		{http.Response{StatusCode: 400}, true},
		{http.Response{StatusCode: 200}, false},
		{http.Response{StatusCode: 429}, true},
	}

	for i, c := range cases {
		fmt.Printf("Running case %d\n", i+1)
		_, err := VerifyAndRetry(&c.res)

		if c.expectingErr {
			if err == nil {
				t.Errorf("Expecting an error in case %d\n", i+1)
			} else {
				fmt.Println(err)
			}
			continue
		}
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "HIBPwned servers should be happy now")
	}))
	defer ts.Close()

	cases := []struct {
		URL          string
		expectingErr bool
	}{
		{":", true},
		{ts.URL, false},
	}

	for i, c := range cases {
		fmt.Printf("Running case %d\n", i+1)
		res, err := Get(c.URL)

		if c.expectingErr {
			if err == nil {
				t.Errorf("Expecting an error in case %d\n", i+1)
			} else {
				fmt.Println(err)
			}
			continue
		}
		if err != nil {
			t.Error(err)
		}
		_, err = ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Error(err)
		}
	}
}

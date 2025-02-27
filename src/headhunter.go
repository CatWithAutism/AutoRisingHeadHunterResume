package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

type HeadHunterCookies struct {
	UserAgent string
	Xsrf      string
	Token     string
}

func RaiseResume(httpClient *http.Client, hunter *HeadHunterCookies) error {
	req, err := http.NewRequest("POST", "https://hh.ru/shards/resume/batch_update", nil)
	if err != nil {
		return nil
	}

	httpClient = &http.Client{}

	req.Header.Set("x-xsrftoken", hunter.Xsrf)
	req.Header.Set("cookie", fmt.Sprintf("_xsrf=%s; hhtoken=%s;", hunter.Xsrf, hunter.Token))
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resume update failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func AuthorizeHeadHunter(username string, password string) (*http.Client, *HeadHunterCookies, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, nil, err
	}

	hunter := HeadHunterCookies{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
	}

	err = hunter.getDefaultCookies()
	if err != nil {
		return nil, nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	formData := url.Values{}
	formData.Set("_xsrf", hunter.Xsrf)
	formData.Set("failUrl", "/account/login?backurl=%2F")
	formData.Set("accountType", "APPLICANT")
	formData.Set("remember", "yes")
	formData.Set("username", username)
	formData.Set("password", password)
	formData.Set("isBot", "false")
	formData.Set("captchaText", "")

	req, err := http.NewRequest("POST", "https://hh.ru/account/login?backurl=%2F", bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("cookie", fmt.Sprintf("_xsrf=%s; hhtoken=%s;", hunter.Xsrf, hunter.Token))
	req.Header.Set("x-xsrftoken", hunter.Xsrf)

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("authorization failed with status code: %d", resp.StatusCode)
	}

	cookies := resp.Header["Set-Cookie"]
	var cookie string
	for _, c := range cookies {
		cookie += c + "\n"
	}

	xsrfRegex := regexp.MustCompile(`_xsrf=.+?;`)
	hhtokenRegex := regexp.MustCompile(`hhtoken=.+?;`)

	xsrfMatch := xsrfRegex.FindString(cookie)
	if xsrfMatch == "" {
		panic("can't authorize")
	}
	hunter.Xsrf = xsrfMatch[6 : len(xsrfMatch)-1]

	hhtokenMatch := hhtokenRegex.FindString(cookie)
	if hhtokenMatch == "" {
		panic("can't authorize")
	}
	hunter.Token = hhtokenMatch[8 : len(hhtokenMatch)-1]

	return client, &hunter, nil
}

func (h *HeadHunterCookies) getDefaultCookies() error {
	url := "https://hh.ru/"
	headers := map[string]string{"user-agent": h.UserAgent}
	response, err := doRequest(http.MethodHead, url, headers, nil)
	if err != nil {
		return err
	}

	cookies := response.Header["Set-Cookie"]
	var cookie string
	for _, c := range cookies {
		cookie += c + "\n"
	}

	xsrfRegex := regexp.MustCompile(`_xsrf=.+?;`)
	hhtokenRegex := regexp.MustCompile(`hhtoken=.+?;`)

	xsrfMatch := xsrfRegex.FindString(cookie)
	if xsrfMatch == "" {
		return fmt.Errorf("xsrf cookie not found")
	}
	h.Xsrf = xsrfMatch[6 : len(xsrfMatch)-1]

	hhtokenMatch := hhtokenRegex.FindString(cookie)
	if hhtokenMatch == "" {
		return fmt.Errorf("hhtoken cookie not found")
	}
	h.Token = hhtokenMatch[8 : len(hhtokenMatch)-1]

	return nil
}

func doRequest(method string, url string, headers map[string]string, body *bytes.Buffer) (*http.Response, error) {
	client := &http.Client{}
	var req *http.Request
	var err error
	if body == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, body)
	}
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

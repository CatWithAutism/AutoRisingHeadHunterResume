package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type HeadHunterClient struct {
	httpClient *http.Client
	notify     func(message string) error
}

func (headhunter *HeadHunterClient) RaiseResume() error {
	const resumeUpdated = "Resume updated."

	req, err := http.NewRequest("POST", "https://hh.ru/shards/resume/batch_update", nil)
	if err != nil {
		return err
	}

	xsrf := GetSpecifiedCookie(headhunter.httpClient, "https", "hh.ru", "_xsrf")
	req.Header.Set("x-xsrftoken", xsrf)

	resp, err := headhunter.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resume update failed with status code: %d", resp.StatusCode)
	}

	if headhunter.notify != nil {
		err := headhunter.notify(resumeUpdated)
		if err != nil {
			println(err.Error())
		}
	}

	err = headhunter.fetchDefaultCookies()
	if err != nil {
		return err
	}

	return nil
}

func AuthorizeHeadHunter(username string, password string, notify func(message string) error) (*HeadHunterClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	headhunter := HeadHunterClient{
		httpClient: &http.Client{Jar: jar},
		notify:     notify,
	}

	err = headhunter.fetchDefaultCookies()
	if err != nil {
		return nil, err
	}

	xsrf := GetSpecifiedCookie(headhunter.httpClient, "https", "hh.ru", "_xsrf")

	formData := url.Values{}
	formData.Set("_xsrf", xsrf)
	formData.Set("failUrl", "/account/login?backurl=%2F")
	formData.Set("accountType", "APPLICANT")
	formData.Set("remember", "yes")
	formData.Set("username", username)
	formData.Set("password", password)
	formData.Set("isBot", "false")
	formData.Set("captchaText", "")

	resp, err := headhunter.httpClient.Post("https://hh.ru/account/login?backurl=%2F", "application/x-www-form-urlencoded", bytes.NewBufferString(formData.Encode()))

	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authorization failed with status code: %d", resp.StatusCode)
	}

	if headhunter.notify != nil {
		err := headhunter.notify("We have authorized on hh.ru")
		if err != nil {
			println(err.Error())
		}
	}

	return &headhunter, nil
}

func (headhunter *HeadHunterClient) fetchDefaultCookies() error {

	resp, err := headhunter.httpClient.Head("https://hh.ru/")
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch default cookies with status code: %d", resp.StatusCode)
	}

	return nil
}

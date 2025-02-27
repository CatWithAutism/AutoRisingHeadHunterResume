package main

import (
	"net/http"
	"net/url"
)

func GetSpecifiedCookie(httpClient *http.Client, scheme string, host string, name string) string {
	cookies := httpClient.Jar.Cookies(&url.URL{
		Scheme: scheme,
		Host:   host,
	})

	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return ""
}

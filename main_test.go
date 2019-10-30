package main

import (
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	_, err := http.Get("https://maxence-sso.workos.dev/login")
	if err != nil {
		t.Log(err)
	}
}

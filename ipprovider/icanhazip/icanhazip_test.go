package icanhazip

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIcanhazipNew(t *testing.T) {
	expectedURL := "https://ipv4.icanhazip.com"
	ifc := New(&http.Client{})
	ifcOriginal := ifc.(*icanhazip)
	if ifcOriginal.url != expectedURL {
		t.Errorf("URL of icanhazip should be %q, but got %q", expectedURL, ifcOriginal.url)
		return
	}
}

func TestForceIPV6(t *testing.T) {
	expectedv6URL := "https://ipv6.icanhazip.com"
	ifc := New(&http.Client{})
	ifc.ForceIPV6()
	ifcOriginal := ifc.(*icanhazip)
	if ifcOriginal.url != expectedv6URL {
		t.Errorf("URL of icanhazip should be %q, but got %q", expectedv6URL, ifcOriginal.url)
		return
	}
}

func TestIcanhazipSuccess(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(`45.45.45.45`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ifc := &icanhazip{
		c:   &http.Client{},
		url: server.URL,
	}

	ip, errGet := ifc.GetIP()
	if errGet != nil {
		t.Errorf("Got error: %s", errGet.Error())
		return
	}

	if ip != "45.45.45.45" {
		t.Errorf("Incorrect IP value. Got %q, but should be %q", ip, "45.45.45.45")
		return
	}
}

func TestIcanhazipNotSuccessCode(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(429)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ifc := &icanhazip{
		c:   &http.Client{},
		url: server.URL,
	}

	_, errGet := ifc.GetIP()
	if errGet == nil {
		t.Errorf("Should be error, but is success")
		return
	}

	if errGet.Error() != "icanhazip: Status code is not in success range: 429" {
		t.Error("Error was, but not about status code")
		return
	}
}

func TestIcanhazipFailedRead(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1")
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ifc := &icanhazip{
		c:   &http.Client{},
		url: server.URL,
	}

	_, errGet := ifc.GetIP()
	if errGet == nil {
		t.Errorf("Should be error, but is success")
		return
	}
}

func TestIcanhazipFailedOnGet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`something wrong`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ifc := &icanhazip{
		c:   &http.Client{},
		url: "http://127.0.0.1:1234",
	}

	_, errGet := ifc.GetIP()
	if errGet == nil {
		t.Errorf("Should be error, but is success")
		return
	}

	if !isMatchingErrorMessage(errGet.Error(), "icanhazip", "connection refused") {
		t.Errorf("Error was, but not related to the request fail: %v", errGet.Error())
		return
	}
}

func isMatchingErrorMessage(message string, prefix, suffix string) bool {
	return strings.HasPrefix(message, prefix) && strings.HasSuffix(message, suffix)
}

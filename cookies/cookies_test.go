package cookies

import (
	"testing"
)

type nameValue map[string]string

var (
	cookiesString = "token=tokenString;murka= murkaString; zhmurka=aga"
)

// TestNewGet checks cookies string parsing and get by key
func TestNewGet(t *testing.T) {
	cookies := New(cookiesString)
	if len(cookies) != 3 {
		t.Fatalf("Expected slice size %d, received %d", 3, len(cookies))
	}

	expectedNewResult := nameValue{
		"token":   "tokenString",
		"murka":   " murkaString",
		"zhmurka": "aga",
	}

	for _, cookie := range cookies {
		if expectedNewResult[cookie.Name] != cookie.Value {
			t.Fatalf("Cookie with name %s, expected: %s, received: %s",
				cookie.Name, expectedNewResult[cookie.Name], cookie.Value)
		}
	}
}

// TestGet tests getting cookie
func TestGet(t *testing.T) {
	cookies := New(cookiesString)
	expectedGetResult := nameValue{
		"token":   "tokenString",
		"murka":   " murkaString",
		"zhmurka": "aga",
		"q":       "",
		"":        "",
	}

	for key, expectedValue := range expectedGetResult {
		receivedCookie := cookies.Get(key)

		if expectedValue == "" {
			if receivedCookie != nil {
				t.Fatalf("Get by name %s, expected value: %s, received not nil cookie: %s",
					key, expectedValue, receivedCookie.Value)
			}
			continue
		}

		if receivedCookie.Value != expectedValue {
			t.Fatalf("Get by name %s, expected value: %s, received value: %s",
				key, expectedValue, receivedCookie.Value)
		}
	}
}

// TestGetValue tests getting cookie value
func TestGetValue(t *testing.T) {
	cookies := New(cookiesString)
	expectedGetResult := nameValue{
		"token":   "tokenString",
		"murka":   " murkaString",
		"zhmurka": "aga",
		"q":       "",
		"":        "",
	}

	for key, expectedValue := range expectedGetResult {
		receivedCookie := cookies.GetValue(key)
		if receivedCookie != expectedValue {
			t.Fatalf("Get by name %s, expected value: %s, received value: %s",
				key, expectedValue, receivedCookie)
		}
	}
}

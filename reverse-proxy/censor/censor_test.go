package censor

import (
	"encoding/json"
	"testing"
)

func TestCensor(t *testing.T) {
	type Message struct {
		Email    string
		Password string
	}
	message := Message{Email: "ab@mail.com", Password: "123"}
	body, _ := json.Marshal(message)
	testCases := []struct {
		desc    string
		body    []byte
		content Message
	}{
		{
			desc:    "Blocked POST request",
			body:    body,
			content: message,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			censoredBody := Censor([]string{"Password"})(string(tC.body))
			parsedBody := Message{}
			json.Unmarshal([]byte(censoredBody), &parsedBody)

			if parsedBody.Email != tC.content.Email {
				t.Errorf("expected email not to change %v, got %v", parsedBody.Email, tC.content.Email)
			}

			if parsedBody.Password != "********" {
				t.Errorf("expected Password to be censored %v, got %v", parsedBody.Password, tC.content.Password)
			}
		})
	}
}

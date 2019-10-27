package command

import (
	"fmt"
	"testing"
)

func TestFillTemplateMessage(t *testing.T) {
	var tpl = "The {{.OwnerNumber}}th {{.Owner}}'s {{.AnimalNumber}}th {{.Animal}}'s sick!"

	t.Run("mix of values", func(t *testing.T) {
		t.Parallel()

		filler := map[string]interface{}{"OwnerNumber": 6,
			"AnimalNumber": 7, "Animal": "Sheep", "Owner": "Sheikh"}
		expectedMsg := fmt.Sprintf("The %dth %s's %dth %s's sick!",
			filler["OwnerNumber"], filler["Owner"], filler["AnimalNumber"], filler["Animal"])

		res, err := FillTemplateMessage(tpl, filler)
		if err != nil {
			t.Fatalf("Call to func returned an unexpected error: %v", err)
		}
		if res != expectedMsg {
			t.Errorf("Expected returned value to be %s, received %s", expectedMsg, res)
		}
	})

	t.Run("inconsistent fields", func(t *testing.T) {
		t.Parallel()

		testCases := []map[string]interface{}{
			// missing fields
			{"OwnerNumber": 17, "Animal": "goat"},
			// extra fields
			{"OwnerNumber": 6, "AnimalNumber": 7,
				"Animal": "Sheep", "Owner": "Sheikh", "origin": "shrute farms"},
			// missing and extra fields
			{"OwnerNumber": 6, "Owner": "Sheikh", "origin": "shrute farms"},
		}

		for _, tc := range testCases {
			if _, err := FillTemplateMessage(tpl, tc); err != nil {
				t.Fatalf("Call to func with %v returned an unexpected error: %v", tc, err)
			}
		}
	})
}

func TestCoalesce(t *testing.T) {
	t.Parallel()

	const (
		whitespace = "      "
		arg        = "house lannister"
	)

	testCases := []struct {
		args     []string
		expected string
	}{
		{[]string{}, ""},
		{[]string{""}, ""},
		{[]string{"", "", "", ""}, ""},
		{[]string{whitespace, "hello-world"}, whitespace},
		{[]string{"", "", whitespace, "", arg}, whitespace},
		{[]string{"", "", arg}, arg},
		{[]string{arg, "foo", whitespace, "gelatto"}, arg},
		{[]string{"", "", arg, "", "hello", "world"}, arg},
	}

	for _, tc := range testCases {
		if res := Coalesce(tc.args...); res != tc.expected {
			t.Errorf("Expected %s, received %s for test case %v", tc.expected, res, tc.args)
		}
	}
}

func TestValidateEmailID(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		testCases := []string{
			"email@domain.com", "firstname.lastname@foobar.com", "email@subdomain.domain.com",
			"firstname+lastname@domain.com", "email@123.123.123.123", "firstname-lastname@domain.com",
			"1234567890@google.com", "email@domain-one.com", "email@domain.co.jp",
		}
		for _, tc := range testCases {
			if err := ValidateEmailID(tc); err != nil {
				t.Errorf("Expected %s to be valid", tc)
			}
		}
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()
		testCases := []string{
			"foobar", "#@%^%#$@#$@#.com", "@domain.com", "Joe Smith <email@domain.com>",
			"email.domain.com", "email@domain@domain.com", "email@-domain.com", "email@domain..com",
		}
		for _, tc := range testCases {
			if err := ValidateEmailID(tc); err == nil {
				t.Errorf("Expected %s to be invalid", tc)
			}
		}
	})
}

package validate

import "testing"

func TestUsername(t *testing.T) {
	cases := []struct {
		in    string
		want  string
		valid bool
	}{
		{"  Alice_99 ", "Alice_99", true},
		{"bob-the-builder", "bob-the-builder", true},
		{"ab", "", false},      // too short
		{"has space", "", false},
		{"bad!char", "", false},
	}
	for _, c := range cases {
		got, err := Username(c.in)
		if c.valid && (err != nil || got != c.want) {
			t.Errorf("Username(%q) = (%q, %v), want (%q, nil)", c.in, got, err, c.want)
		}
		if !c.valid && err == nil {
			t.Errorf("Username(%q) should have failed", c.in)
		}
	}
}

func TestEmail(t *testing.T) {
	if got, err := Email("  Foo@Bar.COM "); err != nil || got != "foo@bar.com" {
		t.Errorf("Email normalize failed: got (%q, %v)", got, err)
	}
	for _, bad := range []string{"no-at", "a@b", "@b.com", "a@b.", "spaces in@x.com"} {
		if _, err := Email(bad); err == nil {
			t.Errorf("Email(%q) should be invalid", bad)
		}
	}
}

func TestPassword(t *testing.T) {
	if err := Password("Password1"); err != nil {
		t.Errorf("valid password rejected: %v", err)
	}
	for _, bad := range []string{"short1A", "alllowercase1", "ALLUPPERCASE1", "NoDigitsHere"} {
		if err := Password(bad); err == nil {
			t.Errorf("Password(%q) should be invalid", bad)
		}
	}
}

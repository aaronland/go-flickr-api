package response

import (
	"strings"
	"testing"
)

func TestUnmarshalCheckLoginJSONResponse(t *testing.T) {

	rsp := `{"user":{"id":"161215698@N03","username":{"_content":"aaronofsfo"},"path_alias":null},"stat":"ok"}`

	fh := strings.NewReader(rsp)

	l, err := UnmarshalCheckLoginJSONResponse(fh)

	if err != nil {
		t.Fatalf("Failed to unmarshal check login response, %v", err)
	}

	if l.User.Id != "161215698@N03" {
		t.Fatalf("Unexpected user ID '%s'", l.User.Id)
	}

	if l.User.Username.Value != "aaronofsfo" {
		t.Fatalf("Unexpected username '%s'", l.User.Username.Value)
	}
}

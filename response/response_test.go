package response

import (
	"strings"
	"testing"
)

func TestUnmarshalResponse(t *testing.T) {

	rsp := `<?xml version="1.0" encoding="utf-8" ?>
<rsp stat="ok"></rsp>`

	fh := strings.NewReader(rsp)

	r, err := UnmarshalResponse(fh)

	if err != nil {
		t.Fatalf("Failed to unmarshal response, %v", err)
	}

	if r.Status != "ok" {
		t.Fatalf("Unexpected status, %s", r.Status)
	}
}

func TestUnmarshalErrorResponse(t *testing.T) {

	rsp := `<?xml version="1.0" encoding="utf-8" ?>
<rsp stat="error"><err code="999" msg="This is an error" /></rsp>`

	fh := strings.NewReader(rsp)

	r, err := UnmarshalResponse(fh)

	if err != nil {
		t.Fatalf("Failed to unmarshal response, %v", err)
	}

	if r.Error == nil {
		t.Fatal("Failed to parse error")
	}
}

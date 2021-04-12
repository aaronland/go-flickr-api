package response

import (
	"strings"
	"testing"
)

func TestUnmarshalUploadResponse(t *testing.T) {

	rsp := `<?xml version="1.0" encoding="utf-8" ?>
<rsp stat="ok">
<photoid secret="b5b7b2d1fc" originalsecret="2d5847f46b">51111590154</photoid>
</rsp>`

	fh := strings.NewReader(rsp)

	up, err := UnmarshalUploadResponse(fh)

	if err != nil {
		t.Fatalf("Failed to unmarshal replace response, %v", err)
	}

	if up.Status != "ok" {
		t.Fatalf("Unexpected status, '%s'", up.Status)
	}

	if up.Photo == nil {
		t.Fatalf("Missing Photo property")
	}

	if up.Photo.Id != 51111590154 {
		t.Fatalf("Unexpected photo ID, '%d'", up.Photo.Id)
	}

	if up.Photo.Secret != "b5b7b2d1fc" {
		t.Fatalf("Unexpected photo secret, '%s'", up.Photo.Secret)
	}

	if up.Photo.OriginalSecret != "2d5847f46b" {
		t.Fatalf("Unexpected photo secret, '%s'", up.Photo.OriginalSecret)
	}

}

func TestUnmarshalUploadTicketResponse(t *testing.T) {

	rsp := `<?xml version="1.0" encoding="utf-8" ?>
<rsp stat="ok">
<ticketid>161192644-72157718911685379</ticketid>
</rsp>`

	fh := strings.NewReader(rsp)

	ut, err := UnmarshalUploadTicketResponse(fh)

	if err != nil {
		t.Fatalf("Failed to unmarshal ticket response, %v", err)
	}

	if ut.TicketId != "161192644-72157718911685379" {
		t.Fatalf("Unexpected ticket Id '%s'", ut.TicketId)
	}
}

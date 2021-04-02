package response

import (
	"encoding/json"
	"io"
)

type CheckTicket struct {
	Uploader *Uploader `json:"uploader"`
}

type Uploader struct {
	Tickets []*UploaderTicket `json:"ticket"`
}

type UploaderTicket struct {
	TicketId string `json:"id"`
	Complete int    `json:"complete"`
	PhotoId  string `json:"photoid"`
	Imported string `json:"imported"`
}

func UnmarshalCheckTicketResponse(fh io.Reader) (*CheckTicket, error) {

	var ct *CheckTicket

	dec := json.NewDecoder(fh)
	err := dec.Decode(&ct)

	if err != nil {
		return nil, err
	}

	return ct, nil
}

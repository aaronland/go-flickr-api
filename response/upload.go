package response

import (
	"encoding/xml"
	"io"
)

type Upload struct {
	XMLName xml.Name `xml:"rsp"`
	Status  string   `xml:"stat,attr"`
	Error   *Error   `xml:"err,omitempty"`
	PhotoId int64    `xml:"photoid"`
}

type Ticket struct {
	XMLName  xml.Name `xml:"rsp"`
	Status   string   `xml:"stat,attr"`
	Error    *Error   `xml:"err,omitempty"`
	TicketId string   `xml:"ticketid"`
}

func UnmarshalUploadResponse(fh io.Reader) (*Upload, error) {

	body, err := io.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	var up *Upload

	err = xml.Unmarshal([]byte(body), &up)

	if err != nil {
		return nil, err
	}

	return up, nil
}

func UnmarshalTicketResponse(fh io.Reader) (*Ticket, error) {

	body, err := io.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	var up *Ticket

	err = xml.Unmarshal([]byte(body), &up)

	if err != nil {
		return nil, err
	}

	return up, nil
}

package response

import (
	"encoding/xml"
	"io"
)

/*

<?xml version="1.0" encoding="utf-8" ?>
<rsp stat="ok">
<photoid>51090254017</photoid>
</rsp>

*/

type Upload struct {
	XMLName xml.Name `xml:"rsp"`
	Status  string   `xml:"stat,attr"`
	Error   *Error   `xml:"err,omitempty"`
	PhotoId int64    `xml:"photoid,omitempty"`
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

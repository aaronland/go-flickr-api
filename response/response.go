package response

import (
	_ "encoding/json"
	"encoding/xml"
	"fmt"
	"io"
)

/*

<?xml version="1.0" encoding="utf-8" ?>
<rsp stat="fail">
	<err code="2" msg="No photo specified" />
</rsp>

*/

type Error struct {
	Code    int    `xml:"code,attr"`
	Message string `xml:"msg,attr"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}

type Response struct {
	XMLName xml.Name `xml:"rsp"`
	Status  string   `xml:"stat,attr"`
	Error   *Error   `xml:err,omitempty"`
}

func UnmarshalResponse(fh io.Reader) (*Response, error) {

	body, err := io.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	var rsp *Response

	err = xml.Unmarshal([]byte(body), &rsp)

	if err != nil {
		return nil, err
	}

	return rsp, nil
}

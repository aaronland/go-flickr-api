package response

import (
	"context"
	"github.com/whosonfirst/go-ioutil"
	"strings"
	"testing"
)

func TestDervivePagination(t *testing.T) {

	ctx := context.Background()
	rsp := `{ "photos": { "page": 1, "pages": "334", "perpage": 3, "total": "1000", 
    "photo": [
      { "id": "51113564244", "owner": "12639178@N07", "secret": "fb9be4c7e4", "server": "65535", "farm": 66, "title": "welcher Weißling? (leider zu weit weg gewesen)", "ispublic": 1, "isfriend": 0, "isfamily": 0 },
      { "id": "51113808616", "owner": "138569268@N04", "secret": "de681b2941", "server": "65535", "farm": 66, "title": "IMG_5633.jpg", "ispublic": 1, "isfriend": 0, "isfamily": 0 },
      { "id": "51114596875", "owner": "30334731@N00", "secret": "74505c2376", "server": "65535", "farm": 66, "title": "St Peter's Church", "ispublic": 1, "isfriend": 0, "isfamily": 0 }
    ] }, "stat": "ok" }`

	r := strings.NewReader(rsp)

	fh, err := ioutil.NewReadSeekCloser(r)

	if err != nil {
		t.Fatalf("Failed to create NewReadSeekCloser, %v", err)
	}

	pg, err := DerivePagination(ctx, fh)

	if err != nil {
		t.Fatalf("Failed to derive pagination, %v", err)
	}

	if pg.Pages != 334 {
		t.Fatalf("Unexpected pages count: %d", pg.Pages)
	}

}

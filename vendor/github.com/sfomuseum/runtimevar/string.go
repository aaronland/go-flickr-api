package runtimevar

import (
	"context"
	"fmt"
	gc "gocloud.dev/runtimevar"
	"gocloud.dev/runtimevar/awsparamstore"
	"github.com/aaronland/go-aws-session"
	_ "log"
	"net/url"
)

func StringVar(ctx context.Context, uri string) (string, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return "", err
	}

	q := u.Query()

	if q.Get("decoder") == "" {
		q.Set("decoder", "string")
		u.RawQuery = q.Encode()
	}

	var v *gc.Variable
	var v_err error
	
	switch u.Scheme {
	case "awsparamstore":

		// https://gocloud.dev/howto/runtimevar/#awsps-ctor
		
		creds := q.Get("credentials")
		region := q.Get("region")
		
		if creds != "" {

			dsn_str := fmt.Sprintf("region=%s credentials=%s", region, creds)
			sess, err := session.NewSessionWithDSN(dsn_str)

			if err != nil {
				return "", err
			}

			v, v_err = awsparamstore.OpenVariable(sess, u.Host, gc.StringDecoder, nil)
		}
		
	default:
		// pass
	}

	if v == nil {
		
		uri = u.String()
		
		v, v_err = gc.OpenVariable(ctx, uri)
	}
	
	if v_err != nil {
		return "", v_err
	}

	defer v.Close()

	snapshot, err := v.Latest(ctx)

	if err != nil {
		return "", err
	}

	return snapshot.Value.(string), nil
}

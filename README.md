# go-flickr-api

Go package for working with the Flickr API

## Important

Work in progress. Documentation to follow.

## Tools

```
$> make cli
go build -mod vendor -o bin/api cmd/api/main.go
go build -mod vendor -o bin/authorize cmd/authorize/main.go
```

### api

```
$> ./bin/api \
	-client-uri 'oauth1://?consumer_key={KEY}&consumer_secret={SECRET}&oauth_token={TOKEN}&oauth_token_secret={SECRET}' \
	-param method=flickr.test.login \

| jq

{
  "user": {
    "id": "123456789@X03",
    "username": {
      "_content": "example"
    },
    "path_alias": null
  },
  "stat": "ok"
}
```

## authorize

```
$> ./bin/authorize \
	-server-uri 'mkcert://localhost:8080' \
	-client-uri 'oauth1://?consumer_key={KEY}&consumer_secret={SECRET}'
	
2021/03/31 22:47:08 Checking whether mkcert is installed. If it is not you may be prompted for your password (in order to install certificate files
2021/03/31 22:47:09 Listening for requests on https://localhost:8080
2021/03/31 22:47:13 Authorize this application https://www.flickr.com/services/oauth/authorize?oauth_token={TOKEN}&perms=read

{"oauth_token":"{TOKEN}","oauth_token_secret":"{SECRET}"}
```

## See also

* https://www.flickr.com/services/api/
* https://www.flickr.com/services/api/auth.oauth.html
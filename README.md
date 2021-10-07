# oclient2

OAuth 2.0 test client

See Medium article:
* OAuth 2.0 in Go

https://levelup.gitconnected.com/oauth-2-0-in-go-846b257d32b4


Note that in this version, the core oauth2 package is now at github.com/exyzzy/oauth2

## To Install:

```
go get github.com/exyzzy/oauth2
go get github.com/exyzzy/oclient2
go install $GOPATH/src/github.com/exyzzy/oclient2
```

## Legacy Notes:

* oauth2.go is the library, services.json is the config file for the services. Everything else in oclient2 is an example of how to use it.
* First you'll need to copy services.json to your client and edit it to match the services for which you have set up api accounts. Look at the curent examples, for these services you will only need to set the client_id and client_secret
For production do not include these in src code, but instead serve them from host env variables. Depending on the api you need you may have to adjust the scope. The redirect_uri is set for localhost, change this to your server when you deploy.

* You may wish to adjust consts: GcPeriod, InitAuthTimeout, MaxState (see oauth2.go)

* See main.go and templates/home.html for an example of how to set up the redirect link and authorization requests.

* See main.go and templates/api.html for an example of how to set up the service api requests.

* Also see github.com/exyzzy/oauth2



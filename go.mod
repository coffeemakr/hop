module github.com/coffeemakr/wedo

go 1.14

require (
	github.com/coffeemakr/go-http-error v1.0.0
	github.com/gorilla/mux v1.7.4
	github.com/spf13/cobra v0.0.6
	github.com/square/go-jose/v3 v3.0.0-20200309005203-643bdf8caec0
	go.mongodb.org/mongo-driver v1.3.1
	golang.org/x/crypto v0.0.0-20200311171314-f7b00557c8c4
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/square/go-jose.v2 v2.4.1
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace github.com/coffeemakr/go-http-error => ../go-http-error

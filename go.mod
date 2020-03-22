module github.com/coffeemakr/wedo

go 1.14

require (
	github.com/coffeemakr/go-http-error v1.0.0
	github.com/go-pg/pg/v9 v9.1.3 // indirect
	github.com/golang/protobuf v1.3.5 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.6.2 // indirect
	github.com/vmihailenco/msgpack/v4 v4.3.10 // indirect
	go.mongodb.org/mongo-driver v1.3.1
	golang.org/x/crypto v0.0.0-20200311171314-f7b00557c8c4
	gopkg.in/square/go-jose.v2 v2.4.1
)

replace github.com/coffeemakr/go-http-error => ../go-http-error

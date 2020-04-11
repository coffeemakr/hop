module github.com/coffeemakr/ruck

go 1.14

require (
	github.com/coffeemakr/go-http-error v1.3.0
	github.com/gorilla/mux v1.7.4
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/square/go-jose/v3 v3.0.0-20200309005203-643bdf8caec0
	go.mongodb.org/mongo-driver v1.3.2
	golang.org/x/crypto v0.0.0-20200406173513-056763e48d71
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace github.com/coffeemakr/go-http-error => ../go-http-error

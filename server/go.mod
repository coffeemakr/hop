module github.com/coffeemakr/ruck/server

go 1.14

require (
	github.com/coffeemakr/go-http-error v1.3.0
	github.com/coffeemakr/ruck v0.0.0
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/mitchellh/mapstructure v1.2.2 // indirect
	github.com/pelletier/go-toml v1.7.0 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.6.3
	github.com/square/go-jose v2.4.1+incompatible // indirect
	github.com/square/go-jose/v3 v3.0.0-20200309005203-643bdf8caec0
	go.mongodb.org/mongo-driver v1.3.2
	golang.org/x/crypto v0.0.0-20200406173513-056763e48d71
	golang.org/x/sys v0.0.0-20200409092240-59c9f1ba88fa // indirect
	gopkg.in/ini.v1 v1.55.0 // indirect
	gopkg.in/square/go-jose.v2 v2.4.1
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/coffeemakr/ruck => ../

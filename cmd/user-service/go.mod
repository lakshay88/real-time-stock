module github.com/lakshay88/real-time-stock/user-service

go 1.19

require (
	github.com/gorilla/mux v1.8.1
	github.com/lakshay88/real-time-stock v0.0.0-00010101000000-000000000000
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/lakshay88/real-time-stock => ../..

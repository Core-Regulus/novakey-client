module novakeyclient

go 1.25.1

replace novakeyclient => .

require (
	github.com/google/uuid v1.6.0
	golang.org/x/crypto v0.41.0
)

require (
	github.com/core-regulus/novakey-types-go v0.1.19
	golang.org/x/sys v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1
)

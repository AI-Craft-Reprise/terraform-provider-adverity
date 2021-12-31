module terraform-provider-adverity

go 1.16

require (
	github.com/fourcast/adverityclient v0.0.1
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/hashicorp/go-version v1.3.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.4.3
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/zclconf/go-cty v1.10.0 // indirect
	golang.org/x/net v0.0.0-20210326060303-6b1517762897 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/fourcast/adverityclient v0.0.1 => ./adverity/adverityclient

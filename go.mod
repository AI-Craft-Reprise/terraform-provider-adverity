module terraform-provider-adverity

go 1.13

require github.com/hashicorp/terraform-plugin-sdk v1.17.2

require (
	github.com/fourcast/adverityclient v0.0.1
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/zclconf/go-cty v1.9.1 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20210502180810-71e4cd670f79 // indirect
)

replace github.com/fourcast/adverityclient v0.0.1 => ./adverity/adverityclient

module terraform-provider-adverity

go 1.13

require github.com/hashicorp/terraform-plugin-sdk v1.17.2

require (
	github.com/fourcast/adverityclient v0.0.1
	github.com/hashicorp/terraform-plugin-docs v0.5.0 // indirect
)

replace github.com/fourcast/adverityclient v0.0.1 => ./adverity/adverityclient

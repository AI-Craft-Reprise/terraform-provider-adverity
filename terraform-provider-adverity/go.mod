module adverityprovider

go 1.13

require (
	example.com/adverityclient v0.0.0
	github.com/hashicorp/terraform-plugin-sdk v1.7.0
)

replace example.com/adverityclient => ./adverityclient

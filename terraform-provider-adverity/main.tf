provider "adverity" {

  instance_url = "sdjsjd"
  token =  "smhdfbshd"
}

resource "adverity_connection" "test" {
  name = "datatap.adverity.com"
  stack = 1
  connection_type_id = 1
}
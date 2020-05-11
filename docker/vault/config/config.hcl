api_addr= "http://127.0.0.1:8200"
ui= true
disable_mlock = true
storage "consul" {
  address = "consul:8500"
  path    = "vault"
}
resource "null_resource" "this" {
  triggers = {
    a_string = random_string.random.result
  }
}

resource "random_string" "random" {
  length           = 16
  special          = true
  override_special = "/@£$"
}
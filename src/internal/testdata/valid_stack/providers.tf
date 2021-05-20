terraform {
  required_providers {
    null = {
      source = "hashicorp/null"
      version = "3.1.0"
    }
    random = {
      source = "hashicorp/random"
      version = "3.1.0"
    }
  }
}

provider "null" {
  # Configuration options
}

provider "random" {
  # Configuration options
}
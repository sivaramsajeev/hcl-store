provider "student" {
  address = "localhost"
  port = "8888"
  token = "secret"
}

resource "student" "gopher" {
  name = "Gopher"
  description = "The cute Gopher"
  subjects = [
    "Go",
    "Kubernetes",
    "Docker",
    "Helm",
    "Terraform"
  ]
}


resource "student" "rustacean" {
  name = "Rustacean"
  description = "The tough Ferris"
  subjects = [
    "Rust",
    "Linux",
    "Windows",
    "AWS"
  ]
}

resource "student" "pythonian" {
  name = "Pythonian"
  description = "The crazy Python"
  subjects = [
    "Python",
    "AI/ML",
    "Datascience",
  ]
}
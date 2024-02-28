variable "gcp_project" {
  type        = string
}

variable "region" {
  type    = string
  default = "us-west1"
}

variable "deletion_protection" {
  type    = bool
  default = true
}


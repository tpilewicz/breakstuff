locals {
  default_name = "${var.environment}-${var.component}"
  tags = {
    Environment = var.environment
    Component   = var.component
  }
}

variable "environment" {}
variable "component" {}

variable "api" {}
variable "api_resource" {}
variable "http_method" {}

variable "lambda_src" {}

variable "funes_table" {}

variable "nb_rows" {
  default = 20
}
variable "nb_cols" {
  default = 20
}

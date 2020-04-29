locals {
  component    = "clausius"
  default_name = "${var.environment}-${local.component}"

  tags = {
    Environment = var.environment
    Component   = local.component
  }
}

variable "environment" {}

variable "vpc_id" {}
variable "funes_subnets" {}
variable "funes_sg_id" {}
variable "funes_url" {}

variable "get_grid_resource" {}
variable "set_cell_resource" {}

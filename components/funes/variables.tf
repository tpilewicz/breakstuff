locals {
  component    = "funes"
  default_name = "${var.environment}-${local.component}"
}

variable environment {}

variable node_type {}
variable node_count {}

variable vpc_id {}
variable subnet_cidr_block {}

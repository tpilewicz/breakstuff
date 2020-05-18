locals {
  component    = "clausius"
  default_name = "${var.environment}-${local.component}"

  troublemaker_name = "${local.default_name}-troublemaker"
  troublemaker_src = "../../components/clausius/src/troublemaker/main"

  nb_rows = 20
  nb_cols = 20

  tags = {
    Environment = var.environment
    Component   = local.component
  }
}

variable "environment" {}

variable "funes_table" {
}

variable "get_grid_resource" {}
variable "set_cell_resource" {}

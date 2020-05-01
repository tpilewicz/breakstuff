locals {
  component    = "show"
  default_name = "${var.environment}-${local.component}"
  tags = {
    Environment = var.environment
    Component   = local.component
  }

  index_template = "index_template.html"
  css_file = "style.css"
  ok_file = "ok.gif"
  broken_file = "broken.gif"

  rendered_index = templatefile(
    "../../components/show/assets/${local.index_template}",
    {
      api_invoke_url = var.api_invoke_url
      get_grid_resource = var.get_grid_resource
      set_cell_resource = var.set_cell_resource
      refresh_seconds = var.refresh_seconds
    }
  )
  index_key = "index.html"
}

variable "environment" {}

variable "api_invoke_url" {}
variable "get_grid_resource" {}
variable "set_cell_resource" {}
variable "refresh_seconds" {}

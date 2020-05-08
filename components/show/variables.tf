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
  favicon_file = "favicon.png"
  about_file = "about.html"
  about_css_file = "about_style.css"

  ok_frame_names = [
      for i in range(var.nb_ok_frames):
      "frames/ok/ok_${i}.png"
    ]
  broken_frame_names = [
    for i in range(var.nb_broken_frames):
    "frames/broken/broken_${i}.png"
  ]
  frame_names = toset(
    concat(local.ok_frame_names, local.broken_frame_names)
  )

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

  subdomain_name = "www.${var.domain_name}"
}

variable "environment" {}

variable "api_invoke_url" {}
variable "get_grid_resource" {}
variable "set_cell_resource" {}
variable "refresh_seconds" {}
variable "domain_name" {}
variable "nb_ok_frames" {}
variable "nb_broken_frames" {}

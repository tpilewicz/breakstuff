locals {
  get_grid_resource = "get_grid"
  set_cell_resource = "set_cell"
}

variable "environment" {
}
variable "refresh_seconds" {
}
variable "domain_name" {
}

module "funes" {
  source = "../components/funes"

  environment = var.environment
  table_name  = "Funes${title(var.environment)}"
}

module "clausius" {
  source = "../components/clausius"

  environment   = var.environment

  funes_table = module.funes.table

  get_grid_resource = local.get_grid_resource
  set_cell_resource = local.set_cell_resource
}

module "show" {
  source = "../components/show"

  environment   = var.environment

  api_invoke_url    = module.clausius.invoke_url
  get_grid_resource = local.get_grid_resource
  set_cell_resource = local.set_cell_resource
  refresh_seconds = var.refresh_seconds
  domain_name = var.domain_name
  nb_ok_frames = 170
  nb_broken_frames = 60
}

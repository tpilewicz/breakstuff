resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/24"
  tags = {
    Name = var.environment
  }
}

module "funes" {
  source = "../components/funes"

  environment       = var.environment
  node_type         = "cache.t2.micro"
  node_count        = 1
  vpc_id            = aws_vpc.main.id
  subnet_cidr_block = "10.0.0.0/28"
}

module "clausius" {
  source = "../components/clausius"

  environment   = var.environment

  vpc_id        = aws_vpc.main.id
  funes_subnets = module.funes.subnet_ids
  funes_sg_id   = module.funes.sg_id
  funes_url     = module.funes.url

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

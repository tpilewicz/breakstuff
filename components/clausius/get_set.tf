# API

resource "aws_api_gateway_rest_api" "api" {
  name        = "${local.default_name}-gateway"
  description = "${local.default_name} gateway"
}

# Deploy the api. Manually changing the `updatedAt` variable will trigger a terraform deployment
resource "aws_api_gateway_deployment" "api" {
  depends_on = [module.get_grid.gateway_integration, module.set_cell.gateway_integration]

  rest_api_id = aws_api_gateway_rest_api.api.id
  stage_name  = "v1"

  variables = {
    "updatedAt" = "2020-04-19"
  }
}

# Lambdas

module "get_grid" {
  source = "./lambda"

  environment   = var.environment
  component     = "${local.component}-get-grid"
  api           = aws_api_gateway_rest_api.api
  api_resource  = var.get_grid_resource
  http_method   = "GET"
  lambda_src    = "../../components/clausius/src/get_grid/main"
  funes_table   = var.funes_table

  nb_rows = local.nb_rows
  nb_cols = local.nb_cols
}

module "set_cell" {
  source = "./lambda"

  environment   = var.environment
  component     = "${local.component}-set-cell"
  api           = aws_api_gateway_rest_api.api
  api_resource  = var.set_cell_resource
  http_method   = "POST"
  lambda_src    = "../../components/clausius/src/set_cell/main"
  funes_table   = var.funes_table

  nb_rows = local.nb_rows
  nb_cols = local.nb_cols
}

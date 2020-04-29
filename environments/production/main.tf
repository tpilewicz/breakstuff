terraform {
  backend "s3" {
    bucket = "tpilewicz-infra-production"
    key    = "breakstuff/terraform.tfstate"
    region = "eu-central-1"
  }

  required_version = "~> 0.12"
}

module "main" {
  source = "../../main"

  environment = local.environment
  refresh_seconds = local.refresh_seconds
}

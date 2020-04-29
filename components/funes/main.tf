resource "aws_elasticache_replication_group" "funes" {
  replication_group_id          = local.default_name
  engine                        = "redis"
  engine_version                = "5.0.6"
  replication_group_description = "Keeps track of who's broken and who's not"
  node_type                     = var.node_type
  number_cache_clusters         = var.node_count
  port                          = 6379
  parameter_group_name          = aws_elasticache_parameter_group.funes.id
  subnet_group_name             = aws_elasticache_subnet_group.funes.name
  security_group_ids            = [aws_security_group.funes.id]
  maintenance_window            = "tue:03:00-tue:07:00"

  tags = {
    Name        = local.default_name
    Environment = var.environment
    Component   = local.component
  }
}

resource "aws_elasticache_parameter_group" "funes" {
  name   = local.default_name
  family = "redis5.0"

  parameter {
    name  = "activerehashing"
    value = "yes"
  }

  parameter {
    name  = "maxmemory-policy"
    value = "noeviction"
  }
}

resource "aws_subnet" "funes" {
  vpc_id     = var.vpc_id
  cidr_block = var.subnet_cidr_block

  tags = {
    Name = local.default_name
  }
}

resource "aws_elasticache_subnet_group" "funes" {
  name       = local.default_name
  subnet_ids = [aws_subnet.funes.id]
}

resource "aws_security_group" "funes" {
  name        = local.default_name
  description = "Inbound traffic on redis"
  vpc_id      = var.vpc_id

  tags = {
    Name        = local.default_name
    Environment = var.environment
    Component   = local.component
  }
}

resource "aws_security_group_rule" "allow_clausius" {
  type              = "ingress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.funes.id
}

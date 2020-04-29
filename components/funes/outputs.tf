output "subnet_ids" {
  value = [aws_subnet.funes.id]
}
output "sg_id" {
  value = aws_security_group.funes.id
}
output "url" {
  value = "redis://${aws_elasticache_replication_group.funes.primary_endpoint_address}:6379/1"
}

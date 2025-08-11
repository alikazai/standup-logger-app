resource "aws_db_subnet_group" "this" {
  name       = "standup-db-subnets"
  subnet_ids = module.vpc.private_subnets
}

resource "aws_security_group" "db" {
  name        = "standup-db-sg"
  description = "DB ingress from ECS"
  #Not sure this will work
  vpc_id = module.vpc.vpc_id
}

# allow ECS tasks to hit 5432
resource "aws_security_group_rule" "db_in" {
  type                     = "ingress"
  from_port                = 5432
  to_port                  = 5432
  protocol                 = "tcp"
  security_group_id        = aws_security_group.db.id
  source_security_group_id = aws_security_group.ecs_tasks.id
}

resource "aws_db_instance" "postgres" {
  identifier              = "standup-db"
  engine                  = "postgres"
  instance_class          = "db.t4g.micro" #keep costs low
  allocated_storage       = 20
  username                = var.db_user
  password                = var.db_password
  db_name                 = var.db_name
  port                    = 5432
  skip_final_snapshot     = true
  publicly_accessible     = false
  vpc_security_group_ids  = [aws_security_group.db.id]
  db_subnet_group_name    = aws_db_subnet_group.this.name
  backup_retention_period = 7
}

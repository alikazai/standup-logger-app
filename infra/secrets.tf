resource "aws_secretsmanager_secret" "db" {
  name = "standup-db"
}

resource "aws_secretsmanager_secret_version" "db" {
  secret_id = aws_secretsmanager_secret.db.id
  secret_string = jsonencode({
    DB_HOST     = aws_db_instance.postgres.address
    DB_PORT     = "5432"
    DB_NAME     = var.db_name
    DB_USER     = var.db_user
    DB_PASSWORD = var.db_password
  })
}

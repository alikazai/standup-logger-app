data "aws_iam_policy_document" "task_exec_assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "task_exec" {
  name               = "standup-task-exec"
  assume_role_policy = data.aws_iam_policy_document.task_exec_assume.json
}

# Allow the EXECUTION role to read your secret for container startup
resource "aws_iam_role_policy" "task_exec_secrets" {
  name = "standup-task-exec-secrets"
  role = aws_iam_role.task_exec.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect   = "Allow",
      Action   = ["secretsmanager:GetSecretValue"],
      Resource = aws_secretsmanager_secret.db.arn
    }]
  })
}

#allow pull from ECR, write logs 
resource "aws_iam_role_policy_attachment" "task_exec_ecr" {
  role       = aws_iam_role.task_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

#app role (read secrets)
resource "aws_iam_role" "task_app" {
  name               = "standup-task-app"
  assume_role_policy = data.aws_iam_policy_document.task_exec_assume.json
}

resource "aws_iam_policy" "read_secrets" {
  name = "standup-read-secrets"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action   = ["secretsmanager:GetSecretValue"],
      Effect   = "Allow",
      Resource = aws_secretsmanager_secret.db.arn
    }]
  })
}

resource "aws_iam_role_policy_attachment" "app_secrets" {
  role       = aws_iam_role.task_app.name
  policy_arn = aws_iam_policy.read_secrets.arn

}

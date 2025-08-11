resource "aws_security_group" "ecs_tasks" {
  name   = "standup-ecs-tasks-sg"
  vpc_id = module.vpc.vpc_id
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_ecs_cluster" "this" {
  name = "standup-cluster"
}

resource "aws_cloudwatch_log_group" "app" {
  name              = "/ecs/standup-app"
  retention_in_days = 14
}

locals {
  container_name = "standup-app"
  image          = "${aws_ecr_repository.app.repository_url}:${var.image_tag}"
}

resource "aws_ecs_task_definition" "app" {
  family                   = "standup-app"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.task_exec.arn
  task_role_arn            = aws_iam_role.task_app.arn

  container_definitions = jsonencode([{
    name  = local.container_name
    image = local.image
    portMappings = [{
      containerPort = 8080,
      protocol      = "tcp",
    }]
    environment = [
      { name = "PORT", value = "8080" },
      #non-secret envs go here
    ]
    secrets = [
      { name = "DB_HOST", valueFrom = "${aws_secretsmanager_secret.db.arn}:DB_HOST::" },
      { name = "DB_PORT", valueFrom = "${aws_secretsmanager_secret.db.arn}:DB_PORT::" },
      { name = "DB_NAME", valueFrom = "${aws_secretsmanager_secret.db.arn}:DB_NAME::" },
      { name = "DB_USER", valueFrom = "${aws_secretsmanager_secret.db.arn}:DB_USER::" },
      { name = "DB_PASSWORD", valueFrom = "${aws_secretsmanager_secret.db.arn}:DB_PASSWORD::" },
    ]
    logConfiguration = {
      logDriver = "awslogs"
      options = {
        awslogs-group         = aws_cloudwatch_log_group.app.name
        awslogs-region        = var.region
        awslogs-stream-prefix = "ecs"
      }
    }
    healthCheck = {
      command     = ["CMD-SHELL", "curl -f http://localhost:8080/healthz || exit 1"]
      interval    = 30,
      timeout     = 5,
      retries     = 3,
      startPeriod = 10
    }
  }])
}

resource "aws_ecs_service" "app" {
  name            = "standup-svc"
  cluster         = aws_ecs_cluster.this.id
  task_definition = aws_ecs_task_definition.app.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    assign_public_ip = false
    security_groups  = [aws_security_group.ecs_tasks.id]
    subnets          = module.vpc.private_subnets
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.app.arn
    container_name   = local.container_name
    container_port   = 8080
  }

  depends_on = [aws_lb_listener.http]
}

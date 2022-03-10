resource "aws_ecs_cluster" "shortly" {
  name = "shortly"
  tags = local.tags
}

resource "aws_ecs_service" "shortly" {
  name                               = "${local.name}-service"
  cluster                            = aws_ecs_cluster.shortly.id
  task_definition                    = aws_ecs_task_definition.shortly.arn
  desired_count                      = var.replicas
  deployment_minimum_healthy_percent = 100
  deployment_maximum_percent         = 200
  launch_type                        = "FARGATE"
  scheduling_strategy                = "REPLICA"

  network_configuration {
    security_groups  = [module.ecs_security_group.security_group_id]
    subnets          = module.vpc.private_subnets
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = module.alb.target_group_arns[0]
    container_name   = "${local.name}-container"
    container_port   = var.port
  }

  lifecycle {
    ignore_changes = [desired_count]
  }

  tags = local.tags
}

resource "aws_ecs_task_definition" "shortly" {
  family                   = "shortly"
  task_role_arn            = aws_iam_role.ecs_task_role.arn
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  requires_compatibilities = ["FARGATE"]
  container_definitions    = jsonencode([
    {
      name        = "${local.name}-container"
      image       = "${var.container_image}:latest"
      essential   = true
      healthCheck = {
        command = ["CMD-SHELL", "wget -O /dev/null -q localhost:${var.port}${var.health_check_path} || exit 1"]
      }
      environment = [
        {
          name  = "SHORTLY_DBUSER",
          value = module.postgres.db_instance_username
        },
        {
          name  = "SHORTLY_DBPASS",
          value = module.postgres.db_instance_password
        },
        {
          name  = "SHORTLY_DBHOST",
          value = module.postgres.db_instance_address
        },
        {
          name  = "SHORTLY_DBPORT",
          value = tostring(module.postgres.db_instance_port)
        },
        {
          name  = "SHORTLY_DBNAME",
          value = module.postgres.db_instance_name
        },
        {
          name  = "SHORTLY_DBSSLMODE",
          value = "require"
        },
        {
          name  = "SHORTLY_MEMCACHESERVERS",
          value = join(",", [for n in aws_elasticache_cluster.shortly.cache_nodes : format("%s:%d", n.address, n.port)])
        },
      ]
      portMappings = [
        {
          protocol      = "tcp"
          containerPort = var.port
          hostPort      = var.port
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options   = {
          awslogs-group         = aws_cloudwatch_log_group.shortly.name
          awslogs-stream-prefix = "ecs"
          awslogs-region        = "us-east-1"
        }
      }
    }
  ])

  tags = local.tags
}

resource "aws_cloudwatch_log_group" "shortly" {
  name = "/ecs/${local.name}-task"
  tags = local.tags
}

module "ecs_security_group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4.0"

  name        = local.name
  description = "Allow traffic from "
  vpc_id      = module.vpc.vpc_id

  ingress_with_cidr_blocks = [
    {
      from_port   = var.port
      to_port     = var.port
      protocol    = "tcp"
      description = "Allow ingress from VPC to service port"
      cidr_blocks = module.vpc.vpc_cidr_block
    },
  ]

  egress_with_cidr_blocks = [
    {
      protocol    = "-1"
      from_port   = 0
      to_port     = 0
      description = "Allow responses"
      cidr_blocks = "0.0.0.0/0"
    }
  ]

  tags = local.tags
}

resource "aws_iam_role" "ecs_task_execution_role" {
  name = "ecs-task-execution-role"

  assume_role_policy = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
   {
     "Action": "sts:AssumeRole",
     "Principal": {
       "Service": "ecs-tasks.amazonaws.com"
     },
     "Effect": "Allow",
     "Sid": ""
   }
 ]
}
EOF
  tags               = local.tags
}

resource "aws_iam_role" "ecs_task_role" {
  name = "ecs-task-role"

  assume_role_policy = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
   {
     "Action": "sts:AssumeRole",
     "Principal": {
       "Service": "ecs-tasks.amazonaws.com"
     },
     "Effect": "Allow",
     "Sid": ""
   }
 ]
}
EOF

  tags = local.tags
}

resource "aws_iam_role_policy_attachment" "ecs-task-execution-role-policy-attachment" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}
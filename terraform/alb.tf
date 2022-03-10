module "alb_security_group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4.0"

  name        = local.name
  description = "Public ALB"
  vpc_id      = module.vpc.vpc_id

  ingress_with_cidr_blocks = [
    {
      from_port   = 80
      to_port     = 80
      protocol    = "tcp"
      description = "Allow external HTTP traffic"
      cidr_blocks = "0.0.0.0/0"
    },
  ]

  egress_with_cidr_blocks = [
    {
      protocol    = "-1"
      from_port   = 0
      to_port     = 0
      description = "Allow responses to requests"
      cidr_blocks = "0.0.0.0/0"
    }
  ]

  tags = local.tags
}

module "alb" {
  source  = "terraform-aws-modules/alb/aws"
  version = "~> 6.0"

  name = "${local.name}-alb"

  load_balancer_type = "application"

  vpc_id          = module.vpc.vpc_id
  subnets         = module.vpc.public_subnets
  security_groups = [module.alb_security_group.security_group_id]

  target_groups = [
    {
      backend_protocol = "HTTP"
      backend_port     = var.port
      target_type      = "ip"
      health_check     = {
        path                = var.health_check_path
        port                = var.port
        healthy_threshold   = 3
        unhealthy_threshold = 2
        timeout             = 2
        interval            = 10
        matcher             = "200-299"
      }
    }
  ]

  http_tcp_listeners = [
    {
      port               = 80
      protocol           = "HTTP"
      target_group_index = 0
    }
  ]

  tags = local.tags
}

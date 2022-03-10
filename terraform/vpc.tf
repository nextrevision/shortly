module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 3.0"

  name = local.name
  cidr = "10.99.0.0/18"

  azs                 = ["${local.region}a", "${local.region}b", "${local.region}c"]
  public_subnets      = ["10.99.0.0/24", "10.99.1.0/24", "10.99.2.0/24"]
  private_subnets     = ["10.99.3.0/24", "10.99.4.0/24", "10.99.5.0/24"]
  database_subnets    = ["10.99.7.0/24", "10.99.8.0/24", "10.99.9.0/24"]
  elasticache_subnets = ["10.99.6.0/26", "10.99.10.64/26", "10.99.6.128/26"]

  enable_nat_gateway     = true
  single_nat_gateway     = true
  one_nat_gateway_per_az = false

  create_database_subnet_group           = true
  create_database_subnet_route_table     = true
  create_database_internet_gateway_route = false

  create_elasticache_subnet_group       = true
  create_elasticache_subnet_route_table = true
  elasticache_subnet_group_tags         = local.tags

  tags = local.tags
}
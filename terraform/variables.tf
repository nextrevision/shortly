variable "port" {
  default     = 8000
  description = "Port to run shortly on"
}

variable "container_image" {
  default     = "docker.io/nextrevision/shortly:9d93c3e77545256e1b9ecc30aaa2738d25c8c164"
  description = "Name of the container image to use including tag"
}

variable "replicas" {
  default     = 2
  description = "Number of shortly containers to run"
}

variable "health_check_path" {
  default     = "/-/health"
  description = "Path to perform health check against"
}

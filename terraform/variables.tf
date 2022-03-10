variable "port" {
  default     = 8000
  description = "Port to run shortly on"
}

variable "container_image" {
  default     = "docker.io/nextrevision/shortly:e612fbfa8bcb7c81d11bf01a9bd1e2f43aff432d"
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

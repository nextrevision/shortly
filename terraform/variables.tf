variable "port" {
  default     = 8000
  description = "Port to run shortly on"
}

variable "container_image" {
  default     = "docker.io/nextrevision/shortly"
  description = "Name of the container image to use without tag"
}

variable "replicas" {
  default     = 2
  description = "Number of shortly containers to run"
}

variable "health_check_path" {
  default     = "/-/health"
  description = "Path to perform health check against"
}

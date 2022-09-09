variable "region" {
  type        = string
  description = "The GCP region to deploy everything to"
  default     = "europe-west3"
}

variable "services" {
  type        = list(string)
  description = "The names of the services to deploy"
  default     = [
    "hello-service"
  ]
}
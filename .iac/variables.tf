variable "gcp_region" {
  type        = string
  description = "The GCP region to deploy everything to"
  default     = "europe-west3"
}

variable "gcp_project" {
  type        = string
  description = "The GCP project to deploy everything to"
  default     = "clubrizer-com"
}

variable "services" {
  type        = list(string)
  description = "The names of the services to deploy"
  default = [
    "hello-service"
  ]
}
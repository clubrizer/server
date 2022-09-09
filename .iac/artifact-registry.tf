###
# Creates an artifact registry repository for the Docker images of all services
###

resource "google_artifact_registry_repository" "server" {
  repository_id = "server"
  location      = var.region
  description   = "Docker images of all services in the Clubrizer Backend"
  format        = "DOCKER"
}
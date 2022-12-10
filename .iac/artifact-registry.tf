###
# Creates an artifact registry repository for the Docker images of all services
###

resource "google_kms_key_ring" "server_registry" {
  name     = "clubrizer-artifact-registry"
  location = "global"
}

resource "google_kms_crypto_key" "server_registry_key" {
  name            = "clubrizer-artifact-registry-key"
  key_ring        = google_kms_key_ring.server_registry.id
  rotation_period = "86400s"

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_artifact_registry_repository" "server" {
  repository_id = "server"
  location      = var.gcp_region
  description   = "Docker images of all services in the Clubrizer Backend"
  format        = "DOCKER"
  kms_key_name  = google_kms_crypto_key.server_registry_key.name
}
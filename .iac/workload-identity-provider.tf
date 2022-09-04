###
# Creates the Workload identity provider.
###

resource "google_iam_workload_identity_pool" "server" {
  workload_identity_pool_id = "server"
  display_name              = "Server"
  description               = "Identity pool for all service accounts used by the Clubrizer Backend"
}

resource "google_iam_workload_identity_pool_provider" "server" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.server.workload_identity_pool_id
  workload_identity_pool_provider_id = "server-provider"
  display_name                       = "Servier Identity provider"
  description                        = "OIDC identity pool provider for all service accounts used by the Clubrizer Backend"

  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.actor"      = "assertion.actor"
    "attribute.repository" = "assertion.repository"
  }
  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

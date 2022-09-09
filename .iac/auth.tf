###
# Creates the Workload identity provider and service accounts that are necessary to deploy soemthing to Cloud Run.
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
  attribute_mapping                  = {
    "google.subject"       = "assertion.sub"
    "attribute.actor"      = "assertion.actor"
    "attribute.repository" = "assertion.repository"
  }
  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

resource "google_service_account" "deploy_cloudrun" {
  account_id   = "deploy-cloudrun"
  display_name = "Deploy to Cloud Run"
  description  = "Used to deploy various services to Cloud Run"
}

resource "google_project_iam_binding" "deploy_cloudrun_run_admin" {
  project = "clubrizer-com"
  role    = "roles/run.admin"
  members = ["serviceAccount:${google_service_account.deploy_cloudrun.email}"]
}

resource "google_project_iam_binding" "deploy_cloudrun_sa_user" {
  project = "clubrizer-com"
  role    = "roles/iam.serviceAccountUser"
  members = ["serviceAccount:${google_service_account.deploy_cloudrun.email}"]
}

resource "google_project_iam_binding" "deploy_cloudrun_registry_admin" {
  project = "clubrizer-com"
  role    = "roles/artifactregistry.admin"
  members = ["serviceAccount:${google_service_account.deploy_cloudrun.email}"]
}

resource "google_service_account_iam_binding" "deploy_cloudrun" {
  service_account_id = google_service_account.deploy_cloudrun.name
  role               = "roles/iam.workloadIdentityUser"

  members = [
    "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.server.name}/attribute.repository/clubrizer/server",
  ]
}

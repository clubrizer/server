###
# Creates the service account that can deploy to Cloud Run.
###

resource "google_service_account" "deploy_cloudrun" {
  account_id   = "deploy-cloudrun"
  display_name = "Deploy to Cloud Run"
  description  = "Used to deploy various services to Cloud Run"
}

resource "google_project_iam_member" "deploy_cloudrun_roles" {
  for_each = toset([
    "roles/run.admin",
    "roles/artifactregistry.admin",
    "roles/artifactregistry.repoAdmin",
    "roles/iam.serviceAccountUser"
  ])

  project = var.gcp_project
  role    = each.key
  member = "serviceAccount:${google_service_account.deploy_cloudrun.email}"
}

resource "google_service_account_iam_binding" "deploy_cloudrun" {
  service_account_id = google_service_account.deploy_cloudrun.name
  role               = "roles/iam.workloadIdentityUser"

  members = [
    "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.server.name}/attribute.repository/clubrizer/server",
  ]
}

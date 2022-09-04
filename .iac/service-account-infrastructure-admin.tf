###
# Creates the service account that can deploy to Cloud Run.
###

resource "google_service_account" "infrastructure_admin" {
  account_id   = "infrastructure-admin"
  display_name = "Infrastructure Admin"
  description  = "Used to deploy the infrastructure in this GCP project"
}

resource "google_project_iam_binding" "infrastructure_admin_roles" {
  for_each = toset([
    "roles/viewer",
    "roles/storage.objectViewer",
    "roles/iam.serviceAccountUser"
  ])

  project = var.gcp_project
  role    = each.key
  members = ["serviceAccount:${google_service_account.infrastructure_admin.email}"]
}

resource "google_service_account_iam_binding" "infrastructure_admin" {
  service_account_id = google_service_account.infrastructure_admin.name
  role               = "roles/iam.workloadIdentityUser"

  members = [
    "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.server.name}/attribute.repository/clubrizer/server",
  ]
}

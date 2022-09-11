###
# Creates the service account that can read infrastructure from GCP.
###

resource "google_service_account" "infrastructure_reader" {
  account_id   = "infrastructure-reader"
  display_name = "Infrastructure Reader"
  description  = "Used to plan the infrastructure in this GCP project"
}

resource "google_project_iam_member" "infrastructure_reader_roles" {
  for_each = toset([
    "roles/viewer",
    "roles/storage.objectViewer",
    "roles/iam.serviceAccountUser"
  ])

  project = var.gcp_project
  role    = each.key
  member = "serviceAccount:${google_service_account.infrastructure_reader.email}"
}

resource "google_service_account_iam_binding" "infrastructure_reader" {
  service_account_id = google_service_account.infrastructure_reader.name
  role               = "roles/iam.workloadIdentityUser"

  members = [
    "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.server.name}/attribute.repository/clubrizer/server",
  ]
}

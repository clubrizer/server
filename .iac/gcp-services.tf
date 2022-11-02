###
# Enables all necessary Google Cloud Project Services.
# See https://console.cloud.google.com/apis/library?project=clubrizer-com for all services.
###

resource "google_project_service" "all" {
  for_each = toset([
    "storage.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
    "compute.googleapis.com",
    "artifactregistry.googleapis.com",
    "run.googleapis.com",
    "pubsub.googleapis.com"
  ])

  service = each.key

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
  disable_on_destroy         = true
}

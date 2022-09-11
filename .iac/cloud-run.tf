###
# Deploys the Cloud Run Services for all backend services.
# Here, the gcr.io/cloudrun/hello image gets deployed, as the image will be set during the service deployment anyways.
###

resource "google_cloud_run_service" "server" {
  for_each = toset(var.services)

  name     = each.key
  location = var.gcp_region
  template {
    metadata {
      labels = {
        environment = "prod"
      }
    }
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  lifecycle {
    ignore_changes = [
      template[0].spec[0].containers
    ]
  }
}

data "google_iam_policy" "server_invoke_public" {
  binding {
    role = "roles/run.invoker"

    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "server_invoke_public" {
  for_each = google_cloud_run_service.server

  location = each.value.location
  project  = each.value.project
  service  = each.value.name

  policy_data = data.google_iam_policy.server_invoke_public.policy_data
}
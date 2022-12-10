# Prerequisites:
#   - Service account
#   - Cloud Resource Manager API must be enabled
# Authentication: Export the location of the service account key to GOOGLE_APPLICATION_CREDENTIALS

resource "google_storage_bucket" "backend" {
  name                        = "clubrizer-com-tfstate"
  location                    = "EUROPE-WEST3"
  storage_class               = "STANDARD"
  uniform_bucket_level_access = true

  versioning {
    enabled = true
  }

  logging {
    log_bucket = "clubrizer-com-access-logs"
  }
}

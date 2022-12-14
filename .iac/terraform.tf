terraform {
  required_version = ">= 1.3.0, < 2.0.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
  }

  backend "gcs" {
    bucket = "clubrizer-com-tfstate"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = "clubrizer-com"
  region  = "europe-west3"
  zone    = "europe-west3-a"
}

provider "google" {
}

data "google_compute_network" "gcp_network" {
  name = "terraform-vpc-${var.tag_name}"
  project = var.google_project
}

resource "google_compute_subnetwork" "gcp_subnetwork" {
  name          = "terraform-${var.tag_name}"
  ip_cidr_range = "10.100.3.0/24"
  region        = var.google_region
  network       = data.google_compute_network.gcp_network.self_link
}

resource "google_compute_instance" "gcp_instance" {
  name         = "terraform-vm-${var.tag_name}"
  machine_type = "n1-standard-1"
  zone         = var.google_zone
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }
  network_interface {
      subnetwork = google_compute_subnetwork.gcp_subnetwork.self_link
  }
  metadata = {
   ssh-keys = "ubuntu:${var.ssh_public_key}"
 }
}

resource "google_compute_interconnect_attachment" "gcp_interconnect" {
  name                     = "terraform-attachment-${var.tag_name}"
  type                     = "PARTNER"
  edge_availability_domain = var.google_interconnect_attachment_edge_availability_domain
  region                   = var.google_region
  router                   = var.google_interconnect_attachment_router
}






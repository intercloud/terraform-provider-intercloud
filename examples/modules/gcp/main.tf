provider "google" {
}

resource "google_compute_network" "gcp_network" {
  name                    = "terraform-vpc-${var.tag_name}"
  project                 = var.google_project
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "gcp_firewall" {
  name    = "terraform-firewall-${var.tag_name}"
  network = google_compute_network.gcp_network.self_link
  allow {
    protocol = "tcp"
    ports    = ["22"]
  }
}

resource "google_compute_subnetwork" "gcp_subnetwork" {
  name          = "terraform-${var.tag_name}"
  ip_cidr_range = "10.100.3.0/24"
  region        = var.google_region
  network       = google_compute_network.gcp_network.self_link
}

resource "google_compute_address" "gcp_public_ip" {
  name    = "terraform-ip-${var.tag_name}"
  project = var.google_project
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
    subnetwork = google_compute_subnetwork.gcp_subnetwork.name
    access_config {
      nat_ip = google_compute_address.gcp_public_ip.address
    }
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






variable "GOOGLE_PROJECT" {
  type     = string
  nullable = false
}

variable "GOOGLE_REGION" {
  type     = string
  nullable = false
  default  = "asia-southeast2"
}

variable "GOOGLE_ZONE" {
  type     = string
  nullable = false
  default  = "asia-southeast2-c"
}

variable "DOCKER_MANAGER_COUNT" {
  type    = number
  default = 1
}

variable "DOCKER_WORKER_COUNT" {
  type    = number
  default = 3
}

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.51.0"
    }
  }
}

provider "google" {
  project = var.GOOGLE_PROJECT
  region  = var.GOOGLE_REGION
  zone    = var.GOOGLE_ZONE
}

resource "google_compute_network" "docker_vpc_network" {
  name                    = "docker-vpc-network"
  auto_create_subnetworks = false
  mtu                     = 1460
}

resource "google_compute_subnetwork" "docker_subnetwork" {
  name          = "docker-subnetwork"
  ip_cidr_range = "10.0.0.0/24"
  region        = var.GOOGLE_REGION
  network       = google_compute_network.docker_vpc_network.name
}

resource "google_compute_firewall" "swarm_firewall" {
  name    = "docker-firewall"
  network = google_compute_network.docker_vpc_network.name

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["22", "8080", "80", "443", "9000", "2377", "7946", "3000"]
  }

  allow {
    protocol = "udp"
    ports    = ["7946", "4789"]
  }

  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_instance" "managers" {
  machine_type = "e2-standard-2"
  name         = "manager-${count.index}"
  count        = var.DOCKER_MANAGER_COUNT
  zone         = var.GOOGLE_ZONE
  tags         = ["http-server"]

  boot_disk {
    auto_delete = true
    device_name = "instance-1"

    initialize_params {
      image = "projects/debian-cloud/global/images/debian-11-bullseye-v20231212"
      size  = 30
      type  = "pd-standard"
    }

    mode = "READ_WRITE"
  }

  network_interface {
    access_config {
      network_tier = "STANDARD"
    }

    subnetwork = google_compute_subnetwork.docker_subnetwork.name
  }

  scheduling {
    automatic_restart  = true
    preemptible        = false
    provisioning_model = "STANDARD"
  }

  shielded_instance_config {
    enable_integrity_monitoring = true
    enable_secure_boot          = true
    enable_vtpm                 = true
  }

  metadata = {
    ssh-keys = "docker:${file("ssh/id_rsa.pub")}"
  }
  
  metadata_startup_script = file("startup-script.sh")
}

resource "google_compute_instance" "workers" {
  machine_type = "e2-standard-2"
  name         = "worker-${count.index}"
  count        = var.DOCKER_WORKER_COUNT
  zone         = var.GOOGLE_ZONE
  tags         = ["http-server"]

  boot_disk {
    auto_delete = true
    device_name = "instance-1"

    initialize_params {
      image = "projects/debian-cloud/global/images/debian-11-bullseye-v20231212"
      size  = 30
      type  = "pd-standard"
    }

    mode = "READ_WRITE"
  }

  network_interface {
    access_config {
      network_tier = "STANDARD"
    }

    subnetwork = google_compute_subnetwork.docker_subnetwork.id
  }

  scheduling {
    automatic_restart  = true
    preemptible        = false
    provisioning_model = "STANDARD"
  }

  shielded_instance_config {
    enable_integrity_monitoring = true
    enable_secure_boot          = true
    enable_vtpm                 = true
  }

  metadata = {
    ssh-keys = "docker:${file("ssh/id_rsa.pub")}"
  }

  metadata_startup_script = file("startup-script.sh")
}

output "manager_ips" {
  value = google_compute_instance.managers[*].network_interface.0.access_config.0.nat_ip
}

output "manager_first_node_private_ip" {
  value = google_compute_instance.managers[0].network_interface.0.network_ip
}

output "manager_first_node_public_ip" {
  value = google_compute_instance.managers[0].network_interface[0].access_config[0].nat_ip
}

output "worker_ips" {
  value = google_compute_instance.workers[*].network_interface.0.access_config.0.nat_ip
}

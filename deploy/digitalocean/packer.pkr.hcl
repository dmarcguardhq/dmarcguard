# Parse DMARC - DigitalOcean Marketplace Image
# Packer template for building the Marketplace-ready Droplet image

packer {
  required_plugins {
    digitalocean = {
      version = ">= 1.0.4"
      source  = "github.com/digitalocean/digitalocean"
    }
  }
}

variable "do_api_token" {
  type        = string
  description = "DigitalOcean API token"
  sensitive   = true
  default     = env("DIGITALOCEAN_TOKEN")
}

variable "parse_dmarc_version" {
  type        = string
  description = "Parse DMARC version to install"
  default     = "latest"
}

variable "image_name" {
  type        = string
  description = "Name for the snapshot"
  default     = "parse-dmarc"
}

variable "droplet_region" {
  type        = string
  description = "Region to create the temporary droplet in"
  default     = "nyc3"
}

variable "droplet_size" {
  type        = string
  description = "Size of the temporary droplet"
  default     = "s-1vcpu-1gb"
}

variable "base_image" {
  type        = string
  description = <<EOT
The DigitalOcean image slug or snapshot ID to use as the base image.
Default is "ubuntu-24-04-x64" (the latest Ubuntu 24.04 LTS image at build time).
For reproducible builds, specify a snapshot ID (e.g., "12345678").
See https://docs.digitalocean.com/products/images/ for details.
EOT
  default     = "ubuntu-24-04-x64"
}

source "digitalocean" "parse-dmarc" {
  api_token     = var.do_api_token
  image         = var.base_image
  region        = var.droplet_region
  size          = var.droplet_size
  ssh_username  = "root"
  snapshot_name = "${var.image_name}-{{timestamp}}"
  snapshot_regions = [
    "nyc1", "nyc2", "nyc3",
    "sfo1", "sfo2", "sfo3",
    "ams2", "ams3",
    "sgp1",
    "lon1",
    "fra1",
    "tor1",
    "blr1",
    "syd1"
  ]
  tags = ["parse-dmarc", "marketplace"]
}

build {
  sources = ["source.digitalocean.parse-dmarc"]

  # Copy scripts to the droplet
  provisioner "file" {
    source      = "scripts/"
    destination = "/tmp/"
  }

  # Run the setup script
  provisioner "shell" {
    environment_vars = [
      "PARSE_DMARC_VERSION=${var.parse_dmarc_version}",
      "DEBIAN_FRONTEND=noninteractive"
    ]
    scripts = [
      "scripts/01-system-setup.sh",
      "scripts/02-install-parse-dmarc.sh",
      "scripts/03-configure-services.sh",
      "scripts/90-cleanup.sh",
      "scripts/99-img-check.sh"
    ]
  }
}

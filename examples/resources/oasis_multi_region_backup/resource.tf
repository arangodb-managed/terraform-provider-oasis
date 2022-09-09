terraform {
  required_version = ">= 0.13.0"
  required_providers {
    oasis = {
      source  = "arangodb-managed/oasis"
      version = ">=2.1.0"
    }
  }
}

provider "oasis" {
  api_key_id     = "" // API Key ID generated in Oasis platform
  api_key_secret = "" // API Key Secret generated in Oasis platform
  organization   = "" // Your Oasis organization where you want to create the resources
}

resource "oasis_deployment" "my_oneshard_deployment" {
	terms_and_conditions_accepted = "true"
	project = "959885701" 
	name = "oasis_multi_region_deployment"
	location {
		region = "gcp-europe-west4"
	}
	version {
		db_version = "3.8.7"
	}
	security {
		disable_foxx_authentication = false
	}
	disk_performance = "dp30"
	configuration {
		model = "oneshard"
		node_size_id = "c4-a8"
		node_disk_size = 20
		maximum_node_disk_size = 40
	}
	notification_settings {
		email_addresses = [
		"test@arangodb.com"
		]
	}
}

resource "oasis_backup" "backup" {
	name = "oasis_backup"
	description = "test backup description update from terraform"
	deployment_id = oasis_deployment.my_oneshard_deployment.id
	upload = true
	auto_deleted_at = 20
	backup_policy_id = "456123"
}

resource "oasis_multi_region_backup" "backup" {
	source_backup_id = oasis_backup.backup.id
	region_id = "gcp-europe-west4"
}
variable "zones" {
	description = "List of zones to configure."
	type        = list
	default     = []
}

variable "nameservers" {
	description = "List of nameservers to configure in the given zones (automatically populated from NS records if not specified)."
	type        = list
	default     = []
}

variable "records" {
	description = "List of records to configure in the given zones."
	type        = any
	default     = []
}

variable "region" {
  type    = string
  default = "eu-west-2"
}

variable "azs" {
  type    = list(string)
  default = ["eu-west-2a", "eu-west-2b"]
}

variable "db_user" {
  type = string
}
variable "db_password" {
  type      = string
  sensitive = true
}
variable "db_name" {
  type = string
}
variable "image_tag" {
  type = string
  #e.g. "v1.0.0"
}
variable "alert_email" {
  type = string
}

variable "suffix" {
  description = "Optional suffix to avoid AWS name collisions"
  type        = string
  default     = ""
}

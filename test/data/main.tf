variable "name" {
  default = "Unknown"
}

output "output_hello" {
  value = var.name
}
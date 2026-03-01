variable "my_key_name" {
  type = string
}

variable "azs" {
  default = [ "ap-southeast-7a" ]
  type = list(string)
}

variable "public_subnets" {
    default = [ "10.0.101.0/24" ]
  type = list(string)
}

variable "project_name" {
  default = "medee"
    type = string
}
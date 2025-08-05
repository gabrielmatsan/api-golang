terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "6.7.0"
    }
  }

  backend "s3" {
    bucket = "terraform-state-bucket-golang"
    key = "terraform.tfstate"
    region = "us-east-1"
    encrypt = true
  }
  
}

provider "aws" {
  region = var.aws_region
}


resource "aws_s3_bucket" "terraform_state" {
  bucket = "terraform-state-bucket-golang"
  region = var.aws_region
  
  versioning {
    enabled = true
  }
}
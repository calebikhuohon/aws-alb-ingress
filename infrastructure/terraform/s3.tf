resource "aws_s3_bucket" "ikh-json-bucket" {
 bucket = var.tag_bucket_name

 tags = {
   Name        = var.tag_bucket_name
   Environment = var.tag_bucket_environment
 }
}

resource "aws_s3_bucket_acl" "ikh-json-bucket-acl" {
 bucket = var.tag_bucket_name
 acl = "private"
}

resource "aws_s3_bucket_versioning" "ikh-json-bucket-version" {
 bucket = var.tag_bucket_name
 versioning_configuration {
   status = "Enabled"
 }
}

data "aws_iam_policy_document" "s3_bucket_policy" {
 statement {
   effect = "Allow"
   principals {

     identifiers = [local.aws_account_id]
     type        = "AWS"
   }
   actions   = ["*"]
   resources = ["${aws_s3_bucket.ikh-json-bucket.arn}/*"]
 }

 statement {
   effect = "Deny"
   principals {
     identifiers = ["*"]
     type        = "AWS"
   }
   actions   = ["*"]
   resources = ["${aws_s3_bucket.ikh-json-bucket.arn}/*"]

   condition {
     test     = "Bool"
     variable = "aws:SecureTransport"
     values = [
       "false",
     ]
   }
 }
}

data "aws_caller_identity" "current" {
}

locals {
 aws_account_id = data.aws_caller_identity.current.account_id
}
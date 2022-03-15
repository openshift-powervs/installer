output "bootstrap_ignition_host" {
  value = ibm_cos_bucket.ignition.s3_endpoint_public
}
output "bootstrap_ignition_bucket" {
  value = ibm_cos_bucket.ignition.bucket_name
}
output "bootstrap_ignition_key" {
  value = ibm_cos_bucket_object.ignition.key
}

data "spicedb_schema" "test" {
}

output "name" {
  value = data.spicedb_schema.test.schema
}
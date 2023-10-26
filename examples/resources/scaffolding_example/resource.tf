resource "spicedb_schema" "test" {
  schema = <<EOF
definition user {}

definition organization {
    permission is_member = member
    relation member : user
}
EOF
}
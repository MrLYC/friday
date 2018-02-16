#!/bin/sh
conf_file=/tmp/migrate_test.yaml

migrate() {
  ${TARGET} confinfo -c ${conf_file}
  ${TARGET} migrate -c ${conf_file} -action rollback || exit $?
  ${TARGET} migrate -c ${conf_file} || exit $?
  ${TARGET} migrate -c ${conf_file} -action rebuild || exit $?
  ${TARGET} migrate -c ${conf_file} -action rollback || exit $?
  ${TARGET} migrate -c ${conf_file} -action list || exit $?
}

sqlite3() {
  conf_file=/tmp/migrate_test.yaml
  dbname=/tmp/migrate_test.db
  cat > ${conf_file} << EOF 
database:
  type: sqlite3
  name: ${dbname}
EOF

  migrate "${conf_file}"

  rm -rf "${dbname}"
}

mysql() {
  cat > ${conf_file} << EOF 
database:
  type: mysql
  name: friday
  host: 127.0.0.1
  port: 3306
  user: root
  password: ""
EOF

  migrate "${conf_file}"
}

postgresql() {
  cat > ${conf_file} << EOF 
database:
  type: postgres
  name: friday
  host: 127.0.0.1
  port: 5432
  user: postgres
  password: my_postgres_password
EOF

  migrate "${conf_file}"
}

sqlite3
mysql
postgresql
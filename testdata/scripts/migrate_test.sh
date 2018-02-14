#!/bin/sh

function migrate() {
  conf_file="$1"
  ${TARGET} confinfo -c ${conf_file}
  ${TARGET} migrate -c ${conf_file} -action rollback || exit $?
  ${TARGET} migrate -c ${conf_file} || exit $?
  ${TARGET} migrate -c ${conf_file} -action rebuild || exit $?
  ${TARGET} migrate -c ${conf_file} -action rollback || exit $?
  ${TARGET} migrate -c ${conf_file} -action list || exit $?
}

function sqlite3() {
  conf_file=/tmp/migrate_test.yaml
  dbname=/tmp/migrate_test.db
  cat > ${conf_file} << EOF 
database:
  type: sqlite3
  name: ${dbname}
EOF

  migrate "${conf_file}"

  rm -rf "${dbname}"
  rm -rf "${conf_file}"
}

sqlite3
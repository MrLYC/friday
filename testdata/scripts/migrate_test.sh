#!/bin/sh

conf_file=/tmp/migrate_test.yaml
dbname=/tmp/migrate_test.db

cat > ${conf_file} << EOF 
database:
  type: sqlite3
  name: ${dbname}
EOF

${TARGET} confinfo -c ${conf_file}
${TARGET} migrate -c ${conf_file} -action rollback || exit $?
${TARGET} migrate -c ${conf_file} || exit $?
${TARGET} migrate -c ${conf_file} -action rollback || exit $?
${TARGET} migrate -c ${conf_file} || exit $?
${TARGET} migrate -c ${conf_file} -action list || exit $?

rm -rf "${dbname}"
rm -rf "${conf_file}"
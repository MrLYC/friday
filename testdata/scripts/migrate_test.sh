#!/bin/sh

dbname=/tmp/migrate_test.db

cat > friday.yaml << EOF 
database:
  type: sqlite3
  name: ${dbname}
EOF

./bin/friday confinfo
./bin/friday migrate -action rollback || exit $?
./bin/friday migrate || exit $?
./bin/friday migrate -action rollback || exit $?
./bin/friday migrate || exit $?
rm -rf "${dbname}"
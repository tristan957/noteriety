#!/usr/bin/env sh

for migration in sql/migrations/*.sql; do
	file="${migration##*/}"
	filename=$(basename $file .sql)
	IFS='-' read -ra migration_parts <<< "$filename"
	if ! sqlite3 noteriety.db "SELECT name FROM migration;" | grep --quiet "${migration_parts[1]}"; then
		echo "Running migration $filename..."
		sqlite3 noteriety.db < $migration
		sqlite3 noteriety.db "INSERT INTO migration (migration_id, name) VALUES (${migration_parts[0]}, '${migration_parts[1]}');"
	fi
done

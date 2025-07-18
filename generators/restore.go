package generators

import (
	"fmt"
	"log"
	"strings"
)

var (
	CommonDatabases = []string{"postgres", "mariadb", "mysql"}
	DatabaseImages  []string
)

func IsCommonDatabaseImage(imageName string) bool {
	lowerName := strings.ToLower(imageName)
	log.Printf("[DEBUG] Checking if image '%s' is a common database image...", imageName)

	for _, db := range CommonDatabases {
		if strings.Contains(lowerName, db) {
			log.Printf("[DEBUG] Match found: image '%s' contains database keyword '%s'", imageName, db)
			return true
		}
	}

	log.Printf("[DEBUG] No common database match found for image '%s'", imageName)
	return false
}

func GenerateBackupRestore(imageName, appName, envFile, dbName string) (restore, backup, clean string) {
	if strings.Contains(strings.ToLower(imageName), "postgres") {
		return PostgresDBBackupRestore(imageName, appName, envFile, dbName)
	} else if strings.Contains(strings.ToLower(imageName), "mariadb") {
		return MariaDBBackupRestore(imageName, envFile, dbName, appName)
	} else if strings.Contains(strings.ToLower(imageName), "mysql") {
		return MariaDBBackupRestore(imageName, envFile, dbName, appName)
	}
	return "", "", ""
}

func MariaDBBackupRestore(imageName, envFile, dbName, appName string) (restore string, backup, clean string) {
	restoreT := fmt.Sprintf(`# Prepare an initialization script that restores the dump file
mkdir -vp initdb.d
mv -v %s.sql initdb.d

# do the bash file to restore and exit once done
cat - >initdb.d/zz_%s_restore.sh <<'EOS'
# Print additional information:
mysql --version
# The script is sourced, override entrypoint args and exit:
set -- true
docker_temp_server_stop
exit 0
EOS

# once we exit we remove initdb.d
trap 'rm -rfv initdb.d/' EXIT

# we start a container to initiate a database and load the dump
# at the end of %s_restore.sh the dump is loaded and
# we exit the container
podman run \
  --rm \
  --interactive \
  --network=none \
  --volume=./initdb.d:/docker-entrypoint-initdb.d:z \
  --volume mysql-data:/var/lib/mysql/:Z \
  %s \
  --replace --name=restore_db \
  "${%s}"`, dbName, dbName, dbName, envFile, imageName)
	backupT := fmt.Sprintf(`
		podman exec %s-app mysqldump \
        --databases %s \
        --default-character-set=utf8mb4 \
        --skip-dump-date \
        --ignore-table=mysql.event \
        --single-transaction \
        --quick \
        --add-drop-table  > %s.sql
		`, appName, dbName, dbName)
	cleanT := fmt.Sprintf(`rm -vf %s.sql`, dbName)
	return restoreT, backupT, cleanT
}

func PostgresDBBackupRestore(imageName, appName, envFile, dbName string) (restore string, backup, clean string) {
	restoreT := fmt.Sprintf(`# Create restore directory
mkdir -vp restore

# Create restore script
cat - >restore/monica_restore.sh <<'EOS'
# Read dump file from standard input:
pg_restore --no-owner --no-privileges -U postgres -d %s
ec=$?
docker_temp_server_stop
exit $ec
EOS

# Run the restore container
podman run \
    --rm \
    --interactive \
    --network=none \
    --volume=./restore:/docker-entrypoint-initdb.d/:Z \
    --volume=postgres-data:/var/lib/postgresql/data:Z \
    --replace --name=restore_db \
    %s \
    "${%s}" < %s.pg_dump

# Clean up
rm -rfv restore/ %s.pg_dump`, dbName, envFile, imageName, dbName, dbName)
	backupT := fmt.Sprintf(`echo "Dumping  postgres database"
podman exec %s-app pg_dump -U postgres --format=c  %s > %s.pg_dump`, appName, dbName, dbName)
	cleanT := fmt.Sprintf(`rm -vf %s.pg_dump`, dbName)
	return restoreT, backupT, cleanT
}


package aws

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)


func AssertPotgresqlDatabaseExists(t *testing.T, dbUrl string, dbPort int32, dbUsername string, dbPassword string, databaseName string) bool {
	output, err := AssertPotgresqlDatabaseExistsE(t, dbUrl, dbPort, dbUsername, dbPassword, databaseName)
	require.NoError(t, err)
	return output
}

func AssertPotgresqlDatabaseExistsE(t *testing.T, dbUrl string, dbPort int32, dbUsername string, dbPassword string, databaseName string) (bool, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", dbUrl, dbPort, dbUsername, dbPassword, databaseName)

	db, connErr := sql.Open("pgx", connectionString)
	if connErr != nil {
		return false, connErr
	}
	defer db.Close()
	return true, nil
}

func AssertPotgresqlSchemaExists(t *testing.T, dbUrl string, dbPort int32, dbUsername string, dbPassword string, databaseName string, expectedSchemaName string) bool {
	output, err := AssertPotgresqlSchemaExistsE(t, dbUrl, dbPort, dbUsername, dbPassword, databaseName, expectedSchemaName)
	if err != nil {
		t.Fatal(err)
	}
	return output
}

func AssertPotgresqlSchemaExistsE(t *testing.T, dbUrl string, dbPort int32, dbUsername string, dbPassword string, databaseName string, expectedSchemaName string) (bool, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", dbUrl, dbPort, dbUsername, dbPassword, databaseName)

	db, connErr := sql.Open("pgx", connectionString)
	if connErr != nil {
		return false, connErr
	}
	defer db.Close()
	var (
		schemaName string
	)
	sqlStatement := `SELECT "schema_name" FROM "information_schema"."schemata" where schema_name=$1`
	row := db.QueryRow(sqlStatement, expectedSchemaName)
	scanErr := row.Scan(&schemaName)
	if scanErr != nil {
		return false, scanErr
	}
	return true, nil
}

func AssertPotgresqlGrantsExists(t *testing.T, dbUrl string, dbPort int32, dbUsername string, dbPassword string, databaseName string, expectedSchemaName string) bool {
	output, err := AssertPotgresqlGrantsExistsE(t, dbUrl, dbPort, dbUsername, dbPassword, databaseName, expectedSchemaName)
	if err != nil {
		t.Fatal(err)
	}
	return output
}

func AssertPotgresqlGrantsExistsE(t *testing.T, dbUrl string, dbPort int32, dbUsername string, dbPassword string, databaseName string, expectedSchemaName string) (bool, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", dbUrl, dbPort, dbUsername, dbPassword, databaseName)

	db, connErr := sql.Open("pgx", connectionString)
	if connErr != nil {
		return false, connErr
	}
	defer db.Close()
	var (
		schemaName string
	)
	sqlStatement := `SELECT grantee AS user, CONCAT(table_schema, '.', table_name) AS table,
			CASE
				WHEN COUNT(privilege_type) = 7 THEN 'ALL'
				ELSE ARRAY_TO_STRING(ARRAY_AGG(privilege_type), ', ')
			END AS grants
		FROM information_schema.role_table_grants
		WHERE grantee = '$1'
		GROUP BY table_name, table_schema, grantee;`
	row := db.QueryRow(sqlStatement, dbUsername)
	scanErr := row.Scan(&schemaName)
	if scanErr != nil {
		return false, scanErr
	}
	return true, nil
}

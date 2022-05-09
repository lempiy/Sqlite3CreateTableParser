package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserNonQuote(t *testing.T) {
	const ddl = `
	CREATE TABLE contact_groups (
	 contact_id integer,
	 group_id integer,
	 PRIMARY KEY (contact_id, group_id),
	 FOREIGN KEY (contact_id) REFERENCES contacts (contact_id)
	 ON DELETE CASCADE ON UPDATE NO ACTION,
	 FOREIGN KEY (group_id) REFERENCES groups (group_id)
	 ON DELETE CASCADE ON UPDATE NO ACTION
	);
	`

	table, errCode := ParseTable(ddl, 0)
	assert.Equal(t, errCode, ERROR_NONE, "Parsing should work")
	assert.Len(t, table.Constraints, 3, "Should have 3 constraints")

}

func TestParserQouteIdentifier(t *testing.T) {
	const ddl = `
	CREATE TABLE "customer" (
		"id"	INTEGER NOT NULL,
		"first_name"	TEXT,
		"last_name"	TEXT,
		PRIMARY KEY("id", "first_name"),
		FOREIGN KEY ("id") REFERENCES "contacts" ("contact_id", "test_field")
	)
	`

	table, errCode := ParseTable(ddl, 0)
	assert.Equal(t, errCode, ERROR_NONE, "Parsing should work")
	assert.Len(t, table.Constraints, 2, "Should have 3 constraints")
}

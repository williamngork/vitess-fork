package schemadiff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSchemaFromQueriesViewWithCTEFail(t *testing.T) {
	queries := []string{"create view v30 as with vcte as (select 1) select * from vcte2"}
	_, err := NewSchemaFromQueries(NewTestEnv(), queries)
	assert.Error(t, err)
	assert.EqualError(t, err, (&ViewDependencyUnresolvedError{View: "v30", MissingReferencedEntities: []string{"dual", "vcte2"}}).Error())
}

func TestNewSchemaFromQueriesViewWithCTE(t *testing.T) {
	tcases := []struct {
		name    string
		queries []string
	}{
		{
			"no table",
			[]string{"create view v20 as with vcte as (select 1) select * from vcte"},
		},
		{
			"with table",
			[]string{
				"create table orders (id int primary key, info int not null)",
				"create view v21 as with vcte as (select * from orders) select * from vcte",
			},
		},
		{
			"with table and column aliasing",
			[]string{
				"create table orders (id int primary key, info int not null)",
				"create view v22 as with vcte as (select id, info as val from orders) select * from vcte",
			},
		},
		{
			"with table and select all from cte",
			[]string{
				"create table orders (id int primary key, info int not null)",
				"create view v22 as with vcte as (select id, info as val from orders) select vcte.* from vcte",
			},
		},
	}
	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewSchemaFromQueries(NewTestEnv(), tc.queries)
			assert.NoError(t, err)
		})
	}
}

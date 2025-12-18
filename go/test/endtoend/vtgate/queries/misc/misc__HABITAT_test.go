package misc

import (
	"github.com/stretchr/testify/require"
	"testing"
	"vitess.io/vitess/go/test/endtoend/utils"
)

// TestTabletTypeRouting tests that the tablet type routing works as intended.
func TestTabletTypeRouting__HABITAT(t *testing.T) {
	// We are gonna configure the routing rules to send the
	// query for a replica tablet in ks_misc.t1 to go to a table that doesn't exist.
	// I know this doesn't make much practical sense, but makes testing really easy.
	routingRules := `{"rules": [
	{
	"from_table": "ks_misc.t1@replica",
	"to_tables": ["uks.unknown"]
	}
]}`
	err := clusterInstance.VtctldClientProcess.ApplyRoutingRules(routingRules)
	require.NoError(t, err)
	defer func() {
		// Clear the routing rules after the test.
		err = clusterInstance.VtctldClientProcess.ApplyRoutingRules("{}")
		require.NoError(t, err)
	}()

	mcmp, closer := start(t)
	defer closer()

	mcmp.Exec("insert into t1(id1, id2) values (0,0)")

	vtConn := mcmp.VtConn
	// We first verify that querying the primary tablet goes to the t1 table.
	utils.Exec(t, vtConn, "use ks_misc@primary")
	utils.AssertMatches(t, vtConn, "select * from ks_misc.t1", `[[INT64(0) INT64(0)]]`)
	// Now we change the connection's target
	utils.Exec(t, vtConn, "use ks_misc@replica")
	// We verify that querying the replica tablet creates an unknown table error.
	_, err = utils.ExecAllowError(t, vtConn, "select * from ks_misc.t1")
	require.ErrorContains(t, err, "table unknown not found")
}

package planbuilder

import (
	"github.com/stretchr/testify/require"
	"vitess.io/vitess/go/test/vschemawrapper"
	"vitess.io/vitess/go/vt/key"
	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
	"vitess.io/vitess/go/vt/vtenv"
	"vitess.io/vitess/go/vt/vtgate/vindexes"
)

// TestForeignKeyPlanning tests the planning of foreign keys in a managed mode by Vitess.
func (s *planTestSuite) TestBypassPlanningKeyrangeTargetFromFile__HABITAT() {
	env := vtenv.NewTestEnv()
	vschema := loadSchema(s.T(), "vschemas/schema__HABITAT.json", true)
	vw, err := vschemawrapper.NewVschemaWrapper(env, vschema, TestBuilder)
	require.NoError(s.T(), err)

	keyRange, _ := key.ParseShardingSpec("-")
	vw.Dest = key.DestinationExactKeyRange{KeyRange: keyRange[0]}

	vw.Vcursor.SetTarget("main")
	vw.Keyspace = &vindexes.Keyspace{Name: "main"}

	s.testFile("bypass_keyrange__HABITAT_cases.json", vw, false)
}

func (s *planTestSuite) TestBypassPlanningShardTargetFromFile__HABITAT() {
	vschema := &vschemawrapper.VSchemaWrapper{
		V: loadSchema(s.T(), "vschemas/schema__HABITAT.json", true),
		Keyspace: &vindexes.Keyspace{
			Name:    "main",
			Sharded: false,
		},
		TabletType_: topodatapb.TabletType_PRIMARY,
		Dest:        key.DestinationShard("-80"),
		Env:         vtenv.NewTestEnv(),
	}

	s.testFile("bypass_shard__HABITAT_cases.json", vschema, false)
}

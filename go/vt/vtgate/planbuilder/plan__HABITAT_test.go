package planbuilder

import (
	"testing"
	"vitess.io/vitess/go/test/vschemawrapper"
	"vitess.io/vitess/go/vt/vtenv"
)

func (s *planTestSuite) TestForeignKeyChecksOn__HABITAT(t *testing.T) {
	vschema := loadSchema(s.T(), "vschemas/schema.json", true)
	s.setFks(vschema)
	fkChecksState := true
	vschemaWrapper := &vschemawrapper.VSchemaWrapper{
		V:                     vschema,
		TestBuilder:           TestBuilder,
		ForeignKeyChecksState: &fkChecksState,
		Env:                   vtenv.NewTestEnv(),
	}

	s.testFile("foreignkey_checks_on__HABITAT_cases.json", vschemaWrapper, false)
}

func (s *planTestSuite) TestForeignKeyPlanning__HABITAT(t *testing.T) {
	vschema := loadSchema(s.T(), "vschemas/schema.json", true)
	s.setFks(vschema)
	vschemaWrapper := &vschemawrapper.VSchemaWrapper{
		V:           vschema,
		TestBuilder: TestBuilder,
		Env:         vtenv.NewTestEnv(),
	}

	s.testFile("foreignkey__HABITAT_cases.json", vschemaWrapper, false)
}

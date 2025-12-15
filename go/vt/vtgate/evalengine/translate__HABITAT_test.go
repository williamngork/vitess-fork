package evalengine

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vitess.io/vitess/go/mysql/collations"
	"vitess.io/vitess/go/sqltypes"
	querypb "vitess.io/vitess/go/vt/proto/query"
	"vitess.io/vitess/go/vt/sqlparser"
	"vitess.io/vitess/go/vt/vtenv"
)

func TestEvaluate__HABITAT(t *testing.T) {
	type testCase struct {
		expression string
		expected   sqltypes.Value
	}

	True := sqltypes.NewInt64(1)
	False := sqltypes.NewInt64(0)
	tests := []testCase{{
		expression: "42",
		expected:   sqltypes.NewInt64(42),
	}, {
		expression: "42.42e0",
		expected:   sqltypes.NewFloat64(42.42),
	}, {
		expression: "40+2",
		expected:   sqltypes.NewInt64(42),
	}, {
		expression: "40-2",
		expected:   sqltypes.NewInt64(38),
	}, {
		expression: "40*2",
		expected:   sqltypes.NewInt64(80),
	}, {
		expression: "40/2",
		expected:   sqltypes.NewDecimal("20.0000"),
	}, {
		expression: ":exp",
		expected:   sqltypes.NewInt64(66),
	}, {
		expression: ":int32_bind_variable",
		expected:   sqltypes.NewInt64(20),
	}, {
		expression: ":uint32_bind_variable",
		expected:   sqltypes.NewUint64(21),
	}, {
		expression: ":uint64_bind_variable",
		expected:   sqltypes.NewUint64(22),
	}, {
		expression: ":string_bind_variable",
		expected:   sqltypes.NewVarChar("bar"),
	}, {
		expression: ":float_bind_variable",
		expected:   sqltypes.NewFloat64(2.2),
	}, {
		expression: "42 in (41, 42)",
		expected:   True,
	}, {
		expression: "42 in (41, 43)",
		expected:   False,
	}, {
		expression: "42 in (null, 41, 43)",
		expected:   NULL,
	}, {
		expression: "(1,2) in ((1,2), (2,3))",
		expected:   True,
	}, {
		expression: "(1,2) = (1,2)",
		expected:   True,
	}, {
		expression: "1 = 'sad'",
		expected:   False,
	}, {
		expression: "(1,2) = (1,3)",
		expected:   False,
	}, {
		expression: "(1,2) = (1,null)",
		expected:   NULL,
	}, {
		expression: "(1,2) in ((4,2), (2,3))",
		expected:   False,
	}, {
		expression: "(1,2) in ((1,null), (2,3))",
		expected:   NULL,
	}, {
		expression: "1 IN ::tuple_bind_variable",
		expected:   True,
	}, {
		expression: "3 IN ::tuple_bind_variable",
		expected:   True,
	}, {
		expression: "4 IN ::tuple_bind_variable",
		expected:   False,
	}, {
		expression: "(1,(1,2,3),(1,(1,2),4),2) = (1,(1,2,3),(1,(1,2),4),2)",
		expected:   True,
	}, {
		expression: "(1,(1,2,3),(1,(1,NULL),4),2) = (1,(1,2,3),(1,(1,2),4),2)",
		expected:   NULL,
	}, {
		expression: "null is null",
		expected:   True,
	}, {
		expression: "true is null",
		expected:   False,
	}, {
		expression: "42 is null",
		expected:   False,
	}, {
		expression: "null is not null",
		expected:   False,
	}, {
		expression: "42 is not null",
		expected:   True,
	}, {
		expression: "true is not null",
		expected:   True,
	}, {
		expression: "null is true",
		expected:   False,
	}, {
		expression: "42 is true",
		expected:   True,
	}, {
		expression: "true is true",
		expected:   True,
	}, {
		expression: "null is false",
		expected:   False,
	}, {
		expression: "42 is false",
		expected:   False,
	}, {
		expression: "false is false",
		expected:   True,
	}, {
		expression: "null is not true",
		expected:   True,
	}, {
		expression: "42 is not true",
		expected:   False,
	}, {
		expression: "true is not true",
		expected:   False,
	}, {
		expression: "null is not false",
		expected:   True,
	}, {
		expression: "42 is not false",
		expected:   True,
	}, {
		expression: "false is not false",
		expected:   False,
	}}

	venv := vtenv.NewTestEnv()
	for _, test := range tests {
		t.Run(test.expression, func(t *testing.T) {
			// Given
			stmt, err := sqlparser.NewTestParser().Parse("select " + test.expression)
			require.NoError(t, err)
			astExpr := stmt.(*sqlparser.Select).SelectExprs.Exprs[0].(*sqlparser.AliasedExpr).Expr
			sqltypesExpr, err := Translate(astExpr, &Config{
				Collation:   venv.CollationEnv().DefaultConnectionCharset(),
				Environment: venv,
			})
			require.Nil(t, err)
			require.NotNil(t, sqltypesExpr)
			env := NewExpressionEnv(context.Background(), map[string]*querypb.BindVariable{
				"exp":                  sqltypes.Int64BindVariable(66),
				"string_bind_variable": sqltypes.StringBindVariable("bar"),
				"int32_bind_variable":  sqltypes.Int32BindVariable(20),
				"uint32_bind_variable": sqltypes.Uint32BindVariable(21),
				"uint64_bind_variable": sqltypes.Uint64BindVariable(22),
				"float_bind_variable":  sqltypes.Float64BindVariable(2.2),
				"tuple_bind_variable": {
					Type: sqltypes.Tuple,
					Values: []*querypb.Value{
						{Type: sqltypes.Int64, Value: []byte("1")},
						{Type: sqltypes.Int64, Value: []byte("2")},
						{Type: sqltypes.Int64, Value: []byte("3")},
					},
				},
			}, NewEmptyVCursor(venv, time.Local))

			// When
			r, err := env.Evaluate(sqltypesExpr)

			// Then
			require.NoError(t, err)
			assert.Equal(t, test.expected, r.Value(collations.MySQL8().DefaultConnectionCharset()), "expected %s", test.expected.String())
		})
	}
}

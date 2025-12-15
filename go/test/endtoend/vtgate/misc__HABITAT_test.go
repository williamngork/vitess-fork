package vtgate

import (
	"testing"

	"vitess.io/vitess/go/test/endtoend/utils"
)

func TestFilterWithINAfterLeftJoin__HABITAT(t *testing.T) {
	conn, closer := start(t)
	defer closer()

	utils.Exec(t, conn, "insert into t1 (id1,id2) values (1, 10)")
	utils.Exec(t, conn, "insert into t1 (id1,id2) values (2, 3)")
	utils.Exec(t, conn, "insert into t1 (id1,id2) values (3, 2)")
	utils.Exec(t, conn, "insert into t1 (id1,id2) values (4, 5)")

	query := "select a.id1 from t1 as a left outer join t2 as b on a.id2 = b.id4 WHERE a.id2 = 12345 AND (b.id4 IS NULL OR b.id4 IN (1))"
	utils.AssertMatches(t, conn, query, `[]`)
}

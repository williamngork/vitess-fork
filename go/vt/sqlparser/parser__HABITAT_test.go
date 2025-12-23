package sqlparser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitStatementToPieces__HABITAT(t *testing.T) {
	testcases := []struct {
		input     string
		output    string
		lenWanted int
	}{{
		// Create procedure with definer current_user.
		input:     "create DEFINER=CURRENT_USER procedure p1 (in country CHAR(3))  begin declare abc DECIMAL(14,2); DECLARE def DECIMAL(14,2); end",
		lenWanted: 1,
	}, {
		// Create procedure with definer current_user().
		input:     "create DEFINER=CURRENT_USER() procedure p1 (in country CHAR(3))  begin declare abc DECIMAL(14,2); DECLARE def DECIMAL(14,2); end",
		lenWanted: 1,
	}, {
		// Create procedure with definer string.
		input:     "create DEFINER='root' procedure p1 (in country CHAR(3))  begin declare abc DECIMAL(14,2); DECLARE def DECIMAL(14,2); end",
		lenWanted: 1,
	}, {
		// Create procedure with definer string at_id.
		input:     "create DEFINER='root'@localhost procedure p1 (in country CHAR(3))  begin declare abc DECIMAL(14,2); DECLARE def DECIMAL(14,2); end",
		lenWanted: 1,
	}, {
		// Create procedure with definer id.
		input:     "create DEFINER=`root` procedure p1 (in country CHAR(3))  begin declare abc DECIMAL(14,2); DECLARE def DECIMAL(14,2); end",
		lenWanted: 1,
	}, {
		// Create procedure with definer id at_id.
		input:     "create DEFINER=`root`@`localhost` procedure p1 (in country CHAR(3))  begin declare abc DECIMAL(14,2); DECLARE def DECIMAL(14,2); end",
		lenWanted: 1,
	},
	}

	parser := NewTestParser()
	for _, tcase := range testcases {
		t.Run(tcase.input, func(t *testing.T) {
			if tcase.output == "" {
				tcase.output = tcase.input
			}

			stmtPieces, err := parser.SplitStatementToPieces(tcase.input)
			require.NoError(t, err)
			if tcase.lenWanted != 0 {
				require.Equal(t, tcase.lenWanted, len(stmtPieces))
			}
			out := strings.Join(stmtPieces, ";")
			require.Equal(t, tcase.output, out)
		})
	}
}

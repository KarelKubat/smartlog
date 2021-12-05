package linebuf

import (
	"testing"
)

func TestMethods(t *testing.T) {
	lb := New()
	for i, test := range []struct {
		addition               string
		wantCompleteBeforeStmt bool
		getStatement           bool
		wantStatement          string
		wantCompleteAfterStmt  bool
		wantFinal              string
	}{
		{
			// first part, no statement yet
			addition:               "hello",
			wantCompleteBeforeStmt: false,
			getStatement:           true,
			wantStatement:          "",
			wantCompleteAfterStmt:  false,
			wantFinal:              "hello",
		},
		{
			// second part, no statement yet
			addition:               " world",
			wantCompleteBeforeStmt: false,
			getStatement:           true,
			wantStatement:          "",
			wantCompleteAfterStmt:  false,
			wantFinal:              "hello world",
		},
		{
			// complete the statement
			addition:               "\n",
			wantCompleteBeforeStmt: true,
			getStatement:           false,
		},
		{
			// extract the statement
			addition:               "",
			wantCompleteBeforeStmt: true,
			getStatement:           true,
			wantStatement:          "hello world\n",
			wantCompleteAfterStmt:  false,
			wantFinal:              "",
		},
		{
			// add two statements and mush
			addition:               "ipsum dolor sit amet\nconsectetur adipiscing elit\nsed do eiusmod tempor",
			wantCompleteBeforeStmt: true,
			getStatement:           true,
			wantStatement:          "ipsum dolor sit amet\n",
			wantCompleteAfterStmt:  true,
			wantFinal:              "consectetur adipiscing elit\nsed do eiusmod tempor",
		},
		{
			// extract next part
			addition:               "",
			wantCompleteBeforeStmt: true,
			getStatement:           true,
			wantStatement:          "consectetur adipiscing elit\n",
			wantCompleteAfterStmt:  false,
			wantFinal:              "sed do eiusmod tempor",
		},
		{
			// nothing more to fetch
			addition:               "",
			wantCompleteBeforeStmt: false,
			getStatement:           true,
			wantStatement:          "",
			wantCompleteAfterStmt:  false,
			wantFinal:              "sed do eiusmod tempor",
		},
	} {
		if test.addition != "" {
			lb.Add([]byte(test.addition), len(test.addition))
		}
		if complete := lb.Complete(); complete != test.wantCompleteBeforeStmt {
			t.Errorf("iteration %v: Complete() before extracting statement = %v, want %v", i, complete, test.wantCompleteBeforeStmt)
		}
		if !test.getStatement {
			continue
		}
		if stmt := string(lb.Statement()); stmt != test.wantStatement {
			t.Errorf("iteration %v: Statement() = %v, want %v", i, stmt, test.wantStatement)
		}
		if complete := lb.Complete(); complete != test.wantCompleteAfterStmt {
			t.Errorf("iteration %v: Complete() after extracting statement = %v, want %v", i, complete, test.wantCompleteAfterStmt)
		}
		if final := string(lb.Bytes()); final != test.wantFinal {
			t.Errorf("iteration %v: final Bytes() = %v, want %v", i, final, test.wantFinal)
		}
	}
}

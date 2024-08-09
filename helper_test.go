package govm

import "testing"

type InOutCase struct {
	In, Expect string
}

func testInOut(t *testing.T, cases []InOutCase, fn func(s string) string) {
	for i, io := range cases {
		result := fn(io.In)
		if result != io.Expect {
			t.Errorf("Test case %d failed, input: %s, expect: %s, got: %s", i, io.In, io.Expect, result)
			return
		}
	}
}

var ToUnixPathTestCases = []InOutCase{
	{`D:\work\code\golearn`, `/d/work/code/golearn`},
	{`/d/work/code/golearn`, `/d/work/code/golearn`},
	{`C:\user\admin\AppData\Local\`, `/c/user/admin/AppData/Local`},
	{``, ``},
}

func TestToUnixPath(t *testing.T) {
	testInOut(t, ToUnixPathTestCases, ToUnixPath)
}

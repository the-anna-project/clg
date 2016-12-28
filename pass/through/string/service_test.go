package string

import (
	"testing"

	"github.com/the-anna-project/context"
)

func Test_Service_Action(t *testing.T) {
	testCases := []struct {
		A        string
		Expected string
	}{
		{
			A:        "",
			Expected: "",
		},
		{
			A:        " ",
			Expected: " ",
		},
		{
			A:        "   ",
			Expected: "   ",
		},
		{
			A:        "foo",
			Expected: "foo",
		},
		{
			A:        "test test",
			Expected: "test test",
		},
		{
			A:        "  test te     st   ",
			Expected: "  test te     st   ",
		},
		{
			A:        "..test,te  /  st / ",
			Expected: "..test,te  /  st / ",
		},
	}

	newService, err := NewService(DefaultServiceConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	for i, testCase := range testCases {
		action := newService.Action().(func(ctx context.Context, s string) string)
		s := action(nil, testCase.A)
		if s != testCase.Expected {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", s)
		}
	}
}

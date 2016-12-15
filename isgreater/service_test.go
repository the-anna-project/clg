package isgreater

import (
	"testing"

	"github.com/the-anna-project/context"
)

func Test_Service_Action(t *testing.T) {
	testCases := []struct {
		A        float64
		B        float64
		Expected bool
	}{
		{
			A:        3.5,
			B:        3.5,
			Expected: false,
		},
		{
			A:        3.5,
			B:        12.5,
			Expected: false,
		},
		{
			A:        35.5,
			B:        14.5,
			Expected: true,
		},
		{
			A:        -3.5,
			B:        7.5,
			Expected: false,
		},
		{
			A:        12.5,
			B:        4.5,
			Expected: true,
		},
		{
			A:        17,
			B:        65,
			Expected: false,
		},
		{
			A:        65,
			B:        17,
			Expected: true,
		},
	}

	newService, err := New(DefaultConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	for i, testCase := range testCases {
		action := newService.Action().(func(ctx context.Context, a, b float64) bool)
		b := action(nil, testCase.A, testCase.B)
		if b != testCase.Expected {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", b)
		}
	}
}

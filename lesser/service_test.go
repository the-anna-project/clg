package lesser

import (
	"testing"

	"github.com/the-anna-project/context"
)

func Test_Service_Action(t *testing.T) {
	testCases := []struct {
		A        float64
		B        float64
		Expected float64
	}{
		{
			A:        3.5,
			B:        3.5,
			Expected: 3.5,
		},
		{
			A:        12.5,
			B:        3.5,
			Expected: 3.5,
		},
		{
			A:        14.5,
			B:        35.5,
			Expected: 14.5,
		},
		{
			A:        7.5,
			B:        -3.5,
			Expected: -3.5,
		},
		{
			A:        4.5,
			B:        12.5,
			Expected: 4.5,
		},
		{
			A:        65,
			B:        17,
			Expected: 17,
		},
		{
			A:        17,
			B:        65,
			Expected: 17,
		},
	}

	newService, err := New(DefaultConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	for i, testCase := range testCases {
		action := newService.Action().(func(ctx context.Context, a, b float64) float64)
		f := action(nil, testCase.A, testCase.B)
		if f != testCase.Expected {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", f)
		}
	}
}

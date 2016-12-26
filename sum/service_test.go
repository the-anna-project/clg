package sum

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
			B:        12.5,
			Expected: 16,
		},
		{
			A:        35.5,
			B:        14.5,
			Expected: 50,
		},
		{
			A:        -3.5,
			B:        7.5,
			Expected: 4,
		},
		{
			A:        12.5,
			B:        4.5,
			Expected: 17,
		},
		{
			A:        36.5,
			B:        6.5,
			Expected: 43,
		},
		{
			A:        36,
			B:        6.5,
			Expected: 42.5,
		},
		{
			A:        99.99,
			B:        12.15,
			Expected: 112.14,
		},
		{
			A:        17,
			B:        65,
			Expected: 82,
		},
	}

	newService, err := NewService(DefaultServiceConfig())
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

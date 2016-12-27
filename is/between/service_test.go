package between

import (
	"testing"

	"github.com/the-anna-project/context"
)

func Test_Service_Action(t *testing.T) {
	testCases := []struct {
		N        float64
		Min      float64
		Max      float64
		Expected bool
	}{
		{
			N:        1,
			Min:      2,
			Max:      4,
			Expected: false,
		},
		{
			N:        2,
			Min:      2,
			Max:      4,
			Expected: true,
		},
		{
			N:        3,
			Min:      2,
			Max:      4,
			Expected: true,
		},
		{
			N:        4,
			Min:      2,
			Max:      4,
			Expected: true,
		},
		{
			N:        5,
			Min:      2,
			Max:      4,
			Expected: false,
		},
		{
			N:        35,
			Min:      -13,
			Max:      518,
			Expected: true,
		},
		{
			N:        -87,
			Min:      -413,
			Max:      -18,
			Expected: true,
		},
		{
			N:        -7,
			Min:      -413,
			Max:      -18,
			Expected: false,
		},
		{
			N:        -987,
			Min:      -413,
			Max:      -18,
			Expected: false,
		},
		{
			N:        1.8,
			Min:      2.34,
			Max:      4.944,
			Expected: false,
		},
		{
			N:        2.334,
			Min:      2.2,
			Max:      4.1,
			Expected: true,
		},
		{
			N:        3.9,
			Min:      2.003,
			Max:      4,
			Expected: true,
		},
		{
			N:        4,
			Min:      2.22,
			Max:      4.83,
			Expected: true,
		},
	}

	newService, err := NewService(DefaultServiceConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	for i, testCase := range testCases {
		action := newService.Action().(func(ctx context.Context, n, min, max float64) bool)
		b := action(nil, testCase.N, testCase.Min, testCase.Max)
		if b != testCase.Expected {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", b)
		}
	}
}

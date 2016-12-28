package float64

import (
	"testing"

	"github.com/the-anna-project/context"
)

func Test_Service_Action(t *testing.T) {
	testCases := []struct {
		A        float64
		Expected float64
	}{
		{
			A:        0,
			Expected: 0,
		},
		{
			A:        0.00,
			Expected: 0.00,
		},
		{
			A:        0.0001,
			Expected: 0.0001,
		},
		{
			A:        0.00305,
			Expected: 0.00305,
		},
		{
			A:        5.87,
			Expected: 5.87,
		},
		{
			A:        17.0,
			Expected: 17.0,
		},
		{
			A:        234,
			Expected: 234,
		},
		{
			A:        -2.34,
			Expected: -2.34,
		},
		{
			A:        -234,
			Expected: -234,
		},
	}

	newService, err := NewService(DefaultServiceConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	for i, testCase := range testCases {
		action := newService.Action().(func(ctx context.Context, f float64) float64)
		f := action(nil, testCase.A)
		if f != testCase.Expected {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", f)
		}
	}
}

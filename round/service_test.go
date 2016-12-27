package round

import (
	"testing"

	"github.com/the-anna-project/context"
)

func Test_Service_Action(t *testing.T) {
	testCases := []struct {
		Float     float64
		Precision int
		Expected  float64
	}{
		{
			Float:     3.5,
			Precision: 0,
			Expected:  4,
		},
		{
			Float:     3.4,
			Precision: 0,
			Expected:  3,
		},
		{
			Float:     3.4,
			Precision: 1,
			Expected:  3.4,
		},
		{
			Float:     3.4,
			Precision: 2,
			Expected:  3.4,
		},
		{
			Float:     3.476,
			Precision: 2,
			Expected:  3.48,
		},
		{
			Float:     -3.476,
			Precision: 2,
			Expected:  -3.48,
		},
		{
			Float:     3,
			Precision: 0,
			Expected:  3,
		},
		{
			Float:     3,
			Precision: 2,
			Expected:  3,
		},
	}

	newService, err := NewService(DefaultServiceConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	for i, testCase := range testCases {
		action := newService.Action().(func(ctx context.Context, f float64, p int) (float64, error))
		f, err := action(nil, testCase.Float, testCase.Precision)
		if err != nil {
			t.Fatal("case", i+1, "expected", nil, "got", err)
		}
		if f != testCase.Expected {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", f)
		}
	}
}

func Test_Service_Action_Error_NegativePrecision(t *testing.T) {
	newService, err := NewService(DefaultServiceConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	action := newService.Action().(func(ctx context.Context, f float64, p int) (float64, error))
	_, err = action(nil, 3.4465, -3)
	if !IsParseFloatSyntax(err) {
		t.Fatal("case", "expected", true, "got", false)
	}
}

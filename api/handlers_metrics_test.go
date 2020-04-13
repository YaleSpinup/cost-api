package api

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/YaleSpinup/cost-api/cloudwatch"
)

func TestParseQuery(t *testing.T) {
	type queryParseTest struct {
		query string
		input cloudwatch.MetricsRequest
		err   error
	}

	tests := []queryParseTest{
		// happy path
		{
			query: "",
			input: cloudwatch.MetricsRequest{
				"start":  "-P1D",
				"end":    "PT0H",
				"period": int64(300),
				"stat":   "Average",
			},
			err: nil,
		},
		{
			query: "start=-PT1H&end=PT0H&period=10s&stat=Maximum",
			input: cloudwatch.MetricsRequest{
				"start":  "-PT1H",
				"end":    "PT0H",
				"period": int64(10),
				"stat":   "Maximum",
			},
			err: nil,
		},
		{
			query: "start=-PT1H",
			input: cloudwatch.MetricsRequest{
				"start":  "-PT1H",
				"end":    "PT0H",
				"period": int64(300),
				"stat":   "Average",
			},
			err: nil,
		},
		{
			query: "end=-PT1H",
			input: cloudwatch.MetricsRequest{
				"start":  "-P1D",
				"end":    "-PT1H",
				"period": int64(300),
				"stat":   "Average",
			},
			err: nil,
		},
		{
			query: "period=10s",
			input: cloudwatch.MetricsRequest{
				"start":  "-P1D",
				"end":    "PT0H",
				"period": int64(10),
				"stat":   "Average",
			},
			err: nil,
		},
		{
			query: "stat=Minimum",
			input: cloudwatch.MetricsRequest{
				"start":  "-P1D",
				"end":    "PT0H",
				"period": int64(300),
				"stat":   "Minimum",
			},
			err: nil,
		},
		// errors
		{
			query: "period=true",
			input: cloudwatch.MetricsRequest{},
			err:   errors.New("failed to parse period as duration: time: invalid duration true"),
		},
	}

	for _, test := range tests {
		input := cloudwatch.MetricsRequest{}
		r := &http.Request{
			URL: &url.URL{
				RawQuery: test.query,
			},
		}

		t.Logf("testing raw query %s", test.query)

		if err := parseQuery(r, input); err != nil {
			if test.err == nil {
				t.Errorf("expected nil error, got %s", err)
			} else if test.err.Error() != err.Error() {
				t.Errorf("expected error %s, got %s", test.err, err)
			}
		} else {
			if !reflect.DeepEqual(input, test.input) {
				t.Errorf("expected %+v, got %+v", test.input, input)
			}
		}
	}
}

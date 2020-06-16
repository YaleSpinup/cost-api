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
				"height": int64(400),
				"width":  int64(600),
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
				"height": int64(400),
				"width":  int64(600),
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
				"height": int64(400),
				"width":  int64(600),
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
				"height": int64(400),
				"width":  int64(600),
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
				"height": int64(400),
				"width":  int64(600),
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
				"height": int64(400),
				"width":  int64(600),
			},
			err: nil,
		},
		{
			query: "height=100",
			input: cloudwatch.MetricsRequest{
				"start":  "-P1D",
				"end":    "PT0H",
				"period": int64(300),
				"stat":   "Average",
				"height": int64(100),
				"width":  int64(600),
			},
			err: nil,
		},
		{
			query: "width=200",
			input: cloudwatch.MetricsRequest{
				"start":  "-P1D",
				"end":    "PT0H",
				"period": int64(300),
				"stat":   "Average",
				"height": int64(400),
				"width":  int64(200),
			},
			err: nil,
		},
		{
			query: "height=2000&width=2000",
			input: cloudwatch.MetricsRequest{
				"start":  "-P1D",
				"end":    "PT0H",
				"period": int64(300),
				"stat":   "Average",
				"height": int64(2000),
				"width":  int64(2000),
			},
			err: nil,
		},
		// errors
		{
			query: "period=true",
			input: cloudwatch.MetricsRequest{},
			err:   errors.New("failed to parse period as duration: time: invalid duration true"),
		},
		{
			query: "height=-100",
			input: cloudwatch.MetricsRequest{},
			err:   errors.New("invalid height -100, value must be >=1 and <= 2000"),
		},
		{
			query: "height=0",
			input: cloudwatch.MetricsRequest{},
			err:   errors.New("invalid height 0, value must be >=1 and <= 2000"),
		},
		{
			query: "height=2001",
			input: cloudwatch.MetricsRequest{},
			err:   errors.New("invalid height 2001, value must be >=1 and <= 2000"),
		},
		{
			query: "width=-100",
			input: cloudwatch.MetricsRequest{},
			err:   errors.New("invalid width -100, value must be >=1 and <= 2000"),
		},
		{
			query: "width=0",
			input: cloudwatch.MetricsRequest{},
			err:   errors.New("invalid width 0, value must be >=1 and <= 2000"),
		},
		{
			query: "width=2001",
			input: cloudwatch.MetricsRequest{},
			err:   errors.New("invalid width 2001, value must be >=1 and <= 2000"),
		},
		{
			query: "height=2001&width=2001",
			input: cloudwatch.MetricsRequest{},
			err:   errors.New("invalid height 2001, value must be >=1 and <= 2000"),
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

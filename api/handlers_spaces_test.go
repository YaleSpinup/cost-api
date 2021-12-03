package api

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func TestParseTime(t *testing.T) {
	// use defaults derived in code
	startResult, endResult, err := parseTime("", "")
	if err != nil {
		t.Errorf("unexpected error from parseTime: %s", err)
	}

	// tests should match defaults from getTimeDefault
	y, m, d := time.Now().Date()
	if d == 1 {
		d = 5
	}

	sTime := fmt.Sprintf("%d-%02d-01", y, m)
	eTime := fmt.Sprintf("%d-%02d-%02d", y, m, d)

	if startResult == sTime {
		t.Logf("got expected default sTime: %s\n", sTime)
	} else {
		t.Errorf("got unexpected sTime: %s\n", sTime)
	}

	if endResult == eTime {
		t.Logf("got expected default eTime: %s\n", eTime)
	} else {
		t.Errorf("got unexpected eTime: %s\n", eTime)
	}

	// negative tests for non-matching defaults from getTimeDefault
	sTime = "2006-01-02"
	eTime = "2006-13-40"

	if startResult != sTime {
		t.Logf("negative test sTime: %s does not match: %s", sTime, startResult)
	} else {
		t.Errorf("got unexpected sTime: %s\n", sTime)
	}
	if endResult != eTime {
		t.Logf("negative test eTime: %s does not match: %s", eTime, endResult)
	} else {
		t.Errorf("got unexpected eTime: %s\n", eTime)
	}

	startTime := "2019-11-01"
	endTime := "2019-11-30"

	startResult, endResult, err = parseTime(startTime, endTime)
	if err != nil {
		t.Errorf("got unexpected error: %s", err)
	}
	if startResult == startTime {
		t.Logf("got expected startResult from getTimeAPI: %s", startResult)
	}
	if endResult == endTime {
		t.Logf("got expected endResult from getTimeAPI: %s", endResult)
	}

	// negative tests for non-matching API inputs from getTimeDefault
	// bad start time fails
	startTime = "2006-01-022"
	endTime = "2006-12-02"

	neg00startResult, neg00endResult, err := parseTime(startTime, endTime)
	if err != nil {
		t.Logf("negative test got expected error: %s", err)
	}
	if neg00startResult == startTime {
		t.Logf("negative test expected neg00_startResult from getTimeAPI: %s", neg00startResult)
	}
	if neg00endResult == endTime {
		t.Logf("negative test got expected neg00_endResult from getTimeAPI: %s", neg00endResult)
	}

	// bad end time fails
	startTime = "2006-01-02"
	endTime = "2006-12-403"

	neg01startResult, neg01endResult, err := parseTime(startTime, endTime)
	if err != nil {
		t.Logf("negative test got expected error: %s", err)
	}
	if neg01startResult == "" {
		t.Logf("negative test expected empty neg01startResult: %s", neg01startResult)
	}
	if neg01endResult == endTime {
		t.Logf("negative test got expected neg01endResult from getTimeAPI: %s", neg01endResult)
	}

	// start after end fails
	startTime = "2006-01-30"
	endTime = "2006-01-01"

	neg02startResult, neg02endResult, err := parseTime(startTime, endTime)
	if err != nil {
		t.Logf("negative test got expected error for start after end : %s", err)
	}
	if neg02startResult == startTime {
		t.Logf("negative test expected neg02startResult from getTimeAPI: %s", neg02startResult)
	}
	if neg02endResult == endTime {
		t.Logf("negative test got expected neg02endResult from getTimeAPI: %s", neg02endResult)
	}

}

func TestInSpace(t *testing.T) {
	spaceIDS := []string{
		"somespace-00012345",
		"yourspace-00012345",
	}

	for _, sid := range spaceIDS {
		expected := &costexplorer.Expression{
			Tags: &costexplorer.TagValues{
				Key: aws.String("spinup:spaceid"),
				Values: []*string{
					aws.String(sid),
				},
			},
		}
		out := inSpace(sid)
		if !awsutil.DeepEqual(expected, out) {
			t.Errorf("expected %s, got %s", awsutil.Prettify(expected), awsutil.Prettify(out))
		}

	}
}

func TestInOrg(t *testing.T) {
	orgs := []string{
		"ss",
		"sstst",
		"cool",
		"yourorg",
	}

	for _, o := range orgs {
		expected := &costexplorer.Expression{
			Or: []*costexplorer.Expression{
				{
					Tags: &costexplorer.TagValues{
						Key: aws.String("yale:org"),
						Values: []*string{
							aws.String(o),
						},
					},
				},
				{
					Tags: &costexplorer.TagValues{
						Key: aws.String("spinup:org"),
						Values: []*string{
							aws.String(o),
						},
					},
				},
			},
		}
		out := inOrg(o)
		if !awsutil.DeepEqual(expected, out) {
			t.Errorf("expected %s, got %s", awsutil.Prettify(expected), awsutil.Prettify(out))
		}

	}
}

func TestNotTryIT(t *testing.T) {
	expected := &costexplorer.Expression{
		Not: &costexplorer.Expression{
			Or: []*costexplorer.Expression{
				{
					Tags: &costexplorer.TagValues{
						Key: aws.String("yale:subsidized"),
						Values: []*string{
							aws.String("true"),
						},
					},
				},
				{
					Tags: &costexplorer.TagValues{
						Key: aws.String("spinup:subsidized"),
						Values: []*string{
							aws.String("true"),
						},
					},
				},
			},
		},
	}

	if out := notTryIT(); !reflect.DeepEqual(expected, out) {
		t.Errorf("expected %s, got %s", awsutil.Prettify(expected), awsutil.Prettify(out))
	}
}

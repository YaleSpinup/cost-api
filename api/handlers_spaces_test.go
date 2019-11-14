package api

import (
	"fmt"
	"testing"
	"time"
)

// getTimeDefault
func TestGetTimeDefault(t *testing.T) {
	// use defaults derived in code
	startTime, endTime, outBool := getTimeDefault()
	if outBool {
		t.Logf("bool of getTimeDefault output: %v\n", outBool)
	}
	if outBool == false {
		t.Errorf("bool of getTimeDefault false, expected true: %v\n", outBool)
	}

	// tests should match defaults from getTimeDefault
	y, m, d := time.Now().Date()
	if d == 1 {
		d = 5
	}

	sTime := fmt.Sprintf("%d-%02d-01", y, m)
	eTime := fmt.Sprintf("%d-%02d-%02d", y, m, d)

	if startTime == sTime {
		t.Logf("got expected default sTime: %s\n", sTime)
	} else {
		t.Errorf("got unexpected sTime: %s\n", sTime)
	}
	if endTime == eTime {
		t.Logf("got expected default eTime: %s\n", eTime)
	} else {
		t.Errorf("got unexpected eTime: %s\n", eTime)
	}

	// negative tests for non-matching defaults from getTimeDefault
	sTime = fmt.Sprint("2006-01-02")
	eTime = fmt.Sprint("2006-13-40")

	if startTime != sTime {
		t.Logf("negative test sTime: %s does not match: %s", sTime, startTime)
	} else {
		t.Errorf("got unexpected sTime: %s\n", sTime)
	}
	if endTime != eTime {
		t.Logf("negative test eTime: %s does not match: %s", eTime, endTime)
	} else {
		t.Errorf("got unexpected eTime: %s\n", eTime)
	}

}

// getTimeAPI
func TestGetTimeAPI(t *testing.T) {
	startTime := "2019-11-01"
	endTime := "2019-11-30"

	startResult, endResult, outBool, err := getTimeAPI(startTime, endTime)
	if err != nil {
		t.Errorf("got unexpected error: %s", err)
	}
	if startResult == startTime {
		t.Logf("got expected startResult from getTimeAPI: %s", startResult)
	}
	if endResult == endTime {
		t.Logf("got expected endResult from getTimeAPI: %s", endResult)
	}
	if outBool == true {
		t.Logf("got expected outBool from getTimeAPI: %v", outBool)
	}

	// negative tests for non-matching API inputs from getTimeDefault
	// bad start time fails
	sTime := fmt.Sprint("2006-01-022")
	eTime := fmt.Sprint("2006-12-02")

	neg00_startResult, neg00_endResult, neg00_outBool, err := getTimeAPI(sTime, eTime)
	if err != nil {
		t.Logf("negative test got expected error: %s", err)
	}
	if neg00_startResult == startTime {
		t.Logf("negative test expected neg00_startResult from getTimeAPI: %s", neg00_startResult)
	}
	if neg00_endResult == endTime {
		t.Logf("negative test got expected neg00_endResult from getTimeAPI: %s", neg00_endResult)
	}
	if neg00_outBool == true {
		t.Logf("negative test got expected neg00_outBool from getTimeAPI: %v", neg00_outBool)
	}

	// bad end time fails
	sTime = fmt.Sprint("2006-01-02")
	eTime = fmt.Sprint("2006-12-403")

	neg01_startResult, neg01_endResult, neg01_outBool, err := getTimeAPI(sTime, eTime)
	if err != nil {
		t.Logf("negative test got expected error: %s", err)
	}
	if neg01_startResult == "" {
		t.Logf("negative test expected empty neg01_startResult: %s", neg01_startResult)
	}
	if neg01_endResult == endTime {
		t.Logf("negative test got expected neg01_endResult from getTimeAPI: %s", neg01_endResult)
	}
	if neg01_outBool == false {
		t.Logf("negative test got expected neg01_bool false: %v", neg01_outBool)
	}

	// start after end fails
	sTime = fmt.Sprint("2006-01-30")
	eTime = fmt.Sprint("2006-01-01")

	neg02_startResult, neg02_endResult, neg02_outBool, err := getTimeAPI(sTime, eTime)
	if err != nil {
		t.Logf("negative test got expected error for start after end : %s", err)
	}
	if neg02_startResult == startTime {
		t.Logf("negative test expected neg02_startResult from getTimeAPI: %s", neg02_startResult)
	}
	if neg02_endResult == endTime {
		t.Logf("negative test got expected neg02_endResult from getTimeAPI: %s", neg02_endResult)
	}
	if neg02_outBool == true {
		t.Logf("negative test got expected neg02_outBool from getTimeAPI: %v", neg02_outBool)
	}
}

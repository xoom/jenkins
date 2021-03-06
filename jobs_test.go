package jenkins

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	jenkinsJobsResponse string = `
{
    "assignedLabels": [
        {}
    ], 
    "description": null, 
    "jobs": [
        {
            "color": "blue", 
            "name": "Jenkins Demo", 
	    "url": "http://build.example.com:8080/job/Jenkins%20Demo/"
        }, 
        {
            "color": "yellow", 
            "name": "cool-service", 
            "url": "http://build.example.com:8080/job/cool-service/"
        } 
    ], 
    "mode": "NORMAL", 
    "nodeDescription": "the master Jenkins node", 
    "nodeName": "", 
    "numExecutors": 2, 
    "overallLoad": {}, 
    "primaryView": {
        "name": "All", 
        "url": "http://platform-jenkins-master.qa.example.com:8080/"
    }, 
    "quietingDown": false, 
    "slaveAgentPort": 0, 
    "unlabeledLoad": {}, 
    "useCrumbs": false, 
    "useSecurity": false, 
    "views": [
        {
            "name": "All", 
            "url": "http://platform-jenkins-master.qa.example.com:8080/"
        }, 
        {
            "name": "DevOps", 
            "url": "http://platform-jenkins-master.qa.example.com:8080/view/DevOps/"
        }
    ]
}
`
)

func TestGetJobsNoError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := *r.URL
		if url.Path != "/api/json/jobs" {
			t.Fatalf("GetJobs() URL path expected to be /api/json/jobs but found %s\n", url.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Fatalf("GetJobs() expected request Accept header to be application/json but found %s\n", r.Header.Get("Accept"))
		}
		if r.Header.Get("Authorization") != "Basic dTpw" {
			t.Fatalf("Want Basic dTpw but got %s\n", r.Header.Get("Authorization"))
		}
		fmt.Fprintln(w, jenkinsJobsResponse)
	}))
	defer testServer.Close()

	url, _ := url.Parse(testServer.URL)
	jenkinsClient := NewClient(url, "u", "p")
	jobs, err := jenkinsClient.GetJobs()
	if err != nil {
		t.Fatalf("GetJobs() not expecting an error, but received: %v\n", err)
	}

	if len(jobs) != 2 {
		t.Fatalf("GetJobs() expected to return map of size 2, but received map of size %d\n", len(jobs))
	}

	expectedJobs := []string{"Jenkins Demo", "cool-service"}
	for _, p := range expectedJobs {
		_, present := jobs[p]
		if !present {
			t.Fatalf("GetJobs() expected to contain %s, but did not\n", p)
		}
	}
	if jobs["Jenkins Demo"].URL != "http://build.example.com:8080/job/Jenkins%20Demo/" {
		t.Fatalf("Want http://build.example.com:8080/job/Jenkins%%20Demo/ but got %s\n", jobs["Jenkins Demo"].URL)
	}
	if jobs["Jenkins Demo"].Color != "blue" {
		t.Fatalf("Want blue but got %s\n", jobs["Jenkins Demo"].Color)
	}

	if jobs["cool-service"].URL != "http://build.example.com:8080/job/cool-service/" {
		t.Fatalf("Want http://build.example.com:8080/job/cool-service/ but got %s\n", jobs["cool-service"].URL)
	}
	if jobs["cool-service"].Color != "yellow" {
		t.Fatalf("Want yellow but got %s\n", jobs["cool-service"].Color)
	}
}

func TestGetJobs500(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := *r.URL
		if url.Path != "/api/json/jobs" {
			t.Fatalf("GetJobs() URL path expected to be /api/json/jobs but found %s\n", url.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Fatalf("GetJobs() expected request Accept header to be application/json but found %s\n", r.Header.Get("Accept"))
		}
		if r.Header.Get("Authorization") != "Basic dTpw" {
			t.Fatalf("Want Basic dTpw but got %s\n", r.Header.Get("Authorization"))
		}
		w.WriteHeader(500)
	}))
	defer testServer.Close()

	url, _ := url.Parse(testServer.URL)
	jenkinsClient := NewClient(url, "u", "p")
	if _, err := jenkinsClient.GetJobs(); err == nil {
		t.Fatalf("GetJobs() expecting an error, but received none\n")
	}
}

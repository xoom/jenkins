package jenkins

import (
	"encoding/xml"
	"net/url"
)

type JobType int

const (
	Maven JobType = iota
	Freestyle
	Unknown
)

type (
	Jenkins interface {
		GetJobs() (map[string]JobDescriptor, error)
		GetJobConfig(jobName string) (JobConfig, error)
		GetJobSummaries() ([]JobSummary, error)
		GetJobSummariesFromFilesystem(root string) ([]JobSummary, error)
		GetLastBuild(jobName string) (LastBuild, error)
		CreateJob(jobName, jobConfigXML string) error
		DeleteJob(jobName string) error
	}

	Client struct {
		baseURL  *url.URL
		userName string
		password string
		Jenkins
	}

	JobDescriptor struct {
		Name  string `json:"name"`
		Color string `json:"color"`
		URL   string `json:"url"`
	}

	Jobs struct {
		Jobs []JobDescriptor `json:"jobs"`
	}

	// Maven project
	JobConfig struct {
		XMLName    xml.Name   `xml:"maven2-moduleset"`
		SCM        Scm        `xml:"scm"`
		Publishers Publishers `xml:"publishers"`
		RootModule RootModule `xml:"rootModule"`
		JobName    string
	}

	// Freestyle project
	FreeStyleJobConfig struct {
		XMLName xml.Name `xml:"project"`
		SCM     Scm      `xml:"scm"`
		JobName string
	}

	// Model of both Maven and Freestyle job types
	JobSummary struct {
		JobDescriptor JobDescriptor
		JobType       JobType
		GitURL        string // the use of this field is deprecated
		Branch        string // the use of this field is deprecated
	}

	Scm struct {
		XMLName           xml.Name          `xml:"scm"`
		Class             string            `xml:"class,attr"`
		UserRemoteConfigs UserRemoteConfigs `xml:"userRemoteConfigs"`
		Branches          Branches          `xml:"branches"`
	}

	Publishers struct {
		XMLName            xml.Name            `xml:"publishers"`
		RedeployPublishers []RedeployPublisher `xml:"hudson.maven.RedeployPublisher"`
	}

	RedeployPublisher struct {
		XMLName xml.Name `xml:"hudson.maven.RedeployPublisher"`
		URL     string   `xml:"url"`
	}

	UserRemoteConfigs struct {
		XMLName          xml.Name           `xml:"userRemoteConfigs"`
		UserRemoteConfig []UserRemoteConfig `xml:"hudson.plugins.git.UserRemoteConfig"`
	}

	UserRemoteConfig struct {
		XMLName xml.Name `xml:"hudson.plugins.git.UserRemoteConfig"`
		URL     string   `xml:"url"`
	}

	Branches struct {
		XMLName xml.Name `xml:"branches"`
		Branch  []Branch `xml:"hudson.plugins.git.BranchSpec"`
	}

	Branch struct {
		XMLName xml.Name `xml:"hudson.plugins.git.BranchSpec"`
		Name    string   `xml:"name"`
	}

	RootModule struct {
		XMLName    xml.Name `xml:"rootModule"`
		GroupID    string   `xml:"groupId"`
		ArtifactID string   `xml:"artifactId"`
	}

	LastBuild struct {
		Result          string `json:"result"`
		TimestampMillis int64  `json:"timestamp"`
		URL             string `json:"url"`
	}
)

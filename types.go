package main

// Metadata represents metatdata structure for network lists
type Metadata struct {
	ID    string             `json:"id"`
	Name  string             `json:"name"`
	Usage []ConfigurationMap `json:"usage"`
}

// ConfigurationMap represents configuration network lists usage structure
type ConfigurationMap struct {
	ConfigID     int              `json:"configId"`
	ConfigName   string           `json:"configName"`
	RatePolicies []RatePolicy     `json:"ratePolicies,omitempty"`
	Policies     []SecurityPolicy `json:"policies,omitempty"`
	MatchTargets []MatchTarget    `json:"match_targets,omitempty"`
}

// RatePolicy represents network lists usage in rate policies structure
type RatePolicy struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Condition string `json:"condition,omitempty"`
}

// MatchTarget represents network lists usage in match targets structure
type MatchTarget struct {
	ID                int64       `json:"id"`
	Hostnames         interface{} `json:"hostnames"`
	Paths             interface{} `json:"paths"`
	NegativePathMatch bool        `json:"negative-path-match"`
	Type              string      `json:"type"`
	SecurityPolicyID  string      `json:"securitypolicyid"`
}

// SecurityPolicy represents network lists usage in security policies structure
type SecurityPolicy struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Action string `json:"action"`
}

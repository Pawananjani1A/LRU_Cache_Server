/*
	Author: Pawananjani Kumar (pawananjani.kumar@goniyo.com)
	CreatedAt: 28 Mar 2024*/

package constants

type Environments string

const (
	EnvTesting     Environments = "test"
	EnvDevelopment Environments = "dev"
	EnvStaging     Environments = "uat"
	EnvBeta        Environments = "beta"
	EnvProd        Environments = "prod"
)

const (
	AppNameSuffixHTTPServer string = "httpserver"
	AppNameSuffixConsumer   string = "consumer"
)

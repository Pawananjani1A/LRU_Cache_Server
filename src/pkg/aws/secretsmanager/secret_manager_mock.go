/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package aws_secretsmanager

import (
	"context"
	GLogger "lruCache/poc/lib/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type mockSecretsManager struct {
	SecretValue *secretsmanager.GetSecretValueOutput
	Err         error
}

func (m *mockSecretsManager) GetSecretValueWithContext(ctx aws.Context, input *secretsmanager.GetSecretValueInput, opts ...request.Option) (*secretsmanager.GetSecretValueOutput, error) {
	return m.SecretValue, m.Err
}

// GetMockSecretManagerClient : input should be a stringified json
func GetMockSecretManagerClient(ctx context.Context, awsRegion string, logger *GLogger.LoggerService, input string, error error) ISecretManager {
	return &secretManager{
		client: &mockSecretsManager{
			SecretValue: &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(input),
			},
			Err: error,
		},
		log: logger,
	}
}

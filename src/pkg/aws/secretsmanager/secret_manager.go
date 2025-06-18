/*
	Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
	CreatedAt: 28 Mar 2024*/

package aws_secretsmanager

import (
	"context"
	"encoding/json"
	"fmt"
	GLogger "lruCache/poc/lib/logger"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
	awstrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/aws/aws-sdk-go/aws"
)

const region = "ap-south-1"

var secretsExpireTime = time.Now().Add(time.Hour * 24)

type ISecretManager interface {
	GetSecret(ctx context.Context, secretArn string) (map[string]interface{}, error)
}

type secretsManagerAPI interface {
	GetSecretValueWithContext(ctx aws.Context, input *secretsmanager.GetSecretValueInput, opts ...request.Option) (*secretsmanager.GetSecretValueOutput, error)
}

type secretManager struct {
	client secretsManagerAPI
	log    *GLogger.LoggerService
}
type secretManagerCache struct {
	data       secretsmanager.GetSecretValueOutput
	expireTime time.Time
}

var secretsCacheMap = make(map[string]secretManagerCache)
var secretsMutex = sync.RWMutex{}

func readSecretsCacheMap(key string) secretManagerCache {
	secretsMutex.RLock()
	defer secretsMutex.RUnlock()
	return secretsCacheMap[key]
}

func writeToSecretsCacheMap(key string, value secretManagerCache) {
	secretsMutex.Lock()
	defer secretsMutex.Unlock()
	secretsCacheMap[key] = value
}

func initializeSecretManager(ctx context.Context, awsRegion string, log *GLogger.LoggerService) *secretsmanager.SecretsManager {
	if awsRegion == "" {
		log.Debug(ctx, "awsRegion_not_found_in_initializeSecretManager_parameters", fmt.Sprintf("using %s as default awsRegion", awsRegion))
		awsRegion = region
	}
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: &awsRegion,
		},
	}))
	awsSession = awstrace.WrapSession(awsSession)
	return secretsmanager.New(awsSession)
}

func GetSecretManagerClient(ctx context.Context, awsRegion string, log *GLogger.LoggerService) ISecretManager {
	smClient := initializeSecretManager(ctx, awsRegion, log)
	return &secretManager{client: smClient, log: log}
}

func (s *secretManager) GetSecret(ctx context.Context, secretArn string) (map[string]interface{}, error) {
	var secretValueOutput secretsmanager.GetSecretValueOutput
	secretsFromCache := readSecretsCacheMap(secretArn)
	isEmptyCache := (secretsFromCache.data.SecretString == nil && secretsFromCache.data.SecretBinary == nil) || secretsFromCache.expireTime.IsZero()
	isExpiredCache := secretsFromCache.expireTime.After(secretsFromCache.expireTime)
	if isEmptyCache || isExpiredCache {
		req := &secretsmanager.GetSecretValueInput{SecretId: &secretArn}
		s.log.Debug(ctx, fmt.Sprintf("secret_fetch_request: %s", req))
		res, err := s.client.GetSecretValueWithContext(ctx, req)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("error_while_fetching_secrets_for_arn=%s", secretArn))
			s.log.Error(ctx, err)
			return nil, err
		}
		s.log.Info(ctx, fmt.Sprintf("successfullly_fetched_secrets_from_arn=%s", secretArn))
		secretValueOutput = *res
		writeToSecretsCacheMap(secretArn, secretManagerCache{
			data:       secretValueOutput,
			expireTime: secretsExpireTime,
		})
	} else {
		s.log.Debug(ctx, "secrets_fetched_from_cache")
		secretValueOutput = secretsFromCache.data
	}
	secrets, err := s.getValueFromSecretString(ctx, secretValueOutput.SecretString)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("error_occurred_while_getting_value_from_secret_string=%s", secretArn))
		s.log.Error(ctx, err)
		return nil, err
	}
	return secrets, nil
}

func (s *secretManager) getValueFromSecretString(ctx context.Context, secretString *string) (map[string]interface{}, error) {
	var secrets = make(map[string]interface{})
	err := json.Unmarshal([]byte(*secretString), &secrets)
	if err != nil {
		err = errors.Wrap(err, "error_while_unmarshalling_SecretString")
		s.log.Error(ctx, err)
		return nil, err
	}
	return secrets, nil
}

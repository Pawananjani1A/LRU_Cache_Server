/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package aws_secretsmanager

import (
	"context"
	GLogger "lruCache/poc/lib/logger"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type MetaDataProvider struct {
}

func (l *MetaDataProvider) AppName(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) UserID(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) CorrelationID(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) LoggerName(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) TraceID(ctx context.Context) string {
	return uuid.NewString()
}
func (l *MetaDataProvider) SpanID(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) Thread(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) ClientDeviceID(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) ClientPlatform(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) ClientVersion(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) AppPackageID(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) ClientSessionID(ctx context.Context) string {
	return ""
}

func getLogger(level GLogger.LogLevel) *GLogger.LoggerService {
	var metaDataProvider GLogger.ILoggingMetaDataProvider = &MetaDataProvider{}
	return GLogger.NewLoggerService(level, metaDataProvider)
}

const awsRegion = "ap-south-1"
const secretArn = "arn:aws:secretsmanager:ap-south-1:111122223333:secret:aes128-1a2b3c"

var secretString = "{\"client_id\": \"secret_client_id\",\"client_secret\": \"secret_client_sec\"}"

func TestSecretManager_GetSecret(t *testing.T) {
	logger := getLogger(GLogger.DEBUG)
	secretManagerClient := GetMockSecretManagerClient(context.Background(), awsRegion, logger, secretString, nil)
	secrets, err := secretManagerClient.GetSecret(context.Background(), secretArn)
	if err != nil {
		logger.Error(context.Background(), "error_occurred_while_fetching_secrets", err)
		return
	}
	logger.Info(context.Background(), "secrets", secrets)
}

func TestSecretManager_GetSecret2(t *testing.T) {
	logger := getLogger(GLogger.DEBUG)
	secretManagerClient := GetMockSecretManagerClient(context.Background(), "", logger, secretString, nil)
	secrets, err := secretManagerClient.GetSecret(context.Background(), secretArn)
	if err != nil {
		logger.Error(context.Background(), "error_occurred_while_fetching_secrets", err)
		return
	}
	logger.Info(context.Background(), "secrets", secrets)
}

func TestSecretManager_GetSecret3(t *testing.T) {
	logger := getLogger(GLogger.DEBUG)
	secretManagerClient := GetMockSecretManagerClient(context.Background(), "", logger, "", errors.New("this is a test error"))
	secrets, err := secretManagerClient.GetSecret(context.Background(), secretArn)
	if err != nil {
		logger.Error(context.Background(), "error_occurred_while_fetching_secrets", err)
		return
	}
	logger.Info(context.Background(), "secrets", secrets)
}

func TestGetValueFromSecretString(t *testing.T) {

	logger := getLogger(GLogger.DEBUG)
	awsSmClient := initializeSecretManager(context.Background(), awsRegion, logger)
	secretManagerClient := secretManager{
		client: awsSmClient,
		log:    logger,
	}
	secrets, err := secretManagerClient.getValueFromSecretString(context.Background(), &secretString)
	if err != nil {
		err = errors.Wrap(err, "error_occurred_while_getting_value_from_secret_string")
		logger.Error(context.Background(), err)
		return
	}
	logger.Info(context.Background(), secrets["client_id"])
	logger.Info(context.Background(), secrets["client_secret"])
}

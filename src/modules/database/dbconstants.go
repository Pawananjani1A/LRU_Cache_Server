/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package database

type DBType string
type AuthMethod string

const (
	MONGO DBType = "mongo"
)

const (
	AwsSecrets AuthMethod = "awssecrets"
	OpenURI    AuthMethod = "openuri"
	AwsIAM     AuthMethod = "awsiam"
)

const (
	DBErrInvalidResultsParamReceived string = "INVALID_RESULTS_RECEIVED"
	DBErrResponseDecodeErr           string = "RESPONSE_DECODE_ERROR"
	DBErrFindQueryError              string = "ERROR_WHILE_RUNNING_FIND_QUERY"
	DBErrBSONUnmarshalError          string = "BSON_UNMARSHALL_ERROR"
	DBErrMongoCursorIterationFailed  string = "MONGO_CURSOR_ITERATION_FAILED"
	DBErrNoDocumentsMatched          string = "NO_DOCUMENTS_MATCHED"
	DBErrUpdateFailed                string = "DB_UPDATE_FAILED"
)

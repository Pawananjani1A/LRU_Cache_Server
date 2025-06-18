/*
	Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
	CreatedAt: 28 Mar 2024*/

package httpclient

import (
	"bytes"
	"context"
	"fmt"
	GLogger "lruCache/poc/lib/logger"
	"lruCache/poc/src/constants"
	"lruCache/poc/src/internal/helpers"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	XCorrelationID string = "x-correlation-id"
	XRequestID     string = "x-request-id"
	XTraceID       string = "x-trace-id"
	XSpanID        string = "x-span-id"
)

type CustomHTTPClient struct {
	client     *http.Client
	logger     *GLogger.LoggerService
	clientName string
}

const (
	defaultTimeout = 10 * time.Second
)

func NewCustomHTTPClient(clientName string, logger *GLogger.LoggerService) *CustomHTTPClient {
	httpClient := &http.Client{Timeout: defaultTimeout}
	//httpClient = httptrace.WrapClient(httpClient)
	return &CustomHTTPClient{
		client:     httpClient,
		logger:     logger,
		clientName: clientName,
	}
}

func (hc *CustomHTTPClient) CallAPI(ctx context.Context, httpMethod string, url string, body []byte, headers http.Header) (*http.Response, []byte, error) {
	var (
		response   *http.Response
		bodyBytes  []byte
		requestErr error
	)
	if headers == nil {
		headers = make(http.Header)
	}
	headers.Add(string(constants.XTraceID), helpers.GetTraceIdFromContext(ctx))
	headers.Add(string(constants.XSpanID), helpers.GetSpanIdFromContext(ctx))
	headers.Add(string(constants.XCorrelationID), helpers.GetDefaultValueFromContext(ctx, string(constants.CorrelationID)))
	headers.Add(string(constants.XRequestID), string(constants.TraceID))
	switch httpMethod {
	case http.MethodGet:
		response, bodyBytes, requestErr = hc.callGETMethod(ctx, url, headers)
	case http.MethodPost:
		response, bodyBytes, requestErr = hc.callPOSTMethod(ctx, url, body, headers)
	default:
		requestErr = fmt.Errorf("invalid/unrecognized HTTP method passed in CallAPI")
	}
	if requestErr != nil {
		hc.logger.Error(ctx, "error_occurred_while_calling_api", "error", requestErr)
	}
	return response, bodyBytes, requestErr
}

func (hc *CustomHTTPClient) callGETMethod(ctx context.Context, url string, headers http.Header) (*http.Response, []byte, error) {
	method := "GET_API"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		hc.logRequestError(ctx, method, url, err)
		return nil, nil, err
	}
	return hc.executeRequest(ctx, method, url, request, headers)
}

func (hc *CustomHTTPClient) callPOSTMethod(ctx context.Context, url string, body []byte, headers http.Header) (*http.Response, []byte, error) {
	method := "POST_API"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		hc.logRequestError(ctx, method, url, err)
		return nil, nil, err
	}

	return hc.executeRequest(ctx, method, url, request, headers)
}

func (hc *CustomHTTPClient) executeRequest(ctx context.Context, method, url string, request *http.Request, headers http.Header) (*http.Response, []byte, error) {
	request.Header = headers
	hc.logger.Debug(ctx, "api_request", map[string]interface{}{
		"url":         url,
		"method":      method,
		"requestBody": request.Body,
		"headers":     request.Header,
	})
	response, err := hc.client.Do(request)
	if err != nil {
		hc.logRequestError(ctx, method, url, err)
		return nil, nil, err
	}
	defer response.Body.Close()
	hc.logger.Debug(ctx, fmt.Sprintf("completed_api_call_%s_for_url=%s", method, url))
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		hc.logger.Error(ctx, map[string]interface{}{
			"statusCode":   response.StatusCode,
			"url":          url,
			"errorMessage": "cannot_decode_response_body",
		}, err)
		return nil, nil, err
	} else {
		hc.logger.Info(ctx, "api_response", map[string]interface{}{
			"statusCode":   response.StatusCode,
			"url":          url,
			"method":       method,
			"responseBody": string(bodyBytes),
		}, err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		hc.logger.Error(ctx, "non_2xx_response_received_from_api", err)
	}
	return response, bodyBytes, nil
}
func (hc *CustomHTTPClient) logRequestError(ctx context.Context, method, url string, err error) {
	if os.IsTimeout(err) {
		hc.logger.Error(ctx, map[string]interface{}{
			"url":          url,
			"method":       method,
			"errorMessage": "request_timed_out",
		}, err)
	} else {
		hc.logger.Error(ctx, map[string]interface{}{
			"url":          url,
			"method":       method,
			"errorMessage": "error_encountered_while_calling",
		}, err)
	}
	hc.logger.Debug(ctx, fmt.Sprintf("url: %s", url))
}

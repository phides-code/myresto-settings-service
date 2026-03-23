package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator"
)

type ResponseStructure struct {
	Data         any     `json:"data"`
	ErrorMessage *string `json:"errorMessage"` // can be string or nil
}

var validate *validator.Validate = validator.New()

var headers = map[string]string{
	"Access-Control-Allow-Origin":  OriginURL,
	"Access-Control-Allow-Headers": "Content-Type, X-CF-Token, x-admin-key",
}

func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("router() received " + req.HTTPMethod + " request")

	headersJSON, _ := json.Marshal(req.Headers)
	log.Printf("router() received %s request, headers=%s", req.HTTPMethod, headersJSON)

	if !localMode {
		awsCfToken := os.Getenv("AWS_CF_TOKEN")

		if awsCfToken == "" {
			return serverError(errors.New("Error reading environment variable"))
		}

		providedCfToken := req.Headers["X-CF-Token"]

		if providedCfToken != awsCfToken {
			return clientError(http.StatusUnauthorized)
		}
	}

	switch req.HTTPMethod {
	case "GET":
		return processGet(ctx, req)
	case "POST":
		return processPost(ctx, req)
	case "OPTIONS":
		return processOptions()
	default:
		log.Println("router() error parsing HTTP method")
		return clientError(http.StatusMethodNotAllowed)
	}
}

func processOptions() (events.APIGatewayProxyResponse, error) {
	additionalHeaders := map[string]string{
		"Access-Control-Allow-Methods": "OPTIONS, POST, GET",
		"Access-Control-Max-Age":       "3600",
	}
	mergedHeaders := mergeHeaders(headers, additionalHeaders)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    mergedHeaders,
	}, nil
}

func processGet(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if strings.HasPrefix(req.Resource, "/settings/validateadminkey") {
		return processValidateAdminKey(req)
	}

	if strings.HasPrefix(req.Resource, "/settings/themesettings") {
		return processGetThemeSettings(ctx)
	}

	if strings.HasPrefix(req.Resource, "/settings/themes") {
		return processGetThemes(ctx)
	}

	return processGetSettings(ctx)
}

func processPost(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if strings.HasPrefix(req.Resource, "/settings/themesettings") {
		return handleAdminOnly(ctx, req, processPostThemeSettings)
	}

	return handleAdminOnly(ctx, req, processPostSettings)
}

func processGetThemeSettings(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	themeSettingsBytes, err := myS3.DownloadFile(ctx, themeSettingsFilename)
	if err != nil {
		return serverError(err)
	}

	var themeSettings ThemeSettings
	err = json.Unmarshal(themeSettingsBytes, &themeSettings)
	if err != nil {
		log.Printf("processGetThemeSettings() Can't unmarshal themeSettings: %v", err)
		return serverError(err)
	}

	err = validate.Struct(&themeSettings)
	if err != nil {
		log.Printf("Invalid body: %v", err)
		return serverError(err)
	}

	response := ResponseStructure{
		Data:         themeSettings,
		ErrorMessage: nil,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println("processGetThemeSettings() error running json.Marshal")
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJson),
		Headers:    headers,
	}, nil
}

func processGetThemes(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	themesBytes, err := myS3.DownloadFile(ctx, themesFilename)
	if err != nil {
		return serverError(err)
	}

	var themes []Theme
	err = json.Unmarshal(themesBytes, &themes)
	if err != nil {
		log.Printf("processGetThemes() Can't unmarshal themes: %v", err)
		return serverError(err)
	}

	for i := range themes {
		if err := validate.Struct(&themes[i]); err != nil {
			log.Printf("processGetThemes() Invalid body: %v", err)
			return serverError(err)
		}
	}

	response := ResponseStructure{
		Data:         themes,
		ErrorMessage: nil,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println("processGetThemes() error running json.Marshal")
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJson),
		Headers:    headers,
	}, nil
}

func processGetSettings(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	settingsBytes, err := myS3.DownloadFile(ctx, settingsFilename)
	if err != nil {
		return serverError(err)
	}

	var settings Settings
	err = json.Unmarshal(settingsBytes, &settings)
	if err != nil {
		log.Printf("processGetSettings() Can't unmarshal settings: %v", err)
		return serverError(err)
	}

	err = validate.Struct(&settings)
	if err != nil {
		log.Printf("Invalid body: %v", err)
		return serverError(err)
	}

	response := ResponseStructure{
		Data:         settings,
		ErrorMessage: nil,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println("processGetSettings() error running json.Marshal")
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJson),
		Headers:    headers,
	}, nil
}

func processPostThemeSettings(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var newThemeSettings ThemeSettings
	err := json.Unmarshal([]byte(req.Body), &newThemeSettings)
	if err != nil {
		log.Printf("processPostThemeSettings() Can't unmarshal body: %v", err)
		return clientError(http.StatusUnprocessableEntity)
	}

	err = validate.Struct(&newThemeSettings)
	if err != nil {
		log.Printf("processPostThemeSettings() Invalid body: %v", err)
		return clientError(http.StatusBadRequest)
	}

	newThemeSettingsBytes, err := json.Marshal(newThemeSettings)
	if err != nil {
		log.Printf("processPostThemeSettings() Can't marshal themeSettings: %v", err)
		return serverError(err)
	}

	fileName, err := myS3.UploadFile(bytes.NewReader(newThemeSettingsBytes), themeSettingsFilename)
	if err != nil {
		return serverError(err)
	}

	response := ResponseStructure{
		Data:         &fileName,
		ErrorMessage: nil,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println("processPostThemeSettings() error running json.Marshal")
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(responseJson),
		Headers:    headers,
	}, nil
}

func processPostSettings(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var newSettings Settings
	err := json.Unmarshal([]byte(req.Body), &newSettings)
	if err != nil {
		log.Printf("processPostSettings() Can't unmarshal body: %v", err)
		return clientError(http.StatusUnprocessableEntity)
	}

	err = validate.Struct(&newSettings)
	if err != nil {
		log.Printf("processPostSettings() Invalid body: %v", err)
		return clientError(http.StatusBadRequest)
	}

	newSettingsBytes, err := json.Marshal(newSettings)
	if err != nil {
		log.Printf("processPostSettings() Can't marshal settings: %v", err)
		return serverError(err)
	}

	fileName, err := myS3.UploadFile(bytes.NewReader(newSettingsBytes), settingsFilename)
	if err != nil {
		return serverError(err)
	}

	response := ResponseStructure{
		Data:         &fileName,
		ErrorMessage: nil,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println("processPostSettings() error running json.Marshal")
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(responseJson),
		Headers:    headers,
	}, nil
}

func processValidateAdminKey(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	adminKey := os.Getenv("ADMIN_KEY")

	if adminKey == "" {
		return serverError(errors.New("processValidateAdminKey(): Error reading environment variable"))
	}

	var providedAdminKey string

	if localMode {
		providedAdminKey = req.Headers["X-Admin-Key"]
	} else {
		providedAdminKey = req.Headers["x-admin-key"]
	}

	validity := providedAdminKey == adminKey

	response := ResponseStructure{
		Data:         validity,
		ErrorMessage: nil,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println("processValidateAdminKey() error running json.Marshal")
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJson),
		Headers:    headers,
	}, nil
}

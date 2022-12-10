package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/josemyduarte/printer"
	"github.com/josemyduarte/printer/internal/aws"
)

var (
	fontURL      = "https://github.com/DevsQuotes/DevsQuotesPrinter/raw/master/assets/FiraSans-Light.ttf"
	fontFileName = "/tmp/FiraSans-Light.ttf"

	bgURL      = "https://github.com/DevsQuotes/DevsQuotesPrinter/raw/master/assets/00-instagram-background.png"
	bgFileName = "/tmp/00-instagram-background.png"
)

func init() {
	if err := downloadFile(fontURL, fontFileName); err != nil {
		panic(fmt.Errorf("couldn't download font: %w", err))
	}

	if err := downloadFile(bgURL, bgFileName); err != nil {
		panic(fmt.Errorf("couldn't download background image: %w", err))
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req struct {
		Text string `json:"text"`
	}

	b := []byte(request.Body)
	err := json.Unmarshal(b, &req)
	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, err
	}

	return aws.Serve(printer.Request{
		BgImgPath: bgFileName,
		FontPath:  fontFileName,
		FontSize:  60,
		Text:      req.Text,
	})
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}

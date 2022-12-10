package aws

import (
	"bytes"
	"encoding/base64"
	"image/jpeg"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/josemyduarte/printer"
)

func Serve(request printer.Request) (events.APIGatewayProxyResponse, error) {
	img, err := printer.TextOnImg(request)
	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}

	imgBuf := new(bytes.Buffer)
	if jpeg.Encode(imgBuf, img, nil) != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		Body:            base64.StdEncoding.EncodeToString(imgBuf.Bytes()),
		IsBase64Encoded: true,
		Headers: map[string]string{
			"Content-Type": "image/png",
		},
	}, nil
}

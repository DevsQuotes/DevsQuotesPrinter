package adapters

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/josemyduarte/printer"
)

var assets = printer.Assets{
	BgImgPath: "../../assets/00-instagram-background.png",
	FontPath:  "../../assets/FiraSans-Light.ttf",
	FontSize:  60,
}

func Serve(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Print(fmt.Sprintf("body:[%s] ", request.Body))

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

	img, err := printer.TextOnImg(
		printer.Request{
			BgImgPath: assets.BgImgPath,
			FontPath:  assets.FontPath,
			FontSize:  assets.FontSize,
			Text:      req.Text,
		},
	)
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
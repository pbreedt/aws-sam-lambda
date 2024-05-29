package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestPayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type ResponsePayload struct {
	Message string `json:"message"`
	From    string `json:"from-ip"`
}

var ErrNon200Response = fmt.Errorf("non 200 response")
var ErrNoIP = fmt.Errorf("no IP returned")

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reply := ResponsePayload{}

	// Part 1: get data from request
	var payload RequestPayload
	err := json.Unmarshal([]byte(request.Body), &payload)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Error: " + err.Error(),
			StatusCode: 400,
		}, nil
	}
	reply.Message = fmt.Sprintf("Hello, %s (aged %d) ", payload.Name, payload.Age)

	// Part 2: get data from APIGatewayProxyRequestContext
	sourceIP := request.RequestContext.Identity.SourceIP
	if sourceIP == "" {
		reply.Message += "from unknown source!"
	} else {
		reply.Message += fmt.Sprintf("from %s!", sourceIP)
	}

	// Part 3: get data from backend service
	resp, err := http.Get("https://checkip.amazonaws.com/")
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	if len(ip) == 0 {
		reply.From = "unknown"
	} else {
		reply.From = strings.TrimSpace(string(ip))
	}

	// Part 4: return APIGatewayProxyResponse
	return events.APIGatewayProxyResponse{
		Body:       reply.String(),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

func (r ResponsePayload) String() string {
	return fmt.Sprintf("{\"message\":\"%s\",\"from\":\"%s\"}", r.Message, r.From)
}

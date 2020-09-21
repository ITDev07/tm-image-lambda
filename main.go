package main


import (
	"fmt"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)


// BodyRequest is our self-made struct to process JSON request from Client
type BodyRequest struct {
	RequestName string `json:"name"`
	MatchTMSID string `json:"matchTMSID"`
	ClientTMSID string `json:"clientTMSID"`
	ValidationKey string `json:"validationKey"`
	AccessToken string `json:"accessToken"`
}

// BodyResponse is our self-made struct to build response for Client
type BodyResponse struct {
	ResponseName string `json:"name"`
	TotalScore float64 `json:"totalScore"`
	City string `json:"city"`
	Firstname string `json:"firstname"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// events.APIGatewayProxyRequest.MultiValueHeaders.
	// events.APIGatewayProxyRequest.MultiValueQueryStringParameters.
	// events.APIGatewayRequestIdentity.AccessKey (populated whether doing an authorization using AWS credentials).

	// BodyRequest will be used to take the json response from client and build it
	bodyRequest := BodyRequest{
		RequestName: "",
	}

	// Unmarshal the json, return 404 if error
	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	// We will build the BodyResponse and send it back in json form
	bodyResponse := BodyResponse{
		ResponseName: bodyRequest.RequestName + " LastName",
	}

	validation_key := bodyRequest.ValidationKey + ".txt"
	access_token := bodyRequest.AccessToken
	success, err := isValidRequest("validation-tm-dev", validation_key, access_token)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	if success == true {
		MatchJSON, _ := getMatchJSON("ids-tm-dev", "BCBC-BBB.json")

		match_score := bodyRequest.MatchTMSID[0:4]
		client_score := bodyRequest.ClientTMSID[0:4]
		cp := decodeScore(client_score)
	    mp := decodeScore(match_score)
	    result := compareScores(cp, mp)

        fmt.Println(result)
		fmt.Println(MatchJSON)

		bodyResponse = BodyResponse{
			ResponseName: bodyRequest.RequestName,
			TotalScore: result.total_score,
			City: MatchJSON["city"],
			Firstname: MatchJSON["firstname"],
		}
	}

	// Marshal the response into json bytes, if error return 404
	response, err := json.Marshal(&bodyResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
}

func main() {
	loadEnv()
	lambda.Start(Handler)
}

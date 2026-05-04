// holds the actual Go server code and the endpoints

package main    //tells go this is the runnable application package

import (
	"encoding/json"
	"log"
	"net/http"
)

type RiskRequest struct {
	AssetID                 int    `json:"assetId"`
	Criticality             string `json:"criticality"`
	CriticalVulnerabilities int    `json:"criticalVulnerabilities"`
	HighVulnerabilities     int    `json:"highVulnerabilities"`
	MediumVulnerabilities   int    `json:"mediumVulnerabilities"`
	LowVulnerabilities      int    `json:"lowVulnerabilities"`
}

type RiskResponse struct {
	AssetID   int    `json:"assetId"`   //asset Id returned backed to the called
	RiskScore int    `json:"riskScore"` //Final calculated risk score
	RiskLevel string `json:"riskLevel"` //final risk label like Low, Medium, High, or Critical
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //set response type
	w.WriteHeader(http.StatusOK)                       //set http status code

	response := map[string]string{ //creates a small response object
		"status": "ok",
	}
	json.NewEncoder(w).Encode(response) //convert response to JSON and send it
}

//function runs when someone call the /calculate-risk endpoint
func calculateRiskHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") //sets the response header

	var request RiskRequest //creates request with type RiskRequest

	// Decode the incoming JSON request into the RiskRequest struct
	err := json.NewDecoder(r.Body).Decode(&request) //Take the incoming JSON and load it into the request struct

	if err != nil { //if JSON body was invalid, go into this block
		w.WriteHeader(http.StatusBadRequest) //send HTTP status 400

		response := map[string]string{ //send error message back to the client as JSON
			"error": "Invalid request body",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate required values before calculating risk
	if request.AssetID <= 0 { //assetId cannot be negative
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "assetId must be greater than 0",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.Criticality != "Low" && request.Criticality != "Medium" && request.Criticality != "High" {
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "criticality must be Low, Medium, or High",
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	if request.CriticalVulnerabilities < 0 || request.HighVulnerabilities < 0 ||
		request.MediumVulnerabilities < 0 || request.LowVulnerabilities < 0 {
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "vulnerability counts cannot be negative",
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	// Calculate the weighted risk score and map it to a risk level
	riskScore := (request.CriticalVulnerabilities * 25) +
		(request.HighVulnerabilities * 15) +
		(request.MediumVulnerabilities * 8) +
		(request.LowVulnerabilities * 3)

	if request.Criticality == "Medium" {
		riskScore += 10
	} else if request.Criticality == "High" {
		riskScore += 20
	}

	if riskScore > 100 {
		riskScore = 100
	}

	riskLevel := "Low"

	if riskScore >= 76 {
		riskLevel = "Critical"
	} else if riskScore >= 51 {
		riskLevel = "High"
	} else if riskScore >= 26 {
		riskLevel = "Medium"
	}

	// Build and return the JSON response with the calculated result
	response := RiskResponse{
		AssetID:   request.AssetID,
		RiskScore: riskScore,
		RiskLevel: riskLevel,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/health", healthHandler)                 //connects the /health URL to that function
	http.HandleFunc("/calculate-risk", calculateRiskHandler)  //connects the /calculate-risk URL to that function
	log.Println("Risk service running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil)) //starts the web server on port 8081
}

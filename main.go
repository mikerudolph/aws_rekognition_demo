package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

// RequestBody is the parsed body from the post request.
type RequestBody struct {
	Image string
}

func main() {
	sess := session.New(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	svc := rekognition.New(sess)

	http.Handle("/", http.FileServer(http.Dir("./html")))

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "Request Body Expected", 400)
			return
		}

		var parsed RequestBody

		err := json.NewDecoder(r.Body).Decode(&parsed)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}

		// Decode the string.
		decodedImage, err := base64.StdEncoding.DecodeString(parsed.Image)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Send request to Rekognition.
		input := &rekognition.DetectLabelsInput{
			Image: &rekognition.Image{
				Bytes: decodedImage,
			},
		}

		result, err := svc.DetectLabels(input)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		output, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(output)
	})

	http.ListenAndServe(":3000", nil)
}

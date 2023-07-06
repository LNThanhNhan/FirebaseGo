package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Response struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}

func MakeSuccessResponse(Data interface{}) []byte {
	response := Response{
		Success: true,
		Message: "Success",
		Data:    []interface{}{Data},
	}
	log.Printf("%+v\n", response)
	jsonData, err := json.Marshal(response)
	if err != nil {
		fmt.Print("Error coverting data to JSON")
	}
	return jsonData
}

func MakeFailResponse(error string) []byte {
	response := Response{
		Success: false,
		Message: error,
		Data:    nil,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		fmt.Print("Error coverting data tog JSON")
	}
	return jsonData
}

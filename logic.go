package main

import (
	"fmt"
	"net/http"
)

var nutritionix string = "https://api.nutritionix.com/v1_1/search/"
var appId string = "c9396ee1"
var appKey string = "3054b427133d447c8bc2826d41ef5574"

func sendAPIRequest(arg http.ResponseWriter, arg2 *http.Request) {
	resp, _ := http.Get(nutritionix + "taco?appId=" + appId + "&appKey=" + appKey)
	fmt.Println(resp.Body)
	// http.Get
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	// "github.com/ramin0/chatbot"

	// Autoload environment variables in .env
	_ "github.com/joho/godotenv/autoload"
	"github.com/omarayousry/chatbot"
)

func getJSONArray(res *http.Response, arrayString string) []map[string]interface{} {
	defer res.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(res.Body).Decode(&data)
	referArray := data[arrayString].([]interface{})
	// var nestedData [len(referArray)]map[string]interface{}
	nestedData := make([]map[string]interface{}, len(referArray))
	for index, dataEntry := range referArray {
		nestedData[index] = dataEntry.(map[string]interface{})
	}
	// st :
	// nestedData := data["results"].([]map[string]interface{})
	return nestedData
}

func getResponse(baseUrl string, params map[string]string, questionString string) (*http.Response, error) {
	var searchString string
	for key, value := range params {
		searchString += key + "=" + value + "&"
	}
	res, err := http.Get(baseUrl + "/" + questionString + "?" + searchString)
	if err != nil {
		return res, fmt.Errorf("THERE WAS A PROBLEM ESTABLISHING CONNECTING USING " + baseUrl + "/" + questionString + "?" + searchString + "The error is" + err.Error())
	}
	return res, nil
}

func chatbotProcess(session chatbot.Session, message string) (string, error) {

	//checks for invalid characters in the user's message
	if strings.ContainsAny(message,
		`123456789,.;()!@#$%^&*[]{}\\|:?><`) {
		return "", fmt.Errorf("Whoops! Please only enter valid answers to the question! No symbols or numbers!")
	}
	if strings.EqualFold(message, "chatbot") {
		return "", fmt.Errorf("This can't be, I'm the one and only %s!", message)
	}
	if reflect.DeepEqual(session["toot"], []string{"THE EPIC MAN"}) {
		fmt.Println("EYYY")
	}
	//
	_, err := getResponse("www.recipepuppy.com/api", map[string]string{"i": "garlic"}, "")
	return err.Error(), nil
	//data := getJSONArray(getResponse("www.recipepuppy.com/api", map[string]string{"i": "garlic"}, ""), "results")

	// return fmt.Sprintf("Hello %s, my name is chatbot. What was yours again?", message), nil
	// resp, _ := http.Get(nutritionix + "taco?appId=" + appId + "&appKey=" + appKey)
	//	res, _ := http.Get("http://www.recipepuppy.com/api/?i=onions,garlic&q=omelet&p=3")
	//	defer res.Body.Close()
	// body, _ := ioutil.ReadAll(resp.Body)
	//	data := getJSONArray(res, "results")
	// var s string
	// for someVar := range body {
	// s += string(someVar)
	// }
	// return s, nil
	// return resp.H, nil

	// return data[0]["title"].(string), nil
}

func main() {
	// Uncomment the following lines to customize the chatbot
	chatbot.WelcomeMessage = `Hello. I am the GUC Sofra AI.
  I will help you in determining the best
  recipes based on the nutritional values you want
  by asking you a series of questions.`
	chatbot.ProcessFunc(chatbotProcess)

	// Use the PORT environment variable
	port := os.Getenv("PORT")
	// Default to 3000 if no PORT environment variable was defined
	if port == "" {
		port = "3000"
	}
	// http.HandleFunc("/", handleDefault)
	fmt.Printf("Listening on port %s...\n", port)
	// http.ListenAndServe(":3000", nil)
	// Start the server
	log.Fatalln(chatbot.Engage(":" + port))
}

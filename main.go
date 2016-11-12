package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	// "github.com/ramin0/chatbot"

	// Autoload environment variables in .env
	_ "github.com/joho/godotenv/autoload"
	"github.com/ramin0/chatbot"
)

func getJSONArray(res *http.Response, arrayString string) []map[string]interface{} {
	defer res.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(res.Body).Decode(&data)
	nestedData := data[arrayString].([]map[string]interface{})
	return nestedData
}

func chatbotProcess(session chatbot.Session, message string) (string, error) {
	// if strings.EqualFold(message, "chatbot") {
	// 	return "", fmt.Errorf("This can't be, I'm the one and only %s!", message)
	// }
	//
	// return fmt.Sprintf("Hello %s, my name is chatbot. What was yours again?", message), nil
	// resp, _ := http.Get(nutritionix + "taco?appId=" + appId + "&appKey=" + appKey)
	res, _ := http.Get("http://www.recipepuppy.com/api/?i=onions,garlic&q=omelet&p=3")
	defer res.Body.Close()
	// body, _ := ioutil.ReadAll(resp.Body)
	data := getJSONArray(res, "results")
	// var s string
	// for someVar := range body {
	// s += string(someVar)
	// }
	// return s, nil
	// return resp.H, nil
	return data[0]["title"].(string), nil
}

func main() {
	// Uncomment the following lines to customize the chatbot
	chatbot.WelcomeMessage = "Hey! What's your name?"
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

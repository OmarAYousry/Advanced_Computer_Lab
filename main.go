package main

import (
	"fmt"
	"log"
	"os"

	// "github.com/ramin0/chatbot"

	// Autoload environment variables in .env
	_ "github.com/joho/godotenv/autoload"
	"github.com/ramin0/chatbot"
)

func chatbotProcess(session chatbot.Session, message string) (string, error) {
	// if strings.EqualFold(message, "chatbot") {
	// 	return "", fmt.Errorf("This can't be, I'm the one and only %s!", message)
	// }
	//
	// return fmt.Sprintf("Hello %s, my name is chatbot. What was yours again?", message), nil

	return "THE ABSOLUTE POTATO", nil
}

func main() {
	// Uncomment the following lines to customize the chatbot
	// chatbot.WelcomeMessage = "What's your name?"
	//potato
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

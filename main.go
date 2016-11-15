package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	// "github.com/ramin0/chatbot"

	// Autoload environment variables in .env
	_ "github.com/joho/godotenv/autoload"
	"github.com/ramin0/chatbot"
)

// var lastNumOfItems int = 0
// var data []map[string]interface{}

func getDetailsForRecipe(rawRecipe map[string]interface{}) string {
	return fmt.Sprintf("<img src=\"%s\"><br/>This recipe is called %s. <br/> The full list of ingredients is: %s.<br/>"+
		"The full recipe is available <a href=\"%s\" target=\"_blank\">here</a>.<br/>",
		rawRecipe["thumbnail"],
		rawRecipe["title"],
		rawRecipe["ingredients"],
		rawRecipe["href"])
}

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

func getResponse(baseUrl string, params []string, questionString string) *http.Response {
	var searchString string
	searchString = strings.Join(params, ",")
	// for key, value := range params {
	// 	searchString += key + "=" + value + "&"
	// }
	res, err := http.Get(baseUrl + "/" + questionString + "?i=" + searchString)
	if err != nil {
		panic("Couldn't establish a connection!")
	}
	return res
}

func chatbotProcess(session chatbot.Session, message string) (string, error) {

	var returnMsg string

	//checks for invalid characters in the user's message
	if strings.ContainsAny(message,
		`;()!@#$%^*[]{}\\|:?><`) {
		return "", fmt.Errorf("Whoops! Please only enter valid answers to the question! No symbols or numbers!")
	}

	if session["phase"] == nil {
		session["name"] = strings.Split(message, " ")
		session["phase"] = []string{"Querying"}
		session["history"] = []string{}
		session["lastNumOfItems"] = []string{}
		// return "Okay, " + session["name"][0] + ". How many ingredients would you like to specify for your dish?", nil
		// } // else if session["phase"][0] == "Number" {
		// 	if len(message) > 1 || !(strings.ContainsAny(message, "123456789")) {
		// 		return "Please choose a number (digit) between 1 and 9 only.", nil
		// 	}
		session["phase"][0] = "Querying"
		// session["number"] = []string{message}
		// returnMsg = "Okay. What is the first ingredient you want?"
		return "Okay, " + session["name"][0] +
			". Please enter the ingredients you want to specify seperated by commas or spaces", nil
	} else if session["phase"][0] == "Querying" && len(session["lastNumOfItems"]) != len(session["history"]) {
		if strings.EqualFold(message, "Yes") || strings.EqualFold(message, "Y") || strings.EqualFold(message, "Yeah") {
			session["lastNumOfItems"] = make([]string, len(session["history"]))
			return "Okay, " + session["name"][0] + ". What would you like to add?", nil
		} else if strings.EqualFold(message, "No") || strings.EqualFold(message, "N") || strings.EqualFold(message, "Nope") {
			session["lastNumOfItems"] = make([]string, len(session["history"]))
			session["phase"][0] = "APIing"
		} else {
			return "", fmt.Errorf("Don't think I quite got that. Please only enter yes or no.")
		}
	} else if session["phase"][0] == "Querying" {
		var items []string
		// just in case the user entered decided to enter a period somewhere or something
		message = strings.Replace(message, ".", "", -1)
		message = strings.Replace(message, "and", "", -1)
		message = strings.Replace(message, "&", "", -1)
		message = strings.Replace(message, "  ", " ", -1)
		if strings.Contains(message, ",") && strings.Contains(message, " ") {
			// assuming user enters something like "Pasta, onions, caramel"
			message = strings.Replace(message, " ", "", -1)
		}
		if strings.Contains(message, ",") {
			items = strings.Split(message, ",")
		} else {
			// default case, can also handle space separation
			items = strings.Split(message, " ")
		}

		returnMsg = "You want "
		// numItems, _ := strconv.Atoi(session["number"][0])
		// numItems -= 1
		// session["number"][0] = strconv.Itoa(numItems)
		// session["history"] = append(session["history"], strings.Split(message, " ")[0])
		// numItems := len(items) + len(session["history"])
		for _, item := range items {
			session["history"] = append(session["history"], item)
		}
		for index, item := range session["history"] {
			if index != len(session["history"])-1 {
				returnMsg += item
				returnMsg += ", "
			} else {
				if len(session["history"]) != 1 {
					returnMsg += "and "
				}
				returnMsg += item
				returnMsg += ". Do you want to add any more ingredients? (Yes/No)"
			}
		}
	}
	if session["phase"][0] == "APIing" {
		data := getJSONArray(getResponse("http://www.recipepuppy.com/api", session["history"], ""), "results")
		if len(data) == 0 {
			returnMsg += "Whoops! I don't seem to have found any recipe matching your entered items. \n Would you like to start over?"
			session["phase"][0] = "Ending"
			session["phase"] = append(session["phase"], "Failure")
			return returnMsg, nil
		} else {
			session["phase"][0] = "Ending"
			session["phase"] = append(session["phase"], "Success")
			returnMsg += getDetailsForRecipe(data[0])
			session["results"] = []string{strconv.Itoa(len(data)), "1"}
		}

		return returnMsg + "Please say 'next' to get more recipes, 'stop' to terminate this session, or 'restart' to start over. ", nil
		// return strings.Join(session["history"], ","), nil

	}
	if session["phase"][0] == "Ending" {
		if session["phase"][1] == "Failure" || session["phase"][1] == "Complete" {
			if strings.EqualFold(message, "yes") || strings.EqualFold(message, "restart") || strings.EqualFold(message, "retry") {
				session["phase"] = nil
				return "A WHOLE NEW WORLD! Sorry, what was your name again?n\n\n", nil
			} else if strings.EqualFold(message, "no") || strings.EqualFold(message, "bye") || strings.EqualFold(message, "goodbye") {
				session["phase"][0] = "Shutdown"
				return "Understood, thank you for using this service. Have a lovely day!", nil
			} else {
				return "I'm sorry, I didn't quite catch that. Please say yes or restart if you want to restart, or no otherwise", nil
			}
		}
		if strings.EqualFold(message, "more") || strings.EqualFold(message, "next") {
			numResults, _ := strconv.Atoi(session["results"][0])
			currentItem, _ := strconv.Atoi(session["results"][1])
			if numResults == currentItem+1 {
				session["phase"][1] = "Complete"
				return "I've already listed everything. Thanks! Please say yes or restart if you want to restart, or no otherwise", nil
			} else {
				currentItem += 1
				session["results"][1] = strconv.Itoa(currentItem)
				data := getJSONArray(getResponse("http://www.recipepuppy.com/api", session["history"], ""), "results")
				return getDetailsForRecipe(data[currentItem]) + "Please say 'next' to get more recipes, 'stop' to end this interaction, or 'restart' to start over", nil
				// return strings.Join(session["history"], ","), nil
			}
		} else if strings.EqualFold(message, "stop") || strings.EqualFold(message, "bye") {
			session["phase"][0] = "Shutdown"
			return "I take it that will be all. Goodbye!", nil
		} else if strings.EqualFold(message, "restart") || strings.EqualFold(message, "redo") {
			session["phase"] = nil
			return "A WHOLE NEW WORLD! Sorry, what was your name again?\n\n\n", nil
		} else {
			return "", fmt.Errorf("Sorry, I didn't quite catch that. Please say 'next' to get more recipes, 'stop' to end this interaction, or 'restart' to start over")
		}
	}
	if session["phase"][0] == "Shutdown" {
		return "", fmt.Errorf("THE CHAT HAS ALREADY BEEN TERMINATED")
	}
	// if strings.EqualFold(message, "chatbot") {
	// 	return "", fmt.Errorf("This can't be, I'm the one and only %s!", message)
	// }
	// if reflect.DeepEqual(session["toot"], []string{"THE EPIC MAN"}) {
	// 	fmt.Println("EYYY")
	// }
	//
	// _, err := getResponse("www.recipepuppy.com/api", map[string]string{"i": "garlic"}, "")
	// return err.Error(), nil

	// return fmt.Sprintf("Hello %s, my name is chatbot. What was yours again?", message), nil
	// resp, _ := http.Get(nutritionix + "taco?appId=" + appId + "&appKey=" + appKey)
	// res, _ := http.Get("http://www.recipepuppy.com/api/?i=onions,garlic&q=omelet&p=3")
	//	defer res.Body.Close()
	// body, _ := ioutil.ReadAll(resp.Body)
	//	data := getJSONArray(res, "results")
	// var s string
	// for someVar := range body {
	// s += string(someVar)
	// }
	// return s, nil
	// return resp.H, nil

	return returnMsg, nil
}

func main() {
	// Uncomment the following lines to customize the chatbot
	chatbot.WelcomeMessage = `Hello. I am the GUC Sofra AI.
  I will help you in determining the best
  recipes based on the nutritional values you want
  by asking you a series of questions.
	First off, could you tell me what your name is please?`
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

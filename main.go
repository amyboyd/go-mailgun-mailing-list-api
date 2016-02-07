package main

import (
	"fmt"
	mailgun "github.com/mailgun/mailgun-go"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Config struct {
	Domain                    string
	ApiKey                    string
	MailingListAddress        string
	HttpPort                  int
	RedirectUrlAfterSubscribe string
}

func main() {
	if len(os.Args) > 1 {
		command := os.Args[1]
		if command == "--help" {
			PrintHelp()
		} else {
			log.Fatal("unknown argument: " + command)
		}
	} else {
		RunApplication()
	}
}

func RunApplication() {
	config := CreateConfigFromEnv()

	PrintConfig(config)

	mg := mailgun.NewMailgun(config.Domain, config.ApiKey, "")

	StartHttpServer(mg, config)
}

func CreateConfigFromEnv() Config {
	httpPort, err := strconv.Atoi(GetConfigVarFromEnv("SUBSCRIBE_HTTP_PORT"))

	if err != nil {
		log.Fatal("HTTP port must be numeric, given: " + GetConfigVarFromEnv("SUBSCRIBE_HTTP_PORT"))
	}

	return Config{
		Domain:                    GetConfigVarFromEnv("MAILGUN_DOMAIN"),
		ApiKey:                    GetConfigVarFromEnv("MAILGUN_API_KEY"),
		MailingListAddress:        GetConfigVarFromEnv("MAILGUN_MAILING_LIST"),
		HttpPort:                  httpPort,
		RedirectUrlAfterSubscribe: GetConfigVarFromEnv("SUBSCRIBE_REDIRECT_URL"),
	}
}

func GetConfigVarFromEnv(envVar string) string {
	value := os.Getenv(envVar)

	if value == "" {
		log.Fatal("System environment not set: " + envVar)
	}

	return value
}

func PrintConfig(config Config) {
	fmt.Println("Running with configuration:")

	fmt.Println("Mailgun domain (from MAILGUN_DOMAIN):", config.Domain)

	// Redact the API key incase the output is redirected to a log file and, for example, sent
	// to a Kibana log server, where the users of the log server shouldn't know the private key.
	redactedApiKey := config.ApiKey[0:10] + strings.Repeat("*", utf8.RuneCountInString(config.ApiKey)-10)
	fmt.Println("Mailgun API key (from MAILGUN_API_KEY):", redactedApiKey)

	fmt.Println("Mailgun mailing list (from MAILGUN_MAILING_LIST):", config.MailingListAddress)

	fmt.Println("HTTP port (from SUBSCRIBE_HTTP_PORT):", config.HttpPort)

	fmt.Println("Redirect URL after subscribe (from SUBSCRIBE_REDIRECT_URL):", config.RedirectUrlAfterSubscribe)
}

func StartHttpServer(mg mailgun.Mailgun, config Config) {
	http.HandleFunc("/subscribe", func(response http.ResponseWriter, request *http.Request) {
		email := request.FormValue("email")

		newMember := mailgun.Member{
			Address:    email,
			Subscribed: mailgun.Subscribed,
		}

		mg.CreateMember(true, config.MailingListAddress, newMember)

		fmt.Println(email + " has been subscribed to the mailing list")

		http.Redirect(response, request, config.RedirectUrlAfterSubscribe, http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/health-check", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(response, "Running")
	})

	http.ListenAndServe(":"+strconv.Itoa(config.HttpPort), nil)
}

func PrintHelp() {
	fmt.Println("You can report issues at: https://github.com/amyboyd/go-mailgun-mailing-list-api/issues")
	fmt.Println("You can download the latest binary from: https://github.com/amyboyd/go-mailgun-mailing-list-api/releases")
	fmt.Println("The source code is available at: https://github.com/amyboyd/go-mailgun-mailing-list-api")
}

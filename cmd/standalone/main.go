package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/justinas/alice"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/labd/terraform-github-registry/internal"
)

var (
	rootCmd = &cobra.Command{
		Use:   "terraform-registry",
		Short: "Run the Terraform Registry",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			internal.InitLogging()
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(cmd); err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
			os.Exit(0)
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(fmt.Errorf("execute: %v", err))
	}
}

func run(cmd *cobra.Command) error {
	r := internal.NewRouter()
	handler := alice.New(
		internal.PanicHandler(),
		internal.LoggingHandler(log.Logger),
	).Then(r)

	// check if certs/registry.pem exists and if it exists then run the server with TLS
	if _, err := os.Stat("certs/registry.test.pem"); err == nil {
		log.Printf("Running service on https://%s\n", "localhost:3000")
		return http.ListenAndServeTLS("127.0.0.1:3000", "certs/registry.test.pem", "certs/registry.test-key.pem", handler)
	}
	log.Printf("Running service on http://%s\n", "localhost:3000")
	return http.ListenAndServe("127.0.0.1:3000", handler)
}

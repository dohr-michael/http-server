package cmd

import (
	"github.com/spf13/cobra"
	"github.com/dohr-michael/http-server/pkg/start"
)

var addr string
var tlsAddr string
var env []string
var certFile string
var keyFile string

var startCmd = &cobra.Command{
	Use:   "start [options] DIRECTORY_TO_SERVE",
	Short: "Start http server to serve files.",
	Long: `
Start http server to serve files.
Will provide additional routes :
- /config.js

This routes provide all environment variable of the server prefixed by HTTP_SERVER_CONFIG_ as map<string, string>
- NO_CACHE_FILES (array with ',' as separator, always /index.html)

Examples:
	http-server start ./static
	http-server start -a :8080 --tls-addr :443 ./static
	http-server start -a :8080 --env "APP_ID=todolist" ./static
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return start.Run(args[0], addr, tlsAddr, certFile, keyFile, env)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&addr, "addr", "a", ":8080", "address to serve as HTTP")
	startCmd.Flags().StringVarP(&tlsAddr, "tls-addr", "", ":8443", "address to serve as HTTPS")
	startCmd.Flags().StringVarP(&certFile, "cert-file", "", "", "Cert File")
	startCmd.Flags().StringVarP(&keyFile, "key-file", "", "", "Key File")
	startCmd.Flags().StringArrayVarP(&env, "env", "", []string{}, "environment variables")
}

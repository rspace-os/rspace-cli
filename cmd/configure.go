package cmd

import (
	_ "embed"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configureCmd)
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Set the URL and API key interactively",
	Long:  "Configure the application with a URL and API key",
	Run: func(cmd *cobra.Command, args []string) {
		doConfigure()
	},
}

func doConfigure() {
	var apiKey string
	var urlS string
	fmt.Println("Please enter API Key:")
	fmt.Scan(&apiKey)
	fmt.Println("Please enter RSpace URL:")
	fmt.Scan(&urlS)
	messageStdErr(fmt.Sprintf("You entered %s, %s", apiKey, urlS))
	checkArgs(apiKey)
	checkUrl(urlS)
	homeDir, e := os.UserHomeDir()
	if e != nil {
		exitWithStdErrMsg("Couldn't find home directory")
	}
	configFile := path.Join(homeDir, DEFAULT_CONFIG_NAME)
	if _, e := os.Stat(configFile); !os.IsNotExist(e) {
		fmt.Printf("There is an existing config file at %s, do you want to overwrite (y/n)?\n ", configFile)
		overwrite := "n"
		fmt.Scanln(&overwrite)
		if overwrite != "y" {
			exitWithStdErrMsg("Not overwritten")
		} else {
			writeConfigFile(urlS, apiKey, configFile, false)
		}
	} else {
		writeConfigFile(urlS, apiKey, configFile, true)

	}

}

func writeConfigFile(urlS, apiKey, configFile string, createFile bool) {
	toWrite := fmt.Sprintf("%s=%s\n%s=%s", APIKEY_ENV_NAME, apiKey, BASE_URL_ENV_NAME, urlS)
	os.WriteFile(configFile, []byte(toWrite), 0644)
}

func checkArgs(apiKey string) {
	if apiKey == "" || len(apiKey) < 8 {
		exitWithStdErrMsg("Invalid API key entered.")
	}
}

func checkUrl(urlS string) {
	_, err := url.Parse(urlS)
	if err != nil {
		exitWithStdErrMsg(fmt.Sprintf(" '%s' is invalid URL syntax", urlS))
	}
}

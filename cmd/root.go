package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// if debug logging is on or off
var debug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "s3-cisbench",
	Short: "s3-csibench is a tool that analyses S3 bucket against CIS benchmark rules",
	Long:  `s3-csibench is a tool that analyses S3 bucket against CIS benchmark rules. Full details can be found at https://github.com/rollwagen/s3-cisbench`,

	// Uncomment the following line if your bare application has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use --help for more information on how to use s3-cisbench.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		setUpLogging(debug)
		return nil
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here, will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.s3-cisbench.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable verbose logging")

	// Cobra also supports local flags, which will only run when this action is called directly.
	// Example rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func setUpLogging(debug bool) {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors:          false,
			FullTimestamp:          true,
			TimestampFormat:        "15:04:05",
			DisableLevelTruncation: true,
			PadLevelText:           true,
		})
	} else {
		logrus.SetLevel(logrus.WarnLevel)
		plainFormatter := new(logrusPlainFormatter)
		logrus.SetFormatter(plainFormatter)
	}
}

type logrusPlainFormatter struct{}

func (f *logrusPlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s\n", entry.Message)), nil
}

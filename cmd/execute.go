package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bpineau/kube-deployments-notifier/config"
	klog "github.com/bpineau/kube-deployments-notifier/pkg/log"
	"github.com/bpineau/kube-deployments-notifier/pkg/run"
)

const appName = "kube-deployments-notifier"

var (
	cfgFile   string
	apiServer string
	kubeConf  string
	dryRun    bool
	logLevel  string
	logOutput string
	logServer string
	endpoint  string
	tokenHdr  string
	tokenVal  string
	filter    string
	healthP   int

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   appName,
		Short: "Notify deployments",
		Long:  "Notify deployments",

		Run: func(cmd *cobra.Command, args []string) {
			config := &config.KdnConfig{
				DryRun:     viper.GetBool("dry-run"),
				Logger:     klog.New(viper.GetString("log.level"), viper.GetString("log.server"), viper.GetString("log.output")),
				Endpoint:   viper.GetString("endpoint"),
				TokenHdr:   viper.GetString("token-header"),
				TokenVal:   viper.GetString("token-value"),
				Filter:     viper.GetString("filter"),
				HealthPort: viper.GetInt("healthcheck-port"),
			}
			config.Init(viper.GetString("api-server"), viper.GetString("kube-config"))
			run.Run(config)
		},
	}
)

// Execute adds all child commands to the root command and sets their flags.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	defaultCfg := "/etc/kdn/" + appName + ".yaml"
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", defaultCfg, "configuration file")

	rootCmd.PersistentFlags().StringVarP(&apiServer, "api-server", "s", "", "kube api server url")
	if err := viper.BindPFlag("api-server", rootCmd.PersistentFlags().Lookup("api-server")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().StringVarP(&kubeConf, "kube-config", "k", "", "kube config path")
	if err := viper.BindPFlag("kube-config", rootCmd.PersistentFlags().Lookup("kube-config")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}
	if err := viper.BindEnv("kube-config", "KUBECONFIG"); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "dry-run mode")
	if err := viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "v", "debug", "log level")
	if err := viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().StringVarP(&logOutput, "log-output", "o", "stderr", "log output")
	if err := viper.BindPFlag("log.output", rootCmd.PersistentFlags().Lookup("log-output")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().StringVarP(&logServer, "log-server", "r", "", "log server (if using syslog)")
	if err := viper.BindPFlag("log.server", rootCmd.PersistentFlags().Lookup("log-server")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "", "API endpoint")
	if err := viper.BindPFlag("endpoint", rootCmd.PersistentFlags().Lookup("endpoint")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().StringVarP(&tokenHdr, "token-header", "t", "", "token header name")
	if err := viper.BindPFlag("token-header", rootCmd.PersistentFlags().Lookup("token-header")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().StringVarP(&tokenVal, "token-value", "a", "", "token header value")
	if err := viper.BindPFlag("token-value", rootCmd.PersistentFlags().Lookup("token-value")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().StringVarP(&filter, "filter", "l", "", "Label filter")
	if err := viper.BindPFlag("filter", rootCmd.PersistentFlags().Lookup("filter")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().IntVarP(&healthP, "healthcheck-port", "p", 0, "port for answering healthchecks")
	if err := viper.BindPFlag("healthcheck-port", rootCmd.PersistentFlags().Lookup("healthcheck-port")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName(appName)

	// all possible config file paths, by priority
	viper.AddConfigPath("/etc/kdn/")
	if home, err := homedir.Dir(); err == nil {
		viper.AddConfigPath(home)
	}
	viper.AddConfigPath(".")

	// prefer the config file path provided by cli flag, if any
	if _, err := os.Stat(cfgFile); !os.IsNotExist(err) {
		viper.SetConfigFile(cfgFile)
	}

	// allow config params through prefixed env variables
	viper.SetEnvPrefix("KDN")
	replacer := strings.NewReplacer("-", "_", ".", "_DOT_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logrus.Info("Using config file: ", viper.ConfigFileUsed())
	}
}

package cmd

import (
	"log"
	"os"
	"strings"
	"time"

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
	resync    int

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
				ResyncIntv: time.Duration(viper.GetInt("resync-interval")) * time.Second,
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

func bindPFlag(key string, cmd string) {
	if err := viper.BindPFlag(key, rootCmd.PersistentFlags().Lookup(cmd)); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	defaultCfg := "/etc/kdn/" + appName + ".yaml"
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", defaultCfg, "configuration file")

	rootCmd.PersistentFlags().StringVarP(&apiServer, "api-server", "s", "", "kube api server url")
	bindPFlag("api-server", "api-server")

	rootCmd.PersistentFlags().StringVarP(&kubeConf, "kube-config", "k", "", "kube config path")
	bindPFlag("kube-config", "kube-config")
	if err := viper.BindEnv("kube-config", "KUBECONFIG"); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "dry-run mode")
	bindPFlag("dry-run", "dry-run")

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "v", "debug", "log level")
	bindPFlag("log.level", "log-level")

	rootCmd.PersistentFlags().StringVarP(&logOutput, "log-output", "o", "stderr", "log output")
	bindPFlag("log.output", "log-output")

	rootCmd.PersistentFlags().StringVarP(&logServer, "log-server", "r", "", "log server (if using syslog)")
	bindPFlag("log.server", "log-server")

	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "", "API endpoint")
	bindPFlag("endpoint", "endpoint")

	rootCmd.PersistentFlags().StringVarP(&tokenHdr, "token-header", "t", "", "token header name")
	bindPFlag("token-header", "token-header")

	rootCmd.PersistentFlags().StringVarP(&tokenVal, "token-value", "a", "", "token header value")
	bindPFlag("token-value", "token-value")

	rootCmd.PersistentFlags().StringVarP(&filter, "filter", "l", "", "Label filter")
	bindPFlag("filter", "filter")

	rootCmd.PersistentFlags().IntVarP(&healthP, "healthcheck-port", "p", 0, "port for answering healthchecks")
	bindPFlag("healthcheck-port", "healthcheck-port")

	rootCmd.PersistentFlags().IntVarP(&resync, "resync-interval", "i", 900, "resync interval in seconds (0 to disable)")
	bindPFlag("resync-interval", "resync-interval")
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

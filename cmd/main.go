package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/djedjethai/generation/internal/agent"
	"github.com/djedjethai/generation/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var environment string
var jaegerEndpoint string
var shards int
var itemsPerShard int
var fileLoggerActive bool
var dbLoggerActive bool
var isTracing bool
var isMetrics bool
var logMode string

type cli struct {
	cfg cfg
}

type cfg struct {
	agent.Config
	ServerTLSConfig config.TLSConfig
	PeerTLSConfig   config.TLSConfig
}

func main() {

	// added
	cli := &cli{}
	cmd := &cobra.Command{
		Use:     "generation",
		PreRunE: cli.setupConfig,
		RunE:    cli.run,
	}
	if err := setupFlags(cmd); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func setupFlags(cmd *cobra.Command) error {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	// service configurations
	// TODO set the "config-file" to be not static
	cmd.Flags().String("config-file", "/home/jerome/Documents/projects/generationProject/generation/devConfig/config.yaml", "Path to config file.")
	dataDir := path.Join(os.TempDir(), "generation")
	cmd.Flags().String("data-dir", dataDir, "Directory to store log and Raft data.")
	cmd.Flags().String("node-name", hostname, "Unique server ID.")
	cmd.Flags().String("bind-addr", "127.0.0.1:8500", "Address to bind Serf on.")
	cmd.Flags().Int("rpc-port", 8400, "Port for RPC clients (and Raft) connections.")
	cmd.Flags().StringSlice("start-join-addrs", nil, "Serf addresses to join.")
	cmd.Flags().Bool("bootstrap", false, "Bootstrap the cluster.")
	// cmd.Flags().String("acl-model-file", "", "Path to ACL model.")
	// cmd.Flags().String("acl-policy-file", "", "Path to ACL policy.")
	cmd.Flags().String("server-tls-cert-file", "/.generation/server.pem", "Path to server tls cert.")
	cmd.Flags().String("server-tls-key-file", "/.generation/server-key.pem", "Path to server tls key.")
	cmd.Flags().String("server-tls-ca-file", "/.generation/ca.pem", "Path to server certificate authority.")
	cmd.Flags().String("peer-tls-cert-file", "/.generation/client.pem", "Path to peer tls cert.")
	cmd.Flags().String("peer-tls-key-file", "/.generation/client-key.pem", "Path to peer tls key.")
	cmd.Flags().String("peer-tls-ca-file", "/.generation/ca.pem", "Path to peer certificate authority.")

	// service options
	// TODO set the environment flag
	cmd.Flags().StringVarP(&environment, "environment", "e", "dev", "set the environment dev or prod")
	cmd.Flags().StringVarP(&jaegerEndpoint, "jaeger", "j", "http://jaeger:14268/api/traces", "the Jaeger end point to connect")
	cmd.Flags().IntVarP(&shards, "shards", "s", 2, "number of shards")
	cmd.Flags().IntVarP(&itemsPerShard, "itemPerShard", "i", 10, "number of shards")
	cmd.Flags().BoolVarP(&dbLoggerActive, "dbLogger", "d", false, "enable the database logging")
	cmd.Flags().BoolVarP(&isTracing, "isTracing", "t", false, "enable Jaeger tracing")
	cmd.Flags().BoolVarP(&isMetrics, "isMetrics", "m", false, "enable Prometheus metrics")
	cmd.Flags().StringVarP(&logMode, "loggerMode", "l", "prod", "logger mode can be prod, development, debug")

	return viper.BindPFlags(cmd.Flags())
}

func (c *cli) setupConfig(cmd *cobra.Command, args []string) error {
	var err error
	configFile, err := cmd.Flags().GetString("config-file")
	if err != nil {
		return err
	}
	log.Println("see the config file path: ", configFile)
	viper.SetConfigFile(configFile)
	if err = viper.ReadInConfig(); err != nil {
		// it's ok if config file doesn't exist
		log.Println("config file has not been found")
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	c.cfg.DataDir = viper.GetString("data-dir")
	log.Println("config file see DataDir: ", c.cfg.DataDir)

	c.cfg.NodeName = viper.GetString("node-name")
	log.Println("config file see NodeName: ", c.cfg.NodeName)

	c.cfg.BindAddr = viper.GetString("bind-addr")
	log.Println("config file see BindAddr: ", c.cfg.BindAddr)

	c.cfg.PortGRPC = viper.GetInt("rpc-port")
	log.Println("config file see RpcPort: ", c.cfg.PortGRPC)

	c.cfg.StartJoinAddrs = viper.GetStringSlice("start-join-addrs")
	log.Println("config file see JoinAddr: ", c.cfg.StartJoinAddrs)

	c.cfg.Bootstrap = viper.GetBool("bootstrap")
	log.Println("config file see Bootstrap: ", c.cfg.Bootstrap)
	// c.cfg.ACLModelFile = viper.GetString("acl-mode-file")
	// c.cfg.ACLPolicyFile = viper.GetString("acl-policy-file")
	c.cfg.ServerTLSConfig.CertFile = viper.GetString("server-tls-cert-file")
	log.Println("config file see ServerTLSConfig CertFile: ", c.cfg.ServerTLSConfig.CertFile)
	c.cfg.ServerTLSConfig.KeyFile = viper.GetString("server-tls-key-file")
	log.Println("config file see ServerTLSConfig KeyFile: ", c.cfg.ServerTLSConfig.KeyFile)
	c.cfg.ServerTLSConfig.CAFile = viper.GetString("server-tls-ca-file")
	log.Println("config file see ServerTLSConfig CaFile: ", c.cfg.ServerTLSConfig.CAFile)
	c.cfg.PeerTLSConfig.CertFile = viper.GetString("peer-tls-cert-file")
	log.Println("config file see PerrTLSConfig CertFile: ", c.cfg.PeerTLSConfig.CertFile)
	c.cfg.PeerTLSConfig.KeyFile = viper.GetString("peer-tls-key-file")
	log.Println("config file see PerrTLSConfig KeyFile: ", c.cfg.PeerTLSConfig.KeyFile)
	c.cfg.PeerTLSConfig.CAFile = viper.GetString("peer-tls-ca-file")
	log.Println("config file see PerrTLSConfig CaFile: ", c.cfg.PeerTLSConfig.CAFile)
	if c.cfg.ServerTLSConfig.CertFile != "" &&
		c.cfg.ServerTLSConfig.KeyFile != "" {
		c.cfg.ServerTLSConfig.Server = true
		c.cfg.Config.ServerTLSConfig, err = config.SetupTLSConfig(
			c.cfg.ServerTLSConfig,
		)
		if err != nil {
			return err
		}
	}
	if c.cfg.PeerTLSConfig.CertFile != "" &&
		c.cfg.PeerTLSConfig.KeyFile != "" {
		c.cfg.PeerTLSConfig.Server = false
		c.cfg.Config.PeerTLSConfig, err = config.SetupTLSConfig(
			c.cfg.PeerTLSConfig,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *cli) run(cmd *cobra.Command, args []string) error {

	// add all var(from flags) to the config first as some are needed after
	c.cfg.FileLoggerActive = fileLoggerActive
	c.cfg.DBLoggerActive = dbLoggerActive
	c.cfg.Shards = shards
	c.cfg.ItemsPerShard = itemsPerShard
	c.cfg.IsTracing = isTracing
	c.cfg.IsMetrics = isMetrics
	c.cfg.JaegerEndpoint = jaegerEndpoint

	var err error
	err = setupSrv(&c.cfg.Config)
	if err != nil {
		os.Exit(1)
	}

	fmt.Println("seee al configsssss: ", c.cfg)
	fmt.Println("seee al configsssss: ", c.cfg.Bootstrap)

	// TODO uncomment here to run the service
	agent, err := agent.New(c.cfg.Config)
	if err != nil {
		return err
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
	return agent.Shutdown()
	// return nil
}

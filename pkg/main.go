package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	pkg "git.solsynth.dev/hypernet/insight/pkg/internal"
	"git.solsynth.dev/hypernet/insight/pkg/internal/gap"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"github.com/fatih/color"

	"git.solsynth.dev/hypernet/insight/pkg/internal/cache"
	"git.solsynth.dev/hypernet/insight/pkg/internal/database"
	"git.solsynth.dev/hypernet/insight/pkg/internal/grpc"

	"git.solsynth.dev/hypernet/insight/pkg/internal/server"
	"git.solsynth.dev/hypernet/insight/pkg/internal/services"
	"github.com/robfig/cron/v3"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func main() {
	// Booting screen
	fmt.Println(color.YellowString(" ___           _       _     _\n|_ _|_ __  ___(_) __ _| |__ | |_\n | || '_ \\/ __| |/ _` | '_ \\| __|\n | || | | \\__ \\ | (_| | | | | |_\n|___|_| |_|___/_|\\__, |_| |_|\\__|\n                 |___/"))
	fmt.Printf("%s v%s\n", color.New(color.FgHiYellow).Add(color.Bold).Sprintf("Hypernet.Insight"), pkg.AppVersion)
	fmt.Printf("The A-word in Hypernet.\n")
	color.HiBlack("=====================================================\n")

	// Configure settings
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.SetConfigName("settings")
	viper.SetConfigType("toml")

	// Load settings
	if err := viper.ReadInConfig(); err != nil {
		log.Panic().Err(err).Msg("An error occurred when loading settings.")
	}

	// Connect to nexus
	if err := gap.InitializeToNexus(); err != nil {
		log.Error().Err(err).Msg("An error occurred when registering service to nexus...")
	}

	// Load keypair
	if reader, err := sec.NewInternalTokenReader(viper.GetString("security.internal_public_key")); err != nil {
		log.Error().Err(err).Msg("An error occurred when reading internal public key for jwt. Authentication related features will be disabled.")
	} else {
		server.IReader = reader
		log.Info().Msg("Internal jwt public key loaded.")
	}

	// Connect to database
	if err := database.NewGorm(); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when connect to database.")
	} else if err := database.RunMigration(database.C); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when running database auto migration.")
	}

	// Connect to ollama
	if err := services.ConnectOllama(); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when connecting to ollama.")
	}

	// Initialize cache
	if err := cache.NewStore(); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when initializing cache.")
	}

	// Configure timed tasks
	quartz := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(&log.Logger)))
	quartz.AddFunc("@every 60m", services.DoAutoDatabaseCleanup)
	quartz.Start()

	// Server
	go server.NewServer().Listen()

	// Grpc Server
	go grpc.NewGrpc().Listen()

	// Messages
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	quartz.Stop()
}

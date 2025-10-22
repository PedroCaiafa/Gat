/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"gat/internal/config"

	"github.com/spf13/cobra"
)

func printWelcome() {
	fmt.Println("╔═══════════════════╗")
	fmt.Println("║      G A T        ║")
	fmt.Println("╚═══════════════════╝")

	//fmt.Printf("Chat Provider: %s\nChat Default Model: %s\n\n", cfgFile.ChatProvider.Provider, cfgFile.ChatProvider.DefaultModel)
	//fmt.Printf("Embedding Provider: %s\nEmbedding Default Model: %s\n\n", cfgFile.EmbeddingProvider.Provider, cfgFile.EmbeddingProvider.DefaultModel)
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Gat",
	Long:  `Creates/checks config files and initializes Qdrant conn and app`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := config.EnsureConfigFile()
		if err != nil {
			return err
		}
		cfgFile, err := config.LoadConfig()
		if err != nil {
			return err
		}

		/*
			-------------------------------	QDRANT SETUP ---------------------------------
		*/
		switch cfgFile.QdrantMode {
		case "disabled":
			fmt.Println("Running without Qdrant - memory features are OFF")
			return nil

		case "external":

		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

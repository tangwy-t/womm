package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/womm/womm/internal/app"
	"github.com/womm/womm/internal/config"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "womm",
	Short: "WOMM - Works On My Machine badge generator",
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP badge server",
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := loadApp()
		if err != nil {
			return err
		}
		defer a.Store.Close()
		fmt.Printf("WOMM server starting on %s:%d\n", a.Config.Server.Host, a.Config.Server.Port)
		return a.Server.ListenAndServe(a.Config)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available badges",
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := loadApp()
		if err != nil {
			return err
		}
		defer a.Store.Close()
		badges := a.Registry.ListAll()
		for _, b := range badges {
			fmt.Printf("[%-8s] %-30s  %s\n", b.Rarity, b.LocalizedName("zh"), b.ID)
		}
		return nil
	},
}

var claimCmd = &cobra.Command{
	Use:   "claim <badge-id>",
	Short: "Claim a declarative badge",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		user, _ := cmd.Flags().GetString("user")
		if user == "" {
			user = "default"
		}
		a, err := loadApp()
		if err != nil {
			return err
		}
		defer a.Store.Close()

		if _, ok := a.Registry.Lookup(args[0]); !ok {
			return fmt.Errorf("badge not found: %s", args[0])
		}
		return a.Store.ClaimBadge(user, args[0])
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show your unlocked badges",
	RunE: func(cmd *cobra.Command, args []string) error {
		user, _ := cmd.Flags().GetString("user")
		if user == "" {
			user = "default"
		}
		a, err := loadApp()
		if err != nil {
			return err
		}
		defer a.Store.Close()
		states, err := a.Store.GetUserBadges(user)
		if err != nil {
			return err
		}
		if len(states) == 0 {
			fmt.Println("No badges unlocked yet.")
			return nil
		}
		for _, s := range states {
			fmt.Printf("  [%.8s] %s\n", s.Source, s.BadgeID)
		}
		return nil
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate <badge-id>",
	Short: "Generate SVG badge file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		theme, _ := cmd.Flags().GetString("theme")
		output, _ := cmd.Flags().GetString("output")
		lang, _ := cmd.Flags().GetString("lang")
		style, _ := cmd.Flags().GetString("style")

		a, err := loadApp()
		if err != nil {
			return err
		}
		defer a.Store.Close()

		b, ok := a.Registry.Lookup(args[0])
		if !ok {
			return fmt.Errorf("badge not found: %s", args[0])
		}
		if theme == "" {
			theme = "pixel"
		}
		if output == "" {
			output = args[0] + ".svg"
		}
		if style == "" {
			style = "github"
		}
		if lang == "" {
			lang = "zh"
		}

		svg, err := a.Renderer.Render(b, theme, style, lang)
		if err != nil {
			return err
		}
		return os.WriteFile(output, []byte(svg), 0644)
	},
}

var tokenCmd = &cobra.Command{
	Use:   "github-token [token]",
	Short: "Set or show GitHub Personal Access Token",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := loadApp()
		if err != nil {
			return err
		}
		defer a.Store.Close()
		if len(args) == 0 {
			if a.Config.GitHub.DefaultToken != "" {
				fmt.Println("Token is configured.")
			} else {
				fmt.Println("No token configured. Set it via womm.toml or: womm github-token <TOKEN>")
			}
			return nil
		}
		a.Config.GitHub.DefaultToken = args[0]
		fmt.Println("Token updated in memory. Persist it in womm.toml for permanent use.")
		return nil
	},
}

var certifyCmd = &cobra.Command{
	Use:   "certify <badge-id>",
	Short: "Attempt to certify a badge via GitHub API",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		user, _ := cmd.Flags().GetString("user")
		if user == "" {
			return fmt.Errorf("--user flag required")
		}
		a, err := loadApp()
		if err != nil {
			return err
		}
		defer a.Store.Close()

		if a.CertEng == nil {
			return fmt.Errorf("GitHub token not configured")
		}
		result, err := a.CertEng.TryCertify(cmd.Context(), user, args[0])
		if err != nil {
			return err
		}
		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}

func loadApp() (*app.App, error) {
	cfg, err := config.Load(configFile)
	if err != nil {
		return nil, err
	}
	return app.New(cfg)
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "womm.toml", "config file path")

	claimCmd.Flags().String("user", "", "GitHub username")
	statusCmd.Flags().String("user", "", "GitHub username")
	certifyCmd.Flags().String("user", "", "GitHub username to certify")

	generateCmd.Flags().String("theme", "", "visual theme (pixel/cyberpunk/glitch/clean)")
	generateCmd.Flags().String("output", "", "output file path")
	generateCmd.Flags().String("lang", "", "language (zh/en)")
	generateCmd.Flags().String("style", "", "template style (github/badge/wide/terminal/stamp)")

	rootCmd.AddCommand(serveCmd, listCmd, claimCmd, statusCmd, generateCmd, tokenCmd, certifyCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

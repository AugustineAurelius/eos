package cmd

import (
	"github.com/AugustineAurelius/fuufu/config"
	"github.com/AugustineAurelius/fuufu/pkg/logger"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	_ "github.com/AugustineAurelius/fuufu/db/postgres/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func createMigrateCMD(manager *config.Manager) *cobra.Command {
	var dir string
	var gooseCMD string

	migrateCMD := &cobra.Command{
		Use: "migrate",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.NewWithManager(manager)

			postgresMasterConfig := manager.LoadPostgres()
			log.Sugar().Infof("get postgres connection url: %s", postgresMasterConfig.DNS())

			db, err := goose.OpenDBWithDriver("postgres", postgresMasterConfig.DNS())
			if err != nil {
				return err
			}
			defer db.Close()

			if err = db.PingContext(cmd.Context()); err != nil {
				log.Panic(err.Error())
			}

			goose.SetDialect("postgres")
			switch gooseCMD {
			case "up":
				if err = goose.UpContext(cmd.Context(), db, "."); err != nil {
					log.Panic(err.Error())
				}
			case "down":
				if err = goose.DownContext(cmd.Context(), db, "."); err != nil {
					log.Panic(err.Error())
				}
			}

			return nil
		},
	}

	migrateCMD.PersistentFlags().StringVar(&gooseCMD, "goose_command", "up", "command for goose migrations")
	migrateCMD.PersistentFlags().StringVar(&dir, "dir", "db/postgres/migration", "path to directory with migrations")

	return migrateCMD
}
package compose

import (
	"fmt"
	"os"

	"github.com/AugustineAurelius/eos/pkg/helpers"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	Postgres, App, Kafka bool
	GoVersion            float64
)

var Cmd = &cobra.Command{
	Use:   "compose [flags]",
	Short: "generate postgres Dockerfile",
	Long:  `generate postgres Dockerfile`,
	Run: func(cmd *cobra.Command, args []string) {

		f, err := os.Create("docker-compose.yaml")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()

		c := compose{
			Version:  "3.9",
			Services: make(map[string]any),
		}

		if Postgres {
			addPostgres(c.Services)
		}
		if Kafka {
			addKafka(c.Services)
		}
		if App {
			addApplication(c.Services)
			addAppDockerfile()
		}

		f.Write(helpers.Must(yaml.Marshal(c)))

	},
}

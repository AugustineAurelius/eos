package compose

import "gopkg.in/yaml.v3"

func addPostgres(m map[string]any) {
	var p map[string]any
	yaml.Unmarshal([]byte(`
    image: postgres:latest   
    command: >
      postgres -c max_connections=1000
               -c shared_buffers=256MB
               -c effective_cache_size=768MB
               -c maintenance_work_mem=64MB
               -c checkpoint_completion_target=0.7
               -c wal_buffers=16MB
               -c default_statistics_target=100
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DATABASE}
      PGDATA: "/var/lib/postgresql/data/pgdata"
    ports:
      - "${POSTGRES_PORT}:5432"`), &p)

	m["postgres"] = p
}

package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateMany2many, downCreateMany2many)
}

func upCreateMany2many(ctx context.Context, tx *sql.Tx) error {
	// Create only the many-to-many junction tables
	// Assumes deployments, domains, containers, scripts, secrets tables already exist

	queries := []string{
		// Deployment-Domains junction table
		`CREATE TABLE IF NOT EXISTS deployment_domains (
			deployment_id INTEGER NOT NULL,
			domain_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (deployment_id, domain_id),
			CONSTRAINT fk_deployment_domains_deployment 
				FOREIGN KEY (deployment_id) REFERENCES deployments(id) ON DELETE CASCADE,
			CONSTRAINT fk_deployment_domains_domain 
				FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
		)`,

		// Deployment-SubDomains junction table
		`CREATE TABLE IF NOT EXISTS deployment_sub_domains (
			deployment_id INTEGER NOT NULL,
			sub_domain_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (deployment_id, sub_domain_id),
			CONSTRAINT fk_deployment_sub_domains_deployment 
				FOREIGN KEY (deployment_id) REFERENCES deployments(id) ON DELETE CASCADE,
			CONSTRAINT fk_deployment_sub_domains_subdomain 
				FOREIGN KEY (sub_domain_id) REFERENCES sub_domains(id) ON DELETE CASCADE
		)`,

		// Deployment-Containers junction table
		`CREATE TABLE IF NOT EXISTS deployment_containers (
			deployment_id INTEGER NOT NULL,
			container_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (deployment_id, container_id),
			CONSTRAINT fk_deployment_containers_deployment 
				FOREIGN KEY (deployment_id) REFERENCES deployments(id) ON DELETE CASCADE,
			CONSTRAINT fk_deployment_containers_container 
				FOREIGN KEY (container_id) REFERENCES containers(id) ON DELETE CASCADE
		)`,

		// Deployment-Scripts junction table
		`CREATE TABLE IF NOT EXISTS deployment_scripts (
			deployment_id INTEGER NOT NULL,
			script_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (deployment_id, script_id),
			CONSTRAINT fk_deployment_scripts_deployment 
				FOREIGN KEY (deployment_id) REFERENCES deployments(id) ON DELETE CASCADE,
			CONSTRAINT fk_deployment_scripts_script 
				FOREIGN KEY (script_id) REFERENCES scripts(id) ON DELETE CASCADE
		)`,

		// Deployment-Secrets junction table
		`CREATE TABLE IF NOT EXISTS deployment_secrets (
			deployment_id INTEGER NOT NULL,
			secret_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (deployment_id, secret_id),
			CONSTRAINT fk_deployment_secrets_deployment 
				FOREIGN KEY (deployment_id) REFERENCES deployments(id) ON DELETE CASCADE,
			CONSTRAINT fk_deployment_secrets_secret 
				FOREIGN KEY (secret_id) REFERENCES secrets(id) ON DELETE CASCADE
		)`,

		// Deployment-Servers junction table (uncomment when servers module is fixed)
		`CREATE TABLE IF NOT EXISTS deployment_servers (
			deployment_id INTEGER NOT NULL,
			server_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (deployment_id, server_id),
			CONSTRAINT fk_deployment_servers_deployment 
				FOREIGN KEY (deployment_id) REFERENCES deployments(id) ON DELETE CASCADE,
			CONSTRAINT fk_deployment_servers_server 
				FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE
		)`,
	}

	// Execute table creation queries
	for _, query := range queries {
		if _, err := tx.ExecContext(ctx, query); err != nil {
			return err
		}
	}

	// Add performance indexes
	indexQueries := []string{
		// Indexes for deployment_domains
		`CREATE INDEX IF NOT EXISTS idx_deployment_domains_deployment_id ON deployment_domains(deployment_id)`,
		`CREATE INDEX IF NOT EXISTS idx_deployment_domains_domain_id ON deployment_domains(domain_id)`,

		// Indexes for deployment_sub_domains
		`CREATE INDEX IF NOT EXISTS idx_deployment_sub_domains_deployment_id ON deployment_sub_domains(deployment_id)`,
		`CREATE INDEX IF NOT EXISTS idx_deployment_sub_domains_subdomain_id ON deployment_sub_domains(sub_domain_id)`,

		// Indexes for deployment_containers
		`CREATE INDEX IF NOT EXISTS idx_deployment_containers_deployment_id ON deployment_containers(deployment_id)`,
		`CREATE INDEX IF NOT EXISTS idx_deployment_containers_container_id ON deployment_containers(container_id)`,

		// Indexes for deployment_scripts
		`CREATE INDEX IF NOT EXISTS idx_deployment_scripts_deployment_id ON deployment_scripts(deployment_id)`,
		`CREATE INDEX IF NOT EXISTS idx_deployment_scripts_script_id ON deployment_scripts(script_id)`,

		// Indexes for deployment_secrets
		`CREATE INDEX IF NOT EXISTS idx_deployment_secrets_deployment_id ON deployment_secrets(deployment_id)`,
		`CREATE INDEX IF NOT EXISTS idx_deployment_secrets_secret_id ON deployment_secrets(secret_id)`,

		// Uncomment when servers module is fixed
		// `CREATE INDEX IF NOT EXISTS idx_deployment_servers_deployment_id ON deployment_servers(deployment_id)`,
		// `CREATE INDEX IF NOT EXISTS idx_deployment_servers_server_id ON deployment_servers(server_id)`,
	}

	// Execute index creation queries
	for _, query := range indexQueries {
		if _, err := tx.ExecContext(ctx, query); err != nil {
			// Continue if index already exists or fails - not critical
			continue
		}
	}

	return nil
}

func downCreateMany2many(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	queries := []string{
		`DROP TABLE IF EXISTS deployment_secrets`,
		`DROP TABLE IF EXISTS deployment_scripts`,
		// `DROP TABLE IF EXISTS deployment_servers`, // Uncomment when servers module is fixed
		`DROP TABLE IF EXISTS deployment_containers`,
		`DROP TABLE IF EXISTS deployment_sub_domains`,
		`DROP TABLE IF EXISTS deployment_domains`,
	}

	// Execute all queries
	for _, query := range queries {
		if _, err := tx.ExecContext(ctx, query); err != nil {
			return err
		}
	}

	return nil
}

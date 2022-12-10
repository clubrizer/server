#!make

define load_env
	set -o allexport && source services/$(1)/cfg/.env && set +o allexport
endef

from-env:
	$(call load_env,$(service)); $(MAKE) $(recipe)
db-migration:
	migrate create -ext sql -dir services/$(service)/internal/db/migrations -seq $(name)
migrate-db:
	echo "Migrating $(service) service DB"
	echo "cockroachdb://$${DB_USER}:$${DB_PASSWORD}@$${DB_URL}"
	migrate \
	-database "cockroachdb://$${DB_USER}:$${DB_PASSWORD}@$${DB_URL}" \
	-path services/$(service)/db/migrations \
	$(to)
mocks:
	cd services/$(service) && rm -rf ./$(dir)/mocks && mockery --dir ./$(dir) --all --with-expecter --exported --keeptree --output ./$(dir)/mocks

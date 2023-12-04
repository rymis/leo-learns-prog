all: backend frontend

test: test-backend test-frontend

backend:
	(cd learnsrv; go build)

test-backend:
	(cd learnsrv/rcs; go test)

frontend:
	@echo "TODO"
	# (cd programming-tasks; npm run build)

test-frontend:
	@echo "TODO"
	# (cd programming-tasks; npm run test)

.PHONY: test backend frontend test-backend test-frontend

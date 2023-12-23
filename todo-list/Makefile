MOCKS_DESTINATION=mocks
.PHONY: mocks
mocks: pkg/data/task.go pkg/handlers/task.go
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; done
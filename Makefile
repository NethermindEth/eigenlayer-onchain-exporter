DEFAULT_GOAL := help 
.PHONY: env

pyenv: ## Activate pyenv
	@source venv/bin/activate

upgrade-pip: ## Upgrade pip version
	@pip3 install --upgrade pip

install-deps-locally: upgrade-pip ## Install dependencies in project folder
	@pip3 install -r requirements.txt -t .

install-deps: upgrade-pip ## Install dependencies
	@pip3 install -r requirements.txt

freeze: ## Freeze dependencies
	@rm -f requirements.txt 
	@pip3 freeze > requirements.txt

test: ## Run unit tests
	@python -m unittest discover

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
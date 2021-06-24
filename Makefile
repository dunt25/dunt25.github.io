GOCMD=go
GOBUILD=$(GOCMD) build

image: build run-image
page: build run-page
link: build run-link
clean: build run-clean

build:
	@echo "Building..."
	@$(GOBUILD)

run-image:
	@echo "Downloading images..."
	@./kuro -image

run-page:
	@echo "Downloading pages..."
	@./kuro -page

run-link:
	@echo "Fetching links..."
	@./kuro -link

run-clean:
	@echo "Cleaning..."
	@./kuro -clean
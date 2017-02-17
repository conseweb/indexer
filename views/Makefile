PWD := $(shell pwd)
PKG := github.com/conseweb/indexer/views
APP := ui

DEV_UI_IMAGE := ckeyer/dev:node

dev-ui:
	docker run --rm -it \
	 --name $(APP)-dev \
	 -p 8080:8080 \
	 -v $(PWD):/opt/$(APP) \
	 -w /opt/$(APP) \
	 $(DEV_UI_IMAGE) bash

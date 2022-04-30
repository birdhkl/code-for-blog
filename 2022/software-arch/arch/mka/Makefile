.PHONY: build-bundle

build-bundle: plugins/so plugins/src
	@echo "# patterns/mka(build-bundle)"
	@$(foreach var,$(shell ls plugins/src),\
		go build -buildmode=plugin -o plugins/so/$(subst .go,.so,$(var)) plugins/src/$(var);\
    )

plugins/so:
	@echo "create folder so"
	@mkdir plugins/so

plugins/src:
	@echo "create folder src"
	@mkdir plugins/src

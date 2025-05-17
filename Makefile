install:
	@sudo go build -o /usr/local/bin/linkt *.go
	@echo "\nlinkt was \033[32msuccessfully\033[0m installed!\n"

uninstall:
	@sudo rm /usr/local/bin/linkt
	@echo "\nlinkt was \033[32msuccessfully\033[0m uninstalled :(\n"

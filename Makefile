install:
	@sudo go build -o /usr/local/bin/linkt *.go
	@echo "\nlinkt was \033[32msuccessfully\033[0m installed!\n"

uninstall:
	@sudo rm /usr/local/bin/linkt
	@echo "\nlinkt was \033[32msuccessfully\033[0m uninstalled :(\n"

help:
	@echo "\nUsage: make <command>\n"
	@echo "Commands:"
	@echo "\tinstall\t\tInstall linkt globally."
	@echo "\tuninstall\tUninstall linkt globally."
	@echo "\thelp\t\tDisplay help for a command.\n"

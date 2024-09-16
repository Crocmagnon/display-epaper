GOK := gok -i display
run:
	${GOK} run
update:
	${GOK} update
logs:
	${GOK} logs -s display-epaper
get:
	${GOK} get --update_all
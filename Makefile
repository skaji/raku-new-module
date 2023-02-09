export APP_RECENT_URL := https://raku.land/recent/json
export APP_TICK := 30
export DEBUG := 1

run:
	go run cmd/raku-new-module/main.go -config-from-env

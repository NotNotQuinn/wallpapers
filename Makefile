# Some variables
OUT_DIR=./bin
EXE=
CLI_BIN=$(OUT_DIR)/wallpaper$(EXE)
# Damn windows... even though I still use it
ifeq ($(OS),Windows_NT)
	EXE=.exe
endif

cli:
	go build -o $(CLI_BIN) ./cmd
debug: cli 
	$(CLI_BIN) random

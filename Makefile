# Some variables
out_dir=./bin
EXE=

# Damn windows... even though I still use it
ifeq ($(OS),Windows_NT)
	EXE=.exe
endif

cli:
	go build -o $(out_dir)/wallpaper$(EXE) .
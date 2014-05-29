solutionRoot = $(abspath ../../../)
prefix = $(addprefix $(solutionRoot), /bin/)

projectName = lobby
outFile = $(projectName)
outFullPath = $(addprefix $(prefix)/, $(outFile))

.PHONY: build
build:
	GOPATH=$(solutionRoot) go build -o $(outFullPath)
	

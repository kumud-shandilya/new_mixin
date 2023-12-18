package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

func upper(data string) string {

	return strings.ToUpper(data)
}

func main() {

	projectPath := "../../../ClassLibraryProjects/FabricApp/FabricApp.csproj"

	// Create a command to run the dotnet command
	cmd := exec.Command("dotnet", "run", "--project", projectPath)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(stdout)

	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", upper(string(data)))
	///_, err = builder.ExecuteSingleStepAction(ctx, m.RuntimeConfig, action)
}

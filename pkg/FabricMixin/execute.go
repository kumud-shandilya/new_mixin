package FabricMixin

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"get.porter.sh/porter/pkg/exec/builder"
	"gopkg.in/yaml.v2"
)

func (m *Mixin) loadAction(ctx context.Context) (*Action, error) {
	var action Action
	err := builder.LoadAction(ctx, m.RuntimeConfig, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &action)
		return &action, err
	})
	return &action, err
}

func upper(data string) string {

	return strings.ToUpper(data)
}

func (m *Mixin) Execute(ctx context.Context) error {
	action, err := m.loadAction(ctx)
	if err != nil {
		return err
	}
	fmt.Println(action)
	///projectPath := "../../../../main.go"

	// Create a command to run the dotnet command
	///cmd := exec.Command("go", "run", projectPath)
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
	return err
}

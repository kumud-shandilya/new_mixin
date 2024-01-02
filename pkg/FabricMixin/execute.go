package FabricMixin

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

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

func (m *Mixin) Execute(ctx context.Context) error {
	action, err := m.loadAction(ctx)
	if err != nil {
		return err
	}
	fmt.Println(action)

	filepath := "/cnab/app/Library/ubuntu.16.04-x64"
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		fmt.Println("File does not exist")
		return err
	}

	fmt.Println("File exists")

	cmd := exec.Command("dotnet", filepath+"/Microsoft.Fabric.Provisioning.Client.dll", "create --token", "\"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjVCM25SeHRRN2ppOGVORGMzRnkwNUtmOTdaRSIsImtpZCI6IjVCM25SeHRRN2ppOGVORGMzRnkwNUtmOTdaRSJ9.eyJhdWQiOiJodHRwczovL2FuYWx5c2lzLndpbmRvd3MubmV0L3Bvd2VyYmkvYXBpIiwiaXNzIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvNzJmOTg4YmYtODZmMS00MWFmLTkxYWItMmQ3Y2QwMTFkYjQ3LyIsImlhdCI6MTcwMzc0NTQ4OCwibmJmIjoxNzAzNzQ1NDg4LCJleHAiOjE3MDM3NDk1ODEsImFjY3QiOjAsImFjciI6IjEiLCJhaW8iOiJBWVFBZS84VkFBQUFGQ2c5Q3hLSnpGZzRzUnJKK3d6U0RrWEF1bWsxcWNKOWMwNWVJanVza05uaVVWSTJ2MnZMYjZINmdxUURNNFZCK0JPVkNyTDkyWmpQYXhYbHQxR01GMUoyZW9ZMHJwSTk2dnVPZFFVdFZxMlRpNGl3d0pwbDRWdnJPc0JDb0l4NncxWkpBNjZlT1dqd0VPc1QrbllXcWJ5SElROHk1ako3Ui9oaWdEUm9DUEU9IiwiYW1yIjpbInJzYSIsIm1mYSJdLCJhcHBpZCI6Ijg3MWMwMTBmLTVlNjEtNGZiMS04M2FjLTk4NjEwYTdlOTExMCIsImFwcGlkYWNyIjoiMCIsImNvbnRyb2xzIjpbImFwcF9yZXMiXSwiY29udHJvbHNfYXVkcyI6WyIwMDAwMDAwOS0wMDAwLTAwMDAtYzAwMC0wMDAwMDAwMDAwMDAiLCIwMDAwMDAwMy0wMDAwLTBmZjEtY2UwMC0wMDAwMDAwMDAwMDAiXSwiZGV2aWNlaWQiOiJlZGM3Y2Y2My03YWRjLTQ2OGYtODhiMi05NDdiNDA4NzU2N2IiLCJmYW1pbHlfbmFtZSI6IlNoYW5kaWx5YSIsImdpdmVuX25hbWUiOiJLdW11ZCIsImlwYWRkciI6IjI0MDQ6ZjgwMTo4MDI4OjE6ZjRhMTo3YjA2OmNhMGY6ODBjIiwibmFtZSI6Ikt1bXVkIFNoYW5kaWx5YSIsIm9pZCI6IjA5ODBkZWVkLTdlNGQtNGIzZC1iYzE4LWU0NDQzZGEyYzBiMCIsIm9ucHJlbV9zaWQiOiJTLTEtNS0yMS0yMTI3NTIxMTg0LTE2MDQwMTI5MjAtMTg4NzkyNzUyNy02OTQ0NDA3MiIsInB1aWQiOiIxMDAzMjAwMkMwNzFBMTk4IiwicmgiOiIwLkFSb0F2NGo1Y3ZHR3IwR1JxeTE4MEJIYlJ3a0FBQUFBQUFBQXdBQUFBQUFBQUFBYUFETS4iLCJzY3AiOiJ1c2VyX2ltcGVyc29uYXRpb24iLCJzaWduaW5fc3RhdGUiOlsiZHZjX21uZ2QiLCJkdmNfY21wIiwiaW5rbm93bm50d2siLCJrbXNpIl0sInN1YiI6IkpqancxajlhV0VLN0J3aVIzQnd1YjRSdDFhVVo5RWk3Z2hBOUN0WW40bkEiLCJ0aWQiOiI3MmY5ODhiZi04NmYxLTQxYWYtOTFhYi0yZDdjZDAxMWRiNDciLCJ1bmlxdWVfbmFtZSI6ImtzaGFuZGlseWFAbWljcm9zb2Z0LmNvbSIsInVwbiI6ImtzaGFuZGlseWFAbWljcm9zb2Z0LmNvbSIsInV0aSI6ImRvZzFoNjA5YlU2ckkzZmZwc2FOQVEiLCJ2ZXIiOiIxLjAiLCJ3aWRzIjpbImI3OWZiZjRkLTNlZjktNDY4OS04MTQzLTc2YjE5NGU4NTUwOSJdLCJ4bXNfY2MiOlsiQ1AxIl19.TjUZ3Pyhu9MjD_8cyiDDdwB4O-ox4LTkUa9KUTkbCln23v0X1IfnU8fL2-pTjWNUvZG-gZtds6vGTNmjKwHEi-ICYCdMETlj4jUYdw2ha9dy220ydBwOeovWh8uobZ6y97UAzF1iqVPLvqKR89qCpRV1o6VzxibfCsLXb6m9t9UwEB-mY61OYyBBUG0-pBIDGiCxsEBkQbd8nEV9YCNN--lSE0-5oNDdiZq9EGYhD4o6UEq1F1V1TS_Rhhx8jeahj8mLYw3w7rVvZAmF84yH7El4NFIBVpUYcwNntEfPd79NPHxWp5s77j_vztbWQg6wxUI6h2vGBzpiJC-kUU1_SA\"", "--payload '{\\\"displayName\\\":\\\"kumud1234567\\\",\\\"description\\\":\\\"any\\\"}'")
	fmt.Println(cmd)
	// File exists, perform further operations
	///projectPath := "../../../../main.go"

	// Create a command to run the dotnet command
	///cmd := exec.Command("go", "run", projectPath)
	// projectPath := "/Library/ubuntu.16.04-x64"
	// // Create a command to run the dotnet command
	// cmd := exec.Command("/bin/bash", projectPath)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
	}

	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}

	data, err := ioutil.ReadAll(stdout)

	if err != nil {
		fmt.Println(err)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s\n", string(data))
	return err
	///_, err = builder.ExecuteSingleStepAction(ctx, m.RuntimeConfig, action)
	//return err
}

package FabricMixin

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"get.porter.sh/porter/pkg/exec/builder"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type InstallAction struct {
	Steps []InstallStep `yaml:"install"`
}

type InstallStep struct {
	InstallArguments `yaml:"FabricMixin"`
}

type InstallArguments struct {
	Token string                 `yaml:"token"`
	Args  map[string]interface{} `yaml:"arguments"`
}

func (m *Mixin) getPayloadData() ([]byte, error) {
	reader := bufio.NewReader(m.In)
	data, err := ioutil.ReadAll(reader)
	return data, errors.Wrap(err, "could not read the payload from STDIN")
}

func (m *Mixin) loadAction(ctx context.Context) (*Action, error) {
	var action Action
	err := builder.LoadAction(ctx, m.RuntimeConfig, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &action)
		return &action, err
	})
	return &action, err
}

func (m *Mixin) Execute(ctx context.Context) error {
	fmt.Fprintln(m.Out, "Starting deployment operations...")

	payload, err := m.getPayloadData()
	if err != nil {
		log.Println(err)
		return err
	}

	var action InstallAction
	err = yaml.Unmarshal(payload, &action)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(action.Steps[0].Token, action.Steps[0].Args)

	filepath := "/usr/local/bin/ubuntu.16.04-x64"
	if _, err := os.Stat(filepath + "/Microsoft.Fabric.Provisioning.Client"); os.IsNotExist(err) {
		fmt.Println("File does not exist")
		return err
	}

	fmt.Println("File exists")
	jsonString, err := json.Marshal(action.Steps[0].Args)

	cmd := exec.Command(filepath+"/Microsoft.Fabric.Provisioning.Client", "create", "--token", action.Steps[0].Token, "--payload", string(jsonString))
	//eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjVCM25SeHRRN2ppOGVORGMzRnkwNUtmOTdaRSIsImtpZCI6IjVCM25SeHRRN2ppOGVORGMzRnkwNUtmOTdaRSJ9.eyJhdWQiOiJodHRwczovL2FuYWx5c2lzLndpbmRvd3MubmV0L3Bvd2VyYmkvYXBpIiwiaXNzIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvNzJmOTg4YmYtODZmMS00MWFmLTkxYWItMmQ3Y2QwMTFkYjQ3LyIsImlhdCI6MTcwNDg5MTA3MiwibmJmIjoxNzA0ODkxMDcyLCJleHAiOjE3MDQ4OTU4MDgsImFjY3QiOjAsImFjciI6IjEiLCJhaW8iOiJBWVFBZS84VkFBQUFTQnlJZUllemd1SXorN1JoSXdrSnNzTUpCTlJTd21KRk9WWmJCd3daNllzU2IrVWJSV1RmV3YvVS9MVjdKdXJRbWhIRFhxbWVPOHVHRElGQ0RjS2pyekVTdEhOUTlvTUdvS2ZERExwYWFJRFI4NTNjc0FNbHJCUEZSd1NDK2pRdUhDYkwwcFRpa0lkUXZRUnRpL1h5UGFkTHVEV29FZlN5NlJVeTBNcU12ZTg9IiwiYW1yIjpbInJzYSIsIm1mYSJdLCJhcHBpZCI6Ijg3MWMwMTBmLTVlNjEtNGZiMS04M2FjLTk4NjEwYTdlOTExMCIsImFwcGlkYWNyIjoiMCIsImNvbnRyb2xzIjpbImFwcF9yZXMiXSwiY29udHJvbHNfYXVkcyI6WyIwMDAwMDAwOS0wMDAwLTAwMDAtYzAwMC0wMDAwMDAwMDAwMDAiLCIwMDAwMDAwMy0wMDAwLTBmZjEtY2UwMC0wMDAwMDAwMDAwMDAiXSwiZGV2aWNlaWQiOiJlZGM3Y2Y2My03YWRjLTQ2OGYtODhiMi05NDdiNDA4NzU2N2IiLCJmYW1pbHlfbmFtZSI6IlNoYW5kaWx5YSIsImdpdmVuX25hbWUiOiJLdW11ZCIsImlwYWRkciI6IjI0MDQ6ZjgwMTo4MDI4OjE6N2ExNToxNGU4OmYwYWI6ZjBjOSIsIm5hbWUiOiJLdW11ZCBTaGFuZGlseWEiLCJvaWQiOiIwOTgwZGVlZC03ZTRkLTRiM2QtYmMxOC1lNDQ0M2RhMmMwYjAiLCJvbnByZW1fc2lkIjoiUy0xLTUtMjEtMjEyNzUyMTE4NC0xNjA0MDEyOTIwLTE4ODc5Mjc1MjctNjk0NDQwNzIiLCJwdWlkIjoiMTAwMzIwMDJDMDcxQTE5OCIsInJoIjoiMC5BUm9BdjRqNWN2R0dyMEdScXkxODBCSGJSd2tBQUFBQUFBQUF3QUFBQUFBQUFBQWFBRE0uIiwic2NwIjoidXNlcl9pbXBlcnNvbmF0aW9uIiwic2lnbmluX3N0YXRlIjpbImR2Y19tbmdkIiwiZHZjX2NtcCIsImlua25vd25udHdrIiwia21zaSJdLCJzdWIiOiJKamp3MWo5YVdFSzdCd2lSM0J3dWI0UnQxYVVaOUVpN2doQTlDdFluNG5BIiwidGlkIjoiNzJmOTg4YmYtODZmMS00MWFmLTkxYWItMmQ3Y2QwMTFkYjQ3IiwidW5pcXVlX25hbWUiOiJrc2hhbmRpbHlhQG1pY3Jvc29mdC5jb20iLCJ1cG4iOiJrc2hhbmRpbHlhQG1pY3Jvc29mdC5jb20iLCJ1dGkiOiJhblc1cGVyaERrdXczWlo3eElKYkFnIiwidmVyIjoiMS4wIiwid2lkcyI6WyJiNzlmYmY0ZC0zZWY5LTQ2ODktODE0My03NmIxOTRlODU1MDkiXSwieG1zX2NjIjpbIkNQMSJdfQ.aAEPekuOlkQhI5AvfhMeIjS_0REAFgvDVBCsZBzQpiIWdTgbCaa6jXI2C0hM4I1Mo9kJMz12XtOyAGF2FCxegSgRC7cJgrgN7fqCK6-NmkK76GRkSLblo3215SeWF900gYNAQDr6dbONod6WQxqv2EO6xzIbZattUkK-G9jldILQ-7QXa5kL1wVRoxj3q-BPkVVq6gwNZEWtymd91q4MGdirx4otFzFp_Gyahb7gh4-fiCtKvIoPSXxVofjr-RbKskTDubfxImHdjjKFulLo8utVZUtiaPwqKbEA85D_tl9yR_eMBSiQcTmxvD8JkaZW8Y-ZSgx_wqBolaXfsje-7A

	// cmd := exec.Command(filepath + "/Microsoft.Fabric.Provisioning.Client create --token \"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjVCM25SeHRRN2ppOGVORGMzRnkwNUtmOTdaRSIsImtpZCI6IjVCM25SeHRRN2ppOGVORGMzRnkwNUtmOTdaRSJ9.eyJhdWQiOiJodHRwczovL2FuYWx5c2lzLndpbmRvd3MubmV0L3Bvd2VyYmkvYXBpIiwiaXNzIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvNzJmOTg4YmYtODZmMS00MWFmLTkxYWItMmQ3Y2QwMTFkYjQ3LyIsImlhdCI6MTcwNDE5MzkzOCwibmJmIjoxNzA0MTkzOTM4LCJleHAiOjE3MDQxOTg5NTEsImFjY3QiOjAsImFjciI6IjEiLCJhaW8iOiJBWVFBZS84VkFBQUFqZ25qa0h4WmtIL3BJaXJaclhCRjRHYnNwdFYzUzF5RG13NlRZMnNEWk9GQlZxUStXSjdqbG5VQzNYVzliNHliMGorZjh6c0szZXVkakk5RlE4SkdMVkhTYU9oSXJKL3gyb2t4eEhOVUY0bW05TVlqRFBNek92bFlpK1B0T3RBaFo1Nmo4U0dLaERxa2VZbmhmcXByMytNSE5GWDY4dEp0NExHMFpuWGpJWTQ9IiwiYW1yIjpbInJzYSIsIm1mYSJdLCJhcHBpZCI6Ijg3MWMwMTBmLTVlNjEtNGZiMS04M2FjLTk4NjEwYTdlOTExMCIsImFwcGlkYWNyIjoiMCIsImNvbnRyb2xzIjpbImFwcF9yZXMiXSwiY29udHJvbHNfYXVkcyI6WyIwMDAwMDAwOS0wMDAwLTAwMDAtYzAwMC0wMDAwMDAwMDAwMDAiLCIwMDAwMDAwMy0wMDAwLTBmZjEtY2UwMC0wMDAwMDAwMDAwMDAiXSwiZGV2aWNlaWQiOiJlZGM3Y2Y2My03YWRjLTQ2OGYtODhiMi05NDdiNDA4NzU2N2IiLCJmYW1pbHlfbmFtZSI6IlNoYW5kaWx5YSIsImdpdmVuX25hbWUiOiJLdW11ZCIsImlwYWRkciI6IjI0MDQ6ZjgwMTo4MDI4OjE6MTkxNToyZDU3OjY2MmY6ZDBmMyIsIm5hbWUiOiJLdW11ZCBTaGFuZGlseWEiLCJvaWQiOiIwOTgwZGVlZC03ZTRkLTRiM2QtYmMxOC1lNDQ0M2RhMmMwYjAiLCJvbnByZW1fc2lkIjoiUy0xLTUtMjEtMjEyNzUyMTE4NC0xNjA0MDEyOTIwLTE4ODc5Mjc1MjctNjk0NDQwNzIiLCJwdWlkIjoiMTAwMzIwMDJDMDcxQTE5OCIsInJoIjoiMC5BUm9BdjRqNWN2R0dyMEdScXkxODBCSGJSd2tBQUFBQUFBQUF3QUFBQUFBQUFBQWFBRE0uIiwic2NwIjoidXNlcl9pbXBlcnNvbmF0aW9uIiwic2lnbmluX3N0YXRlIjpbImR2Y19tbmdkIiwiZHZjX2NtcCIsImlua25vd25udHdrIiwia21zaSJdLCJzdWIiOiJKamp3MWo5YVdFSzdCd2lSM0J3dWI0UnQxYVVaOUVpN2doQTlDdFluNG5BIiwidGlkIjoiNzJmOTg4YmYtODZmMS00MWFmLTkxYWItMmQ3Y2QwMTFkYjQ3IiwidW5pcXVlX25hbWUiOiJrc2hhbmRpbHlhQG1pY3Jvc29mdC5jb20iLCJ1cG4iOiJrc2hhbmRpbHlhQG1pY3Jvc29mdC5jb20iLCJ1dGkiOiJBTnNqMmJsQ1VFR0FOUmpiLXVHMEFRIiwidmVyIjoiMS4wIiwid2lkcyI6WyJiNzlmYmY0ZC0zZWY5LTQ2ODktODE0My03NmIxOTRlODU1MDkiXSwieG1zX2NjIjpbIkNQMSJdfQ.xId5Vk-6cZkMWEdg9Q5GEloRWXeqNaxdbTTbL-ttqayDWWbeOGWr6L8adFZ--RVayaPCSGc-Whp0v7R9jtuVGs2tzGcKfcXBXyXFkdQIdvd5cm9pDbPb6r79VuB7KeBhsKBzi1c5IKDdfv27PJ98YNKlWV2YTRFtZZ2mS6147-6oYpDUm6BQ89NMXJ5oUddQGONBm9IDimanaXYEXIvK0ob90r1fX0CRkguBJNlEAniWHu4J9FvZ0YaX-oQKRhzlZflqE8nPlkpOIv9IvS2Hkwp76V9O1qklHVGiyTQbTWOfcLUhtHtn56gu8iVys_anJQi2fLBMRHMFDVmeQ9D9Nw\" --payload '{\"displayName\":\"kumud12345678\",\"description\":\"any\"}'")
	fmt.Println(cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Output:", string(output))
	//fmt.Println(cmd)
	// // File exists, perform further operations
	// ///projectPath := "../../../../main.go"

	// // Create a command to run the dotnet command
	// ///cmd := exec.Command("go", "run", projectPath)
	// // projectPath := "/Library/ubuntu.16.04-x64"
	// // // Create a command to run the dotnet command
	// // cmd := exec.Command("/bin/bash", projectPath)
	// stdout, err := cmd.StdoutPipe()

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// if err := cmd.Start(); err != nil {
	// 	fmt.Println(err)
	// }

	// data, err := ioutil.ReadAll(stdout)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// if err := cmd.Wait(); err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Printf("%s\n", string(data))
	///return err
	//_, err = builder.ExecuteSingleStepAction(ctx, m.RuntimeConfig, action)
	return err
}

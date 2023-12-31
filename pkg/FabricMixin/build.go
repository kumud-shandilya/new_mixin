package FabricMixin

import (
	"context"
	"fmt"
	"text/template"

	"get.porter.sh/porter/pkg/exec/builder"
	"gopkg.in/yaml.v2"
)

// BuildInput represents stdin passed to the mixin for the build command.
type BuildInput struct {
	Config MixinConfig
}

// MixinConfig represents configuration that can be set on the FabricMixin mixin in porter.yaml
// mixins:
// - FabricMixin:
//	  clientVersion: "v0.0.0"

type MixinConfig struct {
	ClientVersion string `yaml:"clientVersion,omitempty"`
}

type buildConfig struct {
	MixinConfig
}

// This is an example. Replace the following with whatever steps are needed to
// install required components into
const dockerfileLines = `RUN apt-get update && apt-get install wget -y
RUN apt-get update && apt-get install -y gpg
RUN wget -O - https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor -o microsoft.asc.gpg
RUN mv microsoft.asc.gpg /etc/apt/trusted.gpg.d/
RUN wget https://packages.microsoft.com/config/debian/11/prod.list
RUN mv prod.list /etc/apt/sources.list.d/microsoft-prod.list
RUN chown root:root /etc/apt/trusted.gpg.d/microsoft.asc.gpg
RUN chown root:root /etc/apt/sources.list.d/microsoft-prod.list
RUN apt-get update && \
    apt-get install -y dotnet-sdk-7.0
`

// Build will generate the necessary Dockerfile lines
// for an invocation image using this mixin
func (m *Mixin) Build(ctx context.Context) error {

	// Create new Builder.
	var input BuildInput

	err := builder.LoadAction(ctx, m.RuntimeConfig, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &input)
		return &input, err
	})
	if err != nil {
		return err
	}

	suppliedClientVersion := input.Config.ClientVersion

	if suppliedClientVersion != "" {
		m.ClientVersion = suppliedClientVersion
	}

	fmt.Fprintf(m.Out, dockerfileLines)

	tmpl, err := template.New("dockerfile").Parse(dockerfileLines)
	if err != nil {
		return fmt.Errorf("error parsing Dockerfile template for the az mixin: %w", err)
	}

	cfg := buildConfig{MixinConfig: input.Config}

	if err = tmpl.Execute(m.Out, cfg); err != nil {
		return fmt.Errorf("error generating Dockerfile lines for the az mixin: %w", err)
	}

	// Example of pulling and defining a client version for your mixin
	//fmt.Fprintf(m.Out, "\nRUN curl https://get.helm.sh/helm-%s-linux-amd64.tar.gz --output helm3.tar.gz", m.ClientVersion)

	return nil
}

package cli

import (
	"encoding/json"
	gofrogcmd "github.com/jfrog/gofrog/io"
	artifactoryCLI "github.com/jfrog/jfrog-cli-artifactory/artifactory/cli"
	distributionCLI "github.com/jfrog/jfrog-cli-artifactory/distribution/cli"
	evidenceCLI "github.com/jfrog/jfrog-cli-artifactory/evidence/cli"
	"github.com/jfrog/jfrog-cli-artifactory/lifecycle"
	"github.com/jfrog/jfrog-cli-core/v2/common/cliutils"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"os"
	"strings"
)

type jsonFlag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type jsonCommandsStruct struct {
	Name        string     `json:"name"`
	Description string     `json:"usage"`
	Flags       []jsonFlag `json:"flags"`
	Action      string     `json:"action"`
	Arguments   []string   `json:"arguments,omitempty"`
	BuildArgs   string     `json:"build-args,omitempty"`
}

func RunActions(ctx *components.Context) error {
	executableCommand := strings.Split(ctx.ExecutableCommand, " ")
	args := strings.Join(append([]string{}, executableCommand[1:]...), " ")
	command := gofrogcmd.NewCommand(executableCommand[0], args, []string{})
	output, cmdError, _, err := gofrogcmd.RunCmdWithOutputParser(command, true)
	if err != nil {
		log.Error("Error occurred while running command: ", ctx.ExecutableCommand, output, cmdError, err)
	}
	return err
}

func GetJfrogCliArtifactoryApp() components.App {
	data, err := os.ReadFile("/Users/kanishkg/workspace/eco/jfrog-cli-artifactory/cli/main.json")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	var jsonCommands []jsonCommandsStruct
	err = json.Unmarshal(data, &jsonCommands)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	var requiredCommand []components.Command

	for _, command := range jsonCommands {
		newCommand := components.Command{
			Name:        command.Name,
			Description: command.Description,
			Action:      RunActions,
			Flags:       []components.Flag{},
			Category:    command.Action,
		}
		requiredCommand = append(requiredCommand, newCommand)
	}

	app := components.CreateEmbeddedApp(
		"artifactory",
		[]components.Command{},
	)
	app.Subcommands = append(app.Subcommands, components.Namespace{
		Name:        string(cliutils.Ds),
		Description: "Distribution V1 commands.",
		Commands:    distributionCLI.GetCommands(),
		Category:    "Command Namespaces",
	})
	app.Subcommands = append(app.Subcommands, components.Namespace{
		Name:        "evd",
		Description: "Evidence commands.",
		Commands:    evidenceCLI.GetCommands(),
		Category:    "Command Namespaces",
	})
	app.Subcommands = append(app.Subcommands, components.Namespace{
		Name:        string(cliutils.Rt),
		Description: "Artifactory commands.",
		Commands:    artifactoryCLI.GetCommands(),
		Category:    "Command Namespaces",
	})
	app.Commands = append(app.Commands, lifecycle.GetCommands()...)
	app.Commands = append(app.Commands, requiredCommand...)
	return app
}

package cli

import (
	"runtime/debug"
	"time"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/output"
)

// buildInfo descreve o binário em execução.
type buildInfo struct {
	Version string
	Commit  string
	Date    time.Time
}

// resolveVersion usa a versão injetada por -ldflags; senão, a do módulo
// (go install …@vX.Y.Z) e os dados de VCS do build; senão, "dev".
func resolveVersion(injected string) buildInfo {
	bi := buildInfo{Version: injected}
	info, ok := debug.ReadBuildInfo()
	if ok {
		if bi.Version == "" && info.Main.Version != "" && info.Main.Version != "(devel)" {
			bi.Version = info.Main.Version
		}
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				bi.Commit = s.Value
			case "vcs.time":
				bi.Date, _ = time.Parse(time.RFC3339, s.Value)
			}
		}
	}
	if bi.Version == "" {
		bi.Version = "dev"
	}
	return bi
}

func newVersionCmd(injected string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Mostra a versão do ctbz",
		Args:  exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			bi := resolveVersion(injected)
			rec := &output.Record{}
			rec.Add("versao", "Versão", bi.Version).
				Add("commit", "Commit", nilIfEmpty(bi.Commit)).
				Add("data", "Data", output.DateTime{Time: bi.Date})
			return output.Write(cmd.OutOrStdout(), f, rec)
		},
	}
}

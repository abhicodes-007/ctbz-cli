package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// GenMarkdown gera a referência de comandos em dir: um arquivo por comando visível
// e um README.md com o índice. A saída é determinística (sem datas), para que um
// teste possa detectar documentação desatualizada.
func GenMarkdown(root *cobra.Command, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	files := map[string][]byte{"README.md": indexPage(root)}
	for _, c := range visibleCommands(root) {
		files[docFile(c)] = commandPage(c)
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), content, 0o644); err != nil {
			return err
		}
	}
	return nil
}

// visibleCommands percorre a árvore em profundidade, sem comandos ocultos.
func visibleCommands(c *cobra.Command) []*cobra.Command {
	out := []*cobra.Command{c}
	for _, sub := range c.Commands() {
		if sub.IsAvailableCommand() {
			out = append(out, visibleCommands(sub)...)
		}
	}
	return out
}

func docFile(c *cobra.Command) string {
	return strings.ReplaceAll(c.CommandPath(), " ", "_") + ".md"
}

func indexPage(root *cobra.Command) []byte {
	var b bytes.Buffer
	b.WriteString("# Referência de comandos\n\n")
	b.WriteString("Gerada a partir da própria CLI (`go run ./tools/gendocs`); não edite à mão.\n\n")
	b.WriteString("| Comando | Descrição |\n|---|---|\n")
	for _, c := range visibleCommands(root) {
		fmt.Fprintf(&b, "| [`%s`](%s) | %s |\n", c.CommandPath(), docFile(c), c.Short)
	}
	return b.Bytes()
}

func commandPage(c *cobra.Command) []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "# %s\n\n%s\n\n", c.CommandPath(), c.Short)
	if c.Long != "" {
		fmt.Fprintf(&b, "%s\n\n", c.Long)
	}
	if c.Runnable() {
		fmt.Fprintf(&b, "## Uso\n\n```\n%s\n```\n\n", c.UseLine())
	}
	if c.Example != "" {
		fmt.Fprintf(&b, "## Exemplos\n\n```sh\n%s\n```\n\n", c.Example)
	}
	if subs := c.Commands(); c.HasAvailableSubCommands() {
		b.WriteString("## Subcomandos\n\n")
		for _, sub := range subs {
			if sub.IsAvailableCommand() {
				fmt.Fprintf(&b, "- [`%s`](%s): %s\n", sub.CommandPath(), docFile(sub), sub.Short)
			}
		}
		b.WriteString("\n")
	}
	if f := c.NonInheritedFlags(); f.HasAvailableFlags() {
		fmt.Fprintf(&b, "## Flags\n\n```\n%s```\n\n", f.FlagUsages())
	}
	if f := c.InheritedFlags(); f.HasAvailableFlags() {
		fmt.Fprintf(&b, "## Flags globais\n\n```\n%s```\n\n", f.FlagUsages())
	}
	if c.HasParent() {
		p := c.Parent()
		fmt.Fprintf(&b, "Veja também: [`%s`](%s).\n", p.CommandPath(), docFile(p))
	}
	return append(bytes.TrimRight(b.Bytes(), "\n"), '\n')
}

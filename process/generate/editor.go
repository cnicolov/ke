package generate

import (
	"golang.org/x/net/context"
	"kego.io/context/envctx"
	"kego.io/generator"
	"kego.io/kerr"
)

func Editor(ctx context.Context, env *envctx.Env) (source []byte, err error) {

	g := generator.New("main")

	g.Imports.Anonymous("kego.io/system")
	g.Imports.Anonymous("kego.io/system/editors")
	g.Imports.Anonymous(env.Path)
	for p, _ := range env.Aliases {
		g.Imports.Anonymous(p)
	}
	/*
		func main() {
			if err := client.Start(); err != nil {
				console.Error(err.Error())
			}
		}
	*/
	g.Println("func main() {")
	{
		clientStart := generator.Reference("kego.io/editor/client", "Start", env.Path, g.Imports.Add)
		g.Println("if err := ", clientStart, "(); err != nil {")
		{
			consoleError := generator.Reference("kego.io/editor/client/console", "Error", env.Path, g.Imports.Add)
			g.Println(consoleError, "(err.Error())")
		}
		g.Println("}")
	}
	g.Println("}")

	b, err := g.Build()
	if err != nil {
		return nil, kerr.New("CBTOLUQOGL", err, "Build")
	}
	return b, nil
}

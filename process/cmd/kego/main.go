package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"kego.io/kerr"
	"kego.io/process"
	_ "kego.io/system/types"
)

func main() {
	var testFlag = flag.Bool("test", false, "test mode? e.g. don't write the files")
	var pathFlag = flag.String("path", "", "full package path e.g. github.com/foo/bar")
	flag.Parse()
	testMode := *testFlag
	packagePath := *pathFlag

	currentDir, err := filepath.Abs("")
	if err != nil {
		log.Fatalf("Error getting the current working directory:\n%v\n", err.Error())
	}
	typesDir := filepath.Join(currentDir, "types")

	if packagePath == "" {
		path, err := GetPackage(currentDir)
		if err != nil {
			log.Fatalf("Error while getting package path from current working directory. You can specify the full package path with the -path=github.com/foo\n%v\n", err.Error())
		}
		packagePath = path
	}

	imports := map[string]string{}

	if err := process.Scan(currentDir, packagePath, imports); err != nil {
		panic(err)
	}

	mainSource, typesSource, err := process.Generate(packagePath, imports)
	if err != nil {
		panic(err)
	}

	if testMode {
		fmt.Printf("######## Main ########\n%s######## Types ########\n%s", mainSource, typesSource)
		return
	}

	if err = save(currentDir, mainSource); err != nil {
		panic(err)
	}

	if err = save(typesDir, typesSource); err != nil {
		panic(err)
	}

}

func save(dir string, contents []byte) error {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0777); err != nil {
			return kerr.New("BPGOUIYPXO", err, "process/cmd/kego.save", "os.MkdirAll")
		}
	}

	name := "generated.go"
	file := filepath.Join(dir, name)
	backup(file, filepath.Join(dir, fmt.Sprintf("%s.backup", name)))

	output, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return kerr.New("NWLWHSGJWP", err, "process/cmd/kego.save", "os.OpenFile (could not open output file)")
	}
	defer output.Close()

	output.Write(contents)

	return nil
}

func backup(file string, backup string) {

	if _, err := os.Stat(backup); err == nil {
		os.Remove(backup)
	}

	if _, err := os.Stat(file); err == nil {
		os.Rename(file, backup)
	}
}

func GetPackage(dir string) (string, error) {
	gopath := os.Getenv("GOPATH")
	if strings.HasPrefix(gopath, "file://") {
		// This is to fix a bug when running the code in IntelliJ. This can be removed when the
		// bug is fixed: https://github.com/go-lang-plugin-org/go-lang-idea-plugin/issues/1600
		gopath = gopath[7:]
	}
	return getPackage(dir, gopath)
}
func getPackage(dir string, gopathEnv string) (string, error) {
	gopaths := filepath.SplitList(gopathEnv)
	var savedError error
	for _, gopath := range gopaths {
		if strings.HasPrefix(dir, gopath) {
			gosrc := fmt.Sprintf("%s/src", gopath)
			relpath, err := filepath.Rel(gosrc, dir)
			if err != nil {
				savedError = err
				continue
			}
			if relpath == "" {
				continue
			}
			return relpath, nil
		}
	}
	if savedError != nil {
		return "", savedError
	}
	return "", kerr.New("CXOETFPTGM", nil, "process/cmd/kego.getPackage", "Package not found for %s", dir)
}

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/deckarep/golang-set"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
)

const (
	EnableDebug          = true
	OpAdd                = "add"
	OpRun                = "run"
	OpDelete             = "delete"
	OpConfig             = "config"
	OpPush               = "push"
	OpPull               = "pull"
	OpLs                 = "ls"
	BaseDir              = "bm"
	ConfigFile           = "config.json"
	DefaultStorageFolder = "bm"
)

var (
	validConfigSet = mapset.NewSet("DBPath")
)

func init() {
	log.SetFlags(0)
	initBaseDir()
}

func getFullPath(path string) string {
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	return fullPath
}

func push(baseDir string) {
	gitDir := "--git-dir=" + filepath.Join(baseDir, ".git")
	workTree := "--work-tree=" + baseDir
	cmd := exec.Command("git", gitDir, workTree, "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmd = exec.Command("git", gitDir, workTree, "commit", "-m", "\"update\"")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmd = exec.Command("git", gitDir, workTree, "push")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func pull(baseDir string) {
	gitDir := "--git-dir=" + filepath.Join(baseDir, ".git")
	workTree := "--work-tree=" + baseDir
	cmd := exec.Command("git", gitDir, workTree, "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func printConfig(config Config) {
	s := reflect.ValueOf(&config).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%s: %s", typeOfT.Field(i).Name, f.Interface())
	}
}

func convertSetToString(set mapset.Set) string {
	it := set.Iterator()
	items := make([]string, 0, len(it.C))

	for elem := range it.C {
		items = append(items, fmt.Sprintf("%v", elem))
	}
	return fmt.Sprintf("%s", strings.Join(items, ", "))
}

func newError(err error) error {
	if err == nil {
		return nil
	}
	return cli.NewExitError(err.Error(), -1)
}

func newErrorWithText(text string) error {
	return cli.NewExitError(text, -1)
}

func addScript(args []string) error {
	key := args[0]

	if Get(key) != "" {
		flag := false
		prompt := &survey.Confirm{
			Message: "Do you want to override the exists key?",
		}
		survey.AskOne(prompt, &flag, nil)
		if !flag {
			return nil
		}
	}
	bashPath := args[1]
	fullPath := getFullPath(bashPath)
	if fullPath == "" {
		return newErrorWithText("Invalid bash script path: " + bashPath)
	}
	//Debug(fullPath)
	Put(key, filepath.ToSlash(fullPath))
	return nil
}

func runScript(args []string) error {
	key := args[0]
	val := Get(key)
	if val == "" {
		return newErrorWithText("Invalid key: " + key)
	}

	args[0] = val

	ext := filepath.Ext(val)
	var cmd *exec.Cmd
	switch ext {
	case ".py":
		cmd = exec.Command("python", args...)
		break
	case ".jar":
		argsNew := make([]string, 1)
		argsNew[0] = "-jar"
		argsNew = append(argsNew, args...)
		for i, arg := range argsNew {
			if strings.HasPrefix(arg, "_") {
				Debug(i)
				argsNew[i] = strings.Replace(arg, "_", "-", 1)

			}
		}
		Debug(argsNew)
		cmd = exec.Command("java", argsNew...)
		break
	default:
		cmd = exec.Command("bash", args...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	return nil
}

func deleteScript(args []string) {
	key := args[0]
	//if Get(key) != "" {
	//	flag := false
	//	prompt := &survey.Confirm{
	//		Message: "Do you want to delete this key?",
	//	}
	//	survey.AskOne(prompt, &flag, nil)
	//	if !flag {
	//		return
	//	}
	//}

	Delete(key)
	//fullPath := filepath.Join(config.StorePath, key)
	//err := os.RemoveAll(fullPath)
	//if err != nil {
	//	return newError(err)
	//}
}
func handleArgs(op string, args []string) error {

	configPath, err := configFilePath()
	if err != nil {
		return newError(err)
	}
	config, err := LoadConfig(configPath)
	if err != nil {
		return newError(err)
	}

	OpenDB(config)
	InitDBBucket()
	defer CloseDB()

	switch op {
	case OpAdd:
		if len(args) != 2 {
			return newErrorWithText("Valid format: add [key] [script path]")
		}

		addScript(args)
		return nil
		//key := args[0]
		//if FileExists(path.Join(config.StorePath, key)) {
		//	flag := false
		//	prompt := &survey.Confirm{
		//		Message: "Do you want to override the exists key?",
		//	}
		//	survey.AskOne(prompt, &flag, nil)
		//	if !flag {
		//		return nil
		//	}
		//
		//}
		//filePath := args[1]
		//res := checkPath(filePath)
		//if res == -1 {
		//	return newErrorWithText("Invalid template path: " + filePath)
		//}
		//
		//if res == 0 {
		//	err = moveFile(filePath, path.Join(config.StorePath, key), false)
		//} else {
		//	err = moveFile(filePath, path.Join(config.StorePath, key), true)
		//}
		//return newError(err)
	case OpRun:
		if len(args) < 1 {
			return newErrorWithText("Valid format: run [key] [params...]")
		}
		err := runScript(args)
		return err
		//key := args[0]
		//srcPath := path.Join(config.StorePath, key)
		//res := checkPath(srcPath)
		//if res != 1 {
		//	return newErrorWithText("Invalid key: " + key)
		//}
		//currentPath, err := CurrentRunPath()
		//if err != nil {
		//	return newError(err)
		//}
		//err = CopyDir(srcPath, currentPath)
		//return newError(err)
	case OpDelete:
		if len(args) != 1 {
			return newErrorWithText("Valid format: delete [key]")
		}
		deleteScript(args)
		return nil
		//key := args[0]
		//if FileExists(path.Join(config.StorePath, key)) {
		//	flag := false
		//	prompt := &survey.Confirm{
		//		Message: "Do you want to delete this key?",
		//	}
		//	survey.AskOne(prompt, &flag, nil)
		//	if !flag {
		//		return nil
		//	}
		//
		//}
		//fullPath := filepath.Join(config.StorePath, key)
		//err := os.RemoveAll(fullPath)
		//if err != nil {
		//	return newError(err)
		//}
		//break
	case OpConfig:
		if len(args) == 0 {
			printConfig(config)
			return nil
		}
		if len(args) != 2 {
			return newErrorWithText("Valid format: config [type] [value]")
		}
		key := args[0]
		if validConfigSet.Contains(key) {
			value := filepath.FromSlash(args[1])
			absPath, err := filepath.Abs(value)
			if err != nil {
				return newErrorWithText("Invalid config value: " + value)
			}
			err = os.MkdirAll(absPath, 0755)
			if err != nil {
				return newErrorWithText("Invalid config value: " + value)
			}
			config.StorePath = absPath
			configPath, err := configFilePath()
			if err != nil {
				return newError(err)
			}
			err = SaveConfig(configPath, config)
			return newError(err)
		} else {
			return newErrorWithText("Valid configuration types: " + convertSetToString(validConfigSet))
		}
	case OpLs:
		if len(args) == 0 {
			args = []string{""}
		}
		if len(args) != 1 {
			return newErrorWithText("Valid format: ls [prefix]")
		}
		prefix := args[0]
		keys := IterateKey(prefix)
		text := strings.Join(keys, "\n")
		fmt.Println(text)
		//fileList, err := GetFileWithPrefix(config.StorePath, prefix, ignoreFileSet)
		//if err != nil {
		//	return newError(err)
		//}
		//text := strings.Join(fileList, "\n")
		//fmt.Println(text)

	case OpPush:
		push(config.StorePath)
		break
	case OpPull:
		pull(config.StorePath)
		break
	}
	return nil
}

func main() {

	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
    {{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}    {{join .Names ","}}{{"\t"}}{{.ArgsUsage}} {{"\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}
	{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}{{end}}
   {{end}}
`

	app := cli.NewApp()
	app.Name = "bm"
	app.Usage = "A simple bash scripts management tool."
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		{
			Name:      "add",
			Aliases:   []string{"a"},
			Usage:     "Add a bash/py/jar file associated with the specified key to bm.",
			ArgsUsage: "[key] [bash script path]",
			Action: func(c *cli.Context) error {
				return handleArgs(OpAdd, c.Args())
			},
		},

		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage: "Run the target file associated with the specified key. \n" +
				"\t\tIf you want passing parameters to the target file, you should use \"_\" instead of \"-\"",
			ArgsUsage: "[key] [params...]",
			Action: func(c *cli.Context) error {
				return handleArgs(OpRun, c.Args())

			},
		},

		{
			Name:      "delete",
			Aliases:   []string{"d"},
			Usage:     "Delete the key from bm.",
			ArgsUsage: "[key]",
			Action: func(c *cli.Context) error {
				return handleArgs(OpDelete, c.Args())

			},
		},

		{
			Name:      "ls",
			Aliases:   []string{"l"},
			Usage:     "List the keys that begin with prefix.",
			ArgsUsage: "[prefix]",
			Action: func(c *cli.Context) error {
				return handleArgs(OpLs, c.Args())
			},
		},

		{
			Name:      "config",
			Aliases:   []string{"c"},
			Usage:     "Set configurations of bm. The valid configuration is \"DBPath\".",
			ArgsUsage: "[type] [value]",
			Action: func(c *cli.Context) error {
				return handleArgs(OpConfig, c.Args())
			},
		},

		{
			Name:  "push",
			Usage: "Call git push command based on the database directory.",
			Action: func(c *cli.Context) error {
				return handleArgs(OpPush, c.Args())
			},
		},

		{
			Name:  "pull",
			Usage: "Call git pull command based on the database directory.",
			Action: func(c *cli.Context) error {
				return handleArgs(OpPull, c.Args())
			},
		},
	}

	app.Run(os.Args)
}

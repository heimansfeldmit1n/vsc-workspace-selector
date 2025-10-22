package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/BurntSushi/toml"
)

type Project struct {
	Name string `toml:"Name"`
	Path string `toml:"Path"`
}
type Config struct {
	Projects []Project `toml:"projects"`
}

func parse_config_file() Config {
	// Placeholder function to parse the configuration file
	log.Println("Parsing configuration file...")
	// Implementation goes here

	var projects Config

	// Decode the TOML file into the slice
	if _, err := toml.DecodeFile("~/vsc-seclector.config.toml", &projects); err != nil {
		log.Fatal(err)
	}

	return projects
}

func add_new_project(project_path string, project_name string) string {
	// Placeholder function to add a new project
	log.Printf("Adding new project at path: %s", project_path)
	// Implementation goes here
	var new_project Project
	new_project.Name = project_name
	new_project.Path = project_path

	// append toml conf file
	file, err := os.OpenFile("~/vsc-seclector.config.toml", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	if _, err = file.WriteString("\n[[projects]]\nName = \"" + new_project.Name + "\"\nPath = \"" + new_project.Path + "\"\n"); err != nil {
		log.Fatal(err)
	}
	return "New project added"
}

func open_project(project_name string, config Config) string {
	// Placeholder function to open a project

	// get project path from config file
	log.Printf("Opening project at path: %s", project_name)

	for _, p := range config.Projects {
		if p.Name == project_name {
			project_name = p.Path
		}
	}

	// Implementation goes here
	out, err := exec.Command("code", project_name).Output()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	log.Printf("Output: %s\n", out)
	return "Project opened"
}
func main() {

	// Define command-line flags
	project := flag.String("project", "default", "Path to the workspace file")
	new_project_path := flag.String("new-project-path", "", "Path to add a new workspace file")
	new_project := flag.String("new-project", "default", "Map name to project path")
	list_projects := flag.Bool("list-projects", false, "List all available workspaces and file Paths")

	flag.StringVar(project, "p", "", "Path to the workspace file")
	flag.StringVar(new_project_path, "npp", "", "Path to add a new workaspace file")
	flag.StringVar(new_project, "np", "default", "Mapp name to project path")
	flag.BoolVar(list_projects, "ls", false, "List all available workspaces and file Paths")
	flag.Parse()

	// Test if flags can be used together
	if !((len(*new_project) > 0 && len(*new_project_path) > 0) || *list_projects || (len(*project) > 0)) {
		log.Fatalf("Provided Flags can not be used to gether")
		os.Exit(1)
	}

	config := parse_config_file()
	if *list_projects {
		for _, p := range config.Projects {
			fmt.Printf("Project=%s, Path='%s'\n", p.Name, p.Path)
		}
	} else if len(*project) > 0 {
		log.Printf(open_project(*project, config))
	} else if len(*new_project) > 0 && len(*new_project_path) > 0 {
		log.Printf(add_new_project(*new_project_path, *new_project))
	}
}

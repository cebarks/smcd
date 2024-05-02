package smcd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var WorkingDir string

type Server struct {
	Name    string
	Folder  string
	Enabled bool

	Running bool
	Pid     int
}

func (server *Server) Start() {
	log.Printf("Starting server: %s", server.Name)
	//TODO check if server is already running under tmux session name
	log.Println(startCommand(server))
	cmd := exec.Command("/usr/bin/tmux", startCommand(server)...)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Could not start server (%s): %v\n %v", server.Name, err, cmd.Stdout)
	}

	server.Pid = cmd.Process.Pid
	server.Running = true

	log.Printf("Server started: %s", server.Name)
}

func (server *Server) Stop() {
	log.Printf("Stopping server: %s", server.Name)
	err := SendCommand(server, "stop")
	if err != nil {
		log.Printf("Could not stop server: %s\n%#v", server.Name, err)
	}

	log.Println("Server stopped.")
}

func SendCommand(server *Server, command string) error {
	cmd := exec.Command("/usr/bin/tmux", "send-keys", "-t", buildTmuxId(server), command, "ENTER")

	err := cmd.Run()
	if err != nil {
		log.Printf("Error when running command (%s) for server (%s):\n%v", command, server.Name, err)
		return err
	}

	return nil
}

func startCommand(server *Server) []string {
	return []string{"new", "-d",
		"-s", buildTmuxId(server),
		"-c", server.Folder,
		fmt.Sprintf("'%s/start.sh'", server.Folder)}
}

func DiscoverServers() []*Server {
	files, err := os.ReadDir(WorkingDir)
	if err != nil {
		log.Fatalf("Could not list entries in working dir:  %v", err)
	}

	var servers []*Server

	for _, file := range files {
		if file.IsDir() {
			path, err := filepath.Abs(filepath.Join(WorkingDir, file.Name()))
			if err != nil {
				log.Fatalf("Could not get full path:%v", err)
			}

			subFiles, err := os.ReadDir(path)
			if err != nil {
				log.Fatalf("Could not get subentries of working dir:%v", err)
			}

			for _, f := range subFiles {
				if f.Name() == "start.sh" {
					servers = append(servers, &Server{
						Name:   file.Name(),
						Folder: path,
					})
				}
			}
		}
	}

	return servers
}

func buildTmuxId(server *Server) string {
	return fmt.Sprintf("smcd-%s", server.Name)
}

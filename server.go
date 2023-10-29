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
	Name   string
	Folder string
	Pid    int
}

func (server *Server) Start() {
	log.Printf("Starting server: %s", server.Name)
	//TODO check if server is already running under screen identifier
	cmd := exec.Command("/usr/bin/screen", "-dmS", buildScreenId(server), server.startScript())

	err := cmd.Run()
	server.Pid = cmd.Process.Pid

	if err != nil {
		log.Fatalf("Could not start server: %s\n%v", server.Name, err)
	}

	log.Println("Server started.")
}

func (server *Server) Stop() {
	log.Printf("Stopping server: %s", server.Name)
	cmd := exec.Command("/usr/bin/screen", "-S", buildScreenId(server), "-X", "stuff", "stop^M")

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Could not stop server: %s\n%+v", server.Name, err)
	}

	log.Println("Server stopped.")
}

func (server *Server) startScript() string {
	return server.Folder + "/start.sh"
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

func buildScreenId(server *Server) string {
	return fmt.Sprintf("smcd-%s", server.Name)
}

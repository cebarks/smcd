package smcd_test

import (
	"log"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cebarks/smcd"
	"github.com/stretchr/testify/assert"
)

func init() {
	var err error
	smcd.WorkingDir, err = filepath.Abs("test/")
	if err != nil {
		log.Default().Fatalf("error setting up test dir: %v", err)
	}
}

func Test_DiscoverServers(t *testing.T) {
	servers := smcd.DiscoverServers()

	assert.NotEmpty(t, servers, "no servers discovered")
	assert.Equal(t, 2, len(servers))
	assert.Equal(t, "server1", servers[0].Name)
	assert.Equal(t, "server3", servers[1].Name)
	assert.True(t, strings.HasSuffix(servers[0].Folder, "test/server1"))
	assert.True(t, strings.HasSuffix(servers[1].Folder, "test/server3"))
}

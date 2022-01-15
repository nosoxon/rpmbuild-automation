package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"path"
	"reposerver/internal/webserver"
	"reposerver/pkg/cli"
)

const (
	tagComponent = "repo-manager"

	envListenAddress = "LISTEN_ADDRESS"
	envRepoRoot      = "REPO_ROOT"
	envAuthSecret    = "AUTH_SECRET"
)

var (
	ListenAddress string
	RepoRoot      string
	AuthSecret    string
)

func init() {
	cli.Setup(tagComponent)
	ListenAddress = cli.GetEnv(envListenAddress, ":8080")
	RepoRoot      = cli.GetEnv(envRepoRoot,      "/repo")
	AuthSecret    = cli.GetEnv(envAuthSecret,    "TESTING")

	if !cli.IsEnvSet("DEBUG") {
		gin.SetMode(gin.ReleaseMode)
	}

	rpmRoot := path.Join(RepoRoot, "RPMS")
	if _, err := os.Stat(rpmRoot); os.IsNotExist(err) {
		_ = os.MkdirAll(rpmRoot, 0755)
	}
}

func main() {
	log.Info().Str("address", ListenAddress).Str("repo_root", RepoRoot).Msg("starting webserver")
	router := webserver.Router(RepoRoot, AuthSecret)
	router.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	if err := router.Run(ListenAddress); err != nil {
		log.Fatal().Err(err).Msg("failed to start webserver")
	}
}

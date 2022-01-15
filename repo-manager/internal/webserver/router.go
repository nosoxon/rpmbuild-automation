package webserver

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"path"
	"reposerver/internal/createrepo"
)

const headerAuthorization = "X-AUTH"

type repoManager struct {
	repoRoot   string
	authSecret string
}

func (rm *repoManager) postAddPkgs(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		log.Error().Err(err).Msg("could not parse form")
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	if c.GetHeader(headerAuthorization) != rm.authSecret {
		log.Info().Msg("unauthorized addpkgs request")
		c.String(http.StatusUnauthorized, "bad")
		return
	}

	files := form.File["pkg"]
	for _, file := range files {
		dst := path.Join(rm.repoRoot, "RPMS", file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			log.Error().Err(err).Str("file", file.Filename).Msg("could not save to disk")
		} else {
			log.Info().Str("path", dst).Msg("saved to disk")
		}
	}

	result, err := createrepo.Run("--verbose", "--update", rm.repoRoot)
	if err != nil {
		log.Error().Err(err).Msg("could not run createrepo")
		c.JSON(http.StatusInternalServerError, gin.H{"result": result, "err": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func Router(repoRoot, authSecret string) *gin.Engine {
	rm := &repoManager{repoRoot: repoRoot, authSecret: authSecret}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.MaxMultipartMemory = 1 << 28 // 256 MiB

	router.Static("/repo", repoRoot)
	router.POST("/addpkgs", rm.postAddPkgs)

	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404")
	})

	return router
}

// Copyright 2019 Axetroy. All rights reserved. MIT license.
package resource

import (
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/src/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

func File(context *gin.Context) {
	filename := context.Param("filename")
	filePath := path.Join(config.Upload.Path, config.Upload.File.Path, filename)
	if isExistFile := fs.PathExists(filePath); isExistFile == false {
		// if the path not found
		http.NotFound(context.Writer, context.Request)
		return
	}
	http.ServeFile(context.Writer, context.Request, filePath)
}

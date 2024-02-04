package internal

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vuon9/d2m/pkg/api/liquipedia"
)

func NewWebAPI() *gin.Engine {
	app := gin.Default()

	v1 := app.Group("/v1")
	v1.GET("/matches", func(ctx *gin.Context) {
		scheduledMatches, err := liquipedia.NewClient().GetScheduledMatches(ctx)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}

		ctx.JSON(http.StatusOK, scheduledMatches)
	})

	// TODO: Implement this regarding this URL https://liquipedia.net/dota2/Portal:Teams
	// v1.GET("/teams", func(ctx *gin.Context) {
	// })

	v1.GET("/teams/:slug", func(ctx *gin.Context) {
		url := "https://liquipedia.net/dota2/" + ctx.Param("slug")
		fmt.Println(url)
		match, err := liquipedia.NewClient().GetTeamDetailsPage(ctx, url)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}

		ctx.JSON(http.StatusOK, match)
	})

	return app
}

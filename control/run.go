package control

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"meme/middleware"
	"meme/service"
)

func Run(resources embed.FS) {
	//get PORT
	//if not, 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//Template Path
	t := template.Must(template.New("").ParseFS(resources, "views/templates/*"))

	//Gin Instance
	r := gin.Default()

	//set middleware
	r.Use(middleware.RecordUaAndTime)

	//set Template Path
	r.SetHTMLTemplate(t)

	//set Static File Path
	r.StaticFS("/public", http.FS(resources))

	//urls
	//https://patorjk.com/software/taag/#p=display&h=2&f=NScript&t=index%20room
	index := r.Group("")
	{
		index.GET("", Index)

		cmd := index.Group("cmd")
		{
			cmd.POST("exec", Cmd)

			cmd.POST("ssh", Com)

			sql := cmd.Group("sql")
			{
				sql.POST("crud", CRUD)
			}
		}

		meme := index.Group("meme")
		{
			gvg := meme.Group("gvg")
			{
				gvg.GET("", Gvg)

				gvg.POST("stream", GvgStream)

				gvg.GET("ws", GvgHandleConn)
			}
		}

		mst := index.Group(os.Getenv("Master"))
		{
			mst.GET("", Key)
			mst.POST("", Cmd)
		}
	}

	r.GET("sqltest", SqlPage)

	r.GET("reconn", ReConn)
	//NOTE: The machines for [app] have services with 'auto_stop_machines = true' that will be stopped when idling
	///////////////////////////////////////////////////
	memeEngine := r.Group("/meme_")
	{
		gvg := memeEngine.Group("/gvg_")
		{
			gvg.GET("/", func(ctx *gin.Context) {
				ctx.HTML(http.StatusOK, "gvg.html.tmpl", gin.H{})
				service.SetCurrentSub()

				go service.Gvg()

				service.Buffer <- service.GetBuffer()
				<-service.ReqFlg
			})

			gvg.POST("/test", func(ctx *gin.Context) {
				fmt.Println("test")

				service.Buffer <- service.GetBuffer()
				<-service.ReqFlg

				b := make([]byte, 4)
				b[0] = 0
				b[1] = 0
				b[2] = 72
				b[3] = 31
				service.Buffer <- b
				<-service.ReqFlg
			})
		}
	}
	///////////////////////////////////////////////////

	//start
	r.Run(":" + port)
}

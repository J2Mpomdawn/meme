package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"meme/middleware"
	"meme/model"
	"meme/service"
)

func main() {
	/////////////////////////////////////////////////////////////
	/*
	 * json -> golang-struct test
	 */
	var d1 model.Hoge

	dec := json.NewDecoder(bytes.NewBufferString(model.Jsonstr))
	err := dec.Decode(&d1)
	fmt.Println(model.Jsonstr)
	fmt.Printf("%+v %+v \n", d1, err)

	fmt.Println(d1.Data.Castles[0].CastleId, d1.Data.Castles[0].GuildId)
	/////////////////////////////////////////////////////////////

	/////////////////////////////////////////////////////////////
	if os.Getenv("PORT") == "" {
		if err := godotenv.Load("dev.env"); err != nil {
			fmt.Println(err)
		}
	}

	engine := gin.Default()
	engine.Use(middleware.RecordUaAndTime)
	engine.LoadHTMLGlob("./views/*.html")
	engine.Static("/scripts", "./views/scripts")

	memeEngine := engine.Group("/meme")
	{
		gvg := memeEngine.Group("/gvg")
		{
			gvg.GET("/", func(ctx *gin.Context) {
				ctx.HTML(http.StatusOK, "main.html", gin.H{})
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

	engine.Run(":" + os.Getenv("PORT"))
	/////////////////////////////////////////////////////////////
}

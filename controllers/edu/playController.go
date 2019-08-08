package edu

import (
	"edu_api/controllers"
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/models"
	"log"
)

/**
视频播放控制器
 */
type PlayController struct {
	controller controllers.Controller
}

func (play *PlayController) GetPlayList(w rest.ResponseWriter, r *rest.Request)  {
	var playLists []models.Chapter

	playLists, play.controller.Err = play.controller.BaseOrm.GetPlayList(r)

	log.Printf("the playlist is:%v", playLists)
}
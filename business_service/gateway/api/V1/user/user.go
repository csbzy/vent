package user

import "github.com/kataras/iris"

func SetupUserApi(){
	userParty := iris.Party("/api/v1/user")
	userParty.Post("/register",Register)
	userParty.Put("/login",Login)
	userParty.Put("/info",InfoModify)
	userParty.Get("/info/:userID",InfoGet)
}

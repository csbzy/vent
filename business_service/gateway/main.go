package main

import(

)
import "github.com/kataras/iris"

func main(){
	registerApi()
}

func registerApi(){
	iris.API("/api/users/register",RegisterApi{})
}

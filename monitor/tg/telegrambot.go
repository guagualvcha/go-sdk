package tg

import (
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v1"
)

const tgAddr  =  "https://api.telegram.org/"

type TgBot struct {
	Key string
	ChatId string
}

func NewRoRob(key,chatId string) *TgBot{
	return &TgBot{
		Key:key,
		ChatId:chatId,
	}
}

func (t *TgBot) SentMessage(message string){
	params:=map[string]string{"chat_id":t.ChatId,"parse_mode":"html","text":message}

	res,err:=resty.R().SetQueryParams(params).Get(tgAddr + t.Key+"/sendMessage")
	if err!=nil{
		fmt.Println(err)
	}else{
		bz,_:=json.Marshal(res)
		fmt.Println(bz)
	}
}
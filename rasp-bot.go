package main

import (
	"encoding/json"
	"os"
)

type TableRow []string

type ITimeTable interface {
	GroupTimeTable(string) []TableRow            //возвращает полное расписание
	GroupRankTimeTable(string) []TableRow        //зачеты
	GroupExamTimeTable(string) []TableRow        //экзамены
	GroupWeekTimeTable(string) []TableRow        //расписание на текущую неделю с учетом четности
	GroupDayTimeTable(string, string) []TableRow //на день недели
	GroupNearestPair(string) []TableRow          //ближайшая пара
	PrTimeTable(string) []TableRow               //расписание преподавателя
	PrWeekTimeTable(string) []TableRow           //на текущую неделю
	PrRankTimeTable(string) []TableRow           //зачеты
	PrExamTimeTable(string) []TableRow           //экзамены
	PrDayTimeTable(string, string) []TableRow    //на день недели
}

type Config struct {
	Token string `json:"token"`
}

var config Config

var DayOfWeekString = [...]string{"", "Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}

func main() {
	var tt MEPHI_TimeTable
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		panic(err)
	}
	InitializeBotServer(&tt)
}

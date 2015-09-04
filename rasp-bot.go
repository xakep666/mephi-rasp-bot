package main

import (
	"fmt"
)

type TableRow []string

type ITimeTable interface {
	GroupTimeTable(string) ([]TableRow, error)      //возвращает полное расписание
	GroupRankTimeTable(string) ([]TableRow, error)  //зачеты
	GroupExamTimeTable(string) ([]TableRow, error)  //экзамены
	GroupWeekTimeTable(string) ([]TableRow, error)  //расписание на текущую неделю с учетом четности
	GroupTodayTimeTable(string) ([]TableRow, error) //сегодня
	GroupNearestPair(string) ([]TableRow, error)    //ближайшая пара
	PrTimeTable(string) ([]TableRow, error)         //расписание преподавателя
	PrWeekTimeTable(string) ([]TableRow, error)     //на текущую неделю
	PrRankTimeTable(string) ([]TableRow, error)     //зачеты
	PrExamTimeTable(string) ([]TableRow, error)     //экзамены
	PrTodayTimeTable(string) ([]TableRow, error)    //сегодня
}

var DayOfWeekString = [...]string{"", "пн", "вт", "ср", "чт", "пт", "сб", "вс"}

func main() {
	gname := "К03-411"
	var tt MEPHI_TimeTable
	arr, _ := tt.GroupTodayTimeTable(gname)
	for _, rows := range arr {
		for _, cell := range rows {
			fmt.Printf("%s ", cell)
		}
		fmt.Println()
	}
}

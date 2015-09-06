package main

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

var DayOfWeekString = [...]string{"", "Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
var DayOfWeekNum map[string]int

func main() {
	DayOfWeekNum = make(map[string]int)
	DayOfWeekNum["Пн"] = 1
	DayOfWeekNum["Вт"] = 2
	DayOfWeekNum["Ср"] = 3
	DayOfWeekNum["Чт"] = 4
	DayOfWeekNum["Пт"] = 5
	DayOfWeekNum["Сб"] = 6
	DayOfWeekNum["Вс"] = 7
	var tt MEPHI_TimeTable
	InitializeBotServer(&tt)
}

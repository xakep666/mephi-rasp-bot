package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"rasp-bot/Godeps/_workspace/src/github.com/yhat/scrape"
	"rasp-bot/Godeps/_workspace/src/golang.org/x/net/html"
	"rasp-bot/Godeps/_workspace/src/golang.org/x/net/html/atom"
)

/* Структура таблицы:
Дн 	Время 	Н/Ч 	Дисциплина 	Тип зан. 	Группа(ы) 	Преподаватель 	Ауд. 	КВ** 	Прим.
*/

type MEPHI_TimeTable struct {
	root *html.Node
}

func (tt *MEPHI_TimeTable) rootGetter(tt_selector, name, typ_selector string) (err error) {
	resp, err := http.Get("https://eisgateway.mephi.ru/TimeTable/timetableshow.aspx?" +
		tt_selector + "=" + name + "&typ=" + typ_selector)
	if err != nil {
		log.Printf("Error on getting: %s\n", err.Error())
		return
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		log.Printf("Error on parsing: %s\n", err.Error())
		return
	}
	tt.root = root
	return
}

func (tt *MEPHI_TimeTable) GroupTimeTable(gname string) (trs []TableRow) {
	err := tt.rootGetter("gr", gname, "gr")
	if err != nil {
		log.Printf("Error on parsing: %s\n", err.Error())
		return
	}
	i := 0
	for _, row := range scrape.FindAll(tt.root, scrape.ByTag(atom.Tr)) {
		trs = append(trs, make(TableRow, 1))
		for _, cell := range scrape.FindAll(row, scrape.ByTag(atom.Td)) {
			trs[i] = append(trs[i], scrape.Text(cell))
		}
		i++
	}
	err = nil
	return
}

func (tt *MEPHI_TimeTable) GroupRankTimeTable(gname string) (trs []TableRow) {
	err := tt.rootGetter("gr", gname, "grZ")
	if err != nil {
		log.Printf("Error on parsing: %s\n", err.Error())
		return
	}
	i := 0
	for _, row := range scrape.FindAll(tt.root, scrape.ByTag(atom.Tr)) {
		trs = append(trs, make(TableRow, 1))
		for _, cell := range scrape.FindAll(row, scrape.ByTag(atom.Td)) {
			trs[i] = append(trs[i], scrape.Text(cell))
		}
		i++
	}
	err = nil
	return
}

func (tt *MEPHI_TimeTable) GroupExamTimeTable(gname string) (trs []TableRow) {
	err := tt.rootGetter("gr", gname, "grE")
	if err != nil {
		log.Printf("Error on parsing: %s\n", err.Error())
		return
	}
	i := 0
	for _, row := range scrape.FindAll(tt.root, scrape.ByTag(atom.Tr)) {
		trs = append(trs, make(TableRow, 1))
		for _, cell := range scrape.FindAll(row, scrape.ByTag(atom.Td)) {
			trs[i] = append(trs[i], scrape.Text(cell))
		}
		i++
	}
	err = nil
	return
}

func (tt *MEPHI_TimeTable) GroupWeekTimeTable(gname string) (trs []TableRow) {
	err := tt.rootGetter("gr", gname, "gr")
	if err != nil {
		log.Printf("Error on parsing: %s\n", err.Error())
		return
	}
	for _, row := range scrape.FindAll(tt.root, scrape.ByTag(atom.Tr)) {
		rowToAdd := make(TableRow, 1)
		//проверка на четность
		cells := scrape.FindAll(row, scrape.ByTag(atom.Td))
		//ячеек может быть меньше 3х
		if len(cells) >= 3 {
			oddeven := scrape.Text(cells[2])
			dow := scrape.Text(cells[0])
			if oddeven == "/Ч" && !tt.isDayOdd(DayOfWeekNum[dow]) {
				continue
			}
			if oddeven == "Н/" && tt.isDayOdd(DayOfWeekNum[dow]) {
				continue
			}
		}
		for _, cell := range cells {
			//не добавлять н/ч
			txt := scrape.Text(cell)
			if txt == "Н/" || txt == "/Ч" {
				continue
			}
			rowToAdd = append(rowToAdd, txt)
		}
		if len(rowToAdd) >= 3 {
			trs = append(trs, rowToAdd)
		}
	}
	err = nil
	return
}

func (tt *MEPHI_TimeTable) PrTimeTable(pname string) (trs []TableRow) {
	err := tt.rootGetter("prep", pname, "prep")
	if err != nil {
		log.Printf("Error on parsing: %s\n", err.Error())
		return
	}
	i := 0
	for _, row := range scrape.FindAll(tt.root, scrape.ByTag(atom.Tr)) {
		trs = append(trs, make(TableRow, 1))
		for _, cell := range scrape.FindAll(row, scrape.ByTag(atom.Td)) {
			trs[i] = append(trs[i], scrape.Text(cell))
		}
		i++
	}
	err = nil
	return
}

func (tt *MEPHI_TimeTable) PrWeekTimeTable(pname string) (trs []TableRow) {
	err := tt.rootGetter("prep", pname, "prep")
	if err != nil {
		log.Printf("Error on parsing: %s\n", err.Error())
		return
	}
	for _, row := range scrape.FindAll(tt.root, scrape.ByTag(atom.Tr)) {
		rowToAdd := make(TableRow, 1)
		//проверка на четность
		cells := scrape.FindAll(row, scrape.ByTag(atom.Td))
		//ячеек может быть меньше 3х
		if len(cells) >= 3 {
			dow := scrape.Text(cells[0])
			oddeven := scrape.Text(cells[2])
			if oddeven == "/Ч" && !tt.isDayOdd(DayOfWeekNum[dow]) {
				continue
			}
			if oddeven == "Н/" && tt.isDayOdd(DayOfWeekNum[dow]) {
				continue
			}
		}
		for _, cell := range cells {
			//не добавлять н/ч
			txt := scrape.Text(cell)
			if txt == "Н/" || txt == "/Ч" {
				continue
			}
			rowToAdd = append(rowToAdd, txt)
		}
		if len(rowToAdd) >= 3 {
			trs = append(trs, rowToAdd)
		}
	}
	err = nil
	return
}

func (tt *MEPHI_TimeTable) PrRankTimeTable(pname string) (trs []TableRow) {
	err := tt.rootGetter("prep", pname, "prepZ")
	if err != nil {
		log.Printf("Error on parsing: %s\n", err.Error())
		return
	}
	i := 0
	for _, row := range scrape.FindAll(tt.root, scrape.ByTag(atom.Tr)) {
		trs = append(trs, make(TableRow, 1))
		for _, cell := range scrape.FindAll(row, scrape.ByTag(atom.Td)) {
			trs[i] = append(trs[i], scrape.Text(cell))
		}
		i++
	}
	err = nil
	return
}

func (tt *MEPHI_TimeTable) PrExamTimeTable(pname string) (trs []TableRow) {
	err := tt.rootGetter("prep", pname, "prepE")
	if err != nil {
		log.Printf("Error on parsing: %s\n", err.Error())
		return
	}
	i := 0
	for _, row := range scrape.FindAll(tt.root, scrape.ByTag(atom.Tr)) {
		trs = append(trs, make(TableRow, 1))
		for _, cell := range scrape.FindAll(row, scrape.ByTag(atom.Td)) {
			trs[i] = append(trs[i], scrape.Text(cell))
		}
		i++
	}
	err = nil
	return
}

func (tt *MEPHI_TimeTable) GroupNearestPair(gname string) (trs []TableRow) {
	err := tt.rootGetter("gr", gname, "gr")
	if err != nil {
		log.Printf("Error on parsing: %s\n", err.Error())
		return
	}
	trs = make([]TableRow, 1)
	now := time.Now()
	dow := int(now.Weekday())
	//в nrow будет нужный номер строки
	nrow := -1
	rows := scrape.FindAll(tt.root, scrape.ByTag(atom.Tr))
	for i, row := range rows {
		cells := scrape.FindAll(row, scrape.ByTag(atom.Td))
		if len(cells) < 3 {
			continue
		}

		//если день недели не тот, пропускаем
		if scrape.Text(cells[0]) != DayOfWeekString[dow] {
			continue
		}
		//формат времени в ячейке [1] "00:00 - 00:00" (24ч)
		hstart, _ := strconv.Atoi(scrape.Text(cells[1])[0:2])
		mstart, _ := strconv.Atoi(scrape.Text(cells[1])[3:5])
		hend, _ := strconv.Atoi(scrape.Text(cells[1])[8:10])
		mend, _ := strconv.Atoi(scrape.Text(cells[1])[11:13])
		startTime := time.Date(now.Year(), now.Month(), now.Day(), hstart, mstart, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), hend, mend, 0, 0, now.Location())
		//если время "сейчас" между "концом" и "началом" - берем эту строку
		if now.After(startTime) && now.Before(endTime) {
			nrow = i
			break
		}
		//если время "сейчас" после конца - не берем
		if now.After(endTime) {
			continue
		}
		//если время "сейчас" перед началом - берем
		if now.Before(startTime) {
			nrow = i
			break
		}
	}
	//если ничего не нашли, вернем "Ближайших пар на сегодня не найдено"
	if nrow == -1 {
		trs[0] = append(trs[0], "Ближайших пар на сегодня не найдено")
		return
	} else {
		//вернем нужную строку
		for _, cell := range scrape.FindAll(rows[nrow], scrape.ByTag(atom.Td)) {
			trs[0] = append(trs[0], scrape.Text(cell))
		}
		return
	}
}

func (tt *MEPHI_TimeTable) GroupDayTimeTable(gname string, dow string) (trs []TableRow) {
	var tr TableRow
	err := tt.rootGetter("gr", gname, "gr")
	if err != nil {
		log.Printf("Error on parsing: %s\n", err.Error())
		return
	}
	for _, rows := range scrape.FindAll(tt.root, scrape.ByTag(atom.Tr)) {
		cells := scrape.FindAll(rows, scrape.ByTag(atom.Td))
		if len(cells) < 3 {
			continue
		}
		if dow == "Сегодня" {
			dow = DayOfWeekString[int(time.Now().Weekday())]
		}
		if scrape.Text(cells[0]) == dow {
			if len(cells) >= 3 {
				oddeven := scrape.Text(cells[2])
				if oddeven == "/Ч" && !tt.isDayOdd(DayOfWeekNum[dow]) {
					continue
				}
				if oddeven == "Н/" && tt.isDayOdd(DayOfWeekNum[dow]) {
					continue
				}
			}
			for _, cell := range cells {
				txt := scrape.Text(cell)
				if txt == "Н/" || txt == "/Ч" {
					continue
				}
				tr = append(tr, scrape.Text(cell))
			}
		}
		trs = append(trs, tr)
		tr = nil
	}
	return
}

func (tt *MEPHI_TimeTable) PrDayTimeTable(pname string, dow string) (trs []TableRow) {
	var tr TableRow
	err := tt.rootGetter("prep", pname, "prep")
	if err != nil {
		log.Printf("Error on parsing: %s\n", err.Error())
		return
	}
	for _, rows := range scrape.FindAll(tt.root, scrape.ByTag(atom.Tr)) {
		cells := scrape.FindAll(rows, scrape.ByTag(atom.Td))
		if len(cells) < 3 {
			continue
		}
		if dow == "Сегодня" {
			dow = DayOfWeekString[int(time.Now().Weekday())]
		}
		if scrape.Text(cells[0]) == dow {
			if len(cells) >= 3 {
				oddeven := scrape.Text(cells[2])
				if oddeven == "/Ч" && !tt.isDayOdd(DayOfWeekNum[dow]) {
					continue
				}
				if oddeven == "Н/" && tt.isDayOdd(DayOfWeekNum[dow]) {
					continue
				}
			}
			for _, cell := range cells {
				txt := scrape.Text(cell)
				if txt == "Н/" || txt == "/Ч" {
					continue
				}
				tr = append(tr, scrape.Text(cell))
			}
		}
		trs = append(trs, tr)
		tr = nil
	}
	return
}

func (tt MEPHI_TimeTable) isDayOdd(dayInWeekNum int) bool { //день текущей недели
	var delta, s_delta int
	today := time.Now()
	weekBegin := int(today.Day()) - int(today.Weekday()) + 1 //день начала недели
	if weekBegin < 0 {                                       //если раньше 1го числа месяца
		weekBegin *= -1
	}
	day := time.Date(today.Year(), today.Month(), weekBegin+dayInWeekNum, 0, 0, 0, 0, today.Location())
	s1begin := time.Date(today.Year(), 9, 1, 0, 0, 0, 0, today.Location()) //начало 1 семестра
	s2begin := time.Date(today.Year(), 2, 8, 0, 0, 0, 0, today.Location()) //начало 2 семестра
	if day.After(s2begin) {
		delta = int(day.YearDay()) - int(s2begin.YearDay())
		_, d1 := day.ISOWeek()
		_, d2 := s2begin.ISOWeek()
		s_delta = d1 - d2 //сколько воскресений убирать
	} else {
		delta = day.YearDay() - s1begin.YearDay()
		_, d1 := day.ISOWeek()
		_, d2 := s1begin.ISOWeek()
		s_delta = d1 - d2
	}
	return (delta-s_delta)%2 == 0
}

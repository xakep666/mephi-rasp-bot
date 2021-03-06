package main

import (
	"log"
	"os"
	"rasp-bot/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
)

//Разметка клавиатуры выбора дней
var group_kbd_layout = [][]string{
	[]string{"Ближайшая"},
	[]string{"Сегодня", "Пн", "Вт", "Ср"},
	[]string{"Чт", "Пт", "Сб", "Вс"},
	[]string{"Неделя", "Полное", "Зачет", "Экзамен"},
}

var prep_kbd_layout = [][]string{
	[]string{"Сегодня", "Пн", "Вт", "Ср"},
	[]string{"Чт", "Пт", "Сб", "Вс", "Неделя"},
	[]string{"Полное", "Зачет", "Экзамен"},
}

//Подсказка
var help_msg = "Расписание МИФИ\n" +
	"Комманды:\n" +
	"/help - Это сообщение\n" +
	"/group - Расписание группы\n" +
	"/prep - Расписание преподавателя\n" +
	"/cancel - Начать заново"

//Возвращает объект разметки клавиатуры для ответа
func KbdLayout(sel bool) (markup tgbotapi.ReplyKeyboardMarkup) {
	if sel {
		markup.Keyboard = group_kbd_layout
	} else {
		markup.Keyboard = prep_kbd_layout
	}
	markup.OneTimeKeyboard = true
	markup.ResizeKeyboard = true
	return
}

//база для обработки
const (
	NOSEL = iota
	GROUP
	PREP
)

type phases struct {
	groupPrep int
	name      string
}

var process_base map[int]phases //по id чата

//склейка массива строк в одну для отправки
func StringJoiner(rows []TableRow) (newstr string) {
	for _, row := range rows {
		for _, cell := range row {
			newstr += cell //не добавлять для пустых ячеек
			if cell != "" {
				newstr += "|"
			}
		}
		if len(row) != 0 { //не добавлять для пустых строк
			newstr += "\n-----------------\n"
		}
	}
	newstr = "РАСПИСАНИЕ\n" + newstr //telegram не отправляет пустые сообщения (без букв)
	return
}

func HandleRequest(chat_id int, text string, tt ITimeTable) (cfg tgbotapi.MessageConfig) {
	log.Printf("Processing command %s to chat %d", text, chat_id)
	cfg.ReplyMarkup = tgbotapi.ReplyKeyboardHide{true, true} //если по ошибке осталась кастомная клавиатура
	switch text {
	//комманды
	case "/start", "/help":
		{
			cfg.Text = help_msg
		}
	case "/group":
		{
			cfg.Text = "Введите имя группы"
			process_base[chat_id] = phases{GROUP, ""}
		}
	case "/prep":
		{
			cfg.Text = "Введите фамилию преподавателя"
			process_base[chat_id] = phases{PREP, ""}
		}
	case "/cancel":
		{
			cfg.Text = "Выбор отменен"
			process_base[chat_id] = phases{NOSEL, ""}
		}
	//дни недели
	case "Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс", "Сегодня":
		{
			//только если введена первая комманда и имя
			record := process_base[chat_id]
			if record.groupPrep != NOSEL && record.name != "" {
				switch record.groupPrep {
				case GROUP:
					cfg.Text = StringJoiner(tt.GroupDayTimeTable(record.name, text))
				case PREP:
					cfg.Text = StringJoiner(tt.PrDayTimeTable(record.name, text))
				}
			}

		}
	case "Ближайшая":
		{
			//только если выбрана группа
			record := process_base[chat_id]
			if record.groupPrep == GROUP && record.name != "" {
				cfg.Text = StringJoiner(tt.GroupNearestPair(record.name))
			}
		}
	case "Полное":
		{
			record := process_base[chat_id]
			if record.groupPrep != NOSEL && record.name != "" {
				switch record.groupPrep {
				case GROUP:
					cfg.Text = StringJoiner(tt.GroupTimeTable(record.name))
				case PREP:
					cfg.Text = StringJoiner(tt.PrTimeTable(record.name))
				}
			}
		}
	case "Неделя":
		{
			record := process_base[chat_id]
			if record.groupPrep != NOSEL && record.name != "" {
				switch record.groupPrep {
				case GROUP:
					cfg.Text = StringJoiner(tt.GroupWeekTimeTable(record.name))
				case PREP:
					cfg.Text = StringJoiner(tt.PrWeekTimeTable(record.name))
				}
			}
		}
	case "Зачет":
		{
			record := process_base[chat_id]
			if record.groupPrep != NOSEL && record.name != "" {
				switch record.groupPrep {
				case GROUP:
					cfg.Text = StringJoiner(tt.GroupRankTimeTable(record.name))
				case PREP:
					cfg.Text = StringJoiner(tt.PrRankTimeTable(record.name))
				}
			}
		}
	case "Экзамен":
		{
			record := process_base[chat_id]
			if record.groupPrep != NOSEL && record.name != "" {
				switch record.groupPrep {
				case GROUP:
					cfg.Text = StringJoiner(tt.GroupExamTimeTable(record.name))
				case PREP:
					cfg.Text = StringJoiner(tt.PrExamTimeTable(record.name))
				}
			}
		}
	default:
		{
			record := process_base[chat_id]
			switch record.groupPrep {
			case GROUP:
				{
					cfg.Text = "Выберите вариант"
					record.name = text
					cfg.ReplyMarkup = KbdLayout(true)
				}
			case PREP:
				{
					cfg.Text = "Выберите вариант"
					record.name = text
					cfg.ReplyMarkup = KbdLayout(false)
				}
			}
			process_base[chat_id] = record
		}
	}
	cfg.ChatID = chat_id
	return
}

func InitializeBotServer(tt ITimeTable) {
	process_base = make(map[int]phases)
	bot, err := tgbotapi.NewBotAPI(os.Getenv("ACESS_TOKEN"))
	if err != nil {
		log.Panicln("Cannot initialize bot api, " + err.Error())
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	log.Printf("Authorized on account %s\n", bot.Self.UserName)
	err = bot.UpdatesChan(u)
	if err != nil {
		log.Println("Cannot set updates channel, " + err.Error())
	}
	for update := range bot.Updates {
		log.Printf("Message from %s: %s\n", update.Message.From.UserName, update.Message.Text)
		bot.SendMessage(HandleRequest(update.Message.Chat.ID, update.Message.Text, tt))
	}
}

package telegram

import (
	"database/sql"
	query "doMassageBot/internal/db"
	"doMassageBot/internal/utils"
	"fmt"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"strconv"
	"time"
)

const commandStart = "start"

var signMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Зарегистрироваться"),
	))

var mainMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("📝 Записаться"),
		tgbotapi.NewKeyboardButton("🗓 Моё расписание"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("❌ Отменить запись"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("👤 Мой профиль"),
	),
)
var ToBeginning = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🔄 В начало"),
	))

var TimeMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Сегодня"),
		tgbotapi.NewKeyboardButton("Завтра")),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🔄 В начало"),
	),
)

var TimeMenuForFriday = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Сегодня"),
		tgbotapi.NewKeyboardButton("Понедельник")),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🔄 В начало"),
	),
)

var massageMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("💆 Шейно воротниковый массаж"),
		tgbotapi.NewKeyboardButton("🧖 Лечебный массаж")),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🔄 В начало"),
	),
)

var adminMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🗓 Получить расписание")),
)

var adminMenuTime = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("На сегодня"),
		tgbotapi.NewKeyboardButton("На завтра")),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("В начало 🔄"),
	),
)
var adminMenuTimeForFriday = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("На сегодня"),
		tgbotapi.NewKeyboardButton("Понедельник")),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("В начало 🔄"),
	),
)

func (b *Bot) adminHandler(message *tgbotapi.Message, db *sql.DB) error {
	if message.IsCommand() {
		return b.adminHandleCommand(message)
	} else {
		return b.adminHandleText(message, db)
	}
}

func (b *Bot) adminHandleCommand(message *tgbotapi.Message) error {
	cmdText := message.Command()
	if cmdText == "start" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Привет админ!")
		msg.ReplyMarkup = adminMenu
		_, err := b.bot.Send(msg)
		return err
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды :(")
		msg.ReplyMarkup = adminMenu
		_, err := b.bot.Send(msg)
		return err
	}
}

func (b *Bot) adminHandleAdminMenuGetSchedule(message *tgbotapi.Message, db *sql.DB) (err error) {
	switch time.Now().Weekday() {
	case time.Friday:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите день.")
		msg.ReplyMarkup = adminMenuTimeForFriday
		_, err = b.bot.Send(msg)
	case time.Saturday, time.Sunday:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Записи понедельника")
		msg.ReplyMarkup = ToBeginning
		_, err = b.bot.Send(msg)
		objs, err := query.GetAllScheduleForTomorrow(db, "Понедельник")
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
			_, err := b.bot.Send(msg)
			return err
		} else {
			if len(objs) > 0 {
				var text string
				for _, obj := range objs {
					text += fmt.Sprintf("*%s*\n"+"*Дата: * _%v_\n"+"*Время: * _%s_\n"+"*Имя: * _%v_\n"+"*Email: * _%v_\n"+"*Номер телефона: * _%v_\n"+"------------------------------------------------------\n", obj.MType, obj.MDate, obj.MTime, obj.Name, obj.Email, obj.PhoneNum)
				}
				msg := tgbotapi.NewMessage(message.Chat.ID, text)
				msg.ParseMode = "markdown"
				_, err := b.bot.Send(msg)
				return err
			} else {
				msg := tgbotapi.NewMessage(message.Chat.ID, "Еще никто не записался...")
				msg.ReplyMarkup = adminMenu
				_, err := b.bot.Send(msg)
				return err
			}

		}
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите день.")
		msg.ReplyMarkup = adminMenuTime
		_, err = b.bot.Send(msg)

	}
	return
}

func (b *Bot) adminHandleAdminMenuTimeToday(message *tgbotapi.Message, db *sql.DB) error {
	objs, err := query.GetAllScheduleForToday(db)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
		_, err := b.bot.Send(msg)
		return err
	} else {
		if len(objs) > 0 {
			var text string
			for _, obj := range objs {
				text += fmt.Sprintf("*%s*\n"+"*Дата: * _%v_\n"+"*Время: * _%s_\n"+"*Имя: * _%v_\n"+"*Email: * _%v_\n"+"*Номер телефона: * _%v_\n"+"------------------------------------------------------\n", obj.MType, obj.MDate, obj.MTime, obj.Name, obj.Email, obj.PhoneNum)
			}
			msg := tgbotapi.NewMessage(message.Chat.ID, text)
			msg.ParseMode = "markdown"
			_, err := b.bot.Send(msg)
			return err
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Еще никто не записался.")
			msg.ReplyMarkup = adminMenu
			_, err := b.bot.Send(msg)
			return err
		}
	}
}

func (b *Bot) adminHandleAdminMenuTimeTomorrow(message *tgbotapi.Message, db *sql.DB) error {
	objs, err := query.GetAllScheduleForTomorrow(db, message.Text)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
		_, err := b.bot.Send(msg)
		return err
	} else {
		if len(objs) > 0 {
			var text string
			for _, obj := range objs {
				text += fmt.Sprintf("*%s*\n"+"*Дата: * _%v_\n"+"*Время: * _%s_\n"+"*Имя: * _%v_\n"+"*Email: * _%v_\n"+"*Номер телефона: * _%v_\n"+"------------------------------------------------------\n", obj.MType, obj.MDate, obj.MTime, obj.Name, obj.Email, obj.PhoneNum)
			}
			msg := tgbotapi.NewMessage(message.Chat.ID, text)
			msg.ParseMode = "markdown"
			_, err := b.bot.Send(msg)
			return err
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Еще никто не записался...")
			msg.ReplyMarkup = adminMenu
			_, err := b.bot.Send(msg)
			return err
		}

	}
}

func (b *Bot) adminHandleAdminMenuTimeMonday(message *tgbotapi.Message, db *sql.DB) error {
	objs, err := query.GetAllScheduleForMonday(db)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
		_, err := b.bot.Send(msg)
		return err
	} else {
		if len(objs) > 0 {
			var text string
			for _, obj := range objs {
				text += fmt.Sprintf("*%s*\n"+"*Дата: * _%v_\n"+"*Время: * _%s_\n"+"*Имя: * _%v_\n"+"*Email: * _%v_\n"+"*Номер телефона: * _%v_\n"+"------------------------------------------------------\n", obj.MType, obj.MDate, obj.MTime, obj.Name, obj.Email, obj.PhoneNum)
			}
			msg := tgbotapi.NewMessage(message.Chat.ID, text)
			msg.ParseMode = "markdown"
			_, err := b.bot.Send(msg)
			return err
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Еще никто не записался.")
			msg.ReplyMarkup = adminMenu
			_, err := b.bot.Send(msg)
			return err
		}
	}
}

func (b *Bot) adminHandleToBeginning(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Начало")
	msg.ReplyMarkup = adminMenu
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) adminHandleText(message *tgbotapi.Message, db *sql.DB) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды :(")
	msg.ReplyMarkup = adminMenu
	switch message.Text {
	case adminMenu.Keyboard[0][0].Text:
		return b.adminHandleAdminMenuGetSchedule(message, db)
	case adminMenuTime.Keyboard[0][0].Text:
		return b.adminHandleAdminMenuTimeToday(message, db)
	case adminMenuTime.Keyboard[0][1].Text:
		return b.adminHandleAdminMenuTimeTomorrow(message, db)
	case adminMenuTime.Keyboard[1][0].Text:
		return b.adminHandleToBeginning(message)
	case adminMenuTimeForFriday.Keyboard[0][0].Text:
		return b.adminHandleAdminMenuTimeToday(message, db)
	case adminMenuTimeForFriday.Keyboard[0][1].Text:
		return b.adminHandleAdminMenuTimeTomorrow(message, db)
	default:
		_, err := b.bot.Send(msg)
		return err
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message, db *sql.DB) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды.")
	switch message.Command() {
	case commandStart:
		isexist, err := query.CheckIfUserExists(db, message.From.ID)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, err.Error())
			msg.ReplyMarkup = mainMenu
			_, err := b.bot.Send(msg)
			return err
		}
		if isexist {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome!")
			msg.ReplyMarkup = mainMenu
			_, err := b.bot.Send(msg)
			return err
		} else {
			b.uId, err = query.InsertIntoUsers(db, message.From.ID, "", "", "", "", 0)
			if err != nil {
				fmt.Println(err.Error())
				msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
				_, err := b.bot.Send(msg)
				return err
			}
			err := query.UpdateUsername(db, message.From.ID, message.From.UserName)
			if err != nil {
				fmt.Println(err.Error())
				msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
				_, err := b.bot.Send(msg)
				return err
			}
			msg := tgbotapi.NewMessage(message.Chat.ID, "Здравствуйте, вы должны сначала зарегистрироваться.")
			msg.ReplyMarkup = signMenu
			_, err = b.bot.Send(msg)
			return nil
		}
	default:
		_, err := b.bot.Send(msg)
		return err
	}
}

func (b *Bot) handleText(message *tgbotapi.Message, db *sql.DB) error {
	switch message.Text {
	case signMenu.Keyboard[0][0].Text:
		return b.handleSignMenu(message, db)
	case mainMenu.Keyboard[0][0].Text:
		return b.handleMainMenuEnroll(message)
	case mainMenu.Keyboard[0][1].Text:
		return b.handleMainMenuMySchedule(message, db)
	case mainMenu.Keyboard[1][0].Text:
		return b.handleMainMenuCancel(message, db)
	case mainMenu.Keyboard[2][0].Text:
		return b.handleMainMenuMyProfile(message, db)
	case TimeMenu.Keyboard[0][0].Text:
		return b.handleTimeMenuToday(message, db)
	case TimeMenu.Keyboard[0][1].Text:
		return b.handleTimeMenuTomorrow(message, db)
	case TimeMenuForFriday.Keyboard[0][0].Text:
		return b.handleTimeMenuToday(message, db)
	case TimeMenuForFriday.Keyboard[0][1].Text:
		return b.handleTimeMenuTomorrow(message, db)
	case TimeMenu.Keyboard[1][0].Text:
		return b.handleToBeginning(message)
	case massageMenu.Keyboard[0][0].Text:
		return b.handleMassageMenuCollar(message, db)
	case massageMenu.Keyboard[0][1].Text:
		return b.handleMassageMenuMedical(message, db)
	case massageMenu.Keyboard[1][0].Text:
		return b.handleToBeginning(message)
	default:
		return b.handleFinishEntry(message, db)
	}

}

func (b *Bot) handleSignMenu(message *tgbotapi.Message, db *sql.DB) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Ваше имя:")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	_, err := b.bot.Send(msg)
	err = query.UpdateUserStatus(db, message.From.ID, 1)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
		_, err := b.bot.Send(msg)
		return err
	}
	return nil
}

func (b *Bot) handleMainMenuEnroll(message *tgbotapi.Message) (err error) {

	switch time.Now().Weekday() {
	case time.Friday:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите день.")
		msg.ReplyMarkup = TimeMenuForFriday
		_, err = b.bot.Send(msg)
	case time.Saturday, time.Sunday:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Запись не работает")
		msg.ReplyMarkup = ToBeginning
		_, err = b.bot.Send(msg)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите день.")
		msg.ReplyMarkup = TimeMenu
		_, err = b.bot.Send(msg)
	}
	return
}
func (b *Bot) handleMainMenuMySchedule(message *tgbotapi.Message, db *sql.DB) error {
	objs, err := query.GetMySchedule(db, message.From.ID)
	fmt.Println(message.Chat.ID, message.From.ID)
	if err != nil {
		fmt.Println(err.Error())
		msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
		_, err := b.bot.Send(msg)
		return err
	} else {
		if len(objs) > 0 {
			fmt.Println("Its len: ", len(objs))
			for _, obj := range objs {
				text := fmt.Sprintf("*%s*\n"+"*Дата: * _%v_\n"+"*Время: * _%s_\n", obj.MType, obj.MDate, obj.MTime)
				msg := tgbotapi.NewMessage(message.Chat.ID, text)
				msg.ParseMode = "markdown"
				b.bot.Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Вы еще не записались...")
			msg.ReplyMarkup = mainMenu
			_, err := b.bot.Send(msg)
			return err
		}

	}
	return nil

}

func (b *Bot) handleMainMenuCancel(message *tgbotapi.Message, db *sql.DB) error {
	objs, err := query.GetMySchedule(db, message.From.ID)
	if err != nil {
		fmt.Println(err.Error())
		msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
		_, err := b.bot.Send(msg)
		return err
	} else {
		if len(objs) > 0 {
			for _, obj := range objs {
				text := fmt.Sprintf("*%s*\n"+"*Дата: * _%v_\n"+"*Время: * _%s_\n ", obj.MType, obj.MDate, obj.MTime)
				keyboard := tgbotapi.InlineKeyboardMarkup{}
				var row []tgbotapi.InlineKeyboardButton
				btn := tgbotapi.NewInlineKeyboardButtonData("Отменить", strconv.Itoa(obj.Id))
				row = append(row, btn)
				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
				msg := tgbotapi.NewMessage(message.Chat.ID, text)
				msg.ReplyMarkup = keyboard
				msg.ParseMode = "markdown"
				b.bot.Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Вы еще не записались...")
			msg.ReplyMarkup = mainMenu
			_, err := b.bot.Send(msg)
			return err
		}

	}
	return nil
}

func (b *Bot) handleMainMenuMyProfile(message *tgbotapi.Message, db *sql.DB) error {
	obj, err := query.GetMyProfile(db, message.From.ID)
	if err != nil {
		fmt.Println(err.Error())
		msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
		_, err := b.bot.Send(msg)
		return err
	} else {
		text := fmt.Sprintf("*Ваше имя: * _%v_\n"+"*Email: * _%s_\n"+"*Номер телефона: * _%v_\n", obj.Name, obj.Email, obj.PhoneNum)
		keyboard := tgbotapi.InlineKeyboardMarkup{}
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData("Редактировать ", "Edit")
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		msg.ReplyMarkup = keyboard
		msg.ParseMode = "markdown"
		_, err := b.bot.Send(msg)
		return err
	}
}

func (b *Bot) handleTimeMenuToday(message *tgbotapi.Message, db *sql.DB) error {
	err := query.RefreshMassageSchedule(db)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
		_, err = b.bot.Send(msg)
		return err
	}
	b.mId, err = query.InsertIntoSchedule(db, "", TimeMenu.Keyboard[0][0].Text, "00:00", message.From.ID, 0)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
		_, err = b.bot.Send(msg)
		return err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, "Какой вид массажа вам нужен?")
	msg.ReplyMarkup = massageMenu
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleTimeMenuTomorrow(message *tgbotapi.Message, db *sql.DB) error {
	err := query.RefreshMassageSchedule(db)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
		_, err = b.bot.Send(msg)
		return err
	}
	b.mId, err = query.InsertIntoSchedule(db, "", message.Text, "00:00", message.From.ID, 0)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
		_, err = b.bot.Send(msg)
		return err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, "Какой вид массажа вам нужен?")
	msg.ReplyMarkup = massageMenu
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleToBeginning(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Начало")
	msg.ReplyMarkup = mainMenu
	_, err := b.bot.Send(msg)
	return err
}
func (b *Bot) handleMassageMenuCollar(message *tgbotapi.Message, db *sql.DB) error {
	status, err := query.GetScheduleStatus(db, b.mId)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
		_, err := b.bot.Send(msg)
		return err
	} else {
		switch status {
		case 0:
			err := query.UpdateScheduleMType(db, massageMenu.Keyboard[0][0].Text, b.mId, 1)
			if err != nil {
				msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
				_, err := b.bot.Send(msg)
				return err
			} else {
				msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите время:")
				sche, err := query.GetCurrentSchedule(db, b.mId, 1)

				var timeArray []string

				//isToday, err := query.GetScheduleDay(db, b.mId)
				if err != nil {
					msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
					_, err := b.bot.Send(msg)
					return err
				} else {
					if time.Now().Weekday() == sche.DayOfWeek() {
						timeArray, err = query.GenerateTimeCollarToday(db)
					} else if time.Monday == sche.DayOfWeek() {
						timeArray, err = query.GenerateTimeCollarMonday(db)
					} else {
						timeArray, err = query.GenerateTimeCollarTomorrow(db)
					}

					if err != nil {
						msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
						_, err := b.bot.Send(msg)
						return err
					}

					if len(timeArray) == 0 {
						msg := tgbotapi.NewMessage(message.Chat.ID, "Нет свободного места.")
						msg.ReplyMarkup = ToBeginning
						_, err := b.bot.Send(msg)
						return err
					}

					keyboard := tgbotapi.InlineKeyboardMarkup{}
					for _, time := range timeArray {
						var row []tgbotapi.InlineKeyboardButton
						btn := tgbotapi.NewInlineKeyboardButtonData(time, time)
						row = append(row, btn)
						keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
					}
					msg.ReplyMarkup = keyboard
					_, err = b.bot.Send(msg)
					msg := tgbotapi.NewMessage(message.Chat.ID, "🔄")
					msg.ReplyMarkup = ToBeginning
					_, err = b.bot.Send(msg)
					return err

				}

			}

		}

	}
	return nil
}

func (b *Bot) handleMassageMenuMedical(message *tgbotapi.Message, db *sql.DB) error {
	status, err := query.GetScheduleStatus(db, b.mId)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
		_, err := b.bot.Send(msg)
		return err
	} else {
		switch status {
		case 0:
			err := query.UpdateScheduleMType(db, massageMenu.Keyboard[0][1].Text, b.mId, 1)
			if err != nil {
				msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
				_, err := b.bot.Send(msg)
				return err
			}
			msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите время:")
			sche, err := query.GetCurrentSchedule(db, b.mId, 1)

			var timeArray []string

			//isToday, err := query.GetScheduleDay(db, b.mId)
			if err != nil {
				msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
				_, err := b.bot.Send(msg)
				return err
			} else {
				if time.Now().Weekday() == sche.DayOfWeek() {
					timeArray, err = query.GenerateTimeMedicalToday(db)
				} else if time.Monday == sche.DayOfWeek() {
					timeArray, err = query.GenerateTimeMedicalMonday(db)
				} else {
					timeArray, err = query.GenerateTimeMedicalTomorrow(db)
				}

				if err != nil {
					msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
					_, err := b.bot.Send(msg)
					return err
				}

				if len(timeArray) == 0 {
					msg := tgbotapi.NewMessage(message.Chat.ID, "Нет свободного места.")
					msg.ReplyMarkup = ToBeginning
					_, err := b.bot.Send(msg)
					return err
				}

				keyboard := tgbotapi.InlineKeyboardMarkup{}
				for _, time := range timeArray {
					var row []tgbotapi.InlineKeyboardButton
					btn := tgbotapi.NewInlineKeyboardButtonData(time, time)
					row = append(row, btn)
					keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
				}
				msg.ReplyMarkup = keyboard
				_, err = b.bot.Send(msg)
				msg := tgbotapi.NewMessage(message.Chat.ID, "🔄")
				msg.ReplyMarkup = ToBeginning
				_, err = b.bot.Send(msg)
				return err

			}

		}

	}
	return nil
}

func (b *Bot) handleFinishEntry(message *tgbotapi.Message, db *sql.DB) error {
	userStatus, err := query.GetUserStatus(db, message.From.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
		_, err := b.bot.Send(msg)
		return err
	} else {
		switch userStatus {
		case 1:
			fullname := message.Text
			msg := tgbotapi.NewMessage(message.Chat.ID, "Ваш корпоративный адрес электронной почты:")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			_, err := b.bot.Send(msg)
			err = query.UpdateUserStatus(db, message.From.ID, 2)
			if err != nil {
				msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
				_, err := b.bot.Send(msg)
				return err
			}
			err = query.UpdateFullname(db, message.From.ID, fullname)
			if err != nil {
				msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
				_, err := b.bot.Send(msg)
				return err
			}
		case 2:
			email := message.Text
			if utils.IsEmailValid(email) {
				err := query.UpdateEmail(db, message.From.ID, email)
				if err != nil {
					msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
					_, err := b.bot.Send(msg)
					return err
				}
				err = query.UpdateUserStatus(db, message.From.ID, 3)
				if err != nil {
					msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
					_, err := b.bot.Send(msg)
					return err
				}
				msg := tgbotapi.NewMessage(message.Chat.ID, "Номер телефона:")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				btn := tgbotapi.KeyboardButton{
					RequestContact: true,
					Text:           "Поделиться номером телефона",
				}
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{btn})
				_, err = b.bot.Send(msg)
				return err

			} else {
				msg := tgbotapi.NewMessage(message.Chat.ID, "Неверный адрес электронной почты, попробуйте еще раз.")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				_, err := b.bot.Send(msg)
				err = query.UpdateUserStatus(db, message.From.ID, 2)
				if err != nil {
					msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
					_, err := b.bot.Send(msg)
					return err
				}
			}
		case 3:
			if message.Contact != nil {
				phoneNum := message.Contact.PhoneNumber
				err := query.UpdatePhoneNum(db, message.From.ID, phoneNum)
				if err != nil {
					msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
					_, err := b.bot.Send(msg)
					return err
				}
				err = query.UpdateUserStatus(db, message.From.ID, 4)
				if err != nil {
					msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
					_, err := b.bot.Send(msg)
					return err
				}
				err = query.RefreshUserList(db)
				if err != nil {
					msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
					_, err := b.bot.Send(msg)
					return err
				}
				msg := tgbotapi.NewMessage(message.Chat.ID, "Вы успешно зарегистрировались 🎉")
				msg.ReplyMarkup = mainMenu
				_, err = b.bot.Send(msg)
				return err
				//if utils.IsPhoneNumberValid(phoneNum) {
				//err := query.UpdatePhoneNum(db, message.From.ID, phoneNum)
				//if err != nil {
				//	msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
				//	_, err := b.bot.Send(msg)
				//	return err
				//}
				//err = query.UpdateUserStatus(db, message.From.ID, 4)
				//if err != nil {
				//	msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
				//	_, err := b.bot.Send(msg)
				//	return err
				//}
				//err = query.RefreshUserList(db)
				//if err != nil {
				//	msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
				//	_, err := b.bot.Send(msg)
				//	return err
				//}
				//msg := tgbotapi.NewMessage(message.Chat.ID, "Вы успешно зарегистрировались 🎉")
				//msg.ReplyMarkup = mainMenu
				//_, err = b.bot.Send(msg)
				//return err
				//} else {
				//	msg := tgbotapi.NewMessage(message.Chat.ID, "Неверный формат, попробуйте еще раз.")
				//	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				//	_, err := b.bot.Send(msg)
				//	err = query.UpdateUserStatus(db, message.From.ID, 3)
				//	if err != nil {
				//		msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
				//		_, err := b.bot.Send(msg)
				//		return err
				//	}
				//}
			} else {
				msg := tgbotapi.NewMessage(message.Chat.ID, "Неверный формат, попробуйте еще раз.")
				_, err := b.bot.Send(msg)
				err = query.UpdateUserStatus(db, message.From.ID, 3)
				if err != nil {
					msg := tgbotapi.NewMessage(message.Chat.ID, "Database error.")
					_, err := b.bot.Send(msg)
					return err
				}
			}
		}
	}
	return nil
}

func (b *Bot) handleCallbackQuery(callback *tgbotapi.CallbackQuery, db *sql.DB) error {
	fmt.Println(callback.Data)
	if callback.Data == "Edit" {
		err := query.UpdateUserIdStatus(db, callback.From.ID, 1)
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Ваше имя:")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		_, err = b.bot.Send(msg)
		return err
	} else {
		if _, err := strconv.Atoi(callback.Data); err == nil {
			intVar, err := strconv.Atoi(callback.Data)
			if err == nil {
				isCancel, err := query.CancelEntry(db, intVar)
				if err != nil {
					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Database error.")
					_, err = b.bot.Send(msg)
					return err
				} else {
					if isCancel {
						msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Запись отменена.")
						msg.ReplyMarkup = mainMenu
						_, err = b.bot.Send(msg)
						obj, err := query.GetCanceledSchedule(db, b.mId)
						text := fmt.Sprintf("Запись отменена."+"\n*%s*\n"+"*Дата: * _%v_\n"+"*Время: * _%s_\n", obj.MType, obj.MDate, obj.MTime)
						msg = tgbotapi.NewMessage(int64(b.conf.AdminID), text)
						msg.ParseMode = "markdown"
						fmt.Println(obj)
						//msg.ReplyMarkup = mainMenu
						_, err = b.bot.Send(msg)
						return err

					}
				}
			} else {
				fmt.Println("Invalid action.")
			}
		} else {
			status, err := query.GetScheduleStatus(db, b.mId)
			if err != nil {
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Database error.")
				_, err = b.bot.Send(msg)
				return err
			} else {
				switch status {
				case 1:
					time := callback.Data
					err := query.UpdateScheduleTime(db, b.mId, time)
					if err != nil {
						msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Database error.")
						_, err = b.bot.Send(msg)
						return err
					}

					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Вы успешно записаны!")
					_, err = b.bot.Send(msg)
					err = query.UpdateScheduleStatus(db, b.mId, 2)
					//err = sendAlarmToAdmin(b.mId, true)
					//if err != nil {
					//	//TODO
					//}
					obj, err := query.GetCurrentSchedule(db, b.mId, 2)
					text := fmt.Sprintf("*%s*\n"+"*Дата: * _%v_\n"+"*Время: * _%s_\n", obj.MType, obj.MDate, obj.MTime)
					msg = tgbotapi.NewMessage(callback.Message.Chat.ID, text)
					msg.ParseMode = "markdown"
					fmt.Println(obj)
					msg.ReplyMarkup = mainMenu
					_, err = b.bot.Send(msg)
					err = query.UpdateScheduleStatus(db, b.mId, 2)
					text = fmt.Sprintf("Новая запись:"+"\n*%s*\n"+"*Дата: * _%v_\n"+"*Время: * _%s_\n", obj.MType, obj.MDate, obj.MTime)
					msg = tgbotapi.NewMessage(int64(b.conf.AdminID), text)
					msg.ParseMode = "markdown"
					fmt.Println(obj)
					//msg.ReplyMarkup = mainMenu
					_, err = b.bot.Send(msg)
					if err != nil {
						msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Database error.")
						_, err = b.bot.Send(msg)
						return err
					}

				}
			}
		}
	}
	return nil
}

func sendAlarmToAdmin(id int, b bool) error {

	return nil
}

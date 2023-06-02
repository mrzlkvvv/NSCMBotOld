package bot

const (
	REGISTER_DATA_TEMPLATE = "\"Иванов Иван Иванович 123456\", где \"123456\" - это номер паспорта"

	MESSAGE_START = "Здесь вы можете узнать свои результаты экзаменов."

	MESSAGE_HELP = "Чтобы зарегистрироваться, отправьте мне свои данные в формате:\n" + REGISTER_DATA_TEMPLATE + "\n\n/check - проверить результаты\n\n/unregister - удалить аккаунт\n\n/help - показать это сообщение\n\nПо всем вопросам обращайтесь к @KirillMerz"

	MESSAGE_REGISTER_SUCCESS = "Вы были успешно зарегистрированы\nПосмотрим результаты? (/check)"

	MESSAGE_NOT_REGISTERED_ERROR = "Для этого вы должны зарегистрироваться :("

	MESSAGE_UNREGISTER_SUCCESS = "Ваши данные были успешно удалены"

	MESSAGE_DATABASE_ERROR = "В работе базы данных произошла ошибка. Пожалуйста, повторите попытку позже"

	MESSAGE_CHECK_ERROR = "Не удалось получить результаты. Повторите попытку позже"

	MESSAGE_UNKNOWN_COMMAND_ERROR = "Я вас не понимаю...\nПоказать помощь - /help"
)

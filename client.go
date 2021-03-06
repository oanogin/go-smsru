package smsru

import (
	"bytes"
	"net/http"
	"net/url"
)

const API_URL = "https://sms.ru"

var CODES = map[int]string{
	-1:  "Сообщение не найдено",
	100: "Запрос выполнен или сообщение находится в нашей очереди",
	101: "Сообщение передается оператору",
	102: "Сообщение отправлено (в пути)",
	103: "Сообщение доставлено",
	104: "Не может быть доставлено: время жизни истекло",
	105: "Не может быть доставлено: удалено оператором",
	106: "Не может быть доставлено: сбой в телефоне",
	107: "Не может быть доставлено: неизвестная причина",
	108: "Не может быть доставлено: отклонено",
	110: "Сообщение прочитано",
	150: "Не может быть доставлено: не найден маршрут на данный номер",
	200: "Неправильный api_id",
	201: "Не хватает средств на лицевом счету",
	202: "Неправильно указан номер телефона получателя, либо на него нет маршрута",
	203: "Нет текста сообщения",
	204: "Имя отправителя не согласовано с администрацией",
	205: "Сообщение слишком длинное (превышает 8 СМС)",
	206: "Будет превышен или уже превышен дневной лимит на отправку сообщений",
	207: "На этот номер нет маршрута для доставки сообщений",
	208: "Параметр time указан неправильно",
	209: "Вы добавили этот номер (или один из номеров) в стоп-лист",
	210: "Используется GET, где необходимо использовать POST",
	211: "Метод не найден",
	212: "Текст сообщения необходимо передать в кодировке UTF-8 (вы передали в другой кодировке)",
	213: "Указано более 100 номеров в списке получателей",
	220: "Сервис временно недоступен, попробуйте чуть позже",
	230: "Превышен общий лимит количества сообщений на этот номер в день",
	231: "Превышен лимит одинаковых сообщений на этот номер в минуту",
	232: "Превышен лимит одинаковых сообщений на этот номер в день",
	300: "Неправильный token (возможно истек срок действия, либо ваш IP изменился)",
	301: "Неправильный пароль, либо пользователь не найден",
	302: "Пользователь авторизован, но аккаунт не подтвержден (пользователь не ввел код, присланный в регистрационной смс)",
	303: "Код подтверждения неверен",
	304: "Отправлено слишком много кодов подтверждения. Пожалуйста, повторите запрос позднее",
	305: "Слишком много неверных вводов кода, повторите попытку позднее",
	500: "Ошибка на сервере. Повторите запрос.",
	901: "Callback: URL неверный (не начинается на http://)",
	902: "Callback: Обработчик не найден (возможно был удален ранее)",
}

type Client struct {
	APIID    string
	HTTP     *http.Client
	Test     bool
	JSON     bool
	Translit bool
	From     string
}

func NewClient(aid string, testF, jsonF, translitF bool) *Client {
	return &Client{
		APIID:    aid,
		HTTP:     &http.Client{},
		Test:     testF,
		JSON:     jsonF,
		Translit: translitF,
	}
}

func (c *Client) makeRequest(endpoint string, params url.Values) (*bytes.Buffer, error) {
	params.Set("api_id", c.APIID)

	u := API_URL + endpoint + "?" + params.Encode()

	resp, err := c.HTTP.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := new(bytes.Buffer)
	_, err = data.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

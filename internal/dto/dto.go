package dto

//разбить по файлам

type CurrentUser struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

type UserShortResponse struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

type UsersResponse struct {
	Users []UserShortResponse `json:"users"`
}

// Структура ответа с токеном для хендлеров.
type TokenResponse struct {
	Token string `json:"token"`
}

// Структура запроса для хендлеров.
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Структура запроса для хендлеров.
// //Сейчас поля такие же как в LoginRequest, но структура другая, т.к. в дальнейшем полей сюда добавится
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

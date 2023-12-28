package dto

type CategoryDto struct {
	Name string `json:"name"`
}

type StandardResponse struct {
	Status string `json:"status"`
	Msg    string `json:"message"`
}

type CourseDto struct {
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	CategoryID  string `json:"category_id"`
}

type UserDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponseDto struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

package dto

type CreateUserDTO struct {
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=6"`
	RepeatPassword string `json:"repeat_password" validate:"required,eqfield=Password"`
	Username       string `json:"username" validate:"required,alphanum,min=3"`
}

type LoginUserDTO struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserDTO struct {
	Username       string `json:"username" validate:"omitempty,alphanum,min=3"`
	Password       string `json:"password" validate:"omitempty,min=6,eqfield=RepeatPassword"`
	RepeatPassword string `json:"repeat_password" validate:"omitempty,eqfield=Password"`
}

type ChangeEmailDTO struct {
	Password string `json:"password" validate:"required"`
	NewEmail string `json:"new_email" validate:"required,email"`
}

type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,nefield=OldPassword"`
}

type VerifyEmailDTO struct {
	Token string `json:"token" validate:"required"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordDTO struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

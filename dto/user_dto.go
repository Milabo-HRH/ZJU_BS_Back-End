package dto

import "ZJU_BS_Back-End/model"

type UserDto struct {
	Name      string `json:"Name"`
	Email     string `json:"Mail"`
	Privilege string `json:"Privilege"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		Name:      user.Name,
		Email:     user.Mail,
		Privilege: user.Privilege,
	}
}

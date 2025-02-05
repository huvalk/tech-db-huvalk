// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// User Информация о пользователе.
//
//
// swagger:model User
//easyjson:json
type User struct {

	// Описание пользователя.
	About string `json:"about,omitempty"`

	// Почтовый адрес пользователя (уникальное поле).
	// Required: true
	// Format: email
	Email string `json:"email"`

	// Полное имя пользователя.
	// Required: true
	Fullname string `json:"fullname"`

	// Имя пользователя (уникальное поле).
	// Данное поле допускает только латиницу, цифры и знак подчеркивания.
	// Сравнение имени регистронезависимо.
	//
	// Read Only: true
	Nickname string `json:"nickname,omitempty"`
}

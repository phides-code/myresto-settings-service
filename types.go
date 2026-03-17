package main

type Settings struct {
	BannerMessage string          `json:"bannerMessage" validate:"required"`
	Hours         SettingsHours   `json:"hours" validate:"required"`
	Contact       SettingsContact `json:"contact" validate:"required"`
	Social        SettingsSocial  `json:"social" validate:"required"`
}

type SettingsHours struct {
	Mon string `json:"mon"`
	Tue string `json:"tue"`
	Wed string `json:"wed"`
	Thu string `json:"thu"`
	Fri string `json:"fri"`
	Sat string `json:"sat"`
	Sun string `json:"sun"`
}

type SettingsContact struct {
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type SettingsSocial struct {
	Instagram string `json:"instagram"`
	Facebook  string `json:"facebook"`
	Tiktok    string `json:"tiktok"`
}

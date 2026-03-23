package main

type Settings struct {
	BannerMessage  string `json:"bannerMessage" validate:"required"`
	HoursMonday    string `json:"hoursMonday" validate:"required"`
	HoursTuesday   string `json:"hoursTuesday" validate:"required"`
	HoursWednesday string `json:"hoursWednesday" validate:"required"`
	HoursThursday  string `json:"hoursThursday" validate:"required"`
	HoursFriday    string `json:"hoursFriday" validate:"required"`
	HoursSaturday  string `json:"hoursSaturday" validate:"required"`
	HoursSunday    string `json:"hoursSunday" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	Email          string `json:"email" validate:"required"`
	Address        string `json:"address" validate:"required"`
	Instagram      string `json:"instagram" validate:"required"`
	Facebook       string `json:"facebook" validate:"required"`
	Tiktok         string `json:"tiktok" validate:"required"`
}

type ImageSource struct {
	OriginalName string `json:"originalName"`
	UUIDName     string `json:"uuidName"`
}

type Theme struct {
	ThemeName       string `json:"themeName" validate:"required"`
	BackgroundColor string `json:"backgroundColor" validate:"required"`
	TextColor       string `json:"textColor" validate:"required"`
	LinkColor       string `json:"linkColor" validate:"required"`
	AccentColor     string `json:"accentColor" validate:"required"`
	ButtonColor     string `json:"buttonColor" validate:"required"`
	ButtonTextColor string `json:"buttonTextColor" validate:"required"`
}

type ThemeSettings struct {
	BannerImage   ImageSource `json:"bannerImage" validate:"required"`
	SelectedTheme Theme       `json:"selectedTheme" validate:"required"`
}

package main

type Settings struct {
	BannerMessage  string      `json:"bannerMessage" validate:"required"`
	BannerImage    ImageSource `json:"bannerImage" validate:"required"`
	HoursMonday    string      `json:"hoursMonday" validate:"required"`
	HoursTuesday   string      `json:"hoursTuesday" validate:"required"`
	HoursWednesday string      `json:"hoursWednesday" validate:"required"`
	HoursThursday  string      `json:"hoursThursday" validate:"required"`
	HoursFriday    string      `json:"hoursFriday" validate:"required"`
	HoursSaturday  string      `json:"hoursSaturday" validate:"required"`
	HoursSunday    string      `json:"hoursSunday" validate:"required"`
	Phone          string      `json:"phone" validate:"required"`
	Email          string      `json:"email" validate:"required"`
	Address        string      `json:"address" validate:"required"`
	Instagram      string      `json:"instagram" validate:"required"`
	Facebook       string      `json:"facebook" validate:"required"`
	Tiktok         string      `json:"tiktok" validate:"required"`
}

type ImageSource struct {
	OriginalName string `json:"originalName"`
	UUIDName     string `json:"uuidName"`
}

package model

type RequestOtpRegister struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

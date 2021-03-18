package models

type QRCodeBody struct {
	Appid int
	Size  int
	T     string
}

type QRCodeTicketBody struct {
	Code   int    `json:"code"`
	Ticket string `json:"ticket"`
	Msg    string `json:"msg"`
}

type ValidateQRCodeBody struct {
	ReturnCode int    `json:"returnCode"`
	Url        string `json:"url"`
}

package request

type DialogList struct {
	UserID int `uri:"user_id" binding:"required"`
}

type DialogMessage struct {
	UserID int `uri:"user_id" binding:"required"`
}

type SendDialogMessage struct {
	Text string `json:"text" binding:"required"`
}

package command

type Subscribing struct {
	Email string `form:"email" binding:"required"`
}

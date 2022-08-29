package inputs

type Subscribing struct {
	Email string `form:"email" binding:"required"`
}

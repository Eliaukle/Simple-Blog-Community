package vo

type CreateArticleRequest struct {
	// 加上binging用于表单验证
	CategoryId uint   `json:"category_id" binging:"required"`
	Title      string `json:"title" binging:"required"`
	Content    string `json:"content" binging:"required"`
	HeadImage  string `json:"head_image"`
}

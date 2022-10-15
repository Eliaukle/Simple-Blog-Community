package controller

import (
	"blog_server/common"
	"blog_server/model"
	"blog_server/response"
	"blog_server/vo"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

type ArticleController struct {
	DB *gorm.DB
}

type IArticleController interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Show(c *gin.Context)
	List(c *gin.Context)
}

func (a ArticleController) Create(c *gin.Context) {
	var articleRequest vo.CreateArticleRequest
	// 数据验证
	if err := c.ShouldBindJSON(&articleRequest); err != nil {
		response.Fail(c, nil, "数据错误")
		return
	}
	// 获取登录用户
	user, _ := c.Get("user")
	// 创建文章
	article := model.Article{
		UserId:     user.(model.User).ID,
		CategoryId: articleRequest.CategoryId,
		Title:      articleRequest.Title,
		Content:    articleRequest.Content,
		HeadImage:  articleRequest.HeadImage,
	}
	if err := a.DB.Create(&article).Error; err != nil {
		response.Fail(c, nil, "发布失败")
		return
	}
	response.Success(c, gin.H{"id": article.ID}, "发布成功")
}

func (a ArticleController) Update(c *gin.Context) {
	var articleRequest vo.CreateArticleRequest
	// 数据验证
	if err := c.ShouldBindJSON(&articleRequest); err != nil {
		response.Fail(c, nil, "数据错误")
		return
	}
	// 获取path中的id
	articleId := c.Params.ByName("id")
	// 查找文章
	var article model.Article
	if a.DB.Where("id = ?", articleId).First(&article).RecordNotFound() {
		response.Fail(c, nil, "文章不存在")
		return
	}
	// 获取登录用户
	user, _ := c.Get("user")
	userId := user.(model.User).ID
	if userId != article.UserId {
		response.Fail(c, nil, "登录用户不正确")
		return
	}
	// 更新文章
	if err := a.DB.Model(&article).Update(articleRequest).Error; err != nil {
		response.Fail(c, nil, "修改失败")
		return
	}
	response.Success(c, nil, "修改成功")
}

func (a ArticleController) Delete(c *gin.Context) {
	// 获取path中的id
	articleId := c.Params.ByName("id")
	// 查找文章
	var article model.Article
	if a.DB.Where("id = ?", articleId).First(&article).RecordNotFound() {
		response.Fail(c, nil, "文章不存在")
		return
	}
	// 获取登录用户
	user, _ := c.Get("user")
	userId := user.(model.User).ID
	if userId != article.UserId {
		response.Fail(c, nil, "登录用户不正确")
		return
	}
	// 删除文章
	if err := a.DB.Delete(&article).Error; err != nil {
		response.Fail(c, nil, "删除失败")
		return
	}
	response.Success(c, nil, "删除成功")
}

func (a ArticleController) Show(c *gin.Context) {
	// 获取path中的id
	articleId := c.Params.ByName("id")
	// 查找文章
	var article model.Article
	if a.DB.Where("id = ?", articleId).First(&article).RecordNotFound() {
		response.Fail(c, nil, "文章不存在")
		return
	}
	// 展示文章详情
	response.Success(c, gin.H{"article": article}, "查找成功")
}

func (a ArticleController) List(c *gin.Context) {
	// 获取关键词、分类、分页参数
	keyword := c.DefaultQuery("keyword", "")
	categoryId := c.DefaultQuery("categoryId", "0")
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	var query []string
	var args []string
	// 若关键词存在
	if keyword != "" {
		query = append(query, "(title LIKE ? OR content LIKE ?)")
		args = append(args, "%"+keyword+"%")
		args = append(args, "%"+keyword+"%")
	}
	// 若分类存在
	if categoryId != "0" {
		query = append(query, "category_id = ?")
		args = append(args, categoryId)
	}
	// 拼接字符串
	var querystr string
	if len(query) > 0 {
		querystr = strings.Join(query, " AND ")
	}
	// 页面内容
	var article []model.ArticleInfo
	// 文章总数
	var count int
	// 查询文章
	switch len(args) {
	case 0:
		a.DB.Table("articles").Select("id, category_id, title, LEFT(content,80) AS content, head_image, created_at").Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&article)
		a.DB.Model(model.Article{}).Count(&count)
	case 1:
		a.DB.Table("articles").Select("id, category_id, title, LEFT(content,80) AS content, head_image, created_at").Where(querystr, args[0]).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&article)
		a.DB.Model(model.Article{}).Where(querystr, args[0]).Count(&count)
	case 2:
		a.DB.Table("articles").Select("id, category_id, title, LEFT(content,80) AS content, head_image, created_at").Where(querystr, args[0], args[1]).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&article)
		a.DB.Model(model.Article{}).Where(querystr, args[0], args[1]).Count(&count)
	case 3:
		a.DB.Table("articles").Select("id, category_id, title, LEFT(content,80) AS content, head_image, created_at").Where(querystr, args[0], args[1], args[2]).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&article)
		a.DB.Model(model.Article{}).Where(querystr, args[0], args[1], args[2]).Count(&count)
	}
	// 展示文章列表
	response.Success(c, gin.H{"article": article, "count": count}, "查找成功")
}

func NewArticleController() IArticleController {
	db := common.GetDB()
	db.AutoMigrate(model.Article{})
	return ArticleController{DB: db}
}

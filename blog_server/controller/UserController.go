package controller

import (
	"blog_server/common"
	"blog_server/model"
	"blog_server/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

// Register 注册
func Register(c *gin.Context) {
	db := common.GetDB()
	// 获取参数
	var requestUser model.User
	c.Bind(&requestUser)
	userName := requestUser.UserName
	phoneNumber := requestUser.PhoneNumber
	password := requestUser.Password
	// 数据验证
	var user model.User
	db.Where("phone_number = ?", phoneNumber).First(&user)
	if user.ID != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 422,
			"msg":  "用户已存在",
		})
		return
	}
	// 密码加密
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// 创建用户
	newUser := model.User{
		UserName:    userName,
		PhoneNumber: phoneNumber,
		Password:    string(hashedPassword),
		Avatar:      "/images/default_avatar.png",
		Collects:    model.Array{},
		Following:   model.Array{},
		Fans:        0,
	}
	db.Create(&newUser)
	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
	})
}

// Login 登录
func Login(c *gin.Context) {
	db := common.GetDB()
	// 获取参数
	var requestUser model.User
	c.Bind(&requestUser)
	phoneNumber := requestUser.PhoneNumber
	password := requestUser.Password
	// 数据验证
	var user model.User
	db.Where("phone_number =?", phoneNumber).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 422,
			"msg":  "用户不存在",
		})
		return
	}
	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 422,
			"msg":  "密码错误",
		})
		return
	}
	// 发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "系统异常",
		})
		return
	}
	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"token": token},
		"msg":  "登录成功",
	})
}

// GetInfo 登录后获取信息
func GetInfo(c *gin.Context) {
	// 获取上下文中的用户信息
	user, _ := c.Get("user")
	// 返回用户信息
	response.Success(c, gin.H{"id": user.(model.User).ID, "avatar": user.(model.User).Avatar}, "登录获取信息成功")
	//c.JSON(http.StatusOK, gin.H{
	//	"code": 200,
	//	"data": gin.H{"id": user.(model.User).ID, "avatar": user.(model.User).Avatar},
	//	"msg":  "登录获取信息成功",
	//})
}

// GetBriefInfo 获取简要信息
func GetBriefInfo(c *gin.Context) {
	db := common.GetDB()
	// 获取path中的userId
	userId := c.Params.ByName("id")
	// 判断用户身份
	user, _ := c.Get("user")
	//var self bool
	var curUser model.User
	if userId == strconv.Itoa(int(user.(model.User).ID)) {
		//self = true
		curUser = user.(model.User)
	} else {
		//self = false
		db.Where("id =?", userId).First(&curUser)
		if curUser.ID == 0 {
			response.Fail(c, nil, "用户不存在")
			return
		}
	}
	// 返回用户简要信息
	response.Success(c, gin.H{"id": curUser.ID, "name": curUser.UserName, "avatar": curUser.Avatar, "loginId": user.(model.User).ID}, "查找成功")
}

// GetDetailedInfo 获取详细信息
func GetDetailedInfo(c *gin.Context) {
	db := common.GetDB()
	// 获取path中的userId
	userId := c.Params.ByName("id")
	// 判断用户身份
	user, _ := c.Get("user")
	//var self bool
	var curUser model.User
	if userId == strconv.Itoa(int(user.(model.User).ID)) {
		//self = true
		curUser = user.(model.User)
	} else {
		//self = false
		db.Where("id = ?", userId).First(&curUser)
		if curUser.ID == 0 {
			response.Fail(c, nil, "用户不存在")
			return
		}
	}
	// 返回用户详细信息
	var articles, collects []model.ArticleInfo
	var following []model.UserInfo
	var collist, follist []string
	collist = ToStringArray(curUser.Collects)
	follist = ToStringArray(curUser.Following)
	db.Table("articles").Select("id, category_id, title, LEFT(content,80) AS content, head_image, created_at").Where("user_id = ?", userId).Order("created_at desc").Find(&articles)
	db.Table("articles").Select("id, category_id, title, LEFT(content,80) AS content, head_image, created_at").Where("id IN (?)", collist).Order("created_at desc").Find(&collects)
	db.Table("users").Select("id, avatar, user_name").Where("id IN (?)", follist).Find(&following)
	response.Success(c, gin.H{"id": curUser.ID, "name": curUser.UserName, "avatar": curUser.Avatar, "loginId": user.(model.User).ID, "articles": articles, "collects": collects, "following": following, "fans": curUser.Fans}, "查找成功")
}

// ModifyAvatar 修改头像
func ModifyAvatar(c *gin.Context) {
	db := common.GetDB()
	// 获取用户ID
	user, _ := c.Get("user")
	// 获取参数
	var requestUser model.User
	c.Bind(&requestUser)
	avatar := requestUser.Avatar
	// 查找用户
	var curUser model.User
	db.Where("id = ?", user.(model.User).ID).First(&curUser)
	// 更新信息
	if err := db.Model(&curUser).Update("avatar", avatar).Error; err != nil {
		response.Fail(c, nil, "更新失败")
		return
	}
	response.Success(c, nil, "更新成功")
}

// ModifyName 修改用户名
func ModifyName(c *gin.Context) {
	db := common.GetDB()
	// 获取用户ID
	user, _ := c.Get("user")
	// 获取参数
	var requestUser model.User
	c.Bind(&requestUser)
	userName := requestUser.UserName
	// 查找用户
	var curUser model.User
	db.Where("id = ?", user.(model.User).ID).First(&curUser)
	// 更新信息
	if err := db.Model(&curUser).Update("user_name", userName).Error; err != nil {
		response.Fail(c, nil, "更新失败")
		return
	}
	response.Success(c, nil, "更新成功")
}

// Collects 查询收藏
func Collects(c *gin.Context) {
	db := common.GetDB()
	// 获取用户ID
	user, _ := c.Get("user")
	// 获取path中的id
	id := c.Params.ByName("id")
	var curUser model.User
	db.Where("id = ?", user.(model.User).ID).First(&curUser)
	// 判断是否已收藏
	for i := 0; i < len(curUser.Collects); i++ {
		if curUser.Collects[i] == id {
			response.Success(c, gin.H{"collected": true, "index": i}, "查询成功")
			return
		}
	}
	response.Success(c, gin.H{"collected": false}, "查询成功")
}

// NewCollect 新增收藏
func NewCollect(c *gin.Context) {
	db := common.GetDB()
	// 获取用户ID
	user, _ := c.Get("user")
	// 获取path中的id
	id := c.Params.ByName("id")
	// 查找用户
	var curUser model.User
	db.Where("id = ?", user.(model.User).ID).First(&curUser)
	var newCollects []string
	newCollects = append(curUser.Collects, id)
	// 更新收藏夹
	if err := db.Model(&curUser).Update("collects", newCollects).Error; err != nil {
		response.Fail(c, nil, "更新失败")
		return
	}
	response.Success(c, nil, "更新成功")
}

// UnCollect 取消收藏
func UnCollect(c *gin.Context) {
	db := common.GetDB()
	// 获取用户ID
	user, _ := c.Get("user")
	// 获取path中的index
	index, _ := strconv.Atoi(c.Params.ByName("index"))
	// 查找用户
	var curUser model.User
	db.Where("id = ?", user.(model.User).ID).First(&curUser)
	var newCollects []string
	newCollects = append(curUser.Collects[:index], curUser.Collects[index+1:]...)
	// 更新收藏夹
	if err := db.Model(&curUser).Update("collects", newCollects).Error; err != nil {
		response.Fail(c, nil, "更新失败")
		return
	}
	response.Success(c, nil, "更新成功")
}

// Following 查询关注
func Following(c *gin.Context) {
	db := common.GetDB()
	// 获取用户ID
	user, _ := c.Get("user")
	// 获取path中的id
	id := c.Params.ByName("id")
	var curUser model.User
	db.Where("id = ?", user.(model.User).ID).First(&curUser)
	// 判断是否已关注
	for i := 0; i < len(curUser.Following); i++ {
		if curUser.Following[i] == id {
			response.Success(c, gin.H{"followed": true, "index": i}, "查询成功")
			return
		}
	}
	response.Success(c, gin.H{"followed": false}, "查询成功")
}

// NewFollow 新增关注
func NewFollow(c *gin.Context) {
	db := common.GetDB()
	// 获取用户ID
	user, _ := c.Get("user")
	// 获取path中的id
	id := c.Params.ByName("id")
	// 查找用户
	var curUser model.User
	db.Where("id = ?", user.(model.User).ID).First(&curUser)
	//var newFollowing []string
	newFollowing := append(curUser.Following, id)
	// 更新关注列表
	if err := db.Model(&curUser).Update("following", newFollowing).Error; err != nil {
		response.Fail(c, nil, "更新失败")
		return
	}
	// 更新粉丝数
	var followUser model.User
	db.Where("id = ?", id).First(&followUser)
	if err := db.Model(&followUser).Update("fans", followUser.Fans+1).Error; err != nil {
		response.Fail(c, nil, "更新失败")
		return
	}
	response.Success(c, nil, "更新成功")
}

// UnFollow 取消关注
func UnFollow(c *gin.Context) {
	db := common.GetDB()
	// 获取用户ID
	user, _ := c.Get("user")
	// 获取path中的index
	index, _ := strconv.Atoi(c.Params.ByName("index"))
	// 查找用户
	var curUser model.User
	db.Where("id = ?", user.(model.User).ID).First(&curUser)
	//var newFollowing []string
	newFollowing := append(curUser.Following[:index], curUser.Following[index+1:]...)
	followId := curUser.Following[index]
	// 更新关注列表
	if err := db.Model(&curUser).Update("following", newFollowing).Error; err != nil {
		response.Fail(c, nil, "更新失败")
		return
	}
	// 更新粉丝数
	var followUser model.User
	db.Where("id = ?", followId).First(&followUser)
	if err := db.Model(&followUser).Update("fans", followUser.Fans-1).Error; err != nil {
		response.Fail(c, nil, "更新失败")
		return
	}
	response.Success(c, nil, "更新成功")
}

// ToStringArray 将自定义类型转化为字符串数组
func ToStringArray(l []string) (a model.Array) {
	for i := 0; i < len(a); i++ {
		l = append(l, a[i])
	}
	return l
}

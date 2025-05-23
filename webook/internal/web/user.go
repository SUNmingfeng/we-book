package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	time2 "time"
)

const (
	emailRegexpPattern = "^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(.[a-zA-Z0-9_-]+)+$"
	//至少包含字母、数字、特殊字符，1-9位
	passwordRegexpPattern = "^(?=.*\\d)(?=.*[a-zA-Z])(?=.*[^\\da-zA-Z\\s]).{1,9}$"
	bizLogin              = "login"
)

var (
	ErrDuplicateEmail  = service.ErrDuplicateEmail
	ErrUserNotFound    = service.ErrInvaildUserOrPassword
	ErrCodeSendTooMany = service.ErrCodeSendTooMany
)

type UserHandler struct {
	jwtHandler
	emailRexExp    *regexp.Regexp
	passwordPexExp *regexp.Regexp
	svc            service.UserService
	codeSvc        service.CodeService
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		//预编译正则
		emailRexExp:    regexp.MustCompile(emailRegexpPattern, regexp.None),
		passwordPexExp: regexp.MustCompile(passwordRegexpPattern, regexp.None),
		svc:            svc,
		codeSvc:        codeSvc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 使用group分组路由来简化注册
	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp)
	//ug.POST("/login", h.Login)
	ug.POST("/login", h.LoginJWT)
	ug.GET("/profile", h.ProFile)
	ug.POST("/edit", h.Edit)
	ug.POST("login_sms/code/send", h.SendSMSLoginCode)
	ug.POST("login_sms", h.LoginSMS)

}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignReq
	//取出前端数据到req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//对前端输出的一些格式验证
	isEmail, _ := h.emailRexExp.MatchString(req.Email)
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}
	isPassword, _ := h.passwordPexExp.MatchString(req.Password)
	if !isPassword {
		ctx.String(http.StatusOK, "非法密码格式，至少包含字母、数字、特殊字符，1-9位")
		return
	}

	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次密码输入不一致")
		return
	}

	//存储数据
	err := h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		PassWord: req.Password,
	})
	switch err {
	case nil:
		ctx.String(http.StatusOK, fmt.Sprint("注册成功！"))
	case ErrDuplicateEmail:
		ctx.String(http.StatusOK, "注册邮箱冲突，请换一个")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	//取出前端数据到req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		//登录成功后存储session
		sess := sessions.Default(ctx) //获取这个请求的session
		sess.Set("userId", u.Id)      //必须与save搭配使用才能生效
		sess.Options(sessions.Options{
			MaxAge: 900, //有效期
			//HttpOnly: true,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "存储session错误:", err)
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvaildUserOrPassword:
		ctx.String(http.StatusOK, "用户不存在或密码不正确")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	//取出前端数据到req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		h.setJWTToken(ctx, u.Id)
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvaildUserOrPassword:
		ctx.String(http.StatusOK, "用户不存在或密码不正确")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Nickname string `json:"Nickname"`
		Birthday string `json:"Birthday"`
		AboutMe  string `json:"AboutMe"`
	}
	type Result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data any    `json:"data"`
	}
	var req EditReq
	//取出前端数据到req
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "获取信息错误")
		return
	}
	sess := sessions.Default(ctx) //获取当前 HTTP 请求的会话
	userid := sess.Get("userId").(int64)
	if len(req.Nickname) > 8 {
		res := Result{
			Code: 1,
			Msg:  "昵称长度不能超过8个字符",
		}
		ctx.JSON(http.StatusOK, res)
		return
	}
	birthday, err := time2.Parse(time2.DateOnly, req.Birthday)
	if err != nil {
		//ctx.String(http.StatusOK, "系统错误")
		ctx.String(http.StatusOK, "生日格式不正确")
		return
	}
	fmt.Printf("获取到新设置生日：%v", birthday)
	if len(req.AboutMe) > 128 {
		res := Result{
			Code: 1,
			Msg:  "个人简介不能超过128字符",
		}
		ctx.JSON(http.StatusOK, res)
		return
	}
	err = h.svc.UpdateInfo(ctx, domain.User{
		Id:       userid,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.String(http.StatusOK, "更新异常")
		return
	}
	res := Res{
		Code: 0,
		Msg:  "更新成功",
	}
	ctx.JSON(http.StatusOK, res)
}

type Res struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (h *UserHandler) ProFile(ctx *gin.Context) {
	uc := ctx.MustGet("user").(UserClaims)
	//sess := sessions.Default(ctx) //获取这个请求的session
	//userid := sess.Get("userId").(int64)
	//u, err := h.svc.FindById(ctx, userid)
	u, err := h.svc.FindById(ctx, uc.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "获取数据错误")
	}
	type User struct {
		Nickname string `json:"Nickname"`
		Email    string `json:"Email"`
		AboutMe  string `json:"AboutMe"`
		Birthday string `json:"Birthday"`
	}
	log.Println("查到的数据：", u)
	ctx.JSON(http.StatusOK, User{
		Nickname: u.Nickname,
		Email:    u.Email,
		Birthday: u.Birthday.Format(time2.DateOnly),
		AboutMe:  u.AboutMe,
	})
}

func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	println("进入发送验证码方法")
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "手机号码不能为空",
		})
		return
	}

	if len(req.Phone) != 11 {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "手机号码格式不正确",
		})
		return
	}

	err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "短信发送频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
}

func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Res{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Res{
			Code: 4,
			Msg:  "验证码错误，请重新输入",
		})
		return
	}

	u, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Res{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	h.setJWTToken(ctx, u.Id)
	ctx.JSON(http.StatusOK, Res{
		Msg: "登录成功",
	})
	return
}

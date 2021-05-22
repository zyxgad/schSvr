
package kpnmwebpage

import (
	http "net/http"
	os   "os"

	gin  "github.com/gin-gonic/gin"
	util "github.com/zyxgad/go-util/util"
	ksql "github.com/zyxgad/schSvr/handles/sql"
	kses "github.com/zyxgad/schSvr/handles/sql/session"
	krbt "github.com/zyxgad/schSvr/robot"
)


func getClientUuid(cont *gin.Context)(uuid string){
	var (
		err error
	)
	uuid, err = cont.Cookie("client_uuid")
	if err != nil {
		uuid, _ = kses.NewSessionUuid()
		cont.SetCookie("client_uuid", uuid, 60 * 60 * 24 * 30, "/", "", false, true)
		cont.SetCookie("change_uuid_flag", "T", 60 * 60 * 24 * 15, "/", "", false, true)
	}
	return uuid
}

func updateClientUuid(cont *gin.Context)(uuid string){
	var (
		err error
	)
	uuid = getClientUuid(cont)
	_, err = cont.Cookie("change_uuid_flag")
	if err != nil {
		uuid, _ = kses.ChangeSessionUuid(uuid)
		cont.SetCookie("change_uuid_flag", "T", 60 * 60 * 24 * 15, "/", "", false, true)
	}
	cont.SetCookie("client_uuid", uuid, 60 * 60 * 24 * 30, "/", "", false, true)
	return uuid
}

func removeClientUuid(cont *gin.Context){
	uuid, err := cont.Cookie("client_uuid")
	if err != nil {
		return
	}
	kses.RemoveAllSession(uuid)
	cont.SetCookie("client_uuid", "", -1, "/", "", false, true)
}


func getUserByStrid(userid string)(user *SqlUserType){
	if !util.StrIsInt(userid, 10) {
		return nil
	}
	user = getUserById((uint32)(util.StrToInt(userid, 10)))
	if user == nil {
		return nil
	}
	return user
}

func getLoginUser(suuid string)(user *SqlUserType){
	loginuser := kses.GetSession(suuid, "loginuser")
	if loginuser == nil {
		return nil
	}
	return getUserByStrid(loginuser.Value)
}

type userPageSrc int

func (userPageSrc)LoginGetPage(cont *gin.Context){
	suuid := updateClientUuid(cont)
	user := getLoginUser(suuid)
	if user != nil {
		cont.Redirect(http.StatusFound, "/")
		return
	}
	cont.HTML(http.StatusOK, "user/login.html", gin.H{
	})
}

func (userPageSrc)LoginPostPage(cont *gin.Context){
	suuid := updateClientUuid(cont)

	idv := kses.GetSession(suuid, "captcha_id")
	if idv == nil {
		cont.JSON(http.StatusOK, CreateJsonError("CaptchaError", "session no captcha id", ""))
		return
	}
	captcode := cont.PostForm("captcode")
	if ok := krbt.VerifyCaptcha(idv.Value, captcode); !ok {
		cont.JSON(http.StatusOK, CreateJsonError("CaptchaError", "captcha code error", ""))
		return
	}

	var user *SqlUserType = nil
	username := cont.PostForm("username")
	password := cont.PostForm("password")

	if !reg_name.MatchString(username) {
		cont.JSON(http.StatusOK, CreateJsonError("IllegalArgumentException", "username is illegal data", ""))
		return
	}
	user = getUserByName(username)
	if user == nil {
		cont.JSON(http.StatusOK, CreateJsonError("UserNotExistException", "user is not exist", ""))
		return
	}

	user.Password = password
	if !reg_pwd.MatchString(password) || hashUserPwd(user) != user.ShaPwd {
		cont.JSON(http.StatusOK, CreateJsonError("IllegalArgumentException", "the password is wrong", ""))
		return
	}
	user.Password = ""

	if !user.Verified {
		cont.JSON(http.StatusOK, CreateJsonError("UserNotVerifyException", "user is not verified", ""))
		return
	}

	kses.SetSession(&kses.SqlSessionValue{
		Uuid: suuid,
		Key: "loginuser",
		Value: util.JoinObjStr(user.Id),
		Overtime: util.GetTimeAfter(util.TimeDay * 30),
	})
	cont.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (userPageSrc)LogoutPostPage(cont *gin.Context){
	suuid := updateClientUuid(cont)
	loginuser := kses.GetSession(suuid, "loginuser")
	if loginuser != nil {
		kses.RemoveSession(suuid, "loginuser")
	}
	cont.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (userPageSrc)RegisterGetPage(cont *gin.Context){
	suuid := updateClientUuid(cont)
	user := getLoginUser(suuid)
	if user != nil {
		cont.Redirect(http.StatusFound, "/")
		return
	}
	cont.HTML(http.StatusOK, "user/register.html", gin.H{
	})
}

func (userPageSrc)RegisterPostPage(cont *gin.Context){
	suuid := updateClientUuid(cont)

	idv := kses.GetSession(suuid, "captcha_id")
	if idv == nil {
		cont.JSON(http.StatusOK, CreateJsonError("CaptchaError", "session no captcha id", ""))
		return
	}
	captcode := cont.PostForm("captcode")
	if ok := krbt.VerifyCaptcha(idv.Value, captcode); !ok {
		cont.JSON(http.StatusOK, CreateJsonError("CaptchaError", "captcha code error", ""))
		return
	}

	username := cont.PostForm("username")
	password := cont.PostForm("password")
	manager := cont.PostForm("manager")
	if !reg_name.MatchString(username) || !reg_pwd.MatchString(password) || !reg_name.MatchString(manager) {
		cont.JSON(http.StatusOK, CreateJsonError("IllegalArgumentException", "username/password/manager is illegal data", ""))
		return
	}

	var mid uint32
	{
		muser := getUserByName(manager)
		if muser == nil || !muser.Op_v_user {
			cont.JSON(http.StatusOK, CreateJsonError("UserNotFound", "can not found manager user", ""))
			return
		}
		mid = muser.Id
	}

	user := &SqlUserType{
		Username: username,
		Password: password,
		Manager: mid,
	}
	if err := createUser(user); err != nil {
		cont.JSON(http.StatusOK, CreateJsonError("UserExistException", "user is exist", err.Error()))
		return
	}

	cont.JSON(http.StatusOK, gin.H{ "status": "ok" })
}

func (userPageSrc)CaptchaGetPage(cont *gin.Context){
	suuid := updateClientUuid(cont)
	switch cont.Param("mode") {
	case "image":
		id, imgdata, err := krbt.NewCaptcha()
		if err != nil {
			cont.JSON(http.StatusOK, CreateJsonError("NewCaptchaException", err.Error(), ""))
			return
		}
		if idv := kses.GetSession(suuid, "captcha_id"); idv != nil {
			kses.RemoveSession(suuid, "captcha_id")
			krbt.RemoveCaptcha(idv.Value)
		}
		kses.SetSession(&kses.SqlSessionValue{
			Uuid: suuid,
			Key: "captcha_id",
			Value: id,
			Overtime: util.GetTimeAfter(util.TimeMin * 10),
		})
		cont.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data": imgdata,
		})
	default:
		cont.JSON(http.StatusNotFound, CreateJsonError("NoModelException", "don't have the model's partten", ""))
	}
}

func (userPageSrc)InfoGetPage(cont *gin.Context){
	userid := cont.Param("id")
	if !util.StrIsInt(userid, 10) {
		cont.JSON(http.StatusNotFound, CreateJsonError("UseridIllegalException", "userid is must be a number", ""))
		return
	}
	cont.HTML(http.StatusOK, "user/info/info.html", gin.H{
		"userid": userid,
	})
}

func (userPageSrc)InfoResPage(cont *gin.Context){
	userid := cont.Param("id")
	if !util.StrIsInt(userid, 10) {
		cont.JSON(http.StatusNotFound, CreateJsonError("UseridIllegalException", "userid is must be a number", ""))
		return
	}
	user := getUserByStrid(userid)
	if user == nil {
		cont.JSON(http.StatusNotFound, CreateJsonError("UserNotExistException", "user is not exist", ""))
		return
	}
	mode := cont.Param("mode")
	switch mode {
	case "head":
		fd, err := os.Open(util.JoinPath(USER_DATA_PATH, util.JoinObjStr(user.Id), "head.png"))
		if err != nil {
			fd, _ = os.Open(util.JoinPath(RES_PATH, "images", "empty.png"))
			if fd == nil {
				cont.JSON(http.StatusNotFound, CreateJsonError("OpenFileError", "open user head file error", err.Error()))
				return
			}
		}
		cont.Status(http.StatusOK)
		util.MustCopyWR(cont.Writer, fd)
		return
	}
	var response gin.H = gin.H{
		"userid": user.Id,
		"username": user.Username,
		"verified": user.Verified,
		"frozen": user.Frozen,
		"op_v_user": user.Op_v_user,
		"op_v_quest": user.Op_v_quest,
	}
	switch mode {
	case "info":
	case "auth":
		response["auths"] = gin.H{
			"v_user": user.Op_v_user,
			"v_quest": user.Op_v_quest,
		}
	case "children":
		lines, err := sqlUserTable.SqlSearch(ksql.TypeMap{"id": ksql.TYPE_Uint32}, ksql.WhereMap{{"manager", "=", user.Id, ""}}, 0)
		if err != nil {
			cont.JSON(http.StatusOK, CreateJsonError("SqlSelectError", "search user error", err.Error()))
		}
		children := make([]uint32, 0, len(lines))
		for _, l := range lines {
			children = append(children, util.JsonToUint32(l["id"]))
		}
		response["children"] = children
	default:
		cont.JSON(http.StatusNotFound, CreateJsonError("NoModelException", "don't have the model's partten", ""))
		return
	}
	cont.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data": response,
	})
}

func (userPageSrc)InfoSetPage(cont *gin.Context){
	suuid := updateClientUuid(cont)
	user := getLoginUser(suuid)
	if user == nil {
		cont.JSON(http.StatusOK, CreateJsonError("UserNotLogin", "user is not login", ""))
		return
	}
	if !user.Op_v_user {
		cont.JSON(http.StatusOK, CreateJsonError("UserNoPermission", "user no permission", ""))
		return
	}
	userid := cont.Param("id")
	if !util.StrIsInt(userid, 10) {
		cont.JSON(http.StatusOK, CreateJsonError("UseridIllegalException", "userid is must be a number", ""))
		return
	}
	targetuser := getUserByStrid(userid)
	if targetuser == nil {
		cont.JSON(http.StatusOK, CreateJsonError("UserNotExistException", "user is not exist", ""))
		return
	}
	if targetuser.Manager != user.Id {
		cont.JSON(http.StatusOK, CreateJsonError("UserNoPermission", "user no permission", ""))
		return
	}
	var json map[string]interface{}
	if err := cont.ShouldBindJSON(&json); err != nil {
		cont.JSON(http.StatusOK, CreateJsonError("BindJsonError", "bind json body error", err.Error()))
		return
	}
	for k, v := range json {
		switch k {
		case "frozen":
			targetuser.Frozen = util.JsonToBool(v)
		case "verified":
			targetuser.Verified = util.JsonToBool(v)
		case "op_v_user":
			targetuser.Op_v_user = util.JsonToBool(v)
		case "op_v_quest":
			if user.Op_v_quest {
				targetuser.Op_v_quest = util.JsonToBool(v)
			}
		}
	}
	if err := updateUserById(targetuser); err != nil {
		cont.JSON(http.StatusOK, CreateJsonError("SqlUpdateError", "update user data error", err.Error()))
		return
	}
	cont.JSON(http.StatusOK, gin.H{ "status": "ok" })
}

func (userPageSrc)MyInfoGetPage(cont *gin.Context){
	suuid := updateClientUuid(cont)
	user := getLoginUser(suuid)
	if user == nil {
		cont.JSON(http.StatusOK, CreateJsonError("UserNotLogin", "user is not login", ""))
		return
	}
	var response gin.H = gin.H{
		"userid": user.Id,
		"username": user.Username,
	}
	switch cont.Param("mode") {
	case "info":
	case "auth":
		response["auths"] = gin.H{
			"v_user": user.Op_v_user,
			"v_quest": user.Op_v_quest,
		}
	case "children":
		lines, err := sqlUserTable.SqlSearch(ksql.TypeMap{"id": ksql.TYPE_Uint32}, ksql.WhereMap{{"manager", "=", user.Id, ""}}, 0)
		if err != nil {
			cont.JSON(http.StatusOK, CreateJsonError("SqlSelectError", "search user error", err.Error()))
		}
		children := make([]uint32, 0, len(lines))
		for _, l := range lines {
			children = append(children, util.JsonToUint32(l["id"]))
		}
		response["children"] = children
	default:
		cont.JSON(http.StatusNotFound, CreateJsonError("NoModelException", "don't have the model's partten", ""))
		return
	}
	cont.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data": response,
	})
}

func (userPageSrc)SettingPage(cont *gin.Context){
	cont.HTML(http.StatusOK, "user/setting/setting.html", gin.H{
	})
}

func (userPageSrc)SetDataPage(cont *gin.Context){
	suuid := updateClientUuid(cont)
	user := getLoginUser(suuid)
	if user == nil {
		cont.JSON(http.StatusOK, CreateJsonError("UserNotLogin", "user is not login", ""))
		return
	}

	switch cont.Param("mode") {
	case "head":
		imgdata := cont.PostForm("imgdata")
		ok := setUserB64Head(imgdata, user)
		if !ok {
			cont.JSON(http.StatusOK, CreateJsonError("SetUserException", "set user head error", ""))
			return
		}
		cont.JSON(http.StatusOK, gin.H{ "status": "ok" })
	case "password":
		oldpwd := cont.PostForm("oldpwd")
		user.Password = oldpwd
		if !reg_pwd.MatchString(oldpwd) || hashUserPwd(user) != user.ShaPwd {
			cont.JSON(http.StatusOK, CreateJsonError("IllegalArgumentException", "the password is wrong", ""))
			return
		}
		newpwd := cont.PostForm("newpwd")
		if !reg_pwd.MatchString(newpwd) {
			cont.JSON(http.StatusOK, CreateJsonError("IllegalArgumentException", "the new password is illegal data", ""))
			return
		}
		user.Password = newpwd
		err := setUserPwd(user)
		if err != nil {
			cont.JSON(http.StatusOK, CreateJsonError("SetUserException", "set user password error", err.Error()))
			return
		}
		user.Password = ""
		cont.JSON(http.StatusOK, gin.H{ "status": "ok" })
	default:
		cont.JSON(http.StatusNotFound, CreateJsonError("NoModelException", "don't have the model's partten", ""))
	}
}

func (page userPageSrc)Init(){
	userGroup := engine.Group("user");{
		userGroup.GET("/login", page.LoginGetPage)
		userGroup.POST("/login", page.LoginPostPage)
		userGroup.POST("/logout", page.LogoutPostPage)
		userGroup.GET("/register", page.RegisterGetPage)
		userGroup.POST("/register", page.RegisterPostPage)
		userGroup.GET("/captcha/:mode", page.CaptchaGetPage)
		userGroup.GET("/info/:id", page.InfoGetPage)
		userGroup.GET("/info/:id/:mode", page.InfoResPage)
		userGroup.POST("/info/:id/set", page.InfoSetPage)
		userGroup.GET("/myinfo/:mode", page.MyInfoGetPage)
		userGroup.GET("/setting", page.SettingPage)
		userGroup.POST("/setting/set/:mode", page.SetDataPage)
	}
}


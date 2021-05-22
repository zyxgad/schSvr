
package kpnmwebpage

import (
	http "net/http"

	gin  "github.com/gin-gonic/gin"
	util "github.com/zyxgad/go-util/util"
	jwt  "github.com/zyxgad/go-util/jwt"
	kses "github.com/zyxgad/schSvr/handles/sql/session"
	ksql "github.com/zyxgad/schSvr/handles/sql"
)


type questPageSrc int

func (questPageSrc)putGetPage(cont *gin.Context){
	cont.HTML(http.StatusOK, "quest/put.html", gin.H{
	})
}

func (questPageSrc)putPostPage(cont *gin.Context){
	suuid := updateClientUuid(cont)
	user := getLoginUser(suuid)
	if user == nil {
		cont.JSON(http.StatusOK, CreateJsonError("UserNotLogin", "user is not login", ""))
		return
	}
	quest := util.CleanStr(cont.PostForm("question"))
	answer := util.CleanStr(cont.PostForm("answer"))
	if len(quest) == 0 || len(answer) == 0 {
		cont.JSON(http.StatusOK, CreateJsonError("IllegalArgumentException", "question/answer is illegal data", ""))
		return
	}
	sqlQuestTable.SqlInsert(ksql.Map{
		"quest": quest,
		"answer": answer,
		"owner": user.Id,
		"verified": user.Op_v_quest,
	})
	cont.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (questPageSrc)listGetPage(cont *gin.Context){
	suuid := updateClientUuid(cont)
	user := getLoginUser(suuid)
	if user == nil {
		cont.JSON(http.StatusOK, CreateJsonError("UserNotLogin", "user is not login", ""))
		return
	}
	if !user.Op_v_quest {
		cont.JSON(http.StatusOK, CreateJsonError("UserNoPermission", "user no permission", ""))
		return
	}

	lines, err := sqlQuestTable.SqlSearch(ksql.TypeMap{
		"id": ksql.TYPE_Uint32,
		"quest": ksql.TYPE_String,
		"answer": ksql.TYPE_String,
		"owner": ksql.TYPE_Uint32,
		"verified": ksql.TYPE_Bool,
	}, ksql.WhereMap{}, 0)
	if err != nil {
		cont.JSON(http.StatusOK, CreateJsonError("SqlSelectError", "sql select error", err.Error()))
		return
	}
	var data []gin.H = make([]gin.H, 0, len(lines))
	for _, l := range lines {
		data = append(data, gin.H{
			"id": util.JsonToUint32(l["id"]),
			"quest": util.JsonToString(l["quest"]),
			"answer": util.JsonToString(l["answer"]),
			"owner": util.JsonToUint32(l["owner"]),
			"verified": util.JsonToBool(l["verified"]),
		})
	}
	cont.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data": data,
	})
}

func (questPageSrc)setPostPage(cont *gin.Context){
	suuid := updateClientUuid(cont)
	user := getLoginUser(suuid)
	if user == nil {
		cont.JSON(http.StatusOK, CreateJsonError("UserNotLogin", "user is not login", ""))
		return
	}
	if !user.Op_v_quest {
		cont.JSON(http.StatusOK, CreateJsonError("UserNoPermission", "user no permission", ""))
		return
	}
	qid := (uint32)(util.StrToInt(cont.Param("id"), 10))
	if !sqlQuestTable.HasData(ksql.WhereMap{{"id", "=", qid, ""}}) {
		cont.JSON(http.StatusOK, CreateJsonError("NoQuestionFound", "no question id", ""))
		return
	}

	var json map[string]interface{}
	if err := cont.ShouldBindJSON(&json); err != nil {
		cont.JSON(http.StatusOK, CreateJsonError("BindJsonError", "bind json body error", err.Error()))
		return
	}
	for k, v := range json {
		switch k {
		case "question":
			sqlQuestTable.SqlUpdate(ksql.Map{"question": util.JsonToString(v)}, ksql.WhereMap{{"id", "=", qid, ""}})
		case "answer":
			sqlQuestTable.SqlUpdate(ksql.Map{"answer": util.JsonToString(v)}, ksql.WhereMap{{"id", "=", qid, ""}})
		case "verified":
			sqlQuestTable.SqlUpdate(ksql.Map{"verified": util.JsonToBool(v)}, ksql.WhereMap{{"id", "=", qid, ""}})
		}
	}
	cont.JSON(http.StatusOK, gin.H{ "status": "ok" })
}

func (questPageSrc)searchPostPage(cont *gin.Context){
	suuid := updateClientUuid(cont)
	user := getLoginUser(suuid)
	if user == nil {
		cont.JSON(http.StatusOK, CreateJsonError("UserNotLogin", "user is not login", ""))
		return
	}
	qid := (uint32)(util.StrToInt(cont.PostForm("id"), 10))
	lines, err := sqlQuestTable.SqlSearch(ksql.TypeMap{
		"quest": ksql.TYPE_String,
		"owner": ksql.TYPE_Uint32,
		"verified": ksql.TYPE_Bool,
	}, ksql.WhereMap{{"id", "=", qid, ""}}, 1)
	if err != nil || len(lines) != 1 {
		cont.JSON(http.StatusOK, CreateJsonError("LineNotFound", "can not found question id", ""))
		return
	}
	line := lines[0]
	data := gin.H{
		"quest": util.JsonToString(line["quest"]),
		"owner": util.JsonToUint32(line["owner"]),
		"verified": util.JsonToBool(line["verified"]),
	}
	cont.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data": data,
	})
}

func (questPageSrc)matchGetPage(cont *gin.Context){
	cont.HTML(http.StatusOK, "quest/match.html", gin.H{
	})
}

func (questPageSrc)matchPostPage(cont *gin.Context){
	suuid := updateClientUuid(cont)
	user := getLoginUser(suuid)
	if user == nil {
		cont.JSON(http.StatusOK, CreateJsonError("UserNotLogin", "the user is not login", ""))
		return
	}
	switch cont.Param("mode") {
	case "prepare":
		var (
			err error
			lines []ksql.Map
		)
		if tk, err := cont.Cookie("matchid"); err == nil && len(tk) != 0 {
			cont.JSON(http.StatusOK, CreateJsonError("UserPrepared", "user is prepared", ""))
			return
		}
		matchid := cont.PostForm("matchid")
		lines, err = sqlMatchTable.SqlSearch(
			ksql.TypeMap{ "questnum": ksql.TYPE_Uint32 }, ksql.WhereMap{{"id", "=", matchid, ""}}, 1)
		if err != nil || len(lines) != 1 {
			cont.JSON(http.StatusOK, CreateJsonError("NoMatchForUser", "there are no match for user", ""))
			return
		}
		var muserid string
		if lis, err := sqlMatchUserTable.SqlSearch(
			ksql.TypeMap{"id": ksql.TYPE_String},
			ksql.WhereMap{{"matchid", "=", matchid, "AND"}, {"userid", "=", user.Id, ""}}, 1);
			err == nil && len(lis) == 1 {
			muserid = util.JsonToString(lis[0]["id"])
		}else{
			muserid = util.UUID2Hex(util.NewUUID())
			err = sqlMatchUserTable.SqlInsert(ksql.Map{
				"id": muserid,
				"matchid": matchid,
				"userid": user.Id,
			})
			if err != nil {
				cont.JSON(http.StatusOK, CreateJsonError("PrepareError", "prepare error", err.Error()))
				return
			}
		}
		cont.SetCookie("matchid", jwtEncoder.Encode(jwt.SetOutdate(jwt.Json{
			"data": muserid,
		}, util.TimeDay)), 60 * 60 * 24, "/", "", false, true)

		kses.RemoveSession(suuid, "questid")
		kses.SetSession(&kses.SqlSessionValue{
			Uuid: suuid,
			Key: "ans_max",
			Value: util.JoinObjStr(util.JsonToUint32(lines[0]["questnum"])),
			Overtime: util.GetTimeAfter(util.TimeHour * 6),
		})

		cont.JSON(http.StatusOK, gin.H{ "status": "ok" })
	case "getnext":
		var (
			jwtoken string
			jwtdata jwt.Json
			muserid string
			questid uint32
			isout bool
			err error
		)
		jwtoken, _ = cont.Cookie("matchid")
		jwtdata, isout, err = jwtEncoder.Decode(jwtoken)
		if err != nil {
			cont.SetCookie("matchid", "", -1, "/", "", false, true)
			cont.JSON(http.StatusOK, CreateJsonError("ParseJWTError", "parse match token error", err.Error()))
			return
		}
		muserid = util.JsonToString(jwtdata["data"])
		if isout {
			cont.SetCookie("matchid", jwtEncoder.Encode(jwt.SetOutdate(jwt.Json{
				"data": muserid,
			}, util.TimeDay)), 60 * 60 * 24, "/", "", false, true)
		}

		ans_count := sqlMatchAnswerTable.DataCount(ksql.WhereMap{{"muserid", "=", muserid, ""}})
		ans_max := (uint64)(util.StrToInt(kses.GetSession(suuid, "ans_max").Value, 10))

		qidv := kses.GetSession(suuid, "questid")
		answer := util.CleanStr(cont.PostForm("answer"))

		if ans_count < ans_max && qidv != nil {
			if len(answer) == 0 {
				cont.JSON(http.StatusOK, CreateJsonError("IllegalDataException", "answer can't be empty", ""))
				return
			}
			var score int = 0
			questid = (uint32)(util.StrToInt(qidv.Value, 10))
			{
				lines, err := sqlQuestTable.SqlSearch(
					ksql.TypeMap{"answer": ksql.TYPE_String},
					ksql.WhereMap{{"id", "=", questid, "AND"}, {"verified", "=", true, ""}}, 1)
				if err != nil || len(lines) != 1 {
					cont.JSON(http.StatusOK, CreateJsonError("SqlSelectError", "get question answer error", ""))
					return
				}
				std_answer := util.JsonToString(lines[0]["answer"])
				if answer == std_answer {
					score = 1
				}else{
					score = 0
				}
			}
			err := sqlMatchAnswerTable.SqlInsert(ksql.Map{
				"muserid": muserid,
				"questid": questid,
				"answer": answer,
				"score": score,
			})
			if err != nil {
				cont.JSON(http.StatusOK, CreateJsonError("SqlInsertError", "insert answer error", err.Error()))
				return
			}
			ans_count++
		}

		if ans_count >= ans_max {
			kses.RemoveSession(suuid, "ans_max")
			kses.RemoveSession(suuid, "questid")
			cont.JSON(http.StatusOK, CreateJsonError("FinishedMatch", "finish the match", ""))
			return
		}

		{
			questcount := sqlQuestTable.DataCount(ksql.WhereMap{{"verified", "=", true, ""}})
			if questcount == 0 {
				cont.JSON(http.StatusOK, CreateJsonError("GetQuestionError", "no question for match", ""))
				return
			}
			randindex := (uint)(util.GetRandUint64ByRange(0, questcount))
			lines, err := sqlQuestTable.SqlSearchOff(
				ksql.TypeMap{"id": ksql.TYPE_Uint32}, ksql.WhereMap{{"verified", "=", true, ""}}, randindex, 1)
			if err != nil || len(lines) != 1 {
				cont.JSON(http.StatusOK, CreateJsonError("SqlSelectError", "get question error", ""))
				return
			}
			questid = util.JsonToUint32(lines[0]["id"])
		}
		kses.SetSession(&kses.SqlSessionValue{
			Uuid: suuid,
			Key: "questid",
			Value: util.JoinObjStr(questid),
			Overtime: util.GetTimeAfter(util.TimeMin * 60),
		})

		cont.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"questid": questid,
			"count": ans_count + 1,
			"max_count": ans_max,
		})
	case "submit":
		var (
			jwtoken string
			jwtdata jwt.Json
			muserid string
			err error
		)
		jwtoken, _ = cont.Cookie("matchid")
		cont.SetCookie("matchid", "", -1, "/", "", false, true)
		jwtdata, _, err = jwtEncoder.Decode(jwtoken)
		if err != nil {
			cont.JSON(http.StatusOK, CreateJsonError("GetMatchIdError", "parse match token error", err.Error()))
			return
		}

		kses.RemoveSession(suuid, "ans_max")
		kses.RemoveSession(suuid, "questid")

		muserid = util.JsonToString(jwtdata["data"])
		cont.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"muserid": muserid,
		})
	case "check":
		muserid := cont.PostForm("muserid")
		lines, err := sqlMatchAnswerTable.SqlSearch(ksql.TypeMap{
			"questid": ksql.TYPE_Uint32,
			"answer": ksql.TYPE_String,
			"score": ksql.TYPE_Int32,
		}, ksql.WhereMap{{"muserid", "=", muserid, ""}}, 0)
		if err != nil {
			cont.JSON(http.StatusOK, CreateJsonError("SqlSelectError", "sql select error", err.Error()))
			return
		}
		var answers []gin.H = make([]gin.H, 0, len(lines))
		for _, l := range lines {
			answers = append(answers, gin.H{
				"questid": l["questid"],
				"answer": l["answer"],
				"score": l["score"],
			})
		}
		cont.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"answers": answers,
		})
	}
}

func (questPageSrc)matchcheckGetPage(cont *gin.Context){
	muserid := cont.Query("muserid")
	lines, err := sqlMatchUserTable.SqlSearch(
			ksql.TypeMap{"userid": ksql.TYPE_Uint32},
			ksql.WhereMap{{"id", "=", muserid, ""}}, 1)
	if err != nil || len(lines) != 1 {
		cont.JSON(http.StatusOK, CreateJsonError("NoMatchUserFound", "no matchid found", ""))
		return
	}
	userid := util.JsonToUint32(lines[0]["userid"])
	cont.HTML(http.StatusOK, "quest/matchcheck.html", gin.H{
		"muserid": muserid,
		"userid": userid,
	})
}

func (page questPageSrc)Init(){
	questGroup := engine.Group("quest");{
		questGroup.GET("/put", page.putGetPage)
		questGroup.POST("/put", page.putPostPage)
		questGroup.GET("/list", page.listGetPage)
		questGroup.POST("/info/:id/set", page.setPostPage)
		questGroup.POST("/search", page.searchPostPage)
		questGroup.GET("/match", page.matchGetPage)
		questGroup.POST("/match/:mode", page.matchPostPage)
		questGroup.GET("/matchcheck", page.matchcheckGetPage)
	}
}


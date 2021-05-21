

function nextquestion(){
	const ANSWER_INPUT = $("#answer-text");
	const QUEST_BODY = $("#question-body-text");

	const answer = ANSWER_INPUT.val();
	$.ajax({
		url: "/web/quest/match/getnext",
		type: "POST",
		data: {
			answer: answer
		},
		success: function(res){
			if(res.status === "ok"){
				ANSWER_INPUT.val("");
				$("#quest-count").text(res.count + "/" + res.max_count);
				$.ajax({
					url: "/web/quest/search",
					type: "POST",
					data: {
						id: res.questid
					},
					success: function(res){
						if(res.status === "ok"){
							const questdata = res.data;
							QUEST_BODY.html("");
							questdata.quest.split("\n").forEach((item)=>{
								QUEST_BODY.append($(`<div class="question-body-text-line"></div>`).text(item))
							});
							$.ajax({
								url: "/web/user/info/" + questdata.owner,
								type: "POST",
								success: function(res){
									if(res.status === "ok"){
										const userdata = res.data;
										$("#question-body-info-owner").text(userdata.username);
										return;
									}
								}
							})
							return;
						}
						console.log("error res:", res);
					}
				})
				return;
			}
			if(res.error !== undefined){
				switch(res.error){
					case "ParseJWTError":
						alert("获取问题失败: " + res.errorMessage);
						window.location = "/web/quest/match";
						break;
					case "FinishedMatch":
						console.log("finished");
						$.ajax({
							url: "/web/quest/match/submit",
							type: "POST",
							success: function(res){
								if(res.status === "ok"){
									alert("答卷提交完毕");
									window.location = "/web/quest/matchcheck?muserid=" + res.muserid;
									return;
								}
								if(res.error !== undefined){
									alert("提交答卷错误: " + res.errorMessage);
									return;
								}
							}
						});
						return;
				}
				console.log("error:", res);
				return;
			}
			console.log("res:", res);
		}
	});
}

$(document).ready(function(){
	$.ajax({
		url: "/web/user/info",
		type: "POST",
		success: function(res){
			if(res.status !== "ok"){
				alert("您需要先去登录");
				window.location = "/web/user/login";
				return;
			}
		}
	});
	$("#start_page").show();
	$("#start-btn").click(function(){
		const matchid = $("#matchid-input").val();
		if(matchid.length === 0){
			alert("match id error");
			return;
		}
		$.ajax({
			url: "/web/quest/match/prepare",
			type: "POST",
			data: {
				matchid: matchid
			},
			success: function(res){
				if(res.status === "ok"){
					console.log("prepare '" + matchid + "' successed");
					$("#start_page").hide();
					$("#match_page").show();
					nextquestion();
					return;
				}
				if(res.error !== undefined){
					switch(res.error){
						case "UserNotLogin":
							alert("您需要先去登录");
							window.location = "/web/user/login";
							return;
						case "UserPrepared":
							$("#start_page").hide();
							$("#match_page").show();
							nextquestion();
							return;
						case "NoMatchForUser":
							alert("没有匹配的竞赛");
							break;
						default:
							console.log("error:", res);
					}
					return;
				}
				console.log("res:", res);
			}
		});
	});
	$("#submit-btn").click(nextquestion);
	$("#answer-text").keydown(function(event){
		if(event.key == "Enter" || event.keyCode == 13){
			event.preventDefault();
			$("#submit-btn").click();
		}
	});
})
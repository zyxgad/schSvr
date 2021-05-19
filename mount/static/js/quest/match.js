

var questionID = 0;

function nextquestion(){
	const answer = $("#answer-text").val();
	$.ajax({
		url: "/web/quest/match/getnext",
		type: "POST",
		data: {
			answer: answer
		},
		success: function(res){
			if(res.status === "ok"){
				questionID = res.questid;
				$("#quest-count").text(res.count + "/" + res.max_count);
				$.ajax({
					url: "/web/quest/search",
					type: "POST",
					data: {
						id: questionID
					},
					success: function(res){
						if(res.status === "ok"){
							let questdata = res.data;
							console.log("res data:", questdata);
							let qlines = questdata["quest"].split("\n")
							const qnode = $("#question-text");
							qnode.html("");
							for(let i = 0;i < qlines.length;i++){
								qnode.append("<p>" + qlines[i] + "</p>")
							}
							// questdata.owner;
							return;
						}
						console.log("res:", res);
					}
				})
				console.log("next question id:", questionID);
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
})
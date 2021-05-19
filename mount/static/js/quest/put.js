

$(document).ready(function(){
	$("#subject-submit-btn").click(function(){
		var quest = $("#subject-question-text").val();
		var answer = $("#subject-answer-text").val();
		if(quest.length === 0){
			alert("问题不能为空");
			return;
		}
		if(answer.length === 0){
			alert("答案不能为空");
			return;
		}
		$.ajax({
			url: "/web/quest/put",
			type: "POST",
			data: {
				question: quest,
				answer: answer
			},
			success: function(res){
				console.log(res);
				if(res.status === "ok"){
					alert("等待审核中");
					$("#subject-question-text").val("");
					$("#subject-answer-text").val("");
					return;
				}
				if(res.error !== undefined){
					alert("错误:" + res.errorMessage);
					return;
				}
			}
		});
	});

	// $('#subject-question-text').keydown(function(event){
	// 	if(event.key == "Enter" || event.keyCode == 13){
	// 		event.preventDefault();
	// 	}
	// });
	$('#subject-submit-btn').keydown(function(event){if(event.key == "Enter" || event.keyCode == 13){ event.preventDefault(); }});
	$('#subject').keydown(function(event){
		if(event.key == "Enter" || event.keyCode == 13){
			if(event.shiftKey){
				return;
			}
			$('#subject-submit-btn').click();
		}
	});

});




function appendAnswer(item){
	const node = $(`
<div class="match-answers-item">
	<span class="match-answers-item-id"></span>
	<span class="match-answers-item-quest"></span>
	<span class="match-answers-item-answer"></span>
	<span class="match-answers-item-score"></span>
</div>`);
	node.children(".match-answers-item-id:first").text(item.id);
	node.children(".match-answers-item-answer:first").text(item.answer);
	node.children(".match-answers-item-score:first").text(item.score);
	const quest_node = node.children(".match-answers-item-quest:first");
	$.ajax({
		url: "/web/quest/search",
		type: "POST",
		data: {
			id: item.questid
		},
		success: function(res){
			if(res.status === "ok"){
				quest_node.html("");
				res.data["quest"].split("\n").forEach((item)=>{
					quest_node.append("<p>" + item + "</p>")
				});
				return;
			}
			console.log("res:", res);
		}
	});
	$("#match-answers-box").append(node);
}

$(document).ready(function(){
	$.ajax({
		url: "/web/quest/match/check",
		type: "POST",
		data: {
			muserid: muserid
		},
		success: function(res){
			if(res.status === "ok"){
				console.log("anss:", res.answers);
				$("#match-answers-box").html("");
				var score_max = 0;
				var score_count = 0;
				res.answers.forEach((item, ind)=>{
					score_max += 1;
					score_count += item.score;
					appendAnswer({
						id: ind + 1,
						questid: item.questid,
						answer: item.answer,
						score: item.score
					});
				});
				$("#match-info-score").text(score_count + "/" + score_max)
				return;
			}
			console.log("error res:", res);
		}
	});
});

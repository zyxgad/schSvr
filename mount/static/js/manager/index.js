

function makeusercheckbox(disabled, userdata, key, title){
	var vcheckbox = $(`<input type="checkbox" />`);
	vcheckbox.prop("disabled", disabled)
	vcheckbox.prop("checked", userdata[key]);
	vcheckbox.click(function(event){
		var checked = vcheckbox.prop("checked");
		event.preventDefault();
		if(!confirm(`确定将${userdata.username}的'${title}'设为${checked}吗?`)){
			return;
		}
		$.ajax({
			async: true,
			url: "/web/user/info/" + userdata.userid + "/set",
			type: "POST",
			data: `{"${key}":${checked}}`,
			success: function(res){
				if(res.status === "ok"){
					vcheckbox.prop("checked", checked);
					return;
				}
				console.log(title, "set error:", res);
			}
		});
	});
	return vcheckbox;
}

function makequestcheckbox(disabled, questdata, key, title){
	var vcheckbox = $(`<input type="checkbox" />`);
	vcheckbox.prop("disabled", disabled)
	vcheckbox.prop("checked", questdata[key]);
	vcheckbox.click(function(event){
		var checked = vcheckbox.prop("checked");
		event.preventDefault();
		if(!confirm(`确定将问题${questdata.id}的'${title}'设为${checked}吗?`)){
			return;
		}
		$.ajax({
			async: true,
			url: "/web/quest/info/" + questdata.id + "/set",
			type: "POST",
			data: `{"${key}":${checked}}`,
			success: function(res){
				if(res.status === "ok"){
					vcheckbox.prop("checked", checked);
					return;
				}
				console.log(title, "set error:", res);
			}
		});
	});
	return vcheckbox;
}

function flushUser(udata){
	const userauths = udata.auths;
	$("#user-verifies").show();
	$("#user-verifies-body").html("");
	$.ajax({
		async: true,
		url: "/web/user/myinfo/children",
		type: "GET",
		success: function(res){
			if(res.status !== "ok"){
				console.log("error res:", res);
				return;
			}
			res.data.children.sort().forEach((cid)=>{
				var uvitem = $(`<div class="user-verifies-item">
	<span class="user-verifies-item-id"></span>
	<span class="user-verifies-item-name"></span>
	<span class="user-verifies-item-status"></span>
	<span class="user-verifies-item-verified"></span>
	<span class="user-verifies-item-v_user"></span>
	<span class="user-verifies-item-v_quest"></span>
</div>`);
				$("#user-verifies-body").append(uvitem);
				$.ajax({
					url: "/web/user/info/" + cid + "/info",
					type: "GET",
					success: function(res){
						if(res.status !== "ok"){
							return;
						}
						const userdata = res.data;
						uvitem.children(".user-verifies-item-id:first").text(userdata.userid);
						uvitem.children(".user-verifies-item-name:first").text(userdata.username);
						uvitem.children(".user-verifies-item-status:first").text(userdata.frozen?"冻结":"正常");
						uvitem.children(".user-verifies-item-verified:first").html(makeusercheckbox(false, userdata, "verified", "已验证"));
						uvitem.children(".user-verifies-item-v_user:first").html(makeusercheckbox(!userauths["v_user"], userdata, "op_v_user", "允许验证用户"));
						uvitem.children(".user-verifies-item-v_quest:first").html(makeusercheckbox(!userauths["v_quest"], userdata, "op_v_quest", "允许审核问题"));
					}
				});
			});
		}
	});
}

function flushQuest(udata){
	const userauths = udata.auths;
	$("#quest-verifies").show();
	$("#quest-verifies-body").html("");
	$.ajax({
		async: true,
		url: "/web/quest/list",
		type: "GET",
		success: function(res){
			if(res.status !== "ok"){
				console.log("error res:", res);
				return;
			}
			res.data.sort((a, b)=>{return a.id - b.id;}).forEach((questdata)=>{
				var uvitem = $(`<div class="quest-verifies-item">
	<span class="quest-verifies-item-id"></span>
	<span class="quest-verifies-item-user"></span>
	<span class="quest-verifies-item-quest"></span>
	<span class="quest-verifies-item-answer"></span>
	<span class="quest-verifies-item-verified"></span>
</div>`);
				$("#quest-verifies-body").append(uvitem);
				uvitem.children(".quest-verifies-item-id:first").text(questdata.id);
				$.ajax({
					url: "/web/user/info/" + questdata.owner + "/info",
					type: "GET",
					success: function(res){
						if(res.status === "ok"){
							uvitem.children(".quest-verifies-item-user:first").text(res.data.username);
							return;
						}
					}
				})
				let quest_node = uvitem.children(".quest-verifies-item-quest:first");
				quest_node.html("");
				questdata.quest.split("\n").forEach((item)=>{
					quest_node.append($(`<p></p>`).text(item));
				});
				uvitem.children(".quest-verifies-item-answer:first").text(questdata.answer);
				uvitem.children(".quest-verifies-item-verified:first").html(makequestcheckbox(!userauths["v_quest"], questdata, "verified", "已验证"));
			});
		}
	});
}

$(document).ready(function(){
	var userdata = null;
	var userauths = null;
	$.ajax({
		url: "/web/user/myinfo/auth",
		type: "GET",
		success: function(res){
			if(res.status === "ok"){
				userdata = res.data;
				userauths = userdata.auths;
				$("#user-info-name").text(userdata.username);
				Object.keys(userauths).forEach((key)=>{
					$("#user-info-auth").append(
						$(`<span class="user-info-auth-line"></span>`).text(key + ":" + userauths[key]));
				})
				if(userauths["v_user"]){
					flushUser(userdata);
				}
				if(userauths["v_quest"]){
					flushQuest(userdata);
				}
				return;
			}
		}
	});
	$("#user-verifies-flush-btn").click(function(){
		flushUser(userdata);
	})
	$("#quest-verifies-flush-btn").click(function(){
		flushQuest(userdata);
	})
});
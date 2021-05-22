

function makecheckbox(disabled, userdata, key, title){
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

$(document).ready(function(){
	$.ajax({
		url: "/web/user/myinfo/auth",
		type: "GET",
		success: function(res){
			if(res.status === "ok"){
				const userdata = res.data;
				const userauths = userdata.auths;
				$("#user-info-name").text(userdata.username);
				Object.keys(userauths).forEach((key)=>{
					$("#user-info-auth").append(
						$(`<span class="user-info-auth-line"></span>`).text(key + ":" + userauths[key]));
				})
				if(userauths["v_user"]){
					$("#user-verifies").show();
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
										uvitem.children(".user-verifies-item-verified:first").html(makecheckbox(false, userdata, "verified", "已验证"));
										uvitem.children(".user-verifies-item-v_user:first").html(makecheckbox(!userauths["v_user"], userdata, "op_v_user", "允许验证用户"));
										uvitem.children(".user-verifies-item-v_quest:first").html(makecheckbox(!userauths["v_quest"], userdata, "op_v_quest", "允许审核问题"));
									}
								});
							});
						}
					});
				}
				return;
			}
		}
	});
});
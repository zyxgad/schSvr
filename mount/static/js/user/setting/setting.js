
const pwdreg  = /^[A-Za-z][0-9A-Za-z_-]{7,127}$/;

var userdata = null;

function updateUserdata(){
	if(!userdata){
		$("#user-info").hide();
		return;
	}
	$("#user-info-head-img").prop("src", "/web/user/info/" + userdata.userid + "/head");
	$("#user-info-name-cont").text(userdata.username);
	$("#user-info").show();
}

$(document).ready(function(){

	$("#user-info-head-choose-file").bind("change", function(event){
		if(this.files.length == 0){
			return;
		}
		const file = this.files[0];
		const fsize = file.size;
		if(fsize > 1024 * 32){
			alert("头像大小不能大于32Kb");
			return;
		}
		var reader = new FileReader();
		reader.onload = function(){
			var imgdata = this.result;
			$("#user-info-head-img").prop("src", imgdata);
			$("#user-info-head-sure").show();
		}
		reader.readAsDataURL(file);
	});
	$("#user-info-head-sure").click(function(){
		const imgdata = $("#user-info-head-img").prop("src");
		$.ajax({
			url: "/web/user/setting/set/head",
			type: "POST",
			data: {
				"imgdata": imgdata
			},
			success: function(res){
				console.log("res:", res);
			}
		});
		$(this).hide();
	});

	$("#user-info-head-choose-btn").click(function(){
		$("#user-info-head-choose-file").click();
	});
	$("#user-info-changepwd-button").click(function(){
		const oldpwd = $("#user-info-changepwd-originalpwd").val();
		const newpwd = $("#user-info-changepwd-newpwd").val();
		const newpwd2 = $("#user-info-changepwd-repeatnewpwd").val();
		if(!pwdreg.test(oldpwd) || !pwdreg.test(newpwd) || !pwdreg.test(newpwd2)){
			alert("不是合法密码");
			return;
		}
		if(newpwd !== newpwd2){
			alert("两次密码不一致");
			return;
		}
		$.ajax({
			url: "/web/user/setting/set/password",
			type: "POST",
			data: {
				"oldpwd": oldpwd,
				"newpwd": newpwd
			},
			success: function(res){
				if(res.error !== undefined){
					alert("设置密码错误: " + res.errorMessage);
					return
				}
				if(res.status === "ok"){
					alert("密码修改成功");
					return;
				}
			}
		});
	});

	$.ajax({
		url: "/web/user/info",
		type: "POST",
		success: function(res){
			if(res.error !== undefined){
				alert("您还没有登录，请先登录");
				window.location = "/web/user/login";
				return;
			}
			if(res.status === "ok" && res.data !== undefined){
				userdata = res.data;
				updateUserdata();
				return;
			}
		}
	});
})
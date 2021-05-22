
var userdata = null;

function updateUserdata(){
	if(!userdata){
		$("#user-info").hide();
		return;
	}
	$("#user-info-name-cont").text(userdata.username);
	$("#user-info-head-img").prop("src", userdata.userhead);
	$("#user-info").show();
}

$(document).ready(function(){
	$.ajax({
		url: "/web/user/info/" + userid + "/info",
		type: "GET",
		success: function(res){
			if(res.error !== undefined){
				console.log("error:", res);
				$("#error-div").text(res.errorMessage);
				$("#error-div").show();
				return;
			}
			if(res.status === "ok" && res.data !== undefined){
				userdata = res.data;
				updateUserdata();
				return;
			}
		}
	})
})
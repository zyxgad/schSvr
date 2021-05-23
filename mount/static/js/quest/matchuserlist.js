
$(document).ready(function(){
	$.ajax({
		url: "/web/quest/matchlist/" + matchid,
		type: "GET",
		success: function(res){
			if(res.status === "ok"){
				res.data.forEach((item)=>{
					const uvitem = $(`<div class="muser-list-item">
	<span class="muser-list-item-user"></span>
	<span class="muser-list-item-matchid"></span>
</div>`);
					$("#muser-list-body").append(uvitem);
					$.ajax({
						url: "/web/user/info/" + item.userid + "/info",
						type: "GET",
						success: function(res){
							if(res.status === "ok"){
								uvitem.children(".muser-list-item-user:first").text(res.data.username);
								return;
							}
						}
					});
					uvitem.children(".muser-list-item-matchid:first").html(
						$(`<a href="/web/quest/matchcheck?muserid=${item.id}"></a>`).text(item.id));
				});
				return;
			}
			console.log("error res:", res);
		}
	})
});

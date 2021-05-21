
const EMPTY_IMG = 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAABhGlDQ1BJQ0MgcHJvZmlsZQAAKJF9kT1Iw0AcxV9TRZGKgxVEO2SoThZERRy1CkWoEGqFVh1MLv2CJg1Jiouj4Fpw8GOx6uDirKuDqyAIfoC4uTkpukiJ/0sKLWI8OO7Hu3uPu3eAUC8zzeoYBzTdNlOJuJjJropdrwhiGCEMICIzy5iTpCR8x9c9Any9i/Es/3N/jl41ZzEgIBLPMsO0iTeIpzdtg/M+cZgVZZX4nHjMpAsSP3Jd8fiNc8FlgWeGzXRqnjhMLBbaWGljVjQ14iniqKrplC9kPFY5b3HWylXWvCd/YSinryxznWYECSxiCRJEKKiihDJsxGjVSbGQov24j3/I9UvkUshVAiPHAirQILt+8D/43a2Vn5zwkkJxoPPFcT5GgK5doFFznO9jx2mcAMFn4Epv+St1YOaT9FpLix4BfdvAxXVLU/aAyx1g8MmQTdmVgjSFfB54P6NvygL9t0DPmtdbcx+nD0CaukreAAeHwGiBstd93t3d3tu/Z5r9/QBQKXKZByo6SAAAAAZiS0dEAP8A/wD/oL2nkwAAAAlwSFlzAAAuIwAALiMBeKU/dgAAAAd0SU1FB+UDGwwWKoTcD7oAAAAZdEVYdENvbW1lbnQAQ3JlYXRlZCB3aXRoIEdJTVBXgQ4XAAAAC0lEQVQI12NgAAIAAAUAAeImBZsAAAAASUVORK5CYII=';

const SERVER_DATA = {
	_userdata: null,
	set userdata(val){
		if(!val){
			$("#header-user-info-name").text("未登录");
			$("#header-user-info-head img:first").hide();
			$("#header-user-desc .true:first").hide();
			$("#header-user-desc .false:first").show();
			this._userdata = null;
			return;
		}
		$("#header-user-info-name").text(val.username);
		$("#header-user-info-head img:first").prop("src", "/web/user/info/" + val.userid + "/head");
		$("#header-user-info-head img:first").show();
		$("#header-user-desc .false:first").hide();
		$("#header-user-desc .true:first").show();
		this._userdata = val;
	},
	get userdata(){
		return this._userdata;
	}
}

$(document).ready(function(){
	SERVER_DATA.userdata = null;

	var userdescDOM = $("#header-user-desc");
	userdescDOM.hide();
	$("#header-user-info").click(function(){
		userdescDOM.fadeToggle(100);
	});
	$("#header-user-desc-closemount").click(function(){
		userdescDOM.fadeOut(100);
	});
	$("#header-user-info-head img:first").bind("error", function(){
		this.src = EMPTY_IMG;
	});
	$("#header-user-desc-logoutbtn").click(function(){
		if(!confirm("Are you sure to log out?")){
			return;
		}
		$.ajax({
			url: "/web/user/logout",
			type: "POST",
			success: function(res){
				// SERVER_DATA.userdata = null;
				if(res.status === "ok"){
					window.location = "";
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
				return;
			}
			if(res.status === "ok" && res.data !== undefined){
				SERVER_DATA.userdata = res.data;
				return;
			}
		}
	});
});

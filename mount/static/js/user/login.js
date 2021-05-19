
const namereg  = /^[A-Za-z_-][0-9A-Za-z_-]{3,31}$/;
const pwdreg  = /^[A-Za-z][0-9A-Za-z_-]{7,127}$/;

function shakeBox(box){
	box.removeClass("animation-shake-time3");
	setTimeout(()=>{box.addClass("animation-shake-time3");}, 10);
}

function updateCaptcha(){
	$.ajax({
		url: "/web/user/captcha/image",
		type: "GET",
		success: function(res){
			if(res.status === "ok" && res.data !== undefined){
				$('#login-box-captcha-img').prop("src", res.data);
			}
		}
	});
}

$(document).ready(function(){
	function checkUsernameVal(){
		const username = $('#login-box-username-input').val();
		const errbox = $('#login-box-username-error');
		if(!namereg.test(username)){
			errbox.show();
			errbox.text("不是合法用户名");
			shakeBox(errbox);
			return false;
		}
		errbox.hide();
		errbox.text("");
		return true;
	}
	function checkPasswordVal(){
		const password = $('#login-box-password-input').val();
		const errbox = $('#login-box-password-error');
		if(!pwdreg.test(password)){
			errbox.show();
			errbox.text("不是合法密码");
			shakeBox(errbox);
			return false;
		}
		errbox.hide();
		errbox.text("");
		return true;
	}
	function checkCaptchaVal(){
		const captcode = $('#login-box-captcha-input').val();
		const errbox = $('#login-box-captcha-error');
		if(captcode === ""){
			errbox.text("请输入验证码");
			errbox.show();
			shakeBox(errbox);
			return false;
		}
		errbox.hide();
		errbox.text("");
		return true;
	}

	updateCaptcha();
	$('#login-box-username-input').blur(checkUsernameVal);
	$('#login-box-password-input').blur(checkPasswordVal);
	$('#login-box-captcha-input').blur(checkCaptchaVal);
	$('#login-box-captcha-img').click(updateCaptcha);
	$('#login-box-submit').keydown(function(event){ event.preventDefault(); })
	$('#login-box-submit').click(function(){
		let okname = checkUsernameVal(), okpwd = checkPasswordVal(), okcapt = checkCaptchaVal();
		if(!okname || !okpwd || !okcapt){
			return
		}
		const username = $('#login-box-username-input').val();
		const password = $('#login-box-password-input').val();
		const captcode = $('#login-box-captcha-input').val();
		const errorbox = $('#login-box-error');
		$.ajax({
			url: "/web/user/login",
			type: "POST",
			data: {
				username: username,
				password: password,
				captcode: captcode
			},
			success: function(res){
				if(res.status === "ok"){
					console.log("login success");
					errorbox.hide();
					errorbox.text("");
					window.location = "/web/redirectto?url=" + "/web";
					return;
				}
				if(res.error !== undefined){
					switch(res.error){
						case "CaptchaError":
							let errbox = $('#login-box-captcha-error');
							errbox.text("验证码错误");
							errbox.show();
							shakeBox(errbox);
							break;
						case "UserNotVerifyException":
							errorbox.text("请等待管理员审核");
							errorbox.show();
							shakeBox(errorbox);
							break;
						default:
							errorbox.text(res.errorMessage);
							errorbox.show();
							shakeBox(errorbox);
					}
					updateCaptcha();
					return;
				}
				console.log("error res:", res);
			}
		})
	});
	$('#login-box').keyup(function(event){
		if(event.key == "Enter" || event.keyCode == 13){
			$('#login-box-submit').click();
		}
	});
});

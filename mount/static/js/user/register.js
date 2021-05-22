
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
				$('#register-box-captcha-img').prop("src", res.data);
			}
		}
	})
}

$(document).ready(function(){
	function checkUsernameVal(){
		const username = $('#register-box-username-input').val();
		const errbox = $('#register-box-username-error');
		if(!namereg.test(username)){
			errbox.text("不是合法用户名");
			errbox.show();
			shakeBox(errbox);
			return false;
		}
		errbox.hide();
		errbox.text("");
		return true;
	}
	function checkPasswordVal(){
		const password = $('#register-box-password-input').val();
		const errbox = $('#register-box-password-error');
		if(!pwdreg.test(password)){
			errbox.text("不是合法密码");
			errbox.show();
			shakeBox(errbox);
			return false;
		}
		errbox.hide();
		errbox.text("");
		return true;
	}
	function checkPwdagainVal(){
		if(!checkPasswordVal()){
			return false;
		}
		const password = $('#register-box-password-input').val();
		const pwdagain = $('#register-box-pwdagain-input').val();
		const errbox = $('#register-box-pwdagain-error');
		if(password !== pwdagain){
			errbox.text("两次密码不一致");
			errbox.show();
			shakeBox(errbox);
			return false;
		}
		errbox.hide();
		errbox.text("");
		return true;
	}
	function checkManagerVal(){
		const manager = $('#register-box-manager-input').val();
		const errbox = $('#register-box-manager-error');
		if(!namereg.test(manager)){
			errbox.text("不是合法用户名");
			errbox.show();
			shakeBox(errbox);
			return false;
		}
		errbox.hide();
		errbox.text("");
		return true;
	}
	function checkCaptchaVal(){
		const captcode = $('#register-box-captcha-input').val();
		const errbox = $('#register-box-captcha-error');
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
	$('#register-box-username-input').blur(checkUsernameVal);
	$('#register-box-password-input').blur(checkPasswordVal);
	$('#register-box-pwdagain-input').blur(checkPwdagainVal);
	$('#register-box-manager-input').blur(checkManagerVal);
	$('#register-box-captcha-input').blur(checkCaptchaVal);
	$('#register-box-captcha-img').click(updateCaptcha);
	$('#register-box-submit').keydown(function(event){ event.preventDefault(); });
	$('#register-box-submit').click(function(){
		let okname = checkUsernameVal(), okpwd = checkPasswordVal(), okmgr = checkManagerVal(), okcapt = checkCaptchaVal();
		if(!okname || !okpwd || !okmgr || !okcapt){
			return;
		}
		const username = $('#register-box-username-input').val();
		const password = $('#register-box-password-input').val();
		const manager = $('#register-box-manager-input').val();
		const captcode = $('#register-box-captcha-input').val();
		const errorbox = $('#register-box-error');
		$.ajax({
			url: "/web/user/register",
			type: "POST",
			data: {
				username: username,
				password: password,
				manager: manager,
				captcode: captcode
			},
			success: function(res){
				if(res.status === "ok"){
					console.log("register success");
					errorbox.hide();
					errorbox.text("");
					alert("已向管理员提交申请");
					window.location = "/web/redirectto?url=" + "/web/user/login";
					return;
				}
				if(res.error !== undefined){
					switch(res.error){
						case "CaptchaError":
							let errbox = $('#register-box-captcha-error');
							errbox.text("验证码错误");
							errbox.show();
							shakeBox(errbox);
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
		});
	});
	$('#register-box').keyup(function(event){
		if(event.key == "Enter" || event.keyCode == 13){
			$('#register-box-submit').click();
		}
	});
});

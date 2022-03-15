package validators

import "gohub/pkg/captcha"

func ValidateCaptcha(captchaID string, captchaAnswer string, errs map[string][]string) map[string][]string {
	if ok := captcha.NewCaptcha().VerifyCaptcha(captchaID, captchaAnswer); !ok {
		errs["captcha_answer"] = append(errs["captcha_answer"], "验证码不正确")
	}

	return errs
}

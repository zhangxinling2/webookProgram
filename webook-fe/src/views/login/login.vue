<template>
  <div class="loginPage">
    <el-card
      class="loginPanel"
      :body-style="{ height: '100%', width: '100%', boxSizing: 'border-box' }"
    >
      <div class="loginPanelInner">
        <div class="logo">
          <img src="@/assets/images/logo.png" alt="" />
        </div>
        <el-divider direction="vertical" border-style="dashed" class="split" />
        <div class="loginForm">
          <div class="systemName">用户登录</div>
          <el-form
            ref="formRef"
            style="max-width: 600px"
            :model="loginForm"
            status-icon
            :rules="rules"
            size="large"
            label-width="auto"
            class="form"
          >
            <el-form-item label="账号：" prop="account">
              <el-input
                placeholder="请输入账号"
                v-model="loginForm.account"
                autocomplete="off"
                suffix-icon="UserFilled"
              />
            </el-form-item>
            <el-form-item label="密码：" prop="password">
              <el-input
                placeholder="请输入密码"
                v-model="loginForm.password"
                type="password"
                autocomplete="off"
                suffix-icon="Lock"
              />
            </el-form-item>
            <el-form-item label="验证码：" prop="captcha">
              <div style="display: flex; width: 100%">
                <div style="flex: 1">
                  <el-input
                    placeholder="请输入验证码"
                    v-model.number="loginForm.captcha"
                    suffix-icon="Calendar"
                  />
                </div>
                <div class="captchaSrc">
                  <img
                    :src="captchaContent.base64"
                    @click="changeCaptchaSrc"
                  />
                </div>
              </div>
            </el-form-item>
            <div class="registerBtn">
              <el-link type="primary" href="/register" :underline="false">
                去账户？点击注册
              </el-link>
            </div>
            <el-form-item>
              <el-button
                type="primary"
                @click="submitForm(ruleFormRef)"
                class="loginBtn"
              >
                登录
              </el-button>
              <!-- <el-button @click="resetForm(formRef)">登录</el-button> -->
            </el-form-item>
          </el-form>
        </div>
      </div>
    </el-card>
    <div class="footer">@Copyright 通用管理系统 备案信息：陕432432432</div>
  </div>
</template>

<script setup lang="ts">
import { ElLoading, ElMessage, FormInstance, FormRules } from "element-plus";
import { onMounted, reactive, ref } from "vue";
import api from "@/utils/api";
import { AxiosResponse } from "axios";

const loginForm = reactive({
  account: "",
  password: "",
  captcha: "",
});
let captchaContent = reactive({
  captchaId: "",
  base64: "",
  answer: "",
});
const getCaptcha=()=>{
 api.get("/api/captcha/getCaptcha").then((response: AxiosResponse<any>) => {
    if (response.status != 200 || !response.data) {
      if (response.data) {
        ElMessage.error("登陆失败:" + response.data.msg);
      } else {
        ElMessage.error("登陆失败!");
      }
      return;
    }
    captchaContent.captchaId = response.data.captchaId;
    captchaContent.base64 = response.data.base64;
    captchaContent.answer = response.data.answer;
    console.log(captchaContent);
  });
}
onMounted(() => {
 getCaptcha()
});
//校验规则
const rules = reactive<FormRules<typeof loginForm>>({
  account: [{ trigger: "blur" }],
  password: [
    {
      validator: (rule: any, value: any, callback: any) => {
        if (value === "") {
          callback(new Error("密码不能为空"));
        } else if (loginForm.password.length < 6) {
          callback(new Error("密码不能小于6位!"));
        } else {
          callback();
        }
      },
      trigger: "blur",
    },
  ],
  captcha: [{ required: true, message: "请输入验证码", trigger: "blur" }],
});
//表单实例
const ruleFormRef = ref<FormInstance>();
const submitForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.validate((valid) => {
    if (!valid) {
      return;
    }
    const loading = ElLoading.service({
      lock: true,
      text: "正在登录",
      background: "rgba(0,0,0,0.7)",
    });
    api.post("/api/login", loginForm).then((response: AxiosResponse<any>) => {
      if (
        response.status != 200 ||
        !response.data ||
        response.data.status != 200
      ) {
        if (response.data) {
          ElMessage.error("登录失败:" + response.data.msg);
        } else {
          ElMessage.error("登录失败！");
        }
        return;
      }
    });
  });
};
// 验证码路径
const captchaSrc = ref("/captcha/getCaptcha");
/**
 * 刷新验证码
 */
const changecaptchaSrc = () => {
  console.log("click")
  getCaptcha()
};
</script>

<style scoped>
.loginPage {
  width: 100%;
  height: 100%;
  display: flex;
  justify-items: center;
  justify-content: center;
  align-items: center;
  background: linear-gradient(133deg, #1994bb, #19acbb, #19b4bb, #2a89db);
}
.loginPage .loginPanel {
  width: 60%;
  height: 60%;
  min-width: 600px;
  max-width: 1000px;
  min-height: 400px;
  max-height: 500px;
  margin: 0 auto;
}
.loginPage .loginPanel :deep() .el-card__body {
  width: 100%;
  height: 100%;
}

.loginPage .loginPanel .loginPanelInner {
  display: flex;
  width: 100%;
  height: 100%;
}

.loginPage .loginPanel .loginPanelInner .logo {
  width: 40%;
  text-align: center;
  display: flex;
  justify-items: center;
  justify-content: center;
  align-items: center;
}
.loginPage .loginPanel .loginPanelInner .logo img {
  width: 50%;
}
.loginPage .loginPanel .loginPanelInner .split {
  height: 100%;
}
.loginPage .loginPanel .loginPanelInner .loginForm {
  flex: 1;
}
.loginPage .loginPanel .loginPanelInner .loginForm .systemName {
  text-align: center;
  font-size: 30px;
  line-height: 60px;
  margin-bottom: 20px;
  letter-spacing: 3px;
}
.loginPage .loginPanel .loginPanelInner .loginForm .form {
  width: 80%;
}

.loginPage .loginPanel .loginPanelInner .loginForm .form .loginBtn {
  width: 100%;
}
.loginPage .loginPanel .loginPanelInner .loginForm .form .captchaSrc {
    width: 100px;
    height: 100%;
    padding-left: 10px;
}
.loginPage .loginPanel .loginPanelInner .loginForm .form .captchaSrc img {
    width: 100px;
    height: 100%;
    cursor: pointer;
}
.loginPage .loginPanel .loginPanelInner .loginForm .form .captcha {
  width: 100px;
  height: 100%;
  padding-left: 10px;
}
.loginPage .loginPanel .loginPanelInner .loginForm .form .registerBtn {
  text-align: right;
  line-height: 40px;
  margin-bottom: 5px;
}
.footer {
  position: fixed;
  bottom: 0px;
  line-height: 40px;
  text-align: center;
  font-size: 14px;
  color: #000;
}
</style>
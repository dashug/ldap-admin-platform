<template>
  <div class="reset-pass">
    <div class="reset-card">
      <div class="reset-header">
        <h1 class="reset-title">找回密码</h1>
        <p class="reset-desc">输入注册邮箱，接收验证码后重置密码</p>
      </div>
      <el-form ref="form" :model="form" size="medium" class="form-container">
        <el-form-item label="邮箱">
          <div class="input-container">
            <el-input v-model="form.mail" placeholder="请输入个人邮箱" prefix-icon="el-icon-message" />
            <el-button type="primary" @click="sendEmailCode">发送验证码</el-button>
          </div>
        </el-form-item>
        <el-form-item label="验证码" class="code-item">
          <el-input v-model="form.code" placeholder="请输入验证码" prefix-icon="el-icon-mobile-phone" />
        </el-form-item>
        <el-form-item class="reset-item">
          <el-button type="primary" class="submit-btn" @click="resetPass">重置密码</el-button>
        </el-form-item>
      </el-form>
      <div class="reset-footer">
        <router-link to="/login">返回登录</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { emailPass,sendCode } from '@/api/system/user'
import { Message } from 'element-ui'

export default {
  name: 'ChangePass',
  data() {
    return {
      // 查询参数
      form: {
        mail: "",
        code: ""
      }
    }
  },
  methods: {
    // 判断结果
    judgeResult(res){
      if (res.code==0){
          Message({
            showClose: true,
            message: "操作成功",
            type: 'success'
          })
        }
    },

    // 发送邮箱验证码
    async sendEmailCode() {
      console.log('aaaaaaaa',this.form.mail);

      await sendCode({ mail: this.form.mail }).then(res =>{
        this.judgeResult(res)
      })
    },
    // 重置密码
    async resetPass() {
      await emailPass(this.form).then(res =>{
        this.judgeResult(res)
      })
      // 重新登录
      setTimeout(() => {
        this.$router.replace({ path: '/login' })
      }, 1500)
    },
  }
}
</script>

<style scoped lang="scss">
@import "~@/styles/variables.scss";

.reset-pass {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 24px;
  background: linear-gradient(135deg, #1e293b 0%, #334155 100%);
}

.reset-card {
  width: 420px;
  max-width: 100%;
  padding: 40px;
  background: #fff;
  border-radius: $cardRadius;
  box-shadow: 0 24px 48px rgba(0, 0, 0, 0.18);
}

.reset-header {
  text-align: center;
  margin-bottom: 28px;
  .reset-title {
    margin: 0 0 8px;
    font-size: 22px;
    font-weight: 600;
    color: #1e293b;
  }
  .reset-desc {
    margin: 0;
    font-size: 13px;
    color: #64748b;
  }
}

.form-container {
  width: 100%;
}

.input-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.input-container .el-input {
  flex: 1;
  margin-right: 10px;
}

.code-item .el-input {
  width: 100%;
}

.reset-item {
  text-align: center;
  margin-bottom: 0;
}
.submit-btn {
  width: 100%;
  height: 44px;
  border-radius: 8px;
}

.reset-footer {
  text-align: center;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid $borderColor;
  font-size: 13px;
  a {
    color: $themePrimary;
    text-decoration: none;
    &:hover { text-decoration: underline; }
  }
}
</style>

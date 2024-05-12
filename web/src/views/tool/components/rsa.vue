<template>
  <el-form class="page-rsa" label-width="100px">
    <el-form-item label="密钥位数">
      <el-select v-model="bits" class="select-view">
        <el-option :key="item" v-for="item in bitsList" :label="item" :value="item"></el-option>
      </el-select>
      <el-button type="primary" :disabled="generateLoading" @click="generateRsaKey">生成密钥</el-button>
    </el-form-item>
    <el-form-item label="密钥格式">
      <el-select v-model="rsaType" class="select-view">
        <el-option :key="item" v-for="item in rsaTypeList" :label="item" :value="item"></el-option>
      </el-select>
      <el-button type="success" :disabled="encryptLoading" @click="encryptByPublicKey">公钥加密</el-button>
      <el-button type="success" :disabled="decryptLoading" @click="decryptByPrivateKey">私钥解密</el-button>
    </el-form-item>
    <el-form-item label="签名方式">
      <el-select v-model="signType" class="select-view">
        <el-option :key="item" v-for="item in signTypeList" :label="item" :value="item"></el-option>
      </el-select>
      <el-button type="warning" :disabled="signLoading" @click="privateKeySign">私钥签名</el-button>
      <el-button type="warning" :disabled="verifyLoading" @click="publicKeyVerify">公钥验签</el-button>
    </el-form-item>
    <el-form-item label="待加密内容">
      <el-input v-model="content" type="textarea" :rows="3" resize="none" placeholder="待加密的字符串"></el-input>
    </el-form-item>
    <el-form-item label="公钥">
      <el-input v-model="publicKey" type="textarea" :rows="4" resize="none" placeholder="RSA公钥"></el-input>
    </el-form-item>
    <el-form-item label="私钥">
      <el-input v-model="privateKey" type="textarea" :rows="4" resize="none" placeholder="RSA私钥"></el-input>
    </el-form-item>
    <el-form-item label="公钥加密">
      <el-input v-model="encryptContent" type="textarea" :rows="3" resize="none" placeholder="公钥加密后的内容"></el-input>
    </el-form-item>
    <el-form-item label="私钥解密">
      <el-input v-model="decryptContent" type="textarea" :rows="3" resize="none" readonly placeholder="解密公钥加密后的内容"></el-input>
    </el-form-item>
    <el-form-item label="私钥签名">
      <el-input v-model="sign" type="textarea" :rows="2" resize="none" placeholder="使用私钥和原字符串签名"></el-input>
    </el-form-item>
    <el-form-item label="公钥验签">
      <el-input v-model="verify" type="textarea" :rows="2" resize="none" readonly placeholder="使用公钥、签名和原字符串验证私钥签名"></el-input>
    </el-form-item>
  </el-form>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import RsaApi from "@/api/rsa";

const bits = ref(2048);
const bitsList = ref([512, 1024, 2048, 4096]);
const rsaType = ref("PKCS1");
const rsaTypeList = ref(["PKCS1", "PKCS8"]);
const signType = ref("SHA256");
const signTypeList = ref(["MD5", "SHA1", "SHA256", "SHA512"]);
const content = ref("");
const publicKey = ref("");
const privateKey = ref("");
const encryptContent = ref("");
const decryptContent = ref("");
const sign = ref("");
const verify = ref("");
const generateLoading = ref(false);
const encryptLoading = ref(false);
const decryptLoading = ref(false);
const signLoading = ref(false);
const verifyLoading = ref(false);

/**
 * 生成RSA密钥对
 */
const generateRsaKey = () => {
  generateLoading.value = true;
  RsaApi.generate({
    bits: bits.value,
    isPKCS8: rsaType.value === "PKCS8",
  })
    .then((res) => {
      publicKey.value = res.data.publicKey;
      privateKey.value = res.data.privateKey;
    })
    .catch(() => {
      publicKey.value = "";
      privateKey.value = "";
    })
    .finally(() => {
      generateLoading.value = false;
    });
};

/**
 * 公钥加密
 */
const encryptByPublicKey = () => {
  encryptLoading.value = true;
  RsaApi.encrypt({
    message: content.value,
    publicKey: publicKey.value,
  })
    .then((res) => {
      encryptContent.value = res.data;
    })
    .catch(() => {
      encryptContent.value = "";
    })
    .finally(() => {
      encryptLoading.value = false;
    });
};

/**
 * 私钥解密
 */
const decryptByPrivateKey = () => {
  decryptLoading.value = true;
  RsaApi.decrypt({
    message: encryptContent.value,
    privateKey: privateKey.value,
    isPKCS8: rsaType.value === "PKCS8",
  })
    .then((res) => {
      decryptContent.value = res.data;
    })
    .catch(() => {
      decryptContent.value = "";
    })
    .finally(() => {
      decryptLoading.value = false;
    });
};

/**
 * 私钥签名
 */
const privateKeySign = () => {
  signLoading.value = true;
  RsaApi.sign({
    message: content.value,
    privateKey: privateKey.value,
    signType: signType.value,
    isPKCS8: rsaType.value === "PKCS8",
  })
    .then((res) => {
      sign.value = res.data;
    })
    .catch(() => {
      sign.value = "";
    })
    .finally(() => {
      signLoading.value = false;
    });
};

/**
 * 公钥验签
 */
const publicKeyVerify = () => {
  verifyLoading.value = true;
  RsaApi.verify({
    message: content.value,
    publicKey: publicKey.value,
    sign: sign.value,
    signType: signType.value,
  })
    .then((res) => {
      verify.value = res.data ? "验签通过" : "验签失败";
    })
    .catch(() => {
      verify.value = "";
    })
    .finally(() => {
      verifyLoading.value = false;
    });
};
</script>

<style lang="scss">
.page-rsa {
  width: 100%;
  .select-view {
    width: 100px;
    margin: 1px 0;
  }
  .el-button {
    margin: 1px 10px;
  }
  .el-button + .el-button {
    margin-left: 0;
  }
}
</style>

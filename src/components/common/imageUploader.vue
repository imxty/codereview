<template>
  <div class="imgshangc">
    <el-upload action="#" :file-list="filelist" list-type="picture-card" :auto-upload="false"
     :on-preview="handlePictureCardPreview" :on-remove="handleRemove" :on-change="handleDownload">
      <i slot="default" class="el-icon-plus"></i>
      <div slot="file" slot-scope="{file}">
        <div :class="size" class="el-upload el-upload--picture-card">
          <span class="upload-field">
            <span class="upload-icon">+</span>
            <span class="upload-text">上传图片</span>
          </span>
        </div>
      </div>
    </el-upload>
    <div class="cont">
      <el-dialog v-model="dialogVisible">
        <img width="100%" :src="dialogImageUrl" alt="">
      </el-dialog>
    </div>
  </div>
</template>
<script setup>
import { ref, defineEmits } from "vue"
import { elmessage } from '@/utils/popup'
const dialogImageUrl = ref('',)
const dialogVisible = ref(false)
const disabled = ref(false)
const filelist = ref([])
const emit = defineEmits(['onImageUploaded', 'imageBase64'])

const handleRemove = (file) => {
  console.log(file);
}
const handlePictureCardPreview = (file) => {
  dialogImageUrl.value = file.url;
  dialogVisible.value = true;
}
const handleDownload = async (file) => {
  let testFile = file.name.substring(file.name.lastIndexOf('.') + 1).toLowerCase()

  const extension = testFile === 'jpg' || testFile === 'jpeg' || testFile === 'png'

  const isLt2M = (file.size / 1024 / 1024 < 5);
  if (!extension) {
    console.log('请上传png/jpg/jpeg格式的图片');
    elmessage('请上传png/jpg/jpeg格式的图片')
    filelist.value = []
    return false;
  }
  console.log(file.size);
  if (file.size / 1024 / 1024 > 2) {
    console.log('图片不能大于2M')
    elmessage('图片不能大于2M')
    filelist.value = []
    return false;
  }
  const imageBase64 = await getFileBase64(file.raw)
  console.log(imageBase64);
  emit('onImageUploaded', file.raw)
  emit('imageBase64', imageBase64)
  return (extension) && isLt2M

}

const getFileBase64 = async (file) => {
  const info = await new Promise((resolve, reject) => {
    const reader = new FileReader()
    let imgResult = ''
    reader.readAsDataURL(file)
    reader.onload = () => {
      imgResult = reader.result
    }
    reader.onerror = (error) => {
      reject(error)
    }
    reader.onloadend = () => {
      resolve(imgResult)
    }
  })
  return info
}
const props = defineProps({
  imgurl: String
});

const getimgurl = () => {
  if (props.imgurl!=undefined) {
    filelist.value.push({
    url: props.imgurl
  })
  }else{
    return
  }
  
}
getimgurl()
</script>
<style lang="scss">
.el-upload-list__item {
  position: absolute;
}

.upload-field {
  display: flex;
  flex-direction: column;
  align-items: center;

  .upload-icon {
    font-size: 48px;
    color: #C1C7D0
  }

  .upload-text {
    font-size: 14px;
    color: #C1C7D0;
  }
}

.imgshangc {
  .el-dialog {
  }

  .el-dialog__body {
    padding: 70px;
    text-align: center;

    img {
      width: 80%;
    }
  }
}
</style>
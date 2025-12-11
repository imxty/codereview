import { ElMessage } from 'element-plus';


//弹窗组件
export function elmessage (res){
    const txt = res.replaceAll(/^\[errcode:(?<code>\d+)\]/g,'')
    console.log(txt);
    ElMessage({
        showClose: true,
        message: txt,
        type: 'error',
     })
}
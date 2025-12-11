<template>
   <div class="wrap">
      <div class="ownertop">
         <div class="ownerleft">
            <el-button type="primary" class="owbtm" @click.stop="dakaidialog">新建商户</el-button>
            <div @click.stop="dialogurl = true" style="cursor: pointer;flex: 0.8;min-width: 85px;">
               <img src="../../assets/image/icon_invite.png" alt="">
               邀请链接
            </div>
         </div>
         <div class="ownerright">
            <el-select v-model="selectlist" placeholder="方案名称" size="large" style="width: 240px;margin-right: 20px;"
               clearable>
               <el-option key="no" label="未填写" value="null" />
               <el-option v-for="item in selectdata" :key="item.treatment_id" :label="item.treatment_name"
                  :value="item.treatment_id" />
            </el-select>
            <el-input v-model="sminpt" placeholder="请输入商户名称搜索商户" class="ownerinp"
               @blur="sminpt = $event.target.value.trim()" :maxlength="20"></el-input>
            <el-button type="primary" class="owbtm"
               @click.stop="offset = 1, SearchTenants('TENANT_STATUS_UNSET', sminpt, selectlist)">搜索</el-button>
         </div>
      </div>
      <div class="ownercnt">
         <div class="owcnttop">
            <div class="owcnttoplf">
               <div :class="[{ allmer: shanghuzt === '所有商户' }]"
                  @click.stop="offset = 1, SearchTenants('TENANT_STATUS_UNSET')" style="cursor: pointer;">
                  所有商户
               </div>

               <div :class="[{ allmer: shanghuzt === '未认证' }]"
                  @click.stop="offset = 1, SearchTenants('TENANT_STATUS_UNAUTH')" style="cursor: pointer;">
                  未认证
               </div>
               <div :class="[{ allmer: shanghuzt === '未激活' }]"
                  @click.stop="offset = 1, SearchTenants('TENANT_STATUS_USING')" style="cursor: pointer;">
                  使用中
               </div>
               <div :class="[{ allmer: shanghuzt === '已激活' }]"
                  @click.stop="offset = 1, SearchTenants('TENANT_STATUS_PENDING')" style="cursor: pointer;">
                  待续期
               </div>
            </div>
            <div class="owcnttoprg" style="align-items: center;">
               <span v-if="checkList.length > 0">
                  <span class="quanxuan" v-if="biaolie">
                     <el-checkbox v-model="checkAll" :indeterminate="isIndeterminate"
                        @change="handleCheckAllChange">全选</el-checkbox>
                  </span>
                  <span v-if="biaolie">｜</span>
                  <span>
                     已选（{{ checkList.length }}）
                  </span>
                  <span>｜</span>
                  <span style="margin-right: 20px;cursor: pointer;" v-if="biaolie" @click.stop="checkList = []">
                     取消
                  </span>
                  <el-button type="primary" class="owbtm" @click.stop="toprepayment">分配方案</el-button>
                  <el-button type="warning" class="owbtm"
                     @click.stop="offset = 1, dialogconse = true">取消分配方案</el-button>
               </span>
               <img src="../../assets/image/icon_view_category_selected.png" alt="" v-if="biaolie">
               <img src="../../assets/image/icon_view_list.png" alt="" v-if="biaolie" @click.stop="biaoliet(false)"
                  style="cursor: pointer;">
               <img src="../../assets/image/icon_view_category.png" alt="" v-if="biaolie == false"
                  @click.stop="biaoliet(true)" style="cursor: pointer;">
               <img src="../../assets/image/icon_view_list_selected.png" alt="" v-if="biaolie == false">
               <el-tooltip placement="bottom">
                  <template #content>
                     <div class="czzytxt">
                        <div>
                           <h3>商户认证</h3>
                           <p>新建商户或邀请商户自行填写信息后，需要后台审核并对商户进行续期后方可使用</p>
                        </div>
                        <div>
                           <h3>商户续期</h3>
                           <p>商户认证后未开通使用期限或超过使用期限后，将对商户进行停用，请及时联系管理员进行续期</p>
                        </div>
                        <div>
                           <h3>商户地址</h3>
                           <p>{{ tenanturl }}</p>
                        </div>
                     </div>
                  </template>
                  <span class="zczy" style="display: flex;margin-left: 20px;margin-right: 20px;">
                     <img src="../../assets/icons/yiwen.svg" alt="" style="width: 15px;">
                     操作指引
                  </span>
               </el-tooltip>
            </div>
         </div>
         <div class="ownerdata">
            <div v-if="nocommod == false" class="nocom">
               <img src="../../assets/image/icon_noshop.png" alt="">
               <div>暂无商户</div>
            </div>
            <div v-if="nocommod == true">
               <div class="arecom" v-if="biaolie">
                  <template v-if="chonghui">
                     <div class="aretcot aretavc" v-for="(key, i) in commoddata" :key="key" @mouseenter="enter(i)"
                        @mouseleave="leave"
                        :class="[{ yucunsh: key.tenant_status == 'TENANT_STATUS_PRESTORE' || key.tenant_status == 'TENANT_STATUS_PRESTORE_DISABLE' }, { mouseenyr: i == current }]"
                        @click.stop="totenant(key)">
                        <div class="flrmname">
                           <div style="display: flex; align-items: center; cursor: pointer;">
                              {{ key.name }}
                           </div>
                           <!-- v-show="i ==current ||checkList.length>0" -->
                           <div class="checknam" v-show="i == current || checkList.length > 0">
                              <el-checkbox-group v-model="checkList" @change="handleCheckedCitiesChange">
                                 <el-checkbox :label="key" @click.stop="1"><br /></el-checkbox>
                              </el-checkbox-group>
                           </div>
                        </div>
                        <div class="xgshan" @click.stop="xiugai(key)"
                           v-if="key.tenant_review_status != 'TENANT_REVIEW_STATUS_REVIEWING'">修改</div>
                        <div class="flrmess">
                           方案 ：<span>
                              <span v-if="treatments[key.tenant_id]">
                                 {{ treatments[key.tenant_id] }}
                              </span>
                              <span v-else>
                                 未分配方案
                              </span>
                           </span>
                        </div>
                        <div class="flrmess">
                           联系人：{{ key.contact_name }}
                        </div>
                        <div class="flrmess">
                           联系人手机号：{{ key.contact_phone }}
                        </div>
                        <div class="flrmess" style="margin-top: 5px;">
                           <span class="biaoshi guoqi" v-if="key.tenant_status == 'TENANT_STATUS_PENDING'">待续期</span>
                        </div>
                        <div class="usetem under" v-if="key.tenant_review_status == 'TENANT_REVIEW_STATUS_REVIEWING'">
                           <div>
                              审核中
                           </div>
                        </div>
                        <div class="usetem underof"
                           v-if="key.tenant_review_status == 'TENANT_REVIEW_STATUS_REVIEW_FAIL'" @click.stop="1">
                           <el-popover placement="bottom-start" :width="390" trigger="click">
                              <template #reference>
                                 <div>
                                    <el-link type="danger" :underline="false"
                                       style="font-weight: 400; font-size: 12px;">审核失败 > </el-link>
                                 </div>
                              </template>
                              <div class="statusfalt">
                                 <div class="falttop">
                                    <div class="flttop-left">
                                       商户资质审核失败
                                    </div>
                                    <div class="flttop-right">
                                       <el-link :underline="false" class="flttop-right"
                                          @click.stop="xiugai(key)">重新编辑</el-link><span><img
                                             src="../../assets/icons/u82.svg" alt=""
                                             style="width: 6px; height: 7px;margin-left: 3px;"></span>
                                    </div>
                                 </div>
                                 <div class="faltdata">
                                    <p>由于</p>
                                    <div>
                                       {{ key.fail_reason }}
                                       的原因，您的商户资质认证审核失败。请修改后重新提审
                                    </div>
                                 </div>
                              </div>
                           </el-popover>
                        </div>
                     </div>
                  </template>
               </div>
               <div class="datalie" v-if="biaolie == false">
                  <el-table ref="liecheck" :data="commoddata" style="width: 100%"
                     @selection-change="handleSelectionChange" :row-key="getRowKeys" @row-click="totenant">
                     <el-table-column type="selection" width="55" :reserve-selection="true" />
                     <el-table-column prop="name" label="商户名称/商户ID">
                        <template #default="scope">
                           <div>
                              <div class="lienametxt">
                                 {{ scope.row.name }}
                              </div>
                              <div>
                                 {{ scope.row.tenant_id }}
                              </div>
                           </div>
                        </template>
                     </el-table-column>
                     <el-table-column label="状态">
                        <template #default="scope">
                           <div>
                              <div>
                                 <span v-if="scope.row.tenant_status == 'TENANT_STATUS_PENDING'">待续期</span>
                                 <span v-if="scope.row.tenant_status == 'TENANT_STATUS_UNAUTH'">未认证</span>
                                 <span v-if="scope.row.tenant_status == 'TENANT_STATUS_USING'">使用中</span>
                              </div>
                              <div>
                                 <span v-if="scope.row.tenant_status == 'TENANT_REVIEW_STATUS_REVIEWING'">审核中</span>
                                 <span v-if="scope.row.tenant_review_status == 'TENANT_REVIEW_STATUS_REVIEW_FAIL'"
                                    @click.stop="1">
                                    <el-popover placement="bottom-start" :width="390" trigger="click">
                                       <template #reference>
                                          <div>
                                             <el-link type="danger" :underline="false"
                                                style="font-weight: 400; font-size: 12px;">审核失败 > </el-link>
                                          </div>
                                       </template>
                                       <div class="statusfalt">
                                          <div class="falttop">
                                             <div class="flttop-left">
                                                商户资质审核失败
                                             </div>
                                             <div class="flttop-right">
                                                <el-link :underline="false" class="flttop-right"
                                                   @click.stop="xiugai(scope.row)">重新编辑</el-link><span><img
                                                      src="../../assets/icons/u82.svg" alt=""
                                                      style="width: 6px; height: 7px;margin-left: 3px;"></span>
                                             </div>
                                          </div>
                                          <div class="faltdata">
                                             <p>由于</p>
                                             <div>
                                                {{ scope.row.fail_reason }}
                                                的原因，您的商户资质认证审核失败。请修改后重新提审
                                             </div>
                                          </div>
                                       </div>
                                    </el-popover>
                                 </span>
                              </div>

                           </div>
                        </template>
                     </el-table-column>
                     <el-table-column prop="treatments" label="方案">
                        <template #default="scope">
                           <div v-if="treatments[scope.row.tenant_id]">
                              {{ treatments[scope.row.tenant_id] }}
                           </div>
                           <div v-else>
                              --
                           </div>
                        </template>
                     </el-table-column>
                     <el-table-column prop="phone" label="联系人/手机号">
                        <template #default="scope">
                           <div>
                              <div>
                                 {{ scope.row.contact_name }}
                              </div>
                              <div>
                                 {{ scope.row.contact_phone }}
                              </div>
                           </div>
                        </template>
                     </el-table-column>
                     <el-table-column prop="phone" label="开通时间/到期时间">
                        <template #default="scope">
                           <div>
                              <div v-if="scope.row?.timeline?.start_time != ''">
                                 {{ fermitTime(scope.row?.timeline?.start_time, true) }}
                              </div>
                              <div v-else>---</div>
                              <div v-if="scope.row?.timeline?.end_time != ''">
                                 {{ fermitTime(scope.row?.timeline?.end_time, true) }}
                              </div>
                              <div v-else>---</div>

                           </div>
                        </template>
                     </el-table-column>
                     <el-table-column label="操作" :filter-method="filterTag" filter-placement="bottom-end" width="300">
                        <template #default="scope">
                           <div>
                              <el-button type="primary" plain
                                 :disabled="scope.row.tenant_review_status == 'TENANT_REVIEW_STATUS_REVIEWING'"
                                 @click.stop="xiugai(scope.row)">
                                 修改</el-button>
                           </div>
                        </template>
                     </el-table-column>
                  </el-table>
               </div>
            </div>
         </div>
         <div class="fenye">
            <el-pagination v-model:current-page="offset" v-model:page-size="pageSize4" :page-sizes="[10, 30, 50, 100]"
               :small="true" :disabled="disabled" :background="background"
               layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange"
               @current-change="handleCurrentChange" />
         </div>
      </div>
      <el-dialog v-model="dialogVisible" width="50%" :destroy-on-close=true :show-close=false class="changjiantan">
         <div class="newtop">
            <div class="toptxt">
               <span v-if="xiugaitj != '修改'">请填写您的商户信息并提交审核</span>
               <span v-if="xiugaitj == '修改'">请修改您的商户信息并提交审核</span>
            </div>
            <div class="toptxt shenhe" v-if="xiugaitj != '修改'">
               审核通过后才可激活商户进行使用
            </div>
            <div class="newcnt">
               <el-form ref="ruleFormRef" :model="certifyForm" :rules="rules" :size="formSize" label-width="130px">
                  <el-form-item label="商户图标：" width=90 required="true" style="text-align:left">
                     <div class="tableright">
                        <imageUploader @imageBase64="logoImageBase64" @onImageUploaded="logoImageUploaded"
                           :imgurl="certifyForm.logoUrl">
                        </imageUploader>
                        <div class="imgtext">
                           <p>请上传商户门头照片/商户logo；</p>

                           <p>图片不允许涉及政治敏感与色情;</p>

                           <p>图片格式只支持：BMP、JPEG、JPG、GIF、PNG，大小不超过2M</p>
                        </div>
                     </div>
                  </el-form-item>
                  <el-form-item label="商户名称：" required="true" style="text-align:left" prop="name">
                     <el-input placeholder="请输入商户名称" v-model="certifyForm.name" />
                  </el-form-item>
                  <el-form-item label="商户地址：" required="true" style="text-align:left" :error="addressmessage">
                     <div class="dizhi">
                        <div class="dizhishi">
                           <el-select v-model="certifyForm.address.province" placeholder="请选择" @change="checksheng"
                              @blur="addressblur">
                              <el-option v-for="item in provinceOptions()" :key="item.value" :label="item.label"
                                 :value="item.value">
                              </el-option>
                           </el-select>
                           <el-select v-model="certifyForm.address.city" placeholder="请选择" @change="checkshi"
                              @blur="addressblur" style="margin-left: 6px;">
                              <el-option v-for="item in cityOptions()" :key="item.value" :label="item.label"
                                 :value="item.value">
                              </el-option>
                           </el-select>
                           <el-select v-model="certifyForm.address.district" placeholder="请选择" @change="checkxian"
                              @blur="addressblur" style="margin-left: 5px;">
                              <el-option v-for="item in regionOptions()" :key="item.value" :label="item.label"
                                 :value="item.value">
                              </el-option>
                           </el-select>
                        </div>
                        <el-form-item class="el-form-item " style="margin-top: 20px;">
                           <el-input placeholder="具体地址" v-model="certifyForm.address.street" @blur="addressblur" />
                        </el-form-item>
                     </div>
                  </el-form-item>
                  <el-form-item label="营业执照：" required="true" style="text-align:left">
                     <div class="tableright">
                        <imageUploader @imageBase64="yingyeImageBase64" @onImageUploaded="yingyeImageUploaded"
                           :imgurl="certifyForm.businessLicenseUrl" v-if="dtimg == true">
                        </imageUploader>
                        <imageUploader @imageBase64="yingyeImageBase64" @onImageUploaded="yingyeImageUploaded"
                           :imgurl="certifyForm.businessLicenseUrl" v-else>
                        </imageUploader>

                        <div class="imgtext">
                           <p>请上传清晰的营业执照图片，保证以上商户信息与营业执照一致；</p>

                           <p>图片不允许涉及政治敏感与色情;</p>

                           <p>图片格式只支持：BMP、JPEG、JPG、GIF、PNG，大小不超过2M</p>
                        </div>
                     </div>
                  </el-form-item>
                  <el-form-item label="社会信用代码：" required="true" style="text-align:left" prop="socialCreditCode">
                     <el-input placeholder="请输入社会信用代码" v-model="certifyForm.socialCreditCode" />
                  </el-form-item>
                  <el-form-item label="联系人：" required="true" style="text-align:left" prop="contactName">
                     <el-input placeholder="请输入商户联系人姓名" v-model="certifyForm.contactName" />
                  </el-form-item>
                  <el-form-item label="联系人手机号：" required="true" style="text-align:left" :error="phonemessage">
                     <el-input placeholder="请输入联系人手机号，用于商户登录" v-model="userphone.phone" @blur="phoneDate"
                        auto-complete="new-password" />
                  </el-form-item>
                  <el-form-item label="短信验证码：" required="true" style="text-align:left" class="duanxin"
                     v-if="xiugaitj != '修改'" :error="codemessage">
                     <el-input placeholder="请输入验证码" class="no-autofill-pwd" v-model="userphone.smsCode"
                        @blur="checkcodeDate" auto-complete="new-password" />
                     <div class="fasong">
                        <span style="padding-right: 10px; color: #C1C7D0;">|</span>
                        <el-link type="primary" :underline="false" class="huoqu"
                           @click.stop="SendPhone(userphone, 'TEMPLATE_ACTION_BIND_TENANT_PHONE')"
                           :disabled="duanxindj">{{
                              dxtxt }}</el-link>
                     </div>
                  </el-form-item>
                  <el-form-item label="短信验证码：" required="true" style="text-align:left" class="duanxin"
                     v-if="Modifydata?.tenant_status == 'TENANT_STATUS_UNAUTH'" :error="codemessage">
                     <el-input placeholder="请输入验证码" class="no-autofill-pwd" v-model="userphone.smsCode"
                        @blur="checkcodeDate" auto-complete="new-password" />
                     <div class="fasong">
                        <span style="padding-right: 10px; color: #C1C7D0;">|</span>
                        <el-link type="primary" :underline="false" class="huoqu"
                           @click.stop="SendPhone(userphone, 'TEMPLATE_ACTION_BIND_TENANT_PHONE')"
                           :disabled="duanxindj">{{
                              dxtxt }}</el-link>
                     </div>
                  </el-form-item>
                  <el-form-item label="商户登录密码：" required="true" style="text-align:left" v-if="xiugaitj != '修改'"
                     :error="birthDateMessage">
                     <el-input show-password class="no-autofill-pwd" placeholder="请设置商户登录密码"
                        v-model="userInfo.newPlainPassword" @blur="checkBirthDate" auto-complete="new-password" />
                  </el-form-item>
                  <el-form-item label="商户登录密码：" required="true" style="text-align:left"
                     v-if="Modifydata?.tenant_status == 'TENANT_STATUS_UNAUTH'" :error="birthDateMessage">
                     <el-input show-password class="no-autofill-pwd" placeholder="请设置商户登录密码"
                        v-model="userInfo.newPlainPassword" @blur="checkBirthDate" auto-complete="new-password" />
                  </el-form-item>
                  <el-form-item label="确认密码：" required="true" style="text-align:left" v-if="xiugaitj != '修改'"
                     :error="qrpassword">
                     <el-input show-password class="no-autofill-pwd" placeholder="请再次输入商户登录密码"
                        v-model="userInfo.password" @blur="usepassword" auto-complete="new-password" />
                  </el-form-item>
                  <el-form-item label="确认密码：" required="true" style="text-align:left"
                     v-if="Modifydata?.tenant_status == 'TENANT_STATUS_UNAUTH'" :error="qrpassword">
                     <el-input show-password class="no-autofill-pwd" placeholder="请再次输入商户登录密码"
                        v-model="userInfo.password" @blur="usepassword" auto-complete="new-password" />
                  </el-form-item>
               </el-form>
            </div>
         </div>
         <template #footer>
            <span class="dialog-footer">
               <el-button @click.stop="dialogVisible = false" class="quxiaobut">取消</el-button>
               <el-button type="primary" @click.stop="CreateTenant(ruleFormRef)" class="tijiaobut">
                  提交审核
               </el-button>
            </span>
         </template>

      </el-dialog>
      <el-dialog v-model="dialogurl" title="邀请链接" width="30%">
         <div class="diaurl">
            <div class="urltop">
               <div class="urltop-txt">点击复制链接邀请商户填写信息:</div>
               <div style="width: 100%;overflow-wrap:break-word;">
                  {{ jmurl }}
               </div>
               <br/>
               <div>
                  商户地址：{{ tenanturl }}
               </div>
            </div>
         </div>
         <template #footer>
            <span class="dialog-footer-yq">
               <el-button type="primary" @click.stop="copyText(false)" style="background-color: #006be5;" class="yqbut">
                  复制邀请链接
               </el-button>
            </span>
         </template>
      </el-dialog>
      <el-dialog v-model="dialogconse" title="提示" width="30%">
         <div class="diaurl">
            <div class="urltop">
               当前正在取消{{ checkList.length }}个商户的商品分配方案
            </div>
         </div>
         <template #footer>
            <span class="dialog-footer-yq">
               <el-button @click.stop="dialogconse = false" class="quxiaobut" style="height: 40px;">取消</el-button>
               <el-button type="primary" @click.stop="CancelTreatmentToTenants" class="owbtm">
                  确定
               </el-button>
            </span>
         </template>
      </el-dialog>
   </div>
</template>
<script setup>
import { reactive, ref, onMounted, onUnmounted } from "vue";
import imageUploader from '@/components/common/imageUploader.vue';
import { exportExcel } from '@/utils/exportExcel.js'
import { getAreaOptions, TextToCode } from '@/utils/areadDataOptions'
import {
   ListTenantsRequest,
   ActivateTenantRequest,
   DisableTenantRequest,
   SearchTenantsRequest,
   SendPhoneVerificationCode,
   CreateTenantRequest,
   CommitEntityCertificateRequest,
   GetTenantRequest,
   ReCreateTenantRequest,
   UploadImageRequest,
   ListTreatmentsRequest,
   CancelTreatmentToTenantsRequest
} from '@/api/api'
import { getTenantId, getRefreshToken, getuserinfo } from '@/utils/storage'
import { elmessage } from '@/utils/popup'
import { ElMessage } from 'element-plus';
import stort from '@/store/commerindex'
import { useRouter } from 'vue-router';
import { fermitTime, getoneMonth, getDateByDays } from '../../utils/fermitTime'
import { settopname } from '../../utils/sethometop'
import useClipboard from 'vue-clipboard3'
import { rulesz } from '@/utils/rules'
import { useStore } from 'vuex'


//设置vuex对象
const $store = useStore();
const $router = useRouter();
const { toClipboard } = useClipboard()
//商户数据
const commoddata = ref(stort.state.commoddata)
//新建商户弹出框
const dialogVisible = ref(false)
//邀请链接弹出框
const dialogurl = ref(false)
//判断是否有商户
const nocommod = ref(true)
//默认每页几条数据
const pageSize4 = ref(10)
//默认总数据
const total = ref(1)

const dialogconse = ref(false)
const num = ref(1)
const sandiankuang = ref(null)
//接受多选框数据
const checkList = ref([])
const liecheck = ref([])
const datacheck = ref([])

//判断商户展示形式
const biaolie = ref(true)
const timestop = ref()
const searchdata = ref('')
//搜索商户输入框
const sminpt = ref('')

const treatments = ref({})
const selectlist = ref('')
const selectdata = ref('')
//判断当前是什么状态的商户
const shanghuzt = ref('所有商户')
//邀请链接
const jmurl = ref('https://res.jinmuhealth.com' + import.meta.env.VITE_APP_PUBLIC_PATH + 'index.html#/lnvite' + '?o=' + getTenantId() + '&on=' + encodeURI(getuserinfo()?.userinfo.name))

const tenanturl = ref(import.meta.env.VITE_APP_TENRUL_PATH)
console.log(getuserinfo());
const Modifydata = ref(null)
const cities = ref([])
//全选按钮状态
const checkAll = ref(false)
const isIndeterminate = ref(false)
//判断是修改还是添加商户
const xiugaitj = ref('添加')
const current = ref(1)

//
const stencilxg = ref(true)

const dtimg = ref(false)
const chonghui = ref(false)
//
const daletsp = ref(false)
const shanchudata = ref(null)

const stencil_id = ref(null)

const userInfo = ref({
   newPlainPassword: '',
   password: ''
})
const userphone = ref({

   phone: '',
   smsCode: null,
   txId: null,
   templateAction: 'TEMPLATE_ACTION_BIND_TENANT_PHONE'
})
const certifyForm = ref({
   tenantId: null,
   logo: {
      mime: null,
      image: null,
      filename: null
   },
   name: null,
   contactName: null,
   address: {
      province: null,
      city: null,
      district: '',
      street: null
   },

   contactPhone: null,
   businessLicense: {
      mime: null,
      image: null,
      filename: null
   },
   socialCreditCode: null,
   logoUrl: null,
   businessLicenseUrl: null,
   safePhone: null,
   organizationId: null
})

//激活停用弹窗开关
const deactivatesdialog = ref(false)
//判断商户是激活还是停用
const deactivates = ref(null)
//激活停用弹窗数据
const deactivatesdata = ref(null)

const olddata = ref('')

const Merchants = ref([])
const rules = reactive(rulesz)
//表单验证
const formSize = ref('default')
const ruleFormRef = ref()
//密码校验
const birthDateMessage = ref('')
//确认密码验证
const qrpassword = ref('')

const creatcg = ref(true)

const addressmessage = ref('')

const phonemessage = ref('')

const subscriptions = ref([])
//短信验证
const codemessage = ref('')
//短信60秒间隔
const time = ref(60)
const dxtxt = ref('获取验证码')
const duanxindj = ref(false)
const timesphone = ref(null)
//分页偏移量
const offset = ref(1)
const useid = ref(getTenantId())
const handleSizeChange = (val) => {
   isIndeterminate.value = true
   if (shanghuzt.value == '所有商户') {
      getlist()
   } else {
      SearchTenants(searchdata.value)
   }
   checkList.value = [...new Set(checkList.value)]
}
const handleCurrentChange = (val) => {
   isIndeterminate.value = true
   if (shanghuzt.value == '所有商户') {
      getlist()
   } else {
      SearchTenants(searchdata.value)
   }
   checkList.value = [...new Set(checkList.value)]
}
//表单验证事件
const checkBirthDate = () => {
   if (userInfo.value.newPlainPassword == null || userInfo.value.newPlainPassword == '') {
      birthDateMessage.value = '请输入登录密码'
   } else if (userInfo.value.newPlainPassword.length < 8 || userInfo.value.newPlainPassword.length > 18) {
      birthDateMessage.value = '请输入8到18位密码'
   } else {
      birthDateMessage.value = ''
   }
}
const checkcodeDate = () => {
   console.log(userphone.value);
   if (userphone.value.smsCode != null && userphone.value.smsCode != '') {
      codemessage.value = ''
   } else {
      codemessage.value = '请输入短信验证码'
   }
}
const phoneDate = () => {
   console.log(userphone.value.phone);
   if (userphone.value.phone == null && userphone.value.phone == '') {
      phonemessage.value = '请输入手机号码'
   } else if (userphone.value.phone.length != 11) {
      phonemessage.value = '请输入正确的手机号码'
   } else {
      phonemessage.value = ''
   }
}
const usepassword = () => {
   if (userInfo.value.password == null || userInfo.value.password == '') {
      qrpassword.value = "请再次输入密码";
   } else if (userInfo.value.password !== userInfo.value.newPlainPassword) {
      qrpassword.value = "两次输入密码不一致!";
   } else {
      qrpassword.value = '';
   }
};
const addressblur = () => {
   if (certifyForm.value.address.province == null || certifyForm.value.address.city == null || certifyForm.value.address.district == null || certifyForm.value.address.street == null || certifyForm.value.address.city == '' || certifyForm.value.address.province == '' || certifyForm.value.address.street == '') {
      addressmessage.value = '请输入完整地址'
   }
   else {
      addressmessage.value = ''
   }
}
//全选按钮
const handleCheckAllChange = (val) => {
   checkList.value = val ? cities.value : []
   isIndeterminate.value = false
   checkList.value = [...new Set(checkList.value)]
}
const handleCheckedCitiesChange = (value) => {
   const checkedCount = value.length
   checkAll.value = checkedCount === cities.value.length
   isIndeterminate.value = checkedCount > 0 && checkedCount < cities.value.length
   checkList.value = [...new Set(checkList.value)]
}
//获取商户列表数据
const getlist = async () => {
   const rs = {
      organizationId: useid.value,
      pagination: {
         offset: (offset.value - 1) * pageSize4.value,
         size: pageSize4.value
      }
   }
   const add = {
      organizationId: useid.value
   }
   const data = await ListTenantsRequest(rs)

   commoddata.value = data.data.tenants
   subscriptions.value = data.data.subscriptions == undefined ? {} : data.data.subscriptions
   treatments.value = data.data.treatments == undefined ? {} : data.data.treatments
   commoddata.value?.forEach(i => {
      i.timeline = {
         start_time: subscriptions.value[i.tenant_id]?.start_time == undefined ? '' : subscriptions.value[i.tenant_id].start_time,
         end_time: subscriptions.value[i.tenant_id]?.end_time == undefined ? '' : subscriptions.value[i.tenant_id].end_time,
      }
      console.log(i);
   })
   if (commoddata.value) {
      nocommod.value = true
   } else {
      nocommod.value = false
   }
   if (commoddata.value) {
      commoddata.value?.forEach((v) => {
         let tfdata = false
         cities.value.forEach(i => {
            if (i.tenant_id == v.tenant_id) {
               tfdata = true
               return
            }
         })
         if (tfdata == false) {
            cities.value.push(v)
         }
      })
   }
   console.log(cities.value);
   total.value = data.data.total_count
   chonghui.value = true
}

const biaoliet = (data) => {
   console.log('1');
   checkList.value = []
   biaolie.value = data
}

const dakaidialog = () => {
   stencilxg.value = true
   birthDateMessage.value = ''
   //确认密码验证
   qrpassword.value = ''
   codemessage.value = ''
   addressmessage.value = ''
   phonemessage.value = ''
   dtimg.value = false
   dialogVisible.value = true
   xiugaitj.value = '添加'
   certifyForm.value = {
      tenantId: null,
      logo: {
         mime: null,
         image: null,
         filename: null
      },
      name: null,
      contactName: null,
      address: {
         province: null,
         city: null,
         district: '',
         street: null
      },

      contactPhone: null,
      businessLicense: {
         mime: null,
         image: null,
         filename: null
      },
      socialCreditCode: null,
      logoUrl: null,
      businessLicenseUrl: null,
      safePhone: null
   }
   userphone.value = {

      phone: '',
      smsCode: null,
      txId: null,
      templateAction: 'TEMPLATE_ACTION_BIND_TENANT_PHONE'
   }
   userInfo.value = {
      newPlainPassword: '',
      password: ''
   }
   Modifydata.value = null
}
const xiugai = (rs) => {
   console.log(rs);
   dialogVisible.value = true
   xiugaitj.value = '修改'
   dtimg.value = false
   certifyForm.value = {
      tenantId: rs.tenant_id,
      logo: {
         mime: null,
         image: null,
         filename: null
      },
      name: rs.name,
      contactName: rs.contact_name,
      address: {
         province: rs.address.province,
         city: rs.address.city,
         district: rs.address.district == undefined ? '' : rs.address.district,
         street: rs.address.street
      },

      contactPhone: rs.contact_phone,
      businessLicense: {
         mime: null,
         image: null,
         filename: null
      },
      socialCreditCode: rs.social_credit_code,
      logoUrl: rs.logo_url,
      businessLicenseUrl: rs.business_license_url,
      safePhone: rs.safe_phone
   }
   userphone.value =
      userphone.value = {

         phone: rs.contact_phone,
         smsCode: null,
         txId: null,
         templateAction: 'TEMPLATE_ACTION_BIND_TENANT_PHONE'
      }
   userInfo.value = {
      newPlainPassword: '',
      password: ''
   }
   Modifydata.value = rs
}

//一键剔除商户
const SecItemDelete = () => {
   Merchants.value.forEach(e => {
      checkList.value.forEach((v, i, arr) => {
         if (e.tenant_id == v.tenant_id) {
            arr.splice(i, 1)
         }
      })
   })
   if (checkList.value.length == 0) { deactivatesdialog.value = false }

   Merchants.value = []
}
//根据条件搜索商户
const SearchTenants = async (status, tenanname = null, selectlist = null) => {
   searchdata.value = status
   if (status == 'TENANT_STATUS_UNAUTH') {
      shanghuzt.value = '未认证'
   } else if (status == 'TENANT_STATUS_USING') {
      shanghuzt.value = '未激活'
   } else if (status == 'TENANT_STATUS_PENDING') {
      shanghuzt.value = '已激活'
   } else if (status == 'TENANT_STATUS_UNSET') {
      shanghuzt.value = '所有商户'
   }
   chonghui.value = false
   const rs = {
      organizationId: useid.value,
      status: status,
      treatmentId: selectlist,
      tenantName: tenanname,
      pagination: {
         offset: (offset.value - 1) * pageSize4.value,
         size: pageSize4.value
      }

   }
   const data = await SearchTenantsRequest(rs)
   chonghui.value = true
   commoddata.value = data.data.tenants
   subscriptions.value = data.data.subscriptions == undefined ? {} : data.data.subscriptions
   treatments.value = data.data.treatments == undefined ? {} : data.data.treatments
   console.log(treatments.value);

   commoddata.value?.forEach(i => {
      i.timeline = {
         start_time: subscriptions.value[i.tenant_id]?.start_time == undefined ? '' : subscriptions.value[i.tenant_id].start_time,
         end_time: subscriptions.value[i.tenant_id]?.end_time == undefined ? '' : subscriptions.value[i.tenant_id].end_time,
      }
      console.log(i);
   })
   if (commoddata.value) {
      nocommod.value = true
   } else {
      nocommod.value = false
   }
   if (commoddata.value) {
      commoddata.value?.forEach((v) => {
         let tfdata = false
         cities.value.forEach(i => {
            if (i.tenant_id == v.tenant_id) {
               tfdata = true
               return
            }
         })
         if (tfdata == false) {
            cities.value.push(v)
         }
      })
   }
   total.value = data.data.total_count == undefined ? 0 : data.data.total_count
}

//取消分配方案
const CancelTreatmentToTenants = async () => {
   console.log(checkList.value);
   let tenantids = []
   checkList.value.forEach(v => {
      tenantids.push(v.tenant_id)
   })
   const res = {
      organizationId: useid.value,
      tenantIds: tenantids
   }
   const data = await CancelTreatmentToTenantsRequest(res)
   if (data.status === 200) {
      ElMessage({
         message: '取消方案成功',
         type: 'success',
      })
      dialogconse.value = false
      checkList.value = []
      getlist()
   }
}
//处理上传图片
const logoImageBase64 = async (val) => {
   console.log(val);
   const imageInfo = val.split(/data:(.*?);base64,(.*?)/g).filter(Boolean)
   const [mime, image] = imageInfo
   certifyForm.value.logo.mime = mime
   certifyForm.value.logo.image = image
   const data = await UploadImageRequest({ image: certifyForm.value.logo })
   certifyForm.value.logoUrl = data.data.image_url
}
const yingyeImageBase64 = async (val) => {
   console.log(val);
   const imageInfo = val.split(/data:(.*?);base64,(.*?)/g).filter(Boolean)
   const [mime, image] = imageInfo
   certifyForm.value.businessLicense.mime = mime
   certifyForm.value.businessLicense.image = image
   const data = await UploadImageRequest({ image: certifyForm.value.businessLicense })
   certifyForm.value.businessLicenseUrl = data.data.image_url
}

//获取上传图片名称
const logoImageUploaded = (val) => {
   certifyForm.value.logo.filename = val.name.split('.')[0]
   console.log(certifyForm.value);
}
const yingyeImageUploaded = (val) => {
   certifyForm.value.businessLicense.filename = val.name.split('.')[0]
   console.log(certifyForm.value);
}

//点击省重置市和县
const checksheng = () => {
   certifyForm.value.address.city = certifyForm.value.address.district = ''
}

//点击市重置县
const checkshi = () => {
   certifyForm.value.address.district = ''
}

//获取省级地区
const provinceOptions = () => {
   return getAreaOptions()
}

//获取市级地区
const cityOptions = () => {
   if (certifyForm.value.address.province) {
      return getAreaOptions(TextToCode[certifyForm.value.address.province]?.code)
   }
   return []

}

//获取县级地区
const regionOptions = () => {
   if (certifyForm.value.address.city) {
      return getAreaOptions(
         TextToCode[certifyForm.value.address.province][
            certifyForm.value.address.city
         ]?.code
      )
   }
   return []
}

//发送创建短信点击事件
const SendPhone = async (data, tem) => {
   duanxindj.value = true
   if (data.phone != '') {
      if (data.phone.length != 11) {
         elmessage('请输入正确手机号码')
         duanxindj.value = false
      } else {
         const rs = {

            phone: data.phone,
            templateAction: tem,
            language: 'LANGUAGE_SIMPLIFIED_CHINESE'
         }
         await SendPhoneVerificationCode(rs).then((res) => {
            console.log(res);
            if (res?.status === 200) {
               userphone.value.txId = res.data.tx_id
               duanxindj.value = true
               timesphone.value = setInterval(function () {
                  time.value--;
                  if (time.value > 0) {
                     dxtxt.value = '重新发送' + time.value + 's';
                  } else {
                     time.value = 60;//当减到0时赋值为60
                     dxtxt.value = '获取验证码';
                     clearInterval(timesphone.value);//清除定时器
                     duanxindj.value = false
                  }
               }, 1000)
            }
            else {
               elmessage(res.data.detail)
               duanxindj.value = false
            }
         })
      }
   } else {
      elmessage('请输入手机号')
      duanxindj.value = false
   }
}

const CreateTenant = async (formEl) => {
   console.log(formEl);
   await formEl.validate(async (valid, fields) => {
      if (valid) {
         console.log('submit!')
         checkBirthDate()
         checkcodeDate()
         usepassword()
         phoneDate()
         console.log(certifyForm.value);
         if (certifyForm.value.logo.image == null && certifyForm.value.logoUrl == null) {
            elmessage('请上传商户图标')
         } else if (certifyForm.value.businessLicense.image == null && certifyForm.value.businessLicenseUrl == null) {
            elmessage('请上传营业执照')
         } else if (certifyForm.value.address.province == null || certifyForm.value.address.city == null || certifyForm.value.address.district == null || certifyForm.value.address.street == null || certifyForm.value.address.city == '' || certifyForm.value.address.province == '' || certifyForm.value.address.street == '') {
            elmessage('请输入完整地址')
            addressmessage.value = '请输入完整地址'
         } else if (userphone.value.phone.length != 11) {
            phonemessage.value = '请输入正确的手机号码'
         }
         else {
            if (xiugaitj.value == '修改') {
               delete certifyForm.value.businessLicense
               delete certifyForm.value.logo
               certifyForm.value.contactPhone = userphone.value.phone
               if (Modifydata.value.tenant_status == 'TENANT_STATUS_UNAUTH') {
                  if (userphone.value.smsCode == null || userInfo.value.newPlainPassword == '' || userInfo.value.password == '' || userInfo.value.newPlainPassword != userInfo.value.password) {
                     if (!certifyForm.value.businessLicense) {
                        certifyForm.value.businessLicense = {
                           mime: null,
                           image: null,
                           filename: null
                        }
                     }
                     if (!certifyForm.value.logo) {
                        certifyForm.value.logo = {
                           mime: null,
                           image: null,
                           filename: null
                        }
                     }
                     console.log('退出');
                     return
                  } else {
                     if (userphone.value.txId == null) {
                        elmessage('验证码错误请重新获取')
                        return
                     }
                     certifyForm.value.organizationId = getTenantId()
                     certifyForm.value.safePhone = certifyForm.value.contactPhone = userphone.value.phone
                     const rs = {
                        tenant: certifyForm.value,
                        plainPassword: userInfo.value.newPlainPassword,
                        smsCode: userphone.value.smsCode,
                        txId: userphone.value.txId
                     }

                     await ReCreateTenantRequest(rs).then((res) => {
                        console.log(res);
                        if (res.status === 200) {
                           ElMessage({
                              message: '修改成功',
                              type: 'success',
                           })
                           dialogVisible.value = false
                           getlist()
                        }
                        else {
                           elmessage(res.data.detail)
                           if (!certifyForm.value.businessLicense) {
                              certifyForm.value.businessLicense = {
                                 mime: null,
                                 image: null,
                                 filename: null
                              }
                           }
                           if (!certifyForm.value.logo) {
                              certifyForm.value.logo = {
                                 mime: null,
                                 image: null,
                                 filename: null
                              }
                           }
                        }
                     })
                  }

               } else {
                  const rs = {
                     organizationId: getTenantId(),
                     tenant: certifyForm.value,
                  }
                  await CommitEntityCertificateRequest(rs).then((res) => {
                     console.log(res);
                     if (res.status === 200) {
                        ElMessage({
                           message: '修改成功',
                           type: 'success',
                        })
                        dialogVisible.value = false
                        getlist()
                     }
                     else {
                        elmessage(res.data.detail)
                        if (!certifyForm.value.businessLicense) {
                           certifyForm.value.businessLicense = {
                              mime: null,
                              image: null,
                              filename: null
                           }
                        }
                        if (!certifyForm.value.logo) {
                           certifyForm.value.logo = {
                              mime: null,
                              image: null,
                              filename: null
                           }
                        }
                     }
                  })
               }
            } else {
               delete certifyForm.value.businessLicense
               delete certifyForm.value.logo
               if (userphone.value.smsCode == null || userInfo.value.newPlainPassword == '' || userInfo.value.password == '' || userInfo.value.newPlainPassword != userInfo.value.password) {
                  if (!certifyForm.value.businessLicense) {
                     certifyForm.value.businessLicense = {
                        mime: null,
                        image: null,
                        filename: null
                     }
                  }
                  if (!certifyForm.value.logo) {
                     certifyForm.value.logo = {
                        mime: null,
                        image: null,
                        filename: null
                     }
                  }

                  return
               }
               if (userphone.value.txId == null) {
                  elmessage('验证码错误请重新获取')
                  return
               }
               certifyForm.value.safePhone = certifyForm.value.contactPhone = userphone.value.phone
               const rs = {
                  organizationId: localStorage.getItem('hmtenantId'),
                  tenant: certifyForm.value,
                  plainPassword: userInfo.value.newPlainPassword,
                  smsCode: userphone.value.smsCode,
                  txId: userphone.value.txId
               }
               console.log(rs);
               await CreateTenantRequest(rs).then((res) => {
                  console.log(res);
                  if (res.status === 200) {
                     ElMessage({
                        message: '提审成功',
                        type: 'success',
                     })
                     dialogVisible.value = false
                     getlist()
                  }
                  else {
                     elmessage(res.data.detail)
                     if (!certifyForm.value.businessLicense) {
                        certifyForm.value.businessLicense = {
                           mime: null,
                           image: null,
                           filename: null
                        }
                     }
                     if (!certifyForm.value.logo) {
                        certifyForm.value.logo = {
                           mime: null,
                           image: null,
                           filename: null
                        }
                     }
                  }
               })
            }
         }
         certifyForm.value.businessLicense = {
            mime: null,
            image: null,
            filename: null
         }
         if (!certifyForm.value.logo) {
            certifyForm.value.logo = {
               mime: null,
               image: null,
               filename: null
            }
         }
         if (!certifyForm.value.businessLicense) {
            certifyForm.value.businessLicense = {
               mime: null,
               image: null,
               filename: null
            }
         }
      } else {
         console.log('error submit!', fields)
         checkBirthDate()
         checkcodeDate()
         usepassword()
         phoneDate()
         addressblur()

      }
   })
}
//列表多选
const handleSelectionChange = (val) => {
   liecheck.value = val
   const data = []
   liecheck.value.forEach(v => {
      console.log(v);
      data.push(v)
   })
   checkList.value = data

}
const getRowKeys = (row) => {
   return row.tenant_id;
}

const enter = (index) => {
   current.value = index
}
const leave = () => {
   current.value = null
   sandiankuang.value = null
}
const clickdian = (index) => {
   sandiankuang.value = index
}
const toprepayment = () => {
   console.log(checkList.value);
   $router.push({
      name: 'DistributionScheme',
      state: {
         keyword: JSON.stringify(checkList.value)
      }
   })
}

const topname = () => {
   settopname([{ name: '商户管理', url: null }])
   $store.commit('setactiveMenuKey', 1)
}

const copyText = async (data) => {
   let link = {}
   if (data) {
      link = jmurl.value
   } else {
      link = jmurl.value
   }
   try {
      await toClipboard(link)
      ElMessage({
         message: '成功复制到粘贴板',
         type: 'success',
      })
      dialogurl.value = false
   } catch {
      ElMessage({
         message: '复制失败',
         type: 'success',
      })
      dialogurl.value = false
   }
}

const totenant = (data) => {
   if (data.tenant_status != 'TENANT_STATUS_UNAUTH') {
      location.href = 'https://res.jinmuhealth.com/huimai/' + import.meta.env.MODE + '/tenant/index.html#/generalize?tid=' + data.tenant_id + '&oid=' + getTenantId() + '&rtoke=' + getRefreshToken()
   }
}

const closetime = () => {
   clearInterval(timestop.value)
}
function handleKeydown(event) {
   if (event.key !== "ArrowUp" && event.key !== "ArrowDown") {
      event.preventDefault();
   }
}

const ListTreatments = async () => {
   const data = await ListTreatmentsRequest({
      organizationId: getTenantId(),
      isPublish: true
   })
   console.log(data);
   selectdata.value = data.data?.treatments
}
ListTreatments()
onUnmounted(() => {
   clearInterval(timestop.value)
})
topname()
getlist()
</script>

<style lang="scss">
.wraptop {
   position: fixed;
   top: 25px;
   z-index: 999;
   left: 400px;
   display: flex;
   align-items: center;
   font-size: 12px;
   color: rgb(236, 128, 141);
}

.wrap {

   .ownertop,
   .owcnttop {
      display: flex;
      justify-content: space-between;

      .ownerleft {
         display: flex;
         flex: 0.4;
         justify-content: space-between;
         align-items: center;
         font-size: 14px;
         color: #333333;

         div {
            display: flex;

            img {
               width: 17px;
               height: 17px;
               margin-right: 5px;
            }
         }
      }

      .ownerright {
         display: flex;
         flex: 0.5;
         justify-content: flex-end;

         .ownerinp {
            width: 400px;
            height: 40px;
            margin-right: 10px;
         }
      }
   }

   .owcnttop {
      margin-top: 30px;

      .is-active {
         color: #303133;
      }

      .owcnttoplf {
         color: #999999;
         width: 330px;
         display: flex;
         justify-content: space-between;
         align-items: center;
         font-size: 15px;

         .allmer {
            font-weight: 700;
            color: #333333;
         }
      }

      .owcnttoprg {
         display: flex;
         font-size: 14px;
         justify-content: space-between;



         div {
            margin-right: 20px;
         }

         img {
            width: 20px;
            height: 20px;
            margin-right: 5px;
         }
      }
   }
}

.owbtm {
   background-color: rgba(21, 117, 238, 1);
   height: 40px;
   padding: 0 25px;
}

.ownerdata {
   .nocom {
      display: flex;
      height: 70vh;
      flex-direction: column;
      align-items: center;
      justify-content: center;

      img {
         width: 79px;
         height: 79px;
      }

      div {
         color: #666666;
         font-size: 13px;
      }
   }

   .arecom {
      display: flex;
      justify-content: flex-start;
      flex-wrap: wrap;

      .aretavc {
         margin-top: 30px;
         position: relative;
         width: calc((96% - 80px) / 4);
         height: 190px;
         background-color: rgb(255, 255, 255);
         border-radius: 10px;
         padding: 15px;
         overflow: hidden;
         margin-right: 39px;
         border: 1px solid #eeeeee;
         cursor: pointer;

         .xgshan {
            position: absolute;
            font-size: 13px;
            color: #1575EE;
            right: 15px;
         }

         .temptop {
            color: #1575EE;
            font-size: 12px;
            width: 60px;
            height: 24px;
            background-color: #eaf4fb;
            border-radius: 5px;
            display: flex;
            align-items: center;
            justify-content: center;
            position: absolute;
            top: 0px;
            right: 0px;
         }
      }

      .aretavc:nth-of-type(4n+0) {
         margin-right: 0;
      }

      .aretemplate {
         border: 1px solid #1575EE;
      }
   }
}

.flrmname {
   height: 32px;
   font-size: 16px;
   display: flex;
   align-items: center;
   justify-content: space-between;

   .checknam {
      display: flex;
      align-items: center;
      position: relative;

      .tanchushanchu {
         position: absolute;
         width: 150px;
         background-color: #fff;
         border: 1px solid #eeeeee;
         bottom: -80px;
         left: -110px;
         z-index: 99999;

         .sandian {
            height: 40px;
            line-height: 40px;
            padding-left: 20px;
            font-size: 14px;
         }

         .dianshanchu {
            color: rgb(245, 108, 108);
         }

         .dianxiugai {
            color: #c5c5c2;
         }

         .dianxiugai:hover,
         .dianshanchu:hover,
         .kexiugai:hover {
            background-color: rgba(246, 247, 248, 1);
         }
      }
   }

   img {
      width: 26px;
      height: 26px;
   }
}

.flrmess {
   font-size: 12px;
   color: #999999;
}

.usetem {
   width: 100%;
   font-size: 12px;
   color: #1575EE;
   position: absolute;
   bottom: 15px;
   display: flex;
   justify-content: space-between;
   align-items: center;

   .usetembu {
      font-size: 12px;
      width: 65px;
      height: 27px;
      margin-right: 30px;
   }
}

.under {
   color: #C1C7D0;
}

.underof {
   color: #F56C6C;
   z-index: 999;
}

.usestep {
   border: 1px solid #1575EE;
   padding: 0 3px;
}

.dialog-footer button:first-child {
   margin-right: 10px;
}

.el-form-item__label {
   width: 90px;
}

.el-dialog__body {
   padding: 20px;
}

.diaurl {
   border-top: 1px solid #C1C7D0;
   padding: 5%;

   .urlcnt {
      margin-top: 8%;

      div {
         margin-top: 10px;
      }
   }
}

.toptxt {
   text-align: center;
   font-size: 18px;
   font-weight: 400;
   color: #333333;
}

.newcnt {
   margin-top: 30px;
   padding-right: 40px;

   .dizhi {
      width: 100%;

      .dizhishi {
         display: flex;
         justify-content: space-between;
      }
   }

   .el-input {
      height: 40px;
   }

   .duanxin {
      position: relative;

      .fasong {
         display: flex;
         position: absolute;
         right: 20px;
         font-size: 14px;

         .huoqu {
            color: #1575EE;
            font-weight: 400;
         }
      }
   }
}

.ownercnt {
   position: relative;
   min-height: 75vh;

   .fenye {
      margin-top: 30px;
      position: absolute;
      bottom: -45px;
      right: 0px;
   }
}

.wrap .tableright {

   flex: 0.9;
   display: flex;
   align-items: center;

   .imgtext {
      margin-left: 20px;
      font-weight: 400;
      font-style: normal;
      font-size: 14px;
      color: rgb(193, 199, 208);
      line-height: 28px;
   }
}

.newtop {
   padding: 30px;

   .el-dialog__header {
      display: none;
   }

   .shenhe {
      color: #C1C7D0;
      font-size: 14px;
   }

}

.changjiantan {
   .el-dialog__footer {
      text-align: center;
   }

   .el-button {
      width: 103px;
      height: 40px;
   }

   .tijiaobut {
      background-color: #1575ee;
   }
}

.commbut {
   background-color: #1575ee;
   height: 40px;
}

.commer {
   .el-dialog {
      border-radius: 10px;
      height: 200px;
   }

   .el-dialog__body {
      padding: 10px 20px;
   }

   .commerdg {
      padding-top: 30px;
      border-top: 1px solid #eeeeee;
   }

   .pljh {
      color: #3CA1F8;
   }
}

.is-dark {
   width: 400px;
   background-color: #fff !important;
   border: 1px solid #1575EE !important;
}

.el-popper__arrow::before {
   background-color: #fff !important;
}

.czzytxt {
   padding: 10px;
   color: #999999;

   h3 {
      color: #333333;
      margin-bottom: 10px;
      margin-top: 10px;
      ;
   }

}

.el-dialog {
   border-radius: 10px;
}



.elbut {
   width: 103px;
   height: 40px;
}

.quern {
   background-color: #1575ee;
   color: #fff;
}

.shanchutc {
   .el-dialog__body {
      border-top: 1px solid #eee;
      padding-left: 20px;
      padding-top: 20px;
      height: 100px;
   }
}

.yucunsh {
   border: rgba(217, 0, 27, 1) 1px solid !important;
}

.statusfalt {
   padding: 10px;

   .falttop {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .flttop-left {
         font-size: 16px;
         color: rgb(245, 108, 108);
         font-weight: 700;
      }

      .flttop-right {
         font-size: 14px;
         color: rgb(21, 117, 238);
      }
   }

   .faltdata {
      margin-top: 20px;
      margin-bottom: 10px;
   }
}

.deactivates-top {
   font-size: 18px;
   font-weight: 700;
   color: #333;
   height: 60px;
   text-align: center;
   border-bottom: 1px solid #eaf4fb;
}

.deactivates-data {
   display: flex;
   flex-direction: column;
   align-items: center;

   .deactivates_datas {
      display: flex;
      height: 40px;
      align-items: center;
      width: 100%;
      justify-content: center;

      .datas-left {
         width: 100px;
         text-align: right;
         margin-right: 10px;
         flex: 0.5;
         color: #c1c7d0;
      }

      .datas-right {
         flex: 0.5;
      }
   }
}

.datas-footer {
   display: flex;
   justify-content: center;
}

.deacticates_jesi {
   margin-top: 40px;
   padding: 20px;
   text-align: center;
   color: #999999;
   font-size: 12px;
}

.datas_footertxt {
   text-align: center;
   color: #999999;
   font-size: 12px;
   margin-top: 50px;
}

.datatops {
   height: 60px;
   line-height: 60px;
}

.datatable {
   width: 100%;
   padding: 0 30px;
}

.datatabletxt {
   padding: 20px 30px 0;
   display: flex;
}

.yjtc {
   width: 60px;
   height: 20px;
   font-size: 12px;
   font-weight: 400;
}

.rimoney {
   font-size: 16px;
   color: #333333;
   margin-top: 40px;
}

.moneynumber {
   color: #1575EE;
   font-size: 30px;
}

.biaoshi {
   display: inline-block;
   padding: 0 1px;
   border: 1px solid;
}

.guoqi,
.yuqiweij {
   color: rgb(245, 108, 108);
   margin-right: 5px;
}

.yitingy {
   color: rgb(60, 161, 248);
   margin-right: 5px;
}

.dialogbhtop {
   display: flex;
   margin-top: 20px;
   margin-bottom: 20px;
   align-items: center;
   justify-content: space-evenly;

   .dialogxz {
      flex: 0.4;
      text-align: center;
      color: #C1C7D0;
      font-size: 18px;
      border-bottom: 1px solid #333;
      height: 40px;
      cursor: pointer;
   }

   .dialogxuanzhong {
      font-size: 24px;
      color: #1575EE;
      border-bottom: 3px solid #1575EE;
   }

   .ktnxx {
      display: flex;
   }
}

.ktnxx {
   display: flex;
   align-items: center;
   font-size: 20px;
   font-weight: 700;
   color: #333333;
   margin-top: 20px;

   .divlr {
      width: 150px;
      margin-left: 40px;
   }
}

.no-autofill-pwd {
   ::v-deep .el-input__inner {
      -webkit-text-security: disc !important;
   }
}

.dialog-footer-yq {
   display: flex;
   justify-content: center;

   .yqbut {
      width: 170px;
      height: 40px;
   }
}

.mouseenyr {
   border: 1px solid #cfd4db !important;
}

.tybut {
   background-color: #fff;
   border: 1px solid #1575EE;
   color: #1575EE;
}

.jhbut {
   background-color: #006be5;
   border: 1px solid #006be5;
}

.datalie {
   margin-top: 30px;

   .el-table-column--selection {
      .cell {
         width: 50px;
         height: 50px;
         display: flex;

         .el-checkbox {
            width: 50px;
            height: 50px;
            flex: 1;

            .el-checkbox__input {
               width: 50px;
               height: 50px;
               align-items: center;
               justify-content: center;
            }
         }
      }
   }

}
</style>
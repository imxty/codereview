export const rulesz = {
    name: [
        { required: true, message: '商户名不能为空', trigger: 'blur' },
        { min: 3, max: 20, message: '商户名长度数3到20位', trigger: 'blur' },
    ],
    socialCreditCode: [
        {
            required: true,
            message: '社会信用代码不能为空',
            trigger: 'blur',
        },
    ],
    contactName: [
        {
            required: true,
            message: '联系人不能为空',
            trigger: 'blur',
        },
    ],
    contactPhone: [
        {
            required: true,
            message: '联系人手机号不能为空',
            trigger: 'blur',
        }, { min: 11, max: 11, message: '请输入正确的手机号', trigger: 'blur' },
    ],
    newPlainPassword: [
        {
            required: true,
            message: '联系人手机号不能为空',
            trigger: 'blur',
        }
    ],
    type: [
        {
            type: 'array',
            required: true,
            message: 'Please select at least one activity type',
            trigger: 'change',
        },
    ],
    resource: [
        {
            required: true,
            message: 'Please select activity resource',
            trigger: 'change',
        },
    ],
    desc: [
        { required: true, message: 'Please input activity form', trigger: 'blur' },
    ],
}
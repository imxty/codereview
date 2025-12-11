-- ----------------------------
-- Table structure for tenant_entity
-- ----------------------------
CREATE TABLE `tenant_entity` (
    `tenant_id` VARCHAR(20) NOT NULL COMMENT '商户ID',
    `logo_url` VARCHAR(256) NOT NULL COMMENT '商户LOGO地址',
    `safe_phone` VARCHAR(20) NOT NULL COMMENT '安全手机号',
    `store_name` VARCHAR(256) NOT NULL COMMENT '商家名称',
    `province` VARCHAR(50) NOT NULL COMMENT '省',
    `city` VARCHAR(50) NOT NULL COMMENT '市',
    `district` VARCHAR(100) NOT NULL COMMENT '区',
    `street` VARCHAR(256) NOT NULL COMMENT '街道',
    `contact_name` VARCHAR(100) NOT NULL COMMENT '联系人姓名',
    `contact_phone` VARCHAR(20) NOT NULL COMMENT '联系人手机号',
    `social_credit_code` VARCHAR(256) NOT NULL COMMENT '社会信用代码',
    `business_license_url` VARCHAR(256) NOT NULL COMMENT '营业执照',
    `staff_count_quota` INT NOT NULL DEFAULT 6 COMMENT '员工数量上限',
    `report_sharing_status` TINYINT NOT NULL DEFAULT 1 COMMENT '报告共享状态，0是不共享，1是共享',
    `constitution_status` TINYINT NOT NULL DEFAULT 1 COMMENT '体质辩证开关状态，0是不显示，1是显示',
    `overdue` INT NOT NULL DEFAULT 7 COMMENT '复查天数',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`tenant_id`),
    INDEX `idx_tenant_search` (
        `store_name`,
        `contact_name`,
        `contact_phone`
    ),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '租户实体';
-- ----------------------------
-- Table structure for tenant_entity_revision
-- ----------------------------
CREATE TABLE `tenant_entity_revision` (
    `tenant_entity_revision_id` VARCHAR(20) NOT NULL COMMENT '商户信息审核副本ID',
    `tenant_id` VARCHAR(20) NOT NULL COMMENT '商户ID',
    `logo_url` VARCHAR(256) NOT NULL COMMENT '商户LOGO地址',
    `safe_phone` VARCHAR(20) NOT NULL COMMENT '安全手机号',
    `store_name` VARCHAR(256) NOT NULL COMMENT '商家名称',
    `province` VARCHAR(50) NOT NULL COMMENT '省',
    `city` VARCHAR(50) NOT NULL COMMENT '市',
    `district` VARCHAR(100) NOT NULL COMMENT '区',
    `street` VARCHAR(256) NOT NULL COMMENT '街道',
    `contact_name` VARCHAR(100) NOT NULL COMMENT '联系人姓名',
    `contact_phone` VARCHAR(20) NOT NULL COMMENT '联系人手机号',
    `social_credit_code` VARCHAR(256) NOT NULL COMMENT '社会信用代码',
    `business_license_url` VARCHAR(256) NOT NULL COMMENT '营业执照',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号',
    `fail_reason` VARCHAR(256) NULL DEFAULT NULL COMMENT '审核失败原因',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`tenant_entity_revision_id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '商户信息审核副本';
-- ----------------------------
-- Table structure for organization_tenant
-- ----------------------------
CREATE TABLE `organization_tenant` (
    `organization_tenant_id` VARCHAR(20) NOT NULL COMMENT '组织商户关系ID',
    `organization_id` VARCHAR(20) NOT NULL COMMENT '组织ID',
    `tenant_id` VARCHAR(20) NOT NULL COMMENT '商户ID',
    `is_activated` TINYINT NOT NULL DEFAULT 0 COMMENT '0未认证，1已认证',
    `review_status` TINYINT NOT NULL DEFAULT 0 COMMENT '0审核中，1审核成功，2审核失败',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号字段',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`organization_tenant_id`),
    INDEX `idx_tenant_id` (`tenant_id`) USING BTREE,
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '组织商户对应表';
-- ----------------------------
-- Table structure for user
-- ----------------------------
CREATE TABLE `user` (
    `user_id` VARCHAR(20) NOT NULL COMMENT '用户ID',
    `organization_id` VARCHAR(20) NOT NULL COMMENT '组织ID',
    `organization_contact_phone` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '组织联系人手机号',
    `tenant_id` VARCHAR(20) NOT NULL COMMENT '商户ID',
    `username` VARCHAR(100) NOT NULL COMMENT '用户名',
    `role_type` TINYINT NOT NULL DEFAULT 0 COMMENT '角色类型 0商户员工，1组织，2商户，3金姆管理员，4 boss',
    `nickname` VARCHAR(50) NOT NULL COMMENT '昵称',
    `phone` VARCHAR(20) NOT NULL COMMENT '手机号',
    `hashed_password` VARCHAR(100) NOT NULL COMMENT '密码',
    `tenant_limit` INT NOT NULL COMMENT '商户数量限制',
    `privilege_id` VARCHAR(20) NOT NULL COMMENT '权限组ID',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号',
    `is_activated` TINYINT NOT NULL DEFAULT 1 COMMENT '是否状态，离职或被删除为0，存在为1',
    `remark` VARCHAR(256) NOT NULL COMMENT '备注',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`user_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '用户信息';
-- ----------------------------
-- Table structure for subscription_timeline
-- ----------------------------
CREATE TABLE `subscription_timeline` (
    `subscription_timeline_id` VARCHAR(20) NOT NULL COMMENT '订阅时间线ID',
    `tenant_id` VARCHAR(20) NOT NULL COMMENT '商户ID',
    `start_time` DATETIME NOT NULL COMMENT '订阅开始时间',
    `end_time` DATETIME NOT NULL COMMENT '订阅结束时间',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号',
    `years` INT NOT NULL DEFAULT 1 COMMENT '开通年限',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据伪删除时间',
    PRIMARY KEY(`subscription_timeline_id`),
    KEY `idx_tenant_id` (`tenant_id`) USING BTREE,
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '订阅时间线表';
-- ----------------------------
-- Table structure for subscription_period
-- ----------------------------
CREATE TABLE `subscription_period` (
    `subscription_period_id` VARCHAR(20) NOT NULL COMMENT '订阅周期ID,xid',
    `tenant_id` VARCHAR(20) NOT NULL COMMENT '租户ID',
    `organization_id` VARCHAR(20) NOT NULL COMMENT '组织ID',
    `organization_name` VARCHAR(100) NOT NULL COMMENT '组织名称',
    `tenant_name` VARCHAR(100) NOT NULL COMMENT '商户名称',
    `expired_time` DATETIME NOT NULL COMMENT '到期日期',
    `user_id` VARCHAR(20) NOT NULL COMMENT '操作人ID',
    `contact_name` VARCHAR(100) NOT NULL COMMENT '联系人',
    `contact_phone` VARCHAR(20) NOT NULL COMMENT '联系人电话',
    `years` INT NOT NULL DEFAULT 1 COMMENT '开通年限',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号字段',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY(`subscription_period_id`),
    KEY `idx_tenant_id` (`tenant_id`) USING BTREE,
    INDEX `idx_subscription_period_search` (
        `organization_id`,
        `tenant_id`,
        `contact_name`,
        `contact_phone`
    ),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '订阅周期表';
-- ----------------------------
-- Table structure for privilege_policy
-- ----------------------------
CREATE TABLE `privilege_policy` (
    `privilege_id` VARCHAR(20) NOT NULL COMMENT '权限组ID',
    `privilege_name` VARCHAR(256) NOT NULL COMMENT '权限组名称',
    `remark` VARCHAR(256) NOT NULL COMMENT '备注',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号字段',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据伪删除时间',
    PRIMARY KEY(`privilege_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '权限组';
-- ----------------------------
-- Table structure for privilege_preset
-- ----------------------------
CREATE TABLE `privilege_preset` (
    `preset_id` VARCHAR(20) NOT NULL COMMENT '策略集ID',
    `preset_name` VARCHAR(100) NOT NULL COMMENT '策略集名称',
    `policies` TEXT NOT NULL COMMENT '策略',
    `remark` VARCHAR(100) NOT NULL COMMENT '备注',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号字段',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY(`preset_id`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '页面权限策略集表';
-- ----------------------------
-- Table structure for privilege_policy_page
-- ----------------------------
CREATE TABLE `privilege_policy_page` (
    `privilege_id` VARCHAR(20) NOT NULL COMMENT '权限组ID',
    `preset_id` VARCHAR(20) NOT NULL COMMENT '策略集ID',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据伪删除时间',
    PRIMARY KEY(`privilege_id`, `preset_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '页面权限';
-- ----------------------------
-- Table structure for privilege_policy_data
-- ----------------------------
CREATE TABLE `privilege_policy_data` (
    `data_privilege_id` VARCHAR(20) NOT NULL COMMENT '数据权限ID',
    `privilege_id` VARCHAR(20) NOT NULL COMMENT '权限组ID',
    `organization_id` VARCHAR(20) NOT NULL COMMENT '组织ID',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号字段',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据伪删除时间',
    PRIMARY KEY(`data_privilege_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '数据权限';
-- ----------------------------
-- Table structure for tenant_treatment
-- ----------------------------
CREATE TABLE `tenant_treatment` (
    `tenant_treatment_id` VARCHAR(20) NOT NULL COMMENT '商户方案ID',
    `treatment_id` VARCHAR(20) NOT NULL COMMENT '方案ID,xid',
    `tenant_id` VARCHAR(20) NOT NULL COMMENT '商户ID',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号字段',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据伪删除时间',
    PRIMARY KEY(`tenant_treatment_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '商户方案表';
-- ----------------------------
-- Table structure for role_menu
-- ----------------------------
CREATE TABLE `privilege_menu` (
    `privilege_id` VARCHAR(20) NOT NULL COMMENT '权限组ID',
    `menu_id` VARCHAR(20) NOT NULL COMMENT '菜单ID',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号字段',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`privilege_id`, `menu_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '角色菜单表';
-- ----------------------------
-- Table structure for menu    
-- ----------------------------
CREATE TABLE `menu` (
    `menu_id` VARCHAR(20) NOT NULL COMMENT '菜单ID',
    `title` VARCHAR(50) NOT NULL COMMENT '标题',
    `path` VARCHAR(256) NOT NULL COMMENT '路径，点击之后访问的API，sample:"UserAPI,CreateUser"',
    `sort` INT NOT NULL COMMENT '菜单排序',
    `visible` TINYINT NOT NULL COMMENT '是否可见',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号字段',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`menu_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '菜单信息';
-- ----------------------------
-- Table structure for feedback
-- ----------------------------
CREATE TABLE `feedback` (
    `feedback_id` VARCHAR(20) NOT NULL COMMENT '意见反馈ID,xid',
    `tenant_id` VARCHAR(20) NOT NULL COMMENT '商户ID',
    `phone` VARCHAR(20) NOT NULL COMMENT '手机号码',
    `content` VARCHAR(255) NOT NULL COMMENT '反馈意见内容',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`feedback_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '反馈意见表';
-- ----------------------------
-- Table structure for faq
-- ----------------------------
CREATE TABLE `faq` (
    `faq_id` VARCHAR(20) NOT NULL COMMENT '常见问题ID,xid',
    `question` VARCHAR(500) NOT NULL COMMENT '常见问题',
    `answer` TEXT NOT NULL COMMENT '回答',
    PRIMARY KEY (`faq_id`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '常见问题表';
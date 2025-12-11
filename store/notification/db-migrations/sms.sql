-- ----------------------------
-- Table structure for sms
-- ----------------------------
CREATE TABLE `sms` (
    `sms_id` VARCHAR(20) NOT NULL COMMENT '验证码ID,xid',
    `tx_id` VARCHAR(20) NOT NULL COMMENT '验证码上下文ID,xid',
    `phone` VARCHAR(20) NOT NULL COMMENT '手机号',
    `sms_status` TINYINT(4) NOT NULL COMMENT '验证码状态，0表示待定，1表示发送中，2表示发送成功，3表示发送失败',
    `template_action` VARCHAR(255) NOT NULL COMMENT '短信用途',
    `platform_type` VARCHAR(20) NOT NULL COMMENT '运营商 Aliyun,Tencent',
    `template_param` VARCHAR(50) NOT NULL COMMENT '验证码参数',
    `language` VARCHAR(50) NOT NULL COMMENT '语言',
    `serial_number` VARCHAR(20) NOT NULL COMMENT '序列号',
    `sms_error_log` VARCHAR(255) NOT NULL COMMENT '错误信息',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号字段',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY(`sms_id`),
    KEY `idx_phone` (`phone`) USING BTREE,
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '短信验证码信息';

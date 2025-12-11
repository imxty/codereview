-- ----------------------------
-- Table structure for customer
-- ----------------------------
CREATE TABLE `customer` (
    `customer_id` VARCHAR(20) NOT NULL COMMENT '常客ID',
    `organization_id` VARCHAR(50) NOT NULL COMMENT '组织ID',
    `tenant_id` VARCHAR(50) NOT NULL COMMENT '租户ID',
    `staff_id` VARCHAR(20) NOT NULL COMMENT '员工ID',
    `nickname` VARCHAR(50) NOT NULL COMMENT '常客昵称',
    `initial` VARCHAR(1) NOT NULL COMMENT '昵称首字母，大写，如果不是字母则为#',
    `gender` TINYINT NOT NULL COMMENT '0是无效的性别，1是未设置性别，2是男，3是女',
    `phone` VARCHAR(20) NOT NULL COMMENT '手机号码',
    `birthday` DATE NOT NULL DEFAULT '1000-01-01' COMMENT '生日',
    `height` INT NOT NULL COMMENT '身高',
    `weight` INT NOT NULL COMMENT '体重',
    `pmh` VARCHAR(255) NOT NULL COMMENT '既往病史',
    `remark` VARCHAR(255) NOT NULL COMMENT '备注',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '版本号',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`customer_id`),
    INDEX `idx_tenant_id` (`tenant_id`) USING BTREE,
    INDEX `idx_phone` (`phone`) USING BTREE,
    INDEX `idx_initial` (`initial`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '常客表';
-- ----------------------------
-- Table structure for customer_last_status
-- ----------------------------
CREATE TABLE `customer_last_status` (
    `customer_id` VARCHAR(20) NOT NULL COMMENT '常客ID',
    `tenant_id` VARCHAR(50) NOT NULL COMMENT '租户ID',
    `report_id` VARCHAR(20) NOT NULL COMMENT '报告ID',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '版本号',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`customer_id`),
    INDEX `idx_tenant_id` (`tenant_id`) USING BTREE,
    INDEX `idx_updated_at` (`updated_at`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '常客最近测量记录表';

-- ----------------------------
-- Table structure for recommended_product
-- ----------------------------
CREATE TABLE `recommended_product` (
    `recommended_product_id` VARCHAR(20) NOT NULL COMMENT '商品ID,xid',
    `organization_id` VARCHAR(50) NOT NULL COMMENT '组织ID',
    `product_image_url` VARCHAR(500) NOT NULL COMMENT '商品图片地址',
    `product_name` VARCHAR(50) NOT NULL COMMENT '商品名称',
    `product_type` TINYINT NOT NULL COMMENT '2 表示理疗服务,3 表示中成药,4 表示保健品,5 表示营养食品',
    `product_status` TINYINT NOT NULL COMMENT '2 表示使用中,3 表示待使用,4 表示未使用,5 表示已删除',
    `product_introduction` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '商品介绍',
    `remarks` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '商品备注',
    `drug_name` VARCHAR(50) NOT NULL COMMENT '药品名称',
    `drug_classification` TINYINT NOT NULL COMMENT '0 表示处方药, 1 表示非处方药',
    `approved_number` VARCHAR(10) NOT NULL COMMENT '准字号',
    `drug_validity_period` DATETIME NULL DEFAULT NULL COMMENT '药品准效期',
    `symptoms` VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '适用症候的key',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '版本号',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY(`recommended_product_id`),
    INDEX `idx_organization_id` (`organization_id`) USING BTREE,
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '商品表';
-- ----------------------------
-- Table structure for treatment_rev
-- ----------------------------
CREATE TABLE `treatment_rev` (
    `treatment_rev_id` VARCHAR(20) NOT NULL COMMENT '方案修订版ID,xid',
    `treatment_id` VARCHAR(20) NOT NULL COMMENT '方案ID,xid',
    `treatment_name` VARCHAR(20) NOT NULL COMMENT '方案名称',
    `organization_id` VARCHAR(50) NOT NULL COMMENT '组织ID',
    `is_published` TINYINT NOT NULL DEFAULT 0 COMMENT '0 表示未发布, 1 表示已发布,默认为0',
    `treatment_status` TINYINT NOT NULL DEFAULT 0 COMMENT '方案状态, 2 表示草稿,3 表示审核中,4 表示已过审,5 表示未过审',
    `remarks` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '方案备注',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '版本号',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY(`treatment_rev_id`),
    KEY `idx_treatment_id` (`treatment_id`) USING BTREE,
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '方案修订版表';
-- ----------------------------
-- Table structure for treatment
-- ----------------------------
CREATE TABLE `treatment` (
    `treatment_id` VARCHAR(20) NOT NULL COMMENT '方案ID,xid',
    `organization_id` VARCHAR(50) NOT NULL COMMENT '组织ID',
    `latest_rev` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '方案最新版本id',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '版本号',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY(`treatment_id`),
    INDEX `idx_organization_id` (`organization_id`) USING BTREE,
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '方案表';
-- ----------------------------
-- Table structure for treatment_item
-- ----------------------------
CREATE TABLE `treatment_item` (
    `treatment_item_id` VARCHAR(20) NOT NULL COMMENT '方案配置项目ID,xid',
    `treatment_rev_id` VARCHAR(20) NOT NULL COMMENT '方案版本ID,xid',
    `treatment_item_type` TINYINT NOT NULL COMMENT '方案配置类型, 2 表示风险疾病方案,3 表示脏腑辩证方案,4 表示理疗方案',
    `symptom` VARCHAR(50) NOT NULL COMMENT '症候名',
    `recommended_product_id` VARCHAR(20) NOT NULL COMMENT '商品ID,xid',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY(`treatment_item_id`),
    KEY `idx_treatment_rev_id` (`treatment_rev_id`) USING BTREE,
    KEY `idx_recommended_product_id` (`recommended_product_id`) USING BTREE,
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '方案配置项目表';
-- ----------------------------
-- Table structure for recommended_product_statistics
-- ----------------------------
CREATE TABLE `recommended_product_statistics` (
    `recommended_product_statistics_id` VARCHAR(20) NOT NULL COMMENT '推荐商品统计ID,xid',
    `organization_id` VARCHAR(50) NOT NULL COMMENT '组织ID',
    `tenant_id` VARCHAR(50) NOT NULL COMMENT '商户ID',
    `report_id` VARCHAR(20) NOT NULL COMMENT '报告ID,xid',
    `recommended_product_id` VARCHAR(20) NOT NULL COMMENT '商品ID,xid',
    `symptom` VARCHAR(50) NOT NULL COMMENT '症候的key',
    `exposed_at` DATETIME NOT NULL COMMENT '曝光时间',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    PRIMARY KEY(`recommended_product_statistics_id`),
    KEY `idx_tenant_id` (`tenant_id`) USING BTREE,
    KEY `idx_recommended_product_id` (`recommended_product_id`) USING BTREE
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '推荐商品统计表';

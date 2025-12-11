-- ----------------------------
-- Table structure for review_issue
-- ----------------------------
CREATE TABLE `review_issue` (
    `review_issue_id` VARCHAR(20) NOT NULL COMMENT '审核工单ID',
    `submitter_organization_id` VARCHAR(50) NOT NULL COMMENT '提审人ID(组织/商户)',
    `submitter_tenant_id` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '商户ID,只有资质审核才可能有商户ID',
    `submit_time` DATETIME NOT NULL COMMENT '提审时间',
    `reviewer_user_id` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '审核员用户ID',
    `acceptance_time` DATETIME DEFAULT NULL COMMENT '接单时间',
    `review_status` TINYINT NOT NULL DEFAULT 2 COMMENT '审核状态,0待审核，1审核中，2审核失败，3审核成功',
    `target_type` TINYINT NOT NULL COMMENT '审核类型，0商户资质审核，1商品推荐方案审核',
    `is_cancel` TINYINT NOT NULL DEFAULT 0 COMMENT '是否被用户取消，0否，1是',
    `target_rev_id` VARCHAR(20) NOT NULL COMMENT '审核记录目标ID',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号字段',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`review_issue_id`),
    INDEX `idx_deleted_at` (`deleted_at`),
    INDEX `idx_review_search` (`submitter_organization_id`, `review_status`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '审核工单';
-- ----------------------------
-- Table structure for review_result
-- ----------------------------
CREATE TABLE `review_result` (
    `review_result_id` VARCHAR(20) NOT NULL COMMENT '审核结论ID',
    `review_issue_id` VARCHAR(20) NOT NULL COMMENT '审核工单ID',
    `review_time` DATETIME NOT NULL COMMENT '审核时间',
    `reviewer_user_id` VARCHAR(20) NOT NULL COMMENT '审核员用户ID',
    `result` TINYINT NOT NULL COMMENT '0不通过，1通过',
    `comment` VARCHAR(2000) NOT NULL COMMENT '审核意见',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号字段',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`review_result_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '审核结论表';
-- ----------------------------
-- Table structure for review_notification
-- ----------------------------
CREATE TABLE `review_notification` (
    `review_notification_id` VARCHAR(20) NOT NULL COMMENT '审核通知消息ID',
    `organization_id` VARCHAR(50) NOT NULL COMMENT '组织ID',
    `tenant_id` VARCHAR(50) NOT NULL COMMENT '租户ID',
    `notification_title` VARCHAR(200) NOT NULL COMMENT '通知标题',
    `target_type` TINYINT NOT NULL COMMENT '审核类型，0商户资质审核，1商品推荐方案审核',
    `result` TINYINT NOT NULL COMMENT '0不通过，1通过',
    `comment` VARCHAR(2000) NOT NULL COMMENT '审核意见',
    `has_read` TINYINT NOT NULL DEFAULT 0 COMMENT '是否已读，0否，1是',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号字段',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`review_notification_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '审核通知消息';

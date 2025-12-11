-- ----------------------------
-- Table structure for report
-- ----------------------------
CREATE TABLE `report` (
    `report_id` VARCHAR(20) NOT NULL COMMENT '报告ID',
    `organization_id` VARCHAR(50) NOT NULL COMMENT '组织ID',
    `tenant_id` VARCHAR(50) NOT NULL COMMENT '租户ID',
    `staff_id` VARCHAR(20) NOT NULL COMMENT '员工ID',
    `device_id` VARCHAR(20) NOT NULL COMMENT '设备ID',
    `is_customer` TINYINT NOT NULL COMMENT '是否是常客报告，1是，0否',
    `customer_id` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '常客ID',
    `age` VARCHAR(50) NOT NULL COMMENT '年龄',
    `hand` TINYINT NOT NULL COMMENT '测量手，2是左手，3是右手',
    `gender` TINYINT NOT NULL COMMENT '性别，2是男，3是女',
    `remarks` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '普通备注，给测量者看的',
    `staff_remarks` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '员工备注',
    `heart_rate` INT NOT NULL DEFAULT 0 COMMENT '心率，单位bpm',
    `spo` INT NOT NULL COMMENT '血氧值',
    `risks` TEXT NOT NULL COMMENT '风险预估，json格式',
    `c0` INT NOT NULL DEFAULT 0 COMMENT '心包经',
    `c1` INT NOT NULL DEFAULT 0 COMMENT '肝经',
    `c2` INT NOT NULL DEFAULT 0 COMMENT '肾经',
    `c3` INT NOT NULL DEFAULT 0 COMMENT '脾经',
    `c4` INT NOT NULL DEFAULT 0 COMMENT '肺经',
    `c5` INT NOT NULL DEFAULT 0 COMMENT '胃经',
    `c6` INT NOT NULL DEFAULT 0 COMMENT '胆经',
    `c7` INT NOT NULL DEFAULT 0 COMMENT '膀胱经',
    `physical_dialectics` VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '体质辩证,json字符串',
    `dirty_dialectic` VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '脏腑辩证,json字符串',
    `dirty_dialectic_flag` INT NOT NULL DEFAULT 0 COMMENT '脏腑辩证查询判断，用于方便脏腑辨证查找，按照key的先后顺序排序',
    `dietary_advice` VARCHAR(3000) NOT NULL DEFAULT '' COMMENT '调理建议模块,json字符串',
    `measurement_judgment` VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '测量异常模块,json字符串',
    `f0` INT NOT NULL DEFAULT 0 COMMENT '阴虚指数',
    `f1` INT NOT NULL DEFAULT 0 COMMENT '阳虚指数',
    `f2` INT NOT NULL DEFAULT 0 COMMENT '湿气指数',
    `f3` INT NOT NULL DEFAULT 0 COMMENT '血瘀指数',
    `f4` INT NOT NULL DEFAULT 0 COMMENT '阴阳平衡指数',
    `is_report_complete` TINYINT NOT NULL DEFAULT 0 COMMENT '是否获取到金姆报告,0表示未完善，1表示完善',
    `is_stress_state` TINYINT NOT NULL COMMENT '是否是应激态',
    `stress_state` VARCHAR(1000) NOT NULL COMMENT '应激态内容',
    `tenant_treatment_rev_id` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '租户商品推荐方案版本号ID',
    `face_image_url` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '人脸图片URL',
    `tongue_image_url` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '舌面图片URL',
    `tongue_face_report_status` TINYINT NOT NULL DEFAULT 0 COMMENT '舌面分析报告状态，0表示未获取，1表示获取成功，2表示获取失败',
    `tongue_face_report` TEXT NOT NULL COMMENT '舌面分析报告',
    `inquiry_diagnosis` TEXT COMMENT '问诊信息',
    `rev` INT NOT NULL DEFAULT 0 COMMENT '数据版本号字段',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`report_id`),
    INDEX `idx_deleted_at` (`deleted_at`),
    INDEX `idx_staff_report` (`staff_id`),
    INDEX `idx_report_search` (
        `organization_id`,
        `tenant_id`,
        `customer_id`,
        `is_customer`,
        `is_report_complete`,
        `dirty_dialectic_flag`,
        `created_at`
    )
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '健康报告';
-- ----------------------------
-- Table structure for shared_report
-- ----------------------------
CREATE TABLE `shared_report` (
    `token` VARCHAR(50) NOT NULL COMMENT '分享报告token',
    `report_id` BINARY(20) NOT NULL COMMENT '报告ID',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`token`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '分享报告';
-- ----------------------------
-- Table structure for report_address
-- ----------------------------
CREATE TABLE `report_address` (
    `report_id` VARCHAR(20) NOT NULL COMMENT '报告ID,与平台报告ID相同',
    `mac` VARCHAR(30) NOT NULL COMMENT '设备MAC地址',
    `country_code` TINYINT DEFAULT 0 COMMENT '国家编码',
    `country` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '国家',
    `province` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '省',
    `city` VARCHAR(30) NOT NULL DEFAULT '' COMMENT '市',
    `district` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '区',
    `full_address` TEXT NOT NULL COMMENT '完整地址',
    `ip` VARCHAR(30) NOT NULL DEFAULT '' COMMENT 'ip',
    `latitude` VARCHAR(30) NOT NULL DEFAULT '' COMMENT '纬度',
    `longitude` VARCHAR(30) NOT NULL DEFAULT '' COMMENT '经度',
    `remark` VARCHAR(30) NOT NULL DEFAULT '' COMMENT '备注',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    PRIMARY KEY (`report_id`),
    INDEX `idx_mac` (`mac`),
    INDEX `idx_address` (`province`, `city`, `district`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '报告位置信息';
-- ----------------------------
-- Table structure for inquiry_diagnosis
-- ----------------------------
CREATE TABLE `inquiry_diagnosis` (
    `inquiry_id` VARCHAR(20) NOT NULL COMMENT '问题ID',
    `content` VARCHAR(200) NOT NULL COMMENT '内容',
    `is_multiple_selection` TINYINT NOT NULL COMMENT '是否是多选',
    `items` TEXT NOT NULL COMMENT '问题选项，json格式',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`inquiry_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '问诊信息';

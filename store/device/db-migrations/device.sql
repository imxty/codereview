-- ----------------------------
-- Table structure for device
-- ----------------------------
CREATE TABLE `device` (
    `device_id` VARCHAR(20) NOT NULL COMMENT '设备ID,xid',
    `mac` VARCHAR(30) NOT NULL COMMENT '设备MAC地址',
    `sn` VARCHAR(20) NOT NULL COMMENT '设备SN号',
    `model` VARCHAR(100) NOT NULL COMMENT '设备型号',
    `is_available` TINYINT NOT NULL COMMENT '不可用=0,可用=1',
    `remark` VARCHAR(255) NOT NULL COMMENT '备注',
    `created_at` DATETIME NOT NULL COMMENT '数据记录创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '数据记录更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '数据记录伪删除时间',
    PRIMARY KEY (`device_id`),
    UNIQUE KEY `idx_mac` (`mac`) USING BTREE,
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '设备信息表';

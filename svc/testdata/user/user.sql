use huimaibao;
truncate table `tenant_entity`;
-- 已激活
INSERT INTO `tenant_entity`(
        `tenant_id`,
        `logo_url`,
        `safe_phone`,
        `store_name`,
        `province`,
        `city`,
        `district`,
        `street`,
        `contact_name`,
        `contact_phone`,
        `social_credit_code`,
        `business_license_url`,
        `staff_count_quota`,
        `report_sharing_status`,
        `constitution_status`,
        `overdue`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'N1391310003',
        '',
        '13368168788',
        '金姆大药房',
        '江苏省',
        '常州市',
        '天宁区',
        '红梅街道',
        '王先生',
        '13358168788',
        '123456',
        '',
        '6',
        1,
        1,
        20,
        0,
        now(),
        now()
    );
-- 未激活
INSERT INTO `tenant_entity`(
        `tenant_id`,
        `logo_url`,
        `safe_phone`,
        `store_name`,
        `province`,
        `city`,
        `district`,
        `street`,
        `contact_name`,
        `contact_phone`,
        `social_credit_code`,
        `business_license_url`,
        `staff_count_quota`,
        `report_sharing_status`,
        `constitution_status`,
        `overdue`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'N1391310001',
        '',
        '13700000001',
        '金姆大药房',
        '江苏省',
        '常州市',
        '天宁区',
        '红梅街道',
        '王先生',
        '13700000001',
        '123456',
        '',
        '6',
        1,
        1,
        20,
        0,
        now(),
        now()
    );
-- 停用
INSERT INTO `tenant_entity`(
        `tenant_id`,
        `logo_url`,
        `safe_phone`,
        `store_name`,
        `province`,
        `city`,
        `district`,
        `street`,
        `contact_name`,
        `contact_phone`,
        `social_credit_code`,
        `business_license_url`,
        `staff_count_quota`,
        `report_sharing_status`,
        `constitution_status`,
        `overdue`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'N1391310002',
        '',
        '13700000002',
        '金姆大药房',
        '江苏省',
        '常州市',
        '天宁区',
        '红梅街道',
        '王先生',
        '13700000002',
        '123456',
        '',
        '6',
        1,
        1,
        20,
        0,
        now(),
        now()
    );
-- 已激活
INSERT INTO `tenant_entity`(
        `tenant_id`,
        `logo_url`,
        `safe_phone`,
        `store_name`,
        `province`,
        `city`,
        `district`,
        `street`,
        `contact_name`,
        `contact_phone`,
        `social_credit_code`,
        `business_license_url`,
        `staff_count_quota`,
        `report_sharing_status`,
        `constitution_status`,
        `overdue`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'N1391310004',
        '',
        '13700000003',
        '金姆大药房',
        '江苏省',
        '常州市',
        '天宁区',
        '红梅街道',
        '王先生',
        '13700000003',
        '123456',
        '',
        '6',
        0,
        0,
        0,
        0,
        now(),
        now()
    );
-- 未认证
INSERT INTO `tenant_entity`(
        `tenant_id`,
        `logo_url`,
        `safe_phone`,
        `store_name`,
        `province`,
        `city`,
        `district`,
        `street`,
        `contact_name`,
        `contact_phone`,
        `social_credit_code`,
        `business_license_url`,
        `staff_count_quota`,
        `report_sharing_status`,
        `constitution_status`,
        `overdue`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'N1391310005',
        '',
        '13700000004',
        '金姆大药房',
        '江苏省',
        '常州市',
        '天宁区',
        '红梅街道',
        '王先生',
        '13700000004',
        '123456',
        '',
        '6',
        1,
        1,
        20,
        0,
        now(),
        now()
    );
truncate table `organization_tenant`;
INSERT INTO `organization_tenant` (
        `organization_tenant_id`,
        `organization_id`,
        `tenant_id`,
        `is_activated`,
        `review_status`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'c02krl6vvhfg2kd7c10g',
        'N139131',
        'N1391310001',
        1,
        3,
        0,
        now(),
        now()
    );
INSERT INTO `organization_tenant` (
        `organization_tenant_id`,
        `organization_id`,
        `tenant_id`,
        `is_activated`,
        `review_status`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'c02krtmvvhfg43nlt8mg',
        'N139131',
        'N1391310002',
        1,
        3,
        0,
        now(),
        now()
    );
INSERT INTO `organization_tenant` (
        `organization_tenant_id`,
        `organization_id`,
        `tenant_id`,
        `is_activated`,
        `review_status`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'brlgfeript36jkpqd05g',
        'N139131',
        'N1391310003',
        1,
        3,
        0,
        now(),
        now()
    );
INSERT INTO `organization_tenant` (
        `organization_tenant_id`,
        `organization_id`,
        `tenant_id`,
        `is_activated`,
        `review_status`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'brtagdript36udqrctp0',
        'N139131',
        'N1391310004',
        0,
        3,
        0,
        now(),
        now()
    );
INSERT INTO `organization_tenant` (
        `organization_tenant_id`,
        `organization_id`,
        `tenant_id`,
        `is_activated`,
        `review_status`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'brta6hjipt3600qlfadg',
        'N139131',
        'N1391310005',
        0,
        3,
        0,
        now(),
        now()
    );
INSERT INTO `organization_tenant` (
        `organization_tenant_id`,
        `organization_id`,
        `tenant_id`,
        `is_activated`,
        `review_status`,
        `rev`,
        `created_at`,
        `updated_at`,
        `deleted_at`
    )
VALUES (
        'cjfde07ng1s92f92t810',
        'N139131',
        'N139131DDFX',
        1,
        2,
        0,
        '2023-08-18 02:27:12.741',
        '2023-08-18 02:27:12.741',
        NULL
    );
truncate table `user`;
-- 组织 
INSERT INTO `user` (
        `user_id`,
        `organization_id`,
        `tenant_id`,
        `username`,
        `role_type`,
        `nickname`,
        `phone`,
        `hashed_password`,
        `tenant_limit`,
        `privilege_id`,
        `rev`,
        `is_activated`,
        `remark`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'cjens4vng1s2taacdhkg',
        'N139131',
        '',
        'lvguan',
        1,
        '',
        '13700000001',
        '9fcefc0080d894e83ca7d360ce5ccd9ead2c5d8a80a10f9fa9698510aaba865a',
        10,
        'c03ad1mvvhfhdjkuga6g',
        0,
        1,
        '',
        now(),
        now()
    );
-- 组织 
INSERT INTO `user` (
        `user_id`,
        `organization_id`,
        `tenant_id`,
        `username`,
        `role_type`,
        `nickname`,
        `phone`,
        `hashed_password`,
        `tenant_limit`,
        `privilege_id`,
        `rev`,
        `is_activated`,
        `remark`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'cjeo1qvng1s302t0odeg',
        'N254758',
        '',
        'jinmu',
        1,
        '',
        '13700000002',
        '9fcefc0080d894e83ca7d360ce5ccd9ead2c5d8a80a10f9fa9698510aaba865a',
        10,
        'c03ad1mvvhfhdjkuga6g',
        0,
        1,
        '',
        now(),
        now()
    );
-- 员工
INSERT INTO `user` (
        `user_id`,
        `organization_id`,
        `tenant_id`,
        `username`,
        `role_type`,
        `nickname`,
        `phone`,
        `hashed_password`,
        `tenant_limit`,
        `privilege_id`,
        `rev`,
        `is_activated`,
        `remark`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'cjesju7ng1s6r62r75d0',
        '',
        'N1391310003',
        '',
        0,
        '小金',
        '13700000004',
        '9fcefc0080d894e83ca7d360ce5ccd9ead2c5d8a80a10f9fa9698510aaba865a',
        10,
        'c03ad1mvvhfhdjkuga6g',
        0,
        1,
        '',
        now(),
        now()
    );
-- 员工
INSERT INTO `user` (
        `user_id`,
        `organization_id`,
        `tenant_id`,
        `username`,
        `role_type`,
        `nickname`,
        `phone`,
        `hashed_password`,
        `tenant_limit`,
        `privilege_id`,
        `rev`,
        `is_activated`,
        `remark`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'c18qjhbipt3406jsqft0',
        '',
        'N1391310003',
        '',
        0,
        '小金',
        '13710000004',
        '9fcefc0080d894e83ca7d360ce5ccd9ead2c5d8a80a10f9fa9698510aaba865a',
        10,
        'c03ad1mvvhfhdjkuga6g',
        0,
        1,
        '',
        now(),
        now()
    );
-- 员工
INSERT INTO `user` (
        `user_id`,
        `organization_id`,
        `tenant_id`,
        `username`,
        `role_type`,
        `nickname`,
        `phone`,
        `hashed_password`,
        `tenant_limit`,
        `privilege_id`,
        `rev`,
        `is_activated`,
        `remark`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'c077ng6vvhfq3gm6lgk0',
        '',
        'N1391310004',
        '',
        0,
        '小金',
        '13710000001',
        '9fcefc0080d894e83ca7d360ce5ccd9ead2c5d8a80a10f9fa9698510aaba865a',
        2,
        'c03ad1mvvhfhdjkuga6g',
        0,
        0,
        '',
        now(),
        now()
    ),
    (
        'c077ng6vvhfq3gm6lgk1',
        '',
        'N1391310004',
        '',
        0,
        '小金',
        '13710000001',
        '9fcefc0080d894e83ca7d360ce5ccd9ead2c5d8a80a10f9fa9698510aaba865a',
        2,
        'c03ad1mvvhfhdjkuga6g',
        0,
        0,
        '',
        now(),
        now()
    );
INSERT INTO `user` (
        `user_id`,
        `organization_id`,
        `tenant_id`,
        `username`,
        `role_type`,
        `nickname`,
        `phone`,
        `hashed_password`,
        `tenant_limit`,
        `privilege_id`,
        `rev`,
        `remark`,
        `is_activated`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'crsc6kb0bteg717dpqs0',
        '',
        'N1391310004',
        '小金',
        0,
        '小金',
        '13700000001',
        '9fcefc0080d894e83ca7d360ce5ccd9ead2c5d8a80a10f9fa9698510aaba865a',
        2,
        '',
        0,
        '',
        true,
        now(),
        now()
    );
-- 商户
INSERT INTO `user` (
        `user_id`,
        `organization_id`,
        `tenant_id`,
        `username`,
        `role_type`,
        `nickname`,
        `phone`,
        `hashed_password`,
        `tenant_limit`,
        `privilege_id`,
        `rev`,
        `is_activated`,
        `remark`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'cjfde07ng1s92f92t81g',
        'N139131',
        'N1391310003',
        '',
        2,
        '王先生',
        '13358168788',
        '9fcefc0080d894e83ca7d360ce5ccd9ead2c5d8a80a10f9fa9698510aaba865a',
        10,
        '',
        0,
        1,
        '',
        now(),
        now()
    );
-- 审核
INSERT INTO `user` (
        `user_id`,
        `organization_id`,
        `tenant_id`,
        `username`,
        `role_type`,
        `nickname`,
        `phone`,
        `hashed_password`,
        `tenant_limit`,
        `privilege_id`,
        `rev`,
        `is_activated`,
        `remark`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'c02krpmvvhfg3d0f2e00',
        '',
        '',
        'jinmu',
        3,
        '',
        '13358168788',
        '9fcefc0080d894e83ca7d360ce5ccd9ead2c5d8a80a10f9fa9698510aaba865a',
        1,
        'c03ad1mvvhfhdjkuga6g',
        0,
        1,
        '',
        now(),
        now()
    );
INSERT INTO `user` (
        `user_id`,
        `organization_id`,
        `tenant_id`,
        `username`,
        `role_type`,
        `nickname`,
        `phone`,
        `hashed_password`,
        `tenant_limit`,
        `privilege_id`,
        `rev`,
        `is_activated`,
        `remark`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'c02krpmvvhfg3d0f2e02',
        '',
        '',
        'boss',
        4,
        '',
        '13358168789',
        '9fcefc0080d894e83ca7d360ce5ccd9ead2c5d8a80a10f9fa9698510aaba865a',
        1,
        'c03ad1mvvhfhdjkuga6g',
        0,
        1,
        '',
        now(),
        now()
    );
truncate table `tenant_entity_revision`;
INSERT INTO `tenant_entity_revision`(
        `tenant_entity_revision_id`,
        `tenant_id`,
        `logo_url`,
        `safe_phone`,
        `store_name`,
        `province`,
        `city`,
        `district`,
        `street`,
        `contact_name`,
        `contact_phone`,
        `social_credit_code`,
        `business_license_url`,
        `rev`,
        `fail_reason`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'c07nlkmvvhflv1eni9l0',
        'N1391310003',
        '',
        '13358168788',
        '金姆大药房',
        '江苏省',
        '常州市',
        '天宁区',
        '红梅街道',
        '王先生',
        '13358168788',
        '123456',
        '',
        0,
        '',
        now(),
        now()
    );
INSERT INTO `tenant_entity_revision`(
        `tenant_entity_revision_id`,
        `tenant_id`,
        `logo_url`,
        `safe_phone`,
        `store_name`,
        `province`,
        `city`,
        `district`,
        `street`,
        `contact_name`,
        `contact_phone`,
        `social_credit_code`,
        `business_license_url`,
        `rev`,
        `fail_reason`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'c2ll0f8i7qkl01ikbea0',
        'N1391310004',
        '',
        '13358168788',
        '金姆大药房',
        '江苏省',
        '常州市',
        '天宁区',
        '红梅街道',
        '王先生',
        '13358168788',
        '123456',
        '',
        0,
        '',
        now(),
        now()
    );
INSERT INTO `tenant_entity_revision` (
        `tenant_entity_revision_id`,
        `tenant_id`,
        `logo_url`,
        `safe_phone`,
        `store_name`,
        `province`,
        `city`,
        `district`,
        `street`,
        `contact_name`,
        `contact_phone`,
        `social_credit_code`,
        `business_license_url`,
        `rev`,
        `fail_reason`,
        `created_at`,
        `updated_at`,
        `deleted_at`
    )
VALUES (
        'cjfde07ng1s92f92t820',
        'N1391310003',
        '',
        '13358168788',
        '金姆大药房',
        '江苏省',
        '常州市',
        '天宁区',
        '红梅街道',
        '王先生',
        '13358168788',
        '123456',
        '',
        0,
        '',
        now(),
        now(),
        NULL
    );
truncate table `subscription_timeline`;
INSERT INTO `subscription_timeline` (
        `subscription_timeline_id`,
        `tenant_id`,
        `start_time`,
        `end_time`,
        `years`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'c08gi5evvhftnct6v5bg',
        "N1391310003",
        date_sub(UTC_TIMESTAMP(), interval 300 day),
        date_add(UTC_TIMESTAMP(), interval 250 day),
        1,
        0,
        now(),
        now()
    ),
    (
        'c2ll0f8i7qkl01ikbea0',
        "N1391310001",
        date_sub(UTC_TIMESTAMP(), interval 400 day),
        date_sub(UTC_TIMESTAMP(), interval 100 day),
        1,
        0,
        now(),
        now()
    ),
    (
        'c1htsdgi7qkm69d3r31g',
        "N139131DDFX",
        date_sub(UTC_TIMESTAMP(), interval 400 day),
        date_sub(UTC_TIMESTAMP(), interval 100 day),
        1,
        0,
        now(),
        now()
    );
truncate table `subscription_period`;
INSERT INTO `subscription_period`(
        `subscription_period_id`,
        `tenant_id`,
        `organization_id`,
        `organization_name`,
        `tenant_name`,
        `expired_time`,
        `user_id`,
        `contact_name`,
        `contact_phone`,
        `years`,
        `rev`,
        `created_at`,
        `updated_at`
    ) -- 已使用
VALUES (
        'c08grl6vvhfunce65tng',
        "N1391310003",
        'N139131',
        'lvguan',
        '金姆大药房',
        date_add(UTC_TIMESTAMP(), interval 200 day),
        'c02krpmvvhfg3d0f2e00',
        'jinmu',
        '13358168788',
        1,
        0,
        now(),
        now()
    );
truncate table `privilege_policy`;
INSERT INTO `privilege_policy` (
        `privilege_id`,
        `privilege_name`,
        `remark`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'c03ad1mvvhfhdjkuga6g',
        '审核员',
        '',
        0,
        now(),
        now()
    );
truncate table `privilege_preset`;
INSERT INTO `privilege_preset` (
        `preset_id`,
        `preset_name`,
        `policies`,
        `remark`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'c03nueuvvhfhgn8i22pg',
        '1',
        '1',
        '',
        0,
        now(),
        now()
    );
truncate table `privilege_policy_page`;
INSERT INTO `privilege_policy_page` (
        `privilege_id`,
        `preset_id`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'c03ad1mvvhfhdjkuga6g',
        'c03nueuvvhfhgn8i22pg',
        now(),
        now()
    );
truncate table `privilege_policy_data`;
INSERT INTO `privilege_policy_data` (
        `data_privilege_id`,
        `privilege_id`,
        `organization_id`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'c03nuauvvhfhg2ej5avg',
        'c03ad1mvvhfhdjkuga6g',
        'N139131',
        0,
        now(),
        now()
    );
truncate table `tenant_treatment`;

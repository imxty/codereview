use huimaibao;
truncate table `recommended_product`;
INSERT INTO `recommended_product` (
        `recommended_product_id`,
        `organization_id`,
        `product_image_url`,
        `product_name`,
        `product_type`,
        `product_status`,
        `product_introduction`,
        `remarks`,
        `drug_name`,
        `drug_classification`,
        `approved_number`,
        `drug_validity_period`,
        `symptoms`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'cji59svng1s3rjh9sn0g',
        'N139131',
        'a/b/c',
        '草晶华玉米须破壁饮片',
        3,
        2,
        '利水水肿，降血压',
        '',
        '草晶华玉米须破壁饮片',
        0,
        'Y20160328',
        now(),
        'JB0001',
        0,
        now(),
        now()
    ),
    (
        'cji59svng1s3rjh9sn1g',
        'N139131',
        'a/b/c',
        '爱乐维',
        4,
        2,
        '',
        '',
        '爱乐维',
        0,
        'Y20160328',
        '1970-01-02 00:00:01',
        'JB0001,JB0002',
        0,
        now(),
        now()
    ),
    (
        'cji59svng1s3rjh9sn2g',
        'N139131',
        'a/b/c',
        '合生元益生菌冲剂',
        5,
        2,
        '',
        '',
        '合生元益生菌冲剂',
        0,
        'Y20160328',
        '1970-01-02 00:00:01',
        'JB0001',
        0,
        now(),
        now()
    ),
    (
        'cji59t7ng1s3rjh9sn3g',
        'N139131',
        'a/b/c',
        '足疗',
        2,
        3,
        '',
        '',
        '',
        0,
        '',
        '1970-01-02 00:00:01',
        'JB0001',
        0,
        now(),
        now()
    );
INSERT INTO `recommended_product` (
        `recommended_product_id`,
        `organization_id`,
        `product_image_url`,
        `product_name`,
        `product_type`,
        `product_status`,
        `product_introduction`,
        `remarks`,
        `drug_name`,
        `drug_classification`,
        `approved_number`,
        `drug_validity_period`,
        `symptoms`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'cjimtnnng1s04cn7rhh0',
        'N139131',
        'a/b/c',
        '草晶华玉米须破壁饮片',
        3,
        4,
        '利水水肿，降血压',
        '',
        '草晶华玉米须破壁饮片',
        0,
        'Y20160328',
        '2023-08-23 02:28:45.831',
        'JB0001',
        0,
        now(),
        now()
    ),
    (
        'cjimtnnng1s04cn7rhi0',
        'N139131',
        'a/b/c',
        '爱乐维',
        4,
        4,
        '',
        '',
        '爱乐维',
        0,
        'Y20160328',
        '1970-01-02 00:00:01',
        'JB0001,JB0002',
        0,
        now(),
        now()
    ),
    (
        'cjimtnnng1s04cn7rhj0',
        'N139131',
        'a/b/c',
        '合生元益生菌冲剂',
        5,
        4,
        '',
        '',
        '草晶华玉米须破壁饮片',
        0,
        'Y20160328',
        '1970-01-02 00:00:01',
        'JB0001',
        0,
        now(),
        now()
    ),
    (
        'cjimtnnng1s04cn7rhk0',
        'N139131',
        'a/b/c',
        '足疗',
        2,
        4,
        '',
        '',
        '',
        0,
        '',
        '1970-01-02 00:00:01',
        'JB0001',
        0,
        now(),
        now()
    );
truncate table `treatment`;
INSERT INTO `treatment` (
        `treatment_id`,
        `organization_id`,
        `latest_rev`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'cji63inng1s82vrdrcqg',
        'N139131',
        'cji63inng1s82vrdrcr0',
        0,
        now(),
        now()
    );
INSERT INTO `treatment` (
        `treatment_id`,
        `organization_id`,
        `latest_rev`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'cji6j8fng1s1u4usae5g',
        'N139131',
        'cji6j8fng1s1u4usae60',
        0,
        now(),
        now()
    ),
    (
        'cji79tfng1scotmaodgg',
        'N139131',
        'cji79tfng1scotmaodh0',
        0,
        now(),
        now()
    ),
    (
        'cjin7lfng1s4nc430oeg',
        'N139131',
        'cjin7lfng1s4nc430of0',
        0,
        now(),
        now()
    );
truncate table `treatment_rev`;
INSERT INTO `treatment_rev` (
        `treatment_rev_id`,
        `treatment_id`,
        `treatment_name`,
        `remarks`,
        `organization_id`,
        `is_published`,
        `treatment_status`,
        `rev`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'cji63inng1s82vrdrcr0',
        'cji63inng1s82vrdrcqg',
        '养生小达人',
        '',
        'N139131',
        1,
        4,
        0,
        now(),
        now()
    ),
    (
        'cji6j8fng1s1u4usae60',
        'cji6j8fng1s1u4usae5g',
        '养生小达人',
        '',
        'N139131',
        1,
        3,
        0,
        now(),
        now()
    ),
    (
        'cji79tfng1scotmaodh0',
        'cji79tfng1scotmaodgg',
        '养生小达人',
        '',
        'N139131',
        0,
        2,
        0,
        now(),
        now()
    ),
    (
        'cjin7lfng1s4nc430of0',
        'cjin7lfng1s4nc430oeg',
        '养生小达人',
        '',
        'N139131',
        0,
        5,
        0,
        now(),
        now()
    );
truncate table `treatment_item`;
INSERT INTO `treatment_item` (
        `treatment_item_id`,
        `treatment_rev_id`,
        `treatment_item_type`,
        `symptom`,
        `recommended_product_id`,
        `created_at`,
        `updated_at`,
        `deleted_at`
    )
VALUES (
        'cji63inng1s82vrdrcrg',
        'cji63inng1s82vrdrcr0',
        3,
        'Z0001',
        'cji59svng1s3rjh9sn0g',
        '2023-08-22 07:20:42.068',
        '2023-08-22 07:20:42.068',
        NULL
    ),
    (
        'cji63inng1s82vrdrcs0',
        'cji63inng1s82vrdrcr0',
        3,
        'Z0001',
        'cji59svng1s3rjh9sn1g',
        '2023-08-22 07:20:42.068',
        '2023-08-22 07:20:42.068',
        NULL
    ),
    (
        'cji63inng1s82vrdrcsg',
        'cji63inng1s82vrdrcr0',
        2,
        'JB0001',
        'cji59svng1s3rjh9sn0g',
        '2023-08-22 07:20:42.068',
        '2023-08-22 07:20:42.068',
        NULL
    ),
    (
        'cji63inng1s82vrdrct0',
        'cji63inng1s82vrdrcr0',
        2,
        'JB0001',
        'cji59svng1s3rjh9sn1g',
        '2023-08-22 07:20:42.068',
        '2023-08-22 07:20:42.068',
        NULL
    ),
    (
        'cji63inng1s82vrdrctg',
        'cji63inng1s82vrdrcr0',
        4,
        'SQ',
        'cji59svng1s3rjh9sn0g',
        '2023-08-22 07:20:42.068',
        '2023-08-22 07:20:42.068',
        NULL
    ),
    (
        'cji63inng1s82vrdrcu0',
        'cji63inng1s82vrdrcr0',
        4,
        'SQ',
        'cji59svng1s3rjh9sn1g',
        '2023-08-22 07:20:42.068',
        '2023-08-22 07:20:42.068',
        NULL
    );
INSERT INTO `treatment_item` (
        `treatment_item_id`,
        `treatment_rev_id`,
        `treatment_item_type`,
        `symptom`,
        `recommended_product_id`,
        `created_at`,
        `updated_at`,
        `deleted_at`
    )
VALUES (
        'cji6j8fng1s1u4usae6g',
        'cji6j8fng1s1u4usae60',
        3,
        'Z0001',
        'cji59svng1s3rjh9sn0g',
        '2023-08-22 07:54:09.145',
        '2023-08-22 07:54:09.145',
        NULL
    ),
    (
        'cji6j8fng1s1u4usae70',
        'cji6j8fng1s1u4usae60',
        3,
        'Z0001',
        'cji59svng1s3rjh9sn1g',
        '2023-08-22 07:54:09.145',
        '2023-08-22 07:54:09.145',
        NULL
    ),
    (
        'cji6j8fng1s1u4usae7g',
        'cji6j8fng1s1u4usae60',
        2,
        'JB0001',
        'cji59svng1s3rjh9sn0g',
        '2023-08-22 07:54:09.145',
        '2023-08-22 07:54:09.145',
        NULL
    ),
    (
        'cji6j8fng1s1u4usae80',
        'cji6j8fng1s1u4usae60',
        2,
        'JB0001',
        'cji59svng1s3rjh9sn1g',
        '2023-08-22 07:54:09.145',
        '2023-08-22 07:54:09.145',
        NULL
    ),
    (
        'cji6j8fng1s1u4usae8g',
        'cji6j8fng1s1u4usae60',
        4,
        'SQ',
        'cji59svng1s3rjh9sn0g',
        '2023-08-22 07:54:09.145',
        '2023-08-22 07:54:09.145',
        NULL
    ),
    (
        'cji6j8fng1s1u4usae90',
        'cji6j8fng1s1u4usae60',
        4,
        'SQ',
        'cji59svng1s3rjh9sn1g',
        '2023-08-22 07:54:09.145',
        '2023-08-22 07:54:09.145',
        NULL
    ),
    (
        'cjin7lfng1s4nc430ofg',
        'cjin7lfng1s4nc430of0',
        3,
        'Z0001',
        'cji59svng1s3rjh9sn0g',
        '2023-08-23 02:49:57.629',
        '2023-08-23 02:49:57.629',
        NULL
    ),
    (
        'cjin7lfng1s4nc430og0',
        'cjin7lfng1s4nc430of0',
        3,
        'Z0001',
        'cji59svng1s3rjh9sn1g',
        '2023-08-23 02:49:57.629',
        '2023-08-23 02:49:57.629',
        NULL
    ),
    (
        'cjin7lfng1s4nc430ogg',
        'cjin7lfng1s4nc430of0',
        2,
        'JB0001',
        'cji59svng1s3rjh9sn0g',
        '2023-08-23 02:49:57.629',
        '2023-08-23 02:49:57.629',
        NULL
    ),
    (
        'cjin7lfng1s4nc430oh0',
        'cjin7lfng1s4nc430of0',
        2,
        'JB0001',
        'cji59svng1s3rjh9sn1g',
        '2023-08-23 02:49:57.629',
        '2023-08-23 02:49:57.629',
        NULL
    ),
    (
        'cjin7lfng1s4nc430ohg',
        'cjin7lfng1s4nc430of0',
        4,
        'SQ',
        'cji59svng1s3rjh9sn0g',
        '2023-08-23 02:49:57.629',
        '2023-08-23 02:49:57.629',
        NULL
    ),
    (
        'cjin7lfng1s4nc430oi0',
        'cjin7lfng1s4nc430of0',
        4,
        'SQ',
        'cji59svng1s3rjh9sn1g',
        '2023-08-23 02:49:57.629',
        '2023-08-23 02:49:57.629',
        NULL
    );
truncate table `recommended_product_statistics`;
INSERT INTO `recommended_product_statistics` (
        `recommended_product_statistics_id`,
        `organization_id`,
        `tenant_id`,
        `report_id`,
        `recommended_product_id`,
        `symptom`,
        `exposed_at`,
        `created_at`
    )
VALUES (
        'c0q7cj0i7qkp49st6q7g',
        'N139131',
        'N139131003',
        'c0q8r80i7qkp75vgpbo0',
        'cji59svng1s3rjh9sn1g',
        'JB0004',
        date_sub(now(), interval 1 DAY),
        now()
    ),
    (
        'c0q8sl8i7qkp84lq1l40',
        'N139131',
        'N139131003',
        'c0q8sroi7qkp8dppj4cg',
        'cji59svng1s3rjh9sn1g',
        'JB0004',
        date_sub(now(), interval 10 DAY),
        now()
    ),
    (
        'c0q8t0oi7qkp8ookbib0',
        'N139131',
        'N139131003',
        'c0q8sroi7qkp8dppj4cg',
        'cji59svng1s3rjh9sn1g',
        'JB0004',
        now(),
        now()
    ),
    (
        'c0qakf0i7qks6j0khftg',
        'N139131',
        'N139131003',
        'c0q8sroi7qkp8dppj4cg',
        'cji59svng1s3rjh9sn2g',
        'JB0002',
        date_sub(now(), interval 30 DAY),
        now()
    ),
    (
        'c0qakroi7qks74k7ai30',
        'N139131',
        'N139131003',
        'c0q8sroi7qkp8dppj4cg',
        'cji59svng1s3rjh9sn2g',
        'JB0001',
        now(),
        now()
    ),
    (
        'c0qal9gi7qks7l2n5ga0',
        'N139131',
        'N139131003',
        'c0q8sroi7qkp8dppj4cg',
        'cji59svng1s3rjh9sn2g',
        'JB0002',
        date_sub(now(), interval 31 DAY),
        now()
    ),
    (
        'c0qb7j0i7qksoltpoemg',
        'N139131',
        'N139131003',
        'c0q8sroi7qkp8dppj4cg',
        'c0qb8c0i7qksqvdjgpsg',
        'JB0004',
        now(),
        now()
    ),
    (
        'c0qb7j0i7qksolttokjm',
        'N139131',
        'N139131003',
        'c0q8sroi7qkp8dppj4cg',
        'c0qb8c0i7qksqvdjupoo',
        'JB0004',
        date_sub(now(), interval 23 DAY),
        now()
    ),
    (
        'c0qb7j0i7qksolfdokpp',
        'N139131',
        'N139131003',
        'c0q8sroi7qkp8dppj4cg',
        'c0u8pooi7qkobpht46lg',
        'JB0004',
        date_sub(now(), interval 22 DAY),
        now()
    ),
    (
        'c1klhlgi7qkkcbf8f360',
        'N139131',
        'N139131003',
        'c0q8r80i7qkp75vgpbo0',
        'c0ea280i7qkhcdaobad0',
        'JB0004',
        date_sub(now(), interval 23 DAY),
        now()
    ),
    (
        'c1kliogi7qkke5e5ro2g',
        'N139131',
        'N139131003',
        'c0q8r80i7qkp75vgpbo0',
        'c0ea280i7qkhcdaobad0',
        'JB0004',
        date_sub(now(), interval 23 DAY),
        now()
    ),
    (
        'c1klivoi7qkkedp10cs0',
        'N139131',
        'N139131003',
        'c0q8r80i7qkp75vgpbo0',
        'c0ea280i7qkhcdaobad0',
        'JB0004',
        date_sub(now(), interval 16 DAY),
        now()
    ),
    (
        'c22def0i7qkhefumin0g',
        'N139131',
        'N139131003',
        'c22dehgi7qkhev4ug27g',
        'cji59svng1s3rjh9sn1g',
        'JB0001',
        date_sub(now(), interval 16 DAY),
        now()
    ),
    (
        'c22des8i7qkhffhmonag',
        'N139131',
        'N139131003',
        'c22dehgi7qkhev4ug27g',
        'cji59svng1s3rjh9sn1g',
        'JB0002',
        date_sub(now(), interval 10 DAY),
        now()
    ),
    (
        'c235ad0i7qkgshc5es10',
        'N139131',
        'N139131003',
        'c22dehgi7qkhev4ug27g',
        'cji59svng1s3rjh9sn1g',
        'JB0003',
        date_sub(now(), interval 10 DAY),
        now()
    ),
    (
        'c22fta8i7qkiivptk00g',
        'N139131',
        'N139131003',
        'c22dehgi7qkhev4ug27g',
        'cji59svng1s3rjh9sn1g',
        'JB0002',
        date_sub(now(), interval 10 DAY),
        now()
    ),
    (
        'c22ftvgi7qkim1umqgeg',
        'N139131',
        'N139131003',
        'c22dehgi7qkhev4ug27g',
        'cji59svng1s3rjh9sn1g',
        'JB0003',
        date_sub(now(), interval 10 DAY),
        now()
    );

use huimaibao;
truncate table `device`;
INSERT INTO `device` (
        `device_id`,
        `mac`,
        `sn`,
        `model`,
        `is_available`,
        `remark`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'c071jr6vvhfsr0viotng',
        '194055512950961',
        'JMCNUK38D269ED6CD6',
        'JM1300',
        '1',
        '',
        now(),
        now()
    );
INSERT INTO `device` (
        `device_id`,
        `mac`,
        `sn`,
        `model`,
        `is_available`,
        `remark`,
        `created_at`,
        `updated_at`
    )
VALUES (
        'brlgf4bipt36iqdm8v6g',
        '62476371455185',
        'JMCNUK38D269ED6CD1',
        'JM1003',
        '0',
        '',
        now(),
        now()
    );

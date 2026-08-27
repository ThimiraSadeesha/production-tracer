CREATE OR REPLACE PROCEDURE work_order_save(
    IN p_work_order_reference VARCHAR(255),
    IN p_customer_name VARCHAR(255),
    IN p_po_no VARCHAR(255),
    IN p_title VARCHAR(255),
    IN p_quantity BIGINT,
    IN p_delivery_date DATETIME,
    IN p_actor VARCHAR(255),
    IN p_items LONGTEXT
)
BEGIN
    DECLARE v_work_order_id BIGINT;
    DECLARE v_count INT DEFAULT 0;
    DECLARE v_i INT DEFAULT 0;

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            ROLLBACK;
            RESIGNAL;
        END;

    START TRANSACTION;


    INSERT INTO tbl_work_order (work_order_reference, customer_name, po_no, title, quantity,
                                delivery_date, status, created_by, created_at, updated_at)
    VALUES (p_work_order_reference, p_customer_name, p_po_no, p_title, p_quantity,
            p_delivery_date, 'Pending', p_actor, NOW(3), NOW(3))
    ON DUPLICATE KEY UPDATE customer_name = VALUES(customer_name),
                            po_no         = VALUES(po_no),
                            title         = VALUES(title),
                            quantity      = VALUES(quantity),
                            delivery_date = VALUES(delivery_date),
                            updated_by    = p_actor,
                            updated_at    = NOW(3),
                            id            = LAST_INSERT_ID(id);

    SET v_work_order_id = LAST_INSERT_ID();


    DELETE FROM tbl_work_order_item WHERE work_order_id = v_work_order_id;

    SET v_count = COALESCE(JSON_LENGTH(p_items), 0);
    WHILE v_i < v_count
        DO
            INSERT INTO tbl_work_order_item (work_order_id, item_reference, item_code, item_name, uom, stock_uom,
                                             item_group, delivery_date, remarks, quantity, item_status,
                                             created_by, created_at, updated_at)
            VALUES (v_work_order_id,
                    JSON_UNQUOTE(JSON_EXTRACT(p_items, CONCAT('$[', v_i, '].itemReference'))),
                    JSON_UNQUOTE(JSON_EXTRACT(p_items, CONCAT('$[', v_i, '].itemCode'))),
                    JSON_UNQUOTE(JSON_EXTRACT(p_items, CONCAT('$[', v_i, '].itemName'))),
                    JSON_UNQUOTE(JSON_EXTRACT(p_items, CONCAT('$[', v_i, '].uom'))),
                    JSON_UNQUOTE(JSON_EXTRACT(p_items, CONCAT('$[', v_i, '].stockUom'))),
                    JSON_UNQUOTE(JSON_EXTRACT(p_items, CONCAT('$[', v_i, '].itemGroup'))),
                    NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(p_items, CONCAT('$[', v_i, '].deliveryDate'))), 'null'),
                           ''),
                    JSON_UNQUOTE(JSON_EXTRACT(p_items, CONCAT('$[', v_i, '].description'))),
                    CAST(JSON_EXTRACT(p_items, CONCAT('$[', v_i, '].qty')) AS SIGNED),
                    'Pending',
                    p_actor, NOW(3), NOW(3));
            SET v_i = v_i + 1;
        END WHILE;
    COMMIT;

    SELECT v_work_order_id        AS id,
           p_work_order_reference AS work_order_reference,
           v_count                AS item_count;
END

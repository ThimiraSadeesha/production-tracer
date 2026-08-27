DROP PROCEDURE IF EXISTS work_order_save;
CREATE PROCEDURE work_order_save(
    IN work_order_val JSON,
    IN created_by_val VARCHAR(255)
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE total INT DEFAULT 0;
    DECLARE inserted_ids JSON DEFAULT JSON_ARRAY();
    DECLARE v_payload JSON;
    DECLARE v_items JSON;
    DECLARE v_item_i INT DEFAULT 0;
    DECLARE v_item_count INT DEFAULT 0;
    DECLARE v_work_order_id BIGINT;
    DECLARE v_reference VARCHAR(255);
    DECLARE v_customer_name VARCHAR(255);
    DECLARE v_po_no VARCHAR(255);
    DECLARE v_title VARCHAR(255);
    DECLARE v_quantity BIGINT;
    DECLARE v_delivery_date DATETIME;

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            GET DIAGNOSTICS CONDITION 1
                @sqlstate = RETURNED_SQLSTATE,
                @errno = MYSQL_ERRNO,
                @message_text = MESSAGE_TEXT;
            ROLLBACK;
            SELECT CONCAT('Error: [', @sqlstate, '] ', @message_text) AS error_message;
        END;

    SET v_payload = work_order_val;


    SET total = IFNULL(JSON_LENGTH(v_payload), 0);

    START TRANSACTION;

    WHILE i < total
        DO

            SET v_reference = NULLIF(
                    JSON_UNQUOTE(
                            JSON_EXTRACT(
                                    v_payload,
                                    CONCAT('$[', i, '].name')
                            )
                    ),
                    'null'
                              );

            SET v_customer_name = NULLIF(
                    JSON_UNQUOTE(
                            JSON_EXTRACT(
                                    v_payload,
                                    CONCAT('$[', i, '].customer_name')
                            )
                    ),
                    'null'
                                  );

            SET v_po_no = NULLIF(
                    JSON_UNQUOTE(
                            JSON_EXTRACT(
                                    v_payload,
                                    CONCAT('$[', i, '].po_no')
                            )
                    ),
                    'null'
                          );

            SET v_title = NULLIF(
                    JSON_UNQUOTE(
                            JSON_EXTRACT(
                                    v_payload,
                                    CONCAT('$[', i, '].title')
                            )
                    ),
                    'null'
                          );

            SET v_quantity = CAST(
                    COALESCE(
                            NULLIF(
                                    JSON_UNQUOTE(
                                            JSON_EXTRACT(
                                                    v_payload,
                                                    CONCAT('$[', i, '].total_qty')
                                            )
                                    ),
                                    'null'
                            ),
                            0
                    ) AS SIGNED
                             );

            SET v_delivery_date = NULLIF(
                    JSON_UNQUOTE(
                            JSON_EXTRACT(
                                    v_payload,
                                    CONCAT('$[', i, '].delivery_date')
                            )
                    ),
                    'null'
                                  );

            IF v_reference IS NULL OR v_reference = '' THEN
                SIGNAL SQLSTATE '45000'
                    SET MESSAGE_TEXT = 'Work order name is required.';
            END IF;

            INSERT INTO tbl_work_order (work_order_reference,
                                        customer_name,
                                        po_no,
                                        title,
                                        quantity,
                                        delivery_date,
                                        status,
                                        created_by,
                                        created_at,
                                        updated_at)
            VALUES (v_reference,
                    v_customer_name,
                    v_po_no,
                    v_title,
                    v_quantity,
                    v_delivery_date,
                    'Pending',
                    created_by_val,
                    NOW(3),
                    NOW(3))
            ON DUPLICATE KEY UPDATE customer_name = VALUES(customer_name),
                                    po_no         = VALUES(po_no),
                                    title         = VALUES(title),
                                    quantity      = VALUES(quantity),
                                    delivery_date = VALUES(delivery_date),
                                    updated_by    = created_by_val,
                                    updated_at    = NOW(3),
                                    id            = LAST_INSERT_ID(id);

            SET v_work_order_id = LAST_INSERT_ID();

            SET inserted_ids = JSON_ARRAY_APPEND(
                    inserted_ids,
                    '$',
                    v_work_order_id
                               );

            DELETE
            FROM tbl_work_order_item
            WHERE work_order_id = v_work_order_id;

            SET v_items = JSON_EXTRACT(
                    v_payload,
                    CONCAT('$[', i, '].items')
                          );

            SET v_item_count = IFNULL(
                    JSON_LENGTH(v_items),
                    0
                               );

            SET v_item_i = 0;

            WHILE v_item_i < v_item_count
                DO

                    INSERT INTO tbl_work_order_item (work_order_id,
                                                     item_reference,
                                                     item_code,
                                                     item_name,
                                                     uom,
                                                     stock_uom,
                                                     item_group,
                                                     delivery_date,
                                                     remarks,
                                                     quantity,
                                                     item_status,
                                                     created_by,
                                                     created_at,
                                                     updated_at)
                    VALUES (v_work_order_id,
                            NULLIF(
                                    JSON_UNQUOTE(
                                            JSON_EXTRACT(
                                                    v_items,
                                                    CONCAT(
                                                            '$[',
                                                            v_item_i,
                                                            '].custom_sales_order_item_reference'
                                                    )
                                            )
                                    ),
                                    'null'
                            ),
                            NULLIF(
                                    JSON_UNQUOTE(
                                            JSON_EXTRACT(
                                                    v_items,
                                                    CONCAT('$[', v_item_i, '].item_code')
                                            )
                                    ),
                                    'null'
                            ),
                            NULLIF(
                                    JSON_UNQUOTE(
                                            JSON_EXTRACT(
                                                    v_items,
                                                    CONCAT('$[', v_item_i, '].item_name')
                                            )
                                    ),
                                    'null'
                            ),
                            NULLIF(
                                    JSON_UNQUOTE(
                                            JSON_EXTRACT(
                                                    v_items,
                                                    CONCAT('$[', v_item_i, '].uom')
                                            )
                                    ),
                                    'null'
                            ),
                            NULLIF(
                                    JSON_UNQUOTE(
                                            JSON_EXTRACT(
                                                    v_items,
                                                    CONCAT('$[', v_item_i, '].stock_uom')
                                            )
                                    ),
                                    'null'
                            ),
                            NULLIF(
                                    JSON_UNQUOTE(
                                            JSON_EXTRACT(
                                                    v_items,
                                                    CONCAT('$[', v_item_i, '].item_group')
                                            )
                                    ),
                                    'null'
                            ),
                            NULLIF(
                                    JSON_UNQUOTE(
                                            JSON_EXTRACT(
                                                    v_items,
                                                    CONCAT('$[', v_item_i, '].delivery_date')
                                            )
                                    ),
                                    'null'
                            ),
                            NULLIF(
                                    JSON_UNQUOTE(
                                            JSON_EXTRACT(
                                                    v_items,
                                                    CONCAT('$[', v_item_i, '].description')
                                            )
                                    ),
                                    'null'
                            ),
                            CAST(
                                    COALESCE(
                                            NULLIF(
                                                    JSON_UNQUOTE(
                                                            JSON_EXTRACT(
                                                                    v_items,
                                                                    CONCAT('$[', v_item_i, '].qty')
                                                            )
                                                    ),
                                                    'null'
                                            ),
                                            0
                                    ) AS SIGNED
                            ),
                            'Pending',
                            created_by_val,
                            NOW(3),
                            NOW(3));

                    SET v_item_i = v_item_i + 1;
                END WHILE;
            SET i = i + 1;
        END WHILE;
    COMMIT;

    SELECT total        AS savedTotalWorkOrder,
           inserted_ids AS insertedIds;

END;

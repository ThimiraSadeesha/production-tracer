CREATE OR REPLACE PROCEDURE work_order_find(
    IN p_reference VARCHAR(255),
    IN p_customer VARCHAR(255),
    IN p_status VARCHAR(100),
    IN p_cursor BIGINT,
    IN p_limit INT
)
BEGIN
    IF p_limit IS NULL OR p_limit <= 0 THEN
        SET p_limit = 20;
    END IF;

    SELECT wo.id                   AS id,
           wo.work_order_reference AS workOrderReference,
           wo.customer_name        AS customerName,
           wo.po_no                AS poNo,
           wo.title                AS title,
           wo.quantity             AS quantity,
           wo.delivery_date        AS deliveryDate,
           wo.status               AS status
    FROM tbl_work_order wo
    WHERE wo.deleted_at IS NULL
      AND (p_reference IS NULL OR p_reference = '' OR wo.work_order_reference LIKE CONCAT('%', p_reference, '%'))
      AND (p_customer IS NULL OR p_customer = '' OR wo.customer_name LIKE CONCAT('%', p_customer, '%'))
      AND (p_status IS NULL OR p_status = '' OR wo.status LIKE CONCAT('%', p_status, '%'))
      AND (p_cursor = -1 OR wo.id < p_cursor)
    ORDER BY wo.id DESC
    LIMIT p_limit;
END

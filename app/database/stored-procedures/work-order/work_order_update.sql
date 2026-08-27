CREATE OR REPLACE PROCEDURE work_order_update(
    IN p_id BIGINT,
    IN p_customer_name VARCHAR(255),
    IN p_po_no VARCHAR(255),
    IN p_title VARCHAR(255),
    IN p_quantity BIGINT,
    IN p_delivery_date DATETIME,
    IN p_status VARCHAR(100),
    IN p_actor VARCHAR(255)
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            GET DIAGNOSTICS CONDITION 1
                @sqlstate = RETURNED_SQLSTATE,
                @errno = MYSQL_ERRNO,
                @message_text = MESSAGE_TEXT;
            SELECT CONCAT('Error: [', @sqlstate, '] ', @message_text) AS error_message;
            ROLLBACK;
        END;

    START TRANSACTION;

    UPDATE tbl_work_order
    SET customer_name = IFNULL(p_customer_name, customer_name),
        po_no         = IFNULL(p_po_no, po_no),
        title         = IFNULL(p_title, title),
        quantity      = IFNULL(p_quantity, quantity),
        delivery_date = IFNULL(p_delivery_date, delivery_date),
        status        = IFNULL(p_status, status),
        updated_by    = p_actor,
        updated_at    = NOW(3)
    WHERE id = p_id;

    COMMIT;
    SELECT p_id AS workOrderId;
END

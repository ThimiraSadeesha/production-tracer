CREATE OR REPLACE PROCEDURE work_order_get(
    IN p_id BIGINT
)
BEGIN
    SELECT wo.id                   AS id,
           wo.work_order_reference AS workOrderReference,
           wo.customer_name        AS customerName,
           wo.po_no                AS poNo,
           wo.title                AS title,
           wo.quantity             AS quantity,
           wo.delivery_date        AS deliveryDate,
           wo.status               AS status,
           wo.no_of_breakdowns     AS noOfBreakdowns,
           IFNULL(
                   (SELECT JSON_ARRAYAGG(
                                   JSON_OBJECT(
                                           'id', i.id,
                                           'itemReference', i.item_reference,
                                           'itemCode', i.item_code,
                                           'itemName', i.item_name,
                                           'uom', i.uom,
                                           'stockUom', i.stock_uom,
                                           'itemGroup', i.item_group,
                                           'deliveryDate', i.delivery_date,
                                           'remarks', i.remarks,
                                           'quantity', i.quantity,
                                           'itemStatus', i.item_status
                                   )
                           )
                    FROM tbl_work_order_item i
                    WHERE i.work_order_id = wo.id
                      AND i.deleted_at IS NULL),
                   JSON_ARRAY()
           )                       AS items
    FROM tbl_work_order wo
    WHERE wo.id = p_id
      AND wo.deleted_at IS NULL;
END

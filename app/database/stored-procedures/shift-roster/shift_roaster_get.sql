DROP PROCEDURE IF EXISTS get_roasters_by_shift;
CREATE PROCEDURE get_roasters_by_shift(
    IN shift_id_val BIGINT
)
BEGIN

    SELECT s.shift AS shift_name,
           s.start_date,
           s.start_time,
           s.end_date,
           s.end_time,
           JSON_ARRAYAGG(
                   JSON_OBJECT(
                           'operator_id', o.id,
                           'emp_no', o.emp_no,
                           'roasterId', sr.id,
                           'roasterDate', sr.date,
                           'operator_name', o.name,
                           'section', o.section,
                           'machine_id', m.id,
                           'machine_code', m.machine_code,
                           'machine_name', m.machine_name,
                           'capabilities', m.capabilities,
                           'status', sr.status
                   )
           )       AS operators
    FROM tbl_shift_roster sr
             LEFT JOIN tbl_operator o ON sr.operator_id = o.id
             LEFT JOIN tbl_shift s ON sr.shift_id = s.id
             LEFT JOIN tbl_machine m ON sr.machine_id = m.id
    WHERE sr.shift_id = shift_id_val
      AND sr.deleted_at IS NULL
    GROUP BY sr.date, sr.shift_id, s.shift, s.start_date, s.start_time, s.end_date, s.end_time
    ORDER BY sr.date;
END;
-- name: GetPaymentByID :one
SELECT * FROM payments
WHERE id = $1 LIMIT 1;

-- name: GetAllPayments :many
SELECT * FROM payments;

-- name: GetAllPaymentsFromUser :many
SELECT p.* FROM payments p
WHERE p.user_id = $1;

-- name: CreatePayment :one
INSERT INTO payments (
    user_id, name, amount, type, category, date, due_date, paid, paid_at, recurrent, frequency, receipt_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, false, $8, false, $9, $10
)
RETURNING *;

-- name: UpdatePayment :execrows
UPDATE payments
SET
	name = $2,
    amount = $3,
    type = $4,
    category = $5,
    date = $6,
    due_date = $7,
    paid = $8,
    paid_at = $9,
    recurrent = $10,
    frequency = $11,
    receipt_url = $12
WHERE id = $1;

-- name: DeletePayment :execrows
DELETE FROM payments
WHERE id = $1;

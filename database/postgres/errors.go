package postgres

import (
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
)

// getErrorCode извлекает строковый код ошибки из разных драйверов (pgx или pq)
func getErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return string(pqErr.Code)
	}

	return ""
}

// 23505 - Нарушение уникальности (дубликат PRIMARY KEY или UNIQUE constraint)
func IsUniqueViolation(err error) bool {
	return getErrorCode(err) == "23505"
}

// 23503 - Нарушение внешнего ключа (ссылка на несуществующую запись в другой таблице)
func IsForeignKeyViolation(err error) bool {
	return getErrorCode(err) == "23503"
}

// 23502 - Нарушение NOT NULL (попытка вставить NULL в обязательное поле)
func IsNotNullViolation(err error) bool {
	return getErrorCode(err) == "23502"
}

// 23P01 - Нарушение исключения (например, пересечение диапазонов дат)
func IsExclusionViolation(err error) bool {
	return getErrorCode(err) == "23P01"
}

// 23514 - Нарушение CHECK constraint (например, возраст < 0)
func IsCheckViolation(err error) bool {
	return getErrorCode(err) == "23514"
}

// 40P01 - Дедлок (взаимная блокировка транзакций)
func IsDeadlockDetected(err error) bool {
	return getErrorCode(err) == "40P01"
}

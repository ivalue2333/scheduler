package method

import "gorm.io/gen"

type TaskMethod interface {
	/*
		update @@table
		set status = @status
		where id = @id and status = @oldStatus
	*/
	UpdateStatus(id int64, oldStatus, status int32) (gen.RowsAffected, error)
}

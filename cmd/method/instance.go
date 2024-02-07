package method

import (
	"time"

	"gorm.io/gen"
)

type InstanceMethod interface {
	/*
		update @@table
		where id = @id heart_beat_at=@heartBeatAt
	*/
	DoHeartBeat(id int64, heartBeatAt time.Time) (gen.RowsAffected, error)
}

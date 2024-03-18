package method

import (
	"time"

	"gorm.io/gen"
)

type InstanceMethod interface {
	/*
		update @@table
		set heart_beat_at = @heartBeatAt
		where id = @id
	*/
	DoHeartBeat(id int64, heartBeatAt time.Time) (gen.RowsAffected, error)

	/*
		delete from @@table where id = @id
	*/
	PurgeInstance(id int64) (gen.RowsAffected, error)

	/*
		select * from @@table
		where id = id limit 1

	*/
	GetInstance(id int64) (*gen.T, error)

	/*
		select * from @@table
		where id > @id and heart_beat_at < @oldestHeartBeatAt
		order by id limit 20
	*/
	ListOfflineInstances(id int64, oldestHeartBeatAt time.Time) ([]*gen.T, error)
}

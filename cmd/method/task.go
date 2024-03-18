package method

import "gorm.io/gen"

type TaskMethod interface {
	/*
		update @@table
		set status = @status
		where id in @ids and status = @oldStatus and instance_id = @instanceID
	*/
	UpdateStatus(ids []int64, instanceID int64, oldStatus, status int32) (gen.RowsAffected, error)

	/*
		update @@table
		set status = 1, instance_id = 0
		where id in @ids
		and instance_id = @instanceID
		and status = 2
	*/
	PurgeInstanceTask(ids []int64, instanceID int64) (gen.RowsAffected, error)

	/*
		select * from @@table
		where namespace = @namespace
		and category = @category
		and status = 1
		and instance_id in (@instanceID, 0)
		and id > @id
		order by id limit 20
	*/
	ListWaittingTasks(namespace string, category string, instanceID int64, id int64) ([]*gen.T, error)

	/*
		select * from @@table
		where instance_id = @instanceID
		and status = 2
		and id > @id
		order by id limit 20
	*/
	// todo order by id 是否OK
	ListInstanceRunningTasks(id int64, instanceID int64) ([]*gen.T, error)
}

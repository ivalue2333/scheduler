package method

import "gorm.io/gen"

type LockerMethod interface {
	/*
		select * from @@table
		where key = @key limit 1;
	*/
	GetLockInstance(key string) (*gen.T, error)

	/*
		delete from @@table
		where key = @key limit 1;
	*/
	DeleteLockInstance(key string) (*gen.T, error)
}

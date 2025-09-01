package tsdb_cluster

import (
	"info/internal/domain/oracul_speedometers"
	"info/internal/infrastructure/repository/tsdb"
)

type OraculSpeedometersReplicaSet struct {
	*ReplicaSet
}

var _ oracul_speedometers.ReplicaSet = (*OraculSpeedometersReplicaSet)(nil)

func NewOraculSpeedometersReplicaSet(replicaSet *ReplicaSet) *OraculSpeedometersReplicaSet {
	return &OraculSpeedometersReplicaSet{
		ReplicaSet: replicaSet,
	}
}

func (c *OraculSpeedometersReplicaSet) WriteRepo() oracul_speedometers.WriteRepository {
	return tsdb.NewOraculSpeedometersRepository(c.ReplicaSet.WriteRepo())
}

func (c *OraculSpeedometersReplicaSet) ReadRepo() oracul_speedometers.ReadRepository {
	return tsdb.NewOraculSpeedometersRepository(c.ReplicaSet.ReadRepo())
}

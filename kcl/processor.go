package kcl

type CheckpointFunc = func(*string) error

type (
	InitializationInput struct {
		ShardID string
		// Local patch (not upstream): the lease's checkpoint, where KCL will
		// resume delivering: a sequence number, or a sentinel (TRIM_HORIZON,
		// LATEST, AT_TIMESTAMP, SHARD_END). Empty when the daemon did not
		// send one.
		SequenceNumber    string
		SubSequenceNumber int64
	}
	ProcessRecordsInput struct {
		Records    []Record
		Checkpoint CheckpointFunc `json:"-"`
	}
	LeaseLostInput  struct{}
	ShardEndedInput struct {
		Checkpoint CheckpointFunc `json:"-"`
	}
	ShutdownRequestedInput struct {
		Checkpoint CheckpointFunc `json:"-"`
	}
)

type RecordProcessor interface {
	Initialize(*InitializationInput)
	ProcessRecords(*ProcessRecordsInput)
	LeaseLost(*LeaseLostInput)
	ShardEnded(*ShardEndedInput)
	ShutdownRequested(*ShutdownRequestedInput)
}

package ps

type Process interface {
	Pid() uint
	PPid() uint
	Gid() uint
	Uid() uint
	Executable() string
	UserTime() uint
	SysTime() uint
	Threads() uint
	StartTime() uint64
	VirtMem() uint
	PhysMem() uint
}

func Processes() ([]Process, error) {
	return processes()
}

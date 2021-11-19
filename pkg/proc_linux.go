//go:build linux

package ps

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall"
)

type linuxProcess struct {
	pid        uint32
	uid        uint32
	gid        uint32
	comm       string
	state      uint8
	ppid       uint32
	utime      uint
	stime      uint
	numThreads uint
	startTime  uint64
	vsize      uint
	rss        uint
}

var pageSize = os.Getpagesize()

func (p *linuxProcess) Pid() uint {
	return uint(p.pid)
}

func (p *linuxProcess) Gid() uint {
	return uint(p.gid)
}

func (p *linuxProcess) Uid() uint {
	return uint(p.uid)
}

func (p *linuxProcess) PPid() uint {
	return uint(p.ppid)
}

func (p *linuxProcess) Executable() string {
	return p.comm
}

func (p *linuxProcess) UserTime() uint {
	return p.utime
}

func (p *linuxProcess) SysTime() uint {
	return p.stime
}

func (p *linuxProcess) NumThreads() uint {
	return p.numThreads
}

func (p *linuxProcess) StartTime() uint64 {
	return p.startTime
}

func (p *linuxProcess) VirtMem() uint {
	return p.vsize
}

func (p *linuxProcess) PhysMem() uint {
	return uint(pageSize) * p.rss
}

func processes() ([]Process, error) {
	dir, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	results := make([]Process, 0, 32)
	for {
		names, err := dir.Readdirnames(16)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, name := range names {
			if name[0] < '0' || name[0] > '9' {
				continue
			}

			pid, err := strconv.ParseUint(name, 10, 32)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				continue
			}

			var stat syscall.Stat_t
			err = syscall.Stat("/proc/"+name, &stat)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				continue
			}

			p, err := newLinuxProcess(uint32(pid), stat.Uid, stat.Gid)
			if err != nil {
				continue
			}

			results = append(results, p)
		}
	}

	return results, nil
}

func pidStat(pid uint32) ([]string, error) {
	statPath := fmt.Sprintf("/proc/%d/stat", pid)
	stat, err := os.Open(statPath)
	if err != nil {
		return nil, err
	}
	defer stat.Close()

	fields := make([]string, 0, 50)
	scanner := bufio.NewScanner(stat)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		field := scanner.Text()
		fields = append(fields, field)
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return fields, nil
}

func newLinuxProcess(pid, uid, gid uint32) (*linuxProcess, error) {
	fields, err := pidStat(pid)
	if err != nil {
		return nil, err
	}

	ppid, err := strconv.ParseUint(fields[3], 10, 32)
	if err != nil {
		return nil, err
	}

	utime, err := strconv.ParseUint(fields[13], 10, 0)
	if err != nil {
		return nil, err
	}

	stime, err := strconv.ParseUint(fields[14], 10, 0)
	if err != nil {
		return nil, err
	}

	numThreads, err := strconv.ParseUint(fields[19], 10, 0)
	if err != nil {
		return nil, err
	}

	startTime, err := strconv.ParseUint(fields[21], 10, 64)
	if err != nil {
		return nil, err
	}

	vsize, err := strconv.ParseUint(fields[22], 10, 0)
	if err != nil {
		return nil, err
	}

	rss, err := strconv.ParseUint(fields[23], 10, 0)
	if err != nil {
		return nil, err
	}

	return &linuxProcess{
		pid:        pid,
		uid:        uid,
		gid:        gid,
		comm:       fields[1][1 : len(fields[1])-1],
		state:      fields[2][0],
		ppid:       uint32(ppid),
		utime:      uint(utime),
		stime:      uint(stime),
		numThreads: uint(numThreads),
		startTime:  startTime,
		vsize:      uint(vsize),
		rss:        uint(rss),
	}, nil
}

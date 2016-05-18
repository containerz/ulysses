// Copyright (C) 2014 Space Monkey, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package monitor

import "C"

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"syscall"
)

func registerPlatformEnvironment(group *MonitorGroup) {
	group.Chain("rusage", MonitorFunc(
		func(cb func(name string, val float64)) {
			var rusage syscall.Rusage
			err := syscall.Getrusage(syscall.RUSAGE_SELF, &rusage)
			if err != nil {
				logger.Errorf("failed getting rusage data: %s", err)
				return
			}
			MonitorStruct(&rusage, cb)
		}))
	group.Chain("proc.stat", MonitorFunc(
		func(cb func(name string, val float64)) {
			var stat procSelfStat
			err := readProcSelfStat(&stat)
			if err != nil {
				logger.Errorf("failed getting /proc/self/stat data: %s", err)
				return
			}
			MonitorStruct(&stat, cb)
		}))
	group.Chain("proc.statm", MonitorFunc(
		func(cb func(name string, val float64)) {
			var stat procSelfStatm
			err := readProcSelfStatm(&stat)
			if err != nil {
				logger.Errorf("failed getting /proc/self/statm data: %s", err)
				return
			}
			MonitorStruct(&stat, cb)
		}))
}

func fdCount() (count int, err error) {
	f, err := os.Open("/proc/self/fd")
	if err != nil {
		return 0, err
	}
	defer f.Close()
	for {
		names, err := f.Readdirnames(4096)
		count += len(names)
		if err != nil {
			if err == io.EOF {
				return count, nil
			}
			return count, err
		}
	}
}

func openProc() (*os.File, error) {
	return os.Open("/proc/self/exe")
}

type procSelfStat struct {
	Pid                 C.int
	Comm                string
	State               byte
	Ppid                C.int
	Pgrp                C.int
	Session             C.int
	TtyNr               C.int
	Tpgid               C.int
	Flags               C.uint
	Minflt              C.ulong
	Cminflt             C.ulong
	Majflt              C.ulong
	Cmajflt             C.ulong
	Utime               C.ulong
	Stime               C.ulong
	Cutime              C.long
	Cstime              C.long
	Priority            C.long
	Nice                C.long
	NumThreads          C.long
	Itrealvalue         C.long
	Starttime           C.ulonglong
	Vsize               C.ulong
	Rss                 C.long
	Rsslim              C.ulong
	Startcode           C.ulong
	Endcode             C.ulong
	Startstack          C.ulong
	Kstkesp             C.ulong
	Kstkeip             C.ulong
	Signal              C.ulong
	Blocked             C.ulong
	Sigignore           C.ulong
	Sigcatch            C.ulong
	Wchan               C.ulong
	Nswap               C.ulong
	Cnswap              C.ulong
	ExitSignal          C.int
	Processor           C.int
	RtPriority          C.uint
	Policy              C.uint
	DelayAcctBlkioTicks C.ulonglong
	GuestTime           C.ulong
	CguestTime          C.long
}

func readProcSelfStat(s *procSelfStat) error {
	data, err := ioutil.ReadFile("/proc/self/stat")
	if err != nil {
		return err
	}
	_, err = fmt.Sscanf(string(data), "%d %s %c %d %d %d %d %d %d %d %d "+
		"%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d "+
		"%d %d %d %d %d %d %d %d %d %d", &s.Pid, &s.Comm, &s.State, &s.Ppid,
		&s.Pgrp, &s.Session, &s.TtyNr, &s.Tpgid, &s.Flags, &s.Minflt, &s.Cminflt,
		&s.Majflt, &s.Cmajflt, &s.Utime, &s.Stime, &s.Cutime, &s.Cstime,
		&s.Priority, &s.Nice, &s.NumThreads, &s.Itrealvalue, &s.Starttime,
		&s.Vsize, &s.Rss, &s.Rsslim, &s.Startcode, &s.Endcode, &s.Startstack,
		&s.Kstkesp, &s.Kstkeip, &s.Signal, &s.Blocked, &s.Sigignore, &s.Sigcatch,
		&s.Wchan, &s.Nswap, &s.Cnswap, &s.ExitSignal, &s.Processor, &s.RtPriority,
		&s.Policy, &s.DelayAcctBlkioTicks, &s.GuestTime, &s.CguestTime)
	return err
}

type procSelfStatm struct {
	Size     int
	Resident int
	Share    int
	Text     int
	Lib      int
	Data     int
	Dt       int
}

func readProcSelfStatm(s *procSelfStatm) error {
	data, err := ioutil.ReadFile("/proc/self/statm")
	if err != nil {
		return err
	}
	_, err = fmt.Sscanf(string(data), "%d %d %d %d %d %d %d", &s.Size,
		&s.Resident, &s.Share, &s.Text, &s.Lib, &s.Data, &s.Dt)
	return err
}

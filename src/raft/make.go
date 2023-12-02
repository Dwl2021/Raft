package raft

import (
	"math/rand"
	"time"

	"6.824/labrpc"
)

// Make
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.

func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here (2A, 2B, 2C).

	// 对应论文中的初始化
	rf.applyChan = applyCh //2B

	rf.currentTerm = 0
	rf.votedFor = -1
	rf.logs = make([]LogEntry, 0)

	rf.commitIndex = 0
	rf.lastApplied = 0

	rf.nextIndex = make([]int, len(peers))
	rf.matchIndex = make([]int, len(peers))

	rf.status = Follower
	rf.timeout = time.Duration(150+rand.Intn(200)) * time.Millisecond // 随机产生150-350ms
	rf.timer = time.NewTicker(rf.timeout)

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	//fmt.Printf("[ 	Make-func-rf(%v) 	]:  %v\n", rf.me, rf.timeout)
	// start ticker goroutine to start elections
	DPrintf("Now create the raft node [ %d ]!", rf.me)

	go rf.ticker()

	return rf
}

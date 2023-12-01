package raft

import (
	"math/rand"
	"time"
)

// raft的ticker(心脏)
// The ticker go routine starts a new election if this peer hasn't received
// heartsbeats recently.
func (rf *Raft) ticker() {

	for rf.killed() == false {

		// Your code here to check if a leader election should
		// be started and to randomize sleeping time using

		// 当定时器结束进行超时选举
		select {

		case <-rf.timer.C:
			if rf.killed() {
				return
			}
			rf.mu.Lock()
			// 根据自身的status进行一次ticker
			switch rf.status {

			// follower变成竞选者
			case Follower:
				rf.status = Candidate
				fallthrough
			case Candidate:

				// 初始化自身的任期、并把票投给自己
				rf.currentTerm += 1
				rf.votedFor = rf.me
				votedNums := 1 // 统计自身的票数

				// 每轮选举开始时，重新设置选举超时
				rf.overtime = time.Duration(150+rand.Intn(200)) * time.Millisecond // 随机产生200-400ms
				rf.timer.Reset(rf.overtime)

				// 对自身以外的节点进行选举
				for i := 0; i < len(rf.peers); i++ {
					if i == rf.me {
						continue
					}

					voteArgs := RequestVoteArgs{
						Term:         rf.currentTerm,
						CandidateId:  rf.me,
						LastLogIndex: len(rf.logs) - 1,
						LastLogTerm:  0,
					}
					if len(rf.logs) > 0 {
						voteArgs.LastLogTerm = rf.logs[len(rf.logs)-1].Term
					}

					voteReply := RequestVoteReply{}

					go rf.sendRequestVote(i, &voteArgs, &voteReply, &votedNums)
				}
			case Leader:
				// 进行心跳/日志同步
				appendNums := 1 // 对于正确返回的节点数量
				rf.timer.Reset(HeartBeatTimeout)
				// 构造msg
				for i := 0; i < len(rf.peers); i++ {
					if i == rf.me {
						continue
					}
					appendEntriesArgs := AppendEntriesArgs{
						Term:         rf.currentTerm,
						LeaderId:     rf.me,
						PrevLogIndex: 0,
						PrevLogTerm:  0,
						Entries:      nil,
						LeaderCommit: rf.commitIndex,
					}

					appendEntriesReply := AppendEntriesReply{}
					//fmt.Printf("[	ticker(%v) ] : send a election to %v\n", rf.me, i)
					go rf.sendAppendEntries(i, &appendEntriesArgs, &appendEntriesReply, &appendNums)
				}
			}

			rf.mu.Unlock()
		}

	}
}

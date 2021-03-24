// Copyright 2016 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"code.gitea.io/gitea/modules/timeutil"
)

// Vote represents a vote for an issue made by an user.
type Vote struct {
	ID          int64              `xorm:"pk autoincr"`
	UID         int64              `xorm:"UNIQUE(s)"`
	IssueID     int64              `xorm:"UNIQUE(s)"`
	CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
}

// Vote or unvote issue.
func VoteIssue(userID, issueID int64, vote bool) error {
	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	if star {
		if isVoting(sess, userID, issueID) {
			return nil
		}

		if _, err := sess.Insert(&Vote{UID: userID, IssueID: issueID}); err != nil {
			return err
		}
		if _, err := sess.Exec("UPDATE `issue` SET num_votes = num_votes + 1 WHERE id = ?", issueID); err != nil {
			return err
		}
	} else {
		if !isVoting(sess, userID, issueID) {
			return nil
		}

		if _, err := sess.Delete(&Vote{UID: userID, IssueID: issueID}); err != nil {
			return err
		}
		if _, err := sess.Exec("UPDATE `issue` SET num_votes = num_votes - 1 WHERE id = ?", issueID); err != nil {
			return err
		}
	}

	return sess.Commit()
}

// isVoting checks if user has voted given issue.
func isVoting(userID, issueID int64) bool {
	return isVoting(x, userID, issueID)
}

func isVoting(e Engine, userID, issueID int64) bool {
	has, _ := e.Get(&Vote{UID: userID, IssueID: issueID})
	return has
}

// GetVoters returns the users that voted for the issue.
func (issue *Issue) GetVoters(opts ListOptions) ([]*User, error) {
	sess := x.Where("vote.issue_id = ?", issue.ID).
		Join("LEFT", "vote", "`user`.id = vote.uid")
	if opts.Page > 0 {
		sess = opts.setSessionPagination(sess)

		users := make([]*User, 0, opts.PageSize)
		return users, sess.Find(&users)
	}

	users := make([]*User, 0, 8)
	return users, sess.Find(&users)
}

// GetVotersRepos returns the issues the user voted for.
func (u *User) GetVotedIssues(private bool, page, pageSize int, orderBy string) (issues IssueList, err error) {
	if len(orderBy) == 0 {
		orderBy = "updated_unix DESC"
	}
	sess := x.
		Join("INNER", "vote", "vote.issue_id = issue.id").
		Where("vote.uid = ?", u.ID).
		OrderBy(orderBy)

	if page <= 0 {
		page = 1
	}
	sess.Limit(pageSize, (page-1)*pageSize)

	issues = make([]*Issue, 0, pageSize)

	if err = sess.Find(&issues); err != nil {
		return
	}

	if err = issues.loadAttributes(x); err != nil {
		return
	}

	return
}

// GetVotedIssueCount returns the numbers of issues the user voted for.
func (u *User) GetVotedIssueCount(private bool) (int64, error) {
	sess := x.
		Join("INNER", "vote", "vote.issue_id = issue.id").
		Where("vote.uid = ?", u.ID)

	return sess.Count(&Issue{})
}

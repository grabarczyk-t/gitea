// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package migrations

import (
	"code.gitea.io/gitea/modules/timeutil"

	"xorm.io/xorm"
)

func addVoteTable(x *xorm.Engine) error {

	type Vote struct {
		ID          int64              `xorm:"pk autoincr"`
		UID         int64              `xorm:"UNIQUE(s)"`
		IssueID     int64              `xorm:"UNIQUE(s)"`
		CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
	}

	return x.Sync2(new(Vote))
}

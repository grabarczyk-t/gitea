// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package migrations

import (
	"code.gitea.io/gitea/modules/timeutil"

	"xorm.io/xorm"
)

func addIssueVotes(x *xorm.Engine) error {
	type Issue struct {
		NumVotes int `xorm:"NOT NULL DEFAULT 0"`
	}

	if err := x.Sync2(new(Issue)); err != nil {
		return err
	}

	type Vote struct {
		ID          int64              `xorm:"pk autoincr"`
		UID         int64              `xorm:"UNIQUE(s)"`
		IssueID     int64              `xorm:"UNIQUE(s)"`
		CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
	}

	return x.Sync2(new(Vote))
}

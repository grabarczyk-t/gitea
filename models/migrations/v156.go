// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package migrations

import (
	"fmt"

	"xorm.io/xorm"
)

func addNumVotesColumnToIssueTable(x *xorm.Engine) error {
	type Issue struct {
		NumVotes int `xorm:"NOT NULL DEFAULT 0"`
	}

	if err := x.Sync2(new(Webhook)); err != nil {
		return fmt.Errorf("Sync2: %v", err)
	}
	return nil
}

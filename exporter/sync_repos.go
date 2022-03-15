package exporter

import (
	"app/base"
	"app/base/database"
	"app/base/redisdb"
	"app/base/utils"
)

type packageRepoRow struct {
	PackageID string
	Repo      string
}

func syncRepos() error {
	rows, err := database.Db.Table("pkg_repo pr").
		Select("pr.pkg_id AS package_id, 'r:' || c.label AS repo").
		Joins("JOIN repo r ON pr.repo_id = r.id").
		Joins("JOIN content_set c ON c.id = r.content_set_id").
		Rows()
	if err != nil {
		utils.Log("err", err.Error()).Error("Unable to load package repo data from database")
		return err
	}

	var row packageRepoRow
	nSynced, nFailed := 0, 0
	for rows.Next() {
		err = database.Db.ScanRows(rows, &row)
		if err != nil {
			utils.Log("err", err.Error()).Error("Unable to scan single package repo row")
			nFailed += 1
			continue
		}

		err = redisdb.Rdb.SAdd(base.Context, row.Repo, row.PackageID).Err()
		if err != nil {
			utils.Log("err", err.Error(), "packageID", row.PackageID, "repo", row.Repo).
				Error("Unable to store package repo item to Redis")
			nFailed += 1
			continue
		}
		nSynced += 1
	}
	utils.Log("nSynced", nSynced, "nFailed", nFailed).Info("Repos synced")
	return nil
}

package exporter

import (
	"app/base"
	"app/base/database"
	"app/base/redisdb"
	"app/base/utils"
)

type packageRow struct {
	ID      string
	Name    string
	Epoch   int
	Version string
	Release string
	Arch    string
	Errata  string
}

func syncPackages() error {
	rows, err := database.Db.Table("package p").
		Select("p.id AS id, pn.name AS name, epoch, version, release, arch.name AS arch," +
			"COALESCE(e.name, '-') AS errata").
		Joins("JOIN package_name pn ON p.name_id = pn.id").
		Joins("JOIN evr ON p.evr_id = evr.id ").
		Joins("JOIN arch ON p.arch_id = arch.id").
		Joins("LEFT JOIN pkg_errata pe ON p.id = pe.pkg_id").
		Joins("LEFT JOIN errata e ON pe.errata_id = e.id").
		Rows()
	if err != nil {
		utils.Log("err", err.Error()).Error("Unable to load package data from database")
		return err
	}

	var row packageRow
	nSynced, nFailed := 0, 0
	for rows.Next() {
		err = database.Db.ScanRows(rows, &row)
		if err != nil {
			utils.Log("err", err.Error()).Error("Unable to scan single package row")
			nFailed += 1
			continue
		}

		err = addToRedis(row)
		if err != nil {
			nFailed += 1
			continue
		}
		nSynced += 1
	}
	utils.Log("nSynced", nSynced, "nFailed", nFailed).Info("Packages synced")
	return nil
}

func addToRedis(row packageRow) error {
	nevraObj := utils.Nevra{Name: row.Name, Epoch: row.Epoch, Version: row.Version, Release: row.Release, Arch: row.Arch}
	nevra := nevraObj.StringE(false)
	err := redisdb.Rdb.Set(base.Context, row.ID, nevra+" "+row.Errata, 0).Err()
	if err != nil {
		return err
	}

	nameKey := "u:" + row.Name
	err = redisdb.Rdb.SAdd(base.Context, nameKey, row.ID).Err()
	if err != nil {
		return err
	}

	archKey := "a:" + row.Arch
	err = redisdb.Rdb.SAdd(base.Context, archKey, row.ID).Err()
	if err != nil {
		return err
	}
	return nil
}

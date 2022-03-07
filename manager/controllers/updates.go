package controllers

import (
	"app/base"
	"app/base/redisdb"
	"app/base/utils"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/ezamriy/gorpm"
	"github.com/gin-gonic/gin"
)

type UpdatesRequest struct {
	PackageList    []string `json:"package_list"`
	RepositoryList []string `json:"repository_list"`
}

type UpdatesResponse struct {
	UpdateList map[string][][]string `json:"update_list"`
}

// @Summary Show updates
// @Description Show updates
// @ID listAdvisories
// @Security RhIdentity
// @Accept   json
// @Produce  json
// @Success 200 {object} UpdatesResponse
// @Router /api/patch/v1/updates [post]
func UpdatesHandler(c *gin.Context) {
	var request UpdatesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		LogAndRespBadRequest(c, err, "Invalid request body: "+err.Error())
		return
	}

	allNevrasUpdates, evras, counts, err := getNevraErratumUpdates(request.PackageList, request.RepositoryList)
	if err != nil {
		LogAndRespBadRequest(c, err, "Unable find updates")
		return
	}

	updateList := filterAllNevrasUpdates(request.PackageList, allNevrasUpdates, evras, counts)
	var resp = UpdatesResponse{
		UpdateList: updateList,
	}
	c.JSON(http.StatusOK, &resp)
}

func getReposKey(repositoryList []string) ([]string, error) {
	if len(repositoryList) == 0 {
		return []string{}, nil
	}
	repos := make([]string, len(repositoryList))
	for i, repo := range repositoryList {
		repos[i] = fmt.Sprintf("r:%s", repo)
	}
	sort.Strings(repos)
	reposKey := strings.Join(repos, "+")
	status := redisdb.Rdb.Exists(base.Context, reposKey)
	val, err := status.Result()
	if err != nil {
		return nil, err
	}
	if val == 0 {
		redisdb.Rdb.SUnionStore(base.Context, reposKey, repos...)
		redisdb.Rdb.Expire(base.Context, reposKey, 5)
	}
	return []string{reposKey}, nil
}

func getNevraUpdates(nevraStr string, reposKey []string) ([]string, *utils.Nevra) {
	packageUpdatesIDs := []string{"-1"}
	nevra, err := utils.ParseNevra(nevraStr)
	if err != nil {
		utils.Log("err", err.Error(), "nevra", nevra).Debug("Unable to parse nevra")
		return packageUpdatesIDs, nevra
	}
	sinterKeys := []string{"u:" + nevra.Name, "a:" + nevra.Arch}
	sinterKeys = append(sinterKeys, reposKey...)
	packageIDs, err := redisdb.Rdb.SInter(base.Context, sinterKeys...).Result()
	if err != nil {
		utils.Log("err", err.Error(), "nevra").Debug("Unable to get sinter result")
		return packageUpdatesIDs, nevra
	}
	numPackageIDs, err := strarr2intarr(packageIDs)
	if err != nil {
		utils.Log("err", err.Error(), "nevra").Debug("Unable to parse numeric package IDs")
		return packageUpdatesIDs, nevra
	}
	sort.Ints(numPackageIDs)
	sortedPackageIDs := intarr2strarr(numPackageIDs)
	return sortedPackageIDs, nevra
}

func strarr2intarr(sarr []string) ([]int, error) {
	iarr := make([]int, len(sarr))
	for i, s := range sarr {
		num, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		iarr[i] = num
	}
	return iarr, nil
}

func intarr2strarr(iarr []int) []string {
	sarr := make([]string, len(iarr))
	for i, num := range iarr {
		s := strconv.Itoa(num)
		sarr[i] = s
	}
	return sarr
}

func getNevraErratumUpdates(packageList []string, repositoryList []string) (
	[]interface{}, []string, []int, error) {
	updatesIDs := []string{}
	evras := []string{}
	updatesCounts := []int{}

	reposKey, err := getReposKey(repositoryList)
	if err != nil {
		utils.Log("err", err.Error()).Debug("Unable to get repos")
		return nil, nil, nil, nil
	}

	for _, nevraStr := range packageList {
		packageUpdatesIDs, nevra := getNevraUpdates(nevraStr, reposKey)
		updatesIDs = append(updatesIDs, packageUpdatesIDs...)
		updatesCounts = append(updatesCounts, len(packageUpdatesIDs))
		if nevra != nil {
			evras = append(evras, nevra.EVRAString())
		} else {
			evras = append(evras, "")
		}
	}
	if len(updatesIDs) == 0 {
		return nil, evras, updatesCounts, nil
	}

	nevraErratumUpdates, err := redisdb.Rdb.MGet(base.Context, updatesIDs...).Result()
	if err != nil {
		return nil, nil, nil, err
	}
	return nevraErratumUpdates, evras, updatesCounts, nil
}

func filterNevraUpdates(installedEvra string, nevraErratumUpdates []interface{}) [][]string {
	outUpdates := [][]string{}
	for _, nevraErratumUpdate := range nevraErratumUpdates {
		if nevraErratumUpdate == nil {
			continue
		}

		nevraErratumUpdateStr := nevraErratumUpdate.(string)
		nevraErratumArr := strings.Split(nevraErratumUpdateStr, " ")
		if len(nevraErratumArr) != 2 {
			utils.Log("nevraErratum", nevraErratumUpdateStr).Debug("Update parsing failed")
			continue
		}
		nevraStr := nevraErratumArr[0]
		nevra, err := utils.ParseNevra(nevraStr)
		if err != nil {
			utils.Log("err", err.Error(), "nevraToParse", nevraStr).Debug("Unable to parse update nevra")
			continue
		}

		if rpm.Vercmp(nevra.EVRAString(), installedEvra) <= 0 {
			continue
		}

		erratum := nevraErratumArr[1]
		outUpdates = append(outUpdates, []string{nevraStr, erratum})
	}
	return outUpdates
}

func filterAllNevrasUpdates(packageList []string, allNevraErratumUpdates []interface{}, installedEvras []string,
	updatesCounts []int) map[string][][]string {
	updatesList := map[string][][]string{}
	iUpdatesStart := 0
	for i, nevraStr := range packageList {
		updatesCount := updatesCounts[i]
		if updatesCount == 0 {
			updatesList[nevraStr] = [][]string{}
			continue
		}
		packageUpdates := allNevraErratumUpdates[iUpdatesStart : iUpdatesStart+updatesCount]
		iUpdatesStart += updatesCount
		updateResult := filterNevraUpdates(installedEvras[i], packageUpdates)
		updatesList[nevraStr] = updateResult
	}
	return updatesList
}

package controller

import (
	"comment-chuli/format"
	"comment-chuli/lib"
	"comment-chuli/model"
	"comment-chuli/setup"
	"comment-chuli/tool"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LiveStartController struct {
	BaseController
}

func (live *LiveStartController) Start(c *gin.Context) {

	sign := c.Query("sign")
	if sign != "1mw9ug4SHisq5PYC" {
		return
	}

	startTime := time.Now()
	limit := 100
	returnData := make(map[string]interface{})

	var async model.LiveCommentAsyncV2
	rows := async.GetList(limit)
	if len(rows) == 0 {
		c.AsciiJSON(http.StatusOK, format.ResponseData{
			Status: http.StatusOK,
			Data:   returnData,
			Msg:    "没有任务队列",
		})
		return
	}

	// 构造id数组
	var cidSlice []int
	var cidStrSlice []string
	var successCidSlice []int
	successNum := 0
	for _, row := range rows {
		cidSlice = append(cidSlice, row.Cid)
		cidStrSlice = append(cidStrSlice, strconv.Itoa(row.Cid))
	}

	// 初始化日志
	cidLen := len(cidSlice)
	setup.Logger.Printf("处理条数=%d,", cidLen)
	var log = fmt.Sprintf("任务数量%d，处理情况\n", cidLen)

	// 更新处理状态
	async.UpdateHangByFilename(cidSlice)

	// 创建通道用于存储结果，使用goroutine等待所有请求完成
	commentResultChan := make(chan *format.CommentResult, limit)
	var wg sync.WaitGroup
	for _, item := range rows {
		wg.Add(1)
		go live.Node(&wg, commentResultChan, item.Cid)
	}
	wg.Wait()
	close(commentResultChan)

	// 按批更新评论数
	filenameNum := make(map[string]map[string]int)
	for resItem := range commentResultChan {

		if resItem.Cid == 0 {
			continue
		}

		// 接口返回成功记录数+1
		successNum++

		if resItem.Status != "normal" {
			continue
		}

		// 请求结果记录
		successCidSlice = append(successCidSlice, resItem.Cid)
		log += fmt.Sprintf("success:%d:%d\n", resItem.Cid, resItem.Duration)

		if _, ok := filenameNum[resItem.Filename]; !ok {
			filenameNum[resItem.Filename] = make(map[string]int)
		}
		filenameNum[resItem.Filename]["num"]++
		if resItem.ParentId == 0 {
			filenameNum[resItem.Filename]["root_normal"]++
		}
	}

	if len(filenameNum) > 0 {
		var numWg sync.WaitGroup
		for k, v := range filenameNum {
			numWg.Add(1)
			go live.updateNum(&numWg, k, v["num"], v["root_normal"])
		}
		numWg.Wait()
	}

	//fmt.Println(cidSlice, filenameNum, log)

	// 评论通过的批量写源
	live.writeQueue(successCidSlice)

	// 删除队列数据
	async.DeleteHangByFilename(cidSlice)

	// 写入日志
	var asyncLog model.LiveCommentAsyncLog
	statusValue := "error"
	if len(cidSlice) == successNum {
		statusValue = "success"
	}
	cidStr := strings.Join(cidStrSlice, ",")
	asyncLog.SaveData("node_v2", "", time.Since(startTime).Milliseconds(), statusValue, "")

	// 返回结果
	returnData["cid_str"] = cidStr
	c.AsciiJSON(http.StatusOK, format.ResponseData{
		Status: http.StatusOK,
		Data:   returnData,
		Msg:    "成功",
	})

}

func (live *LiveStartController) Node(wg *sync.WaitGroup, commentResult chan *format.CommentResult, cid int) {
	defer wg.Done()

	// 防止协程崩溃
	defer func() {
		if err := recover(); err != nil {
			var tmpInsertLog model.TmpInsertLog
			tmpInsertLog.SaveData("go-start", fmt.Sprintf("并发请求异常cid=%d", cid), 0, "", 0)
			setup.Logger.Printf("err-并发请求异常: %v\n", err)
		}
	}()
	startTime := time.Now()

	defaultReturn := &format.CommentResult{}

	url := fmt.Sprintf("%s/live/node/%d", setup.Config.Server.NodeUrl, cid)
	ret, reqCode := lib.CurlReq("GET", url, time.Millisecond*2000, nil, nil, "评论处理")
	if ret == nil || reqCode == http.StatusGatewayTimeout {
		var tmpInsertLog model.TmpInsertLog
		tmpInsertLog.SaveData("go-start-check-timeout", fmt.Sprintf("cid=%d", cid), 0, "", 0)
		setup.Logger.Printf("err,cid=%d请求超时", cid)
		commentResult <- defaultReturn
		return
	}
	retMap, err := ret.Map()
	if err != nil {
		setup.Logger.Printf("err,cid=%d返回格式错误", cid)
		commentResult <- defaultReturn
		return
	}

	status, err := retMap["status"].(json.Number).Int64()
	if err != nil {
		setup.Logger.Printf("err,cid=%d返回内部状态码错误", cid)
		commentResult <- defaultReturn
		return
	}

	msg, msgOk := retMap["msg"].(string)
	if !msgOk {
		setup.Logger.Printf("err,cid=%d返回错误信息错误", cid)
		commentResult <- defaultReturn
		return
	}

	if status != http.StatusOK {
		setup.Logger.Printf("err,cid=%d接口请求失败,状态码%d,%s", cid, reqCode, msg)
		commentResult <- defaultReturn
		return
	} else {
		data, ok := retMap["data"].(map[string]interface{})
		if !ok {
			setup.Logger.Printf("err,cid=%d返回详细数据错误", cid)
			commentResult <- defaultReturn
			return
		}

		room, _ := data["room"].(json.Number).Int64()
		ParentId, _ := data["parent_id"].(json.Number).Int64()
		commentResult <- &format.CommentResult{
			Cid:      cid,
			Room:     room,
			ParentId: ParentId,
			Status:   data["status"].(string),
			Filename: data["filename"].(string),
			Duration: time.Since(startTime).Milliseconds(),
		}
	}

}

func (live *LiveStartController) updateNum(wg *sync.WaitGroup, filename string, num int, rootNormal int) {

	defer wg.Done()

	// 防止协程崩溃
	defer func() {
		if err := recover(); err != nil {
			var tmpInsertLog model.TmpInsertLog
			tmpInsertLog.SaveData("go-start-updateNum", fmt.Sprintf("并发请求异常filename=%s", filename), 0, "", 0)
			setup.Logger.Printf("err,updateNum并发请求异常: %v\n", err)
		}
	}()

	var commentNum model.LiveCommentNum
	commentNum.OpCommentNum(filename, num, rootNormal)
}

func (live *LiveStartController) writeQueue(cidSlice []int) {

	if len(cidSlice) == 0 {
		return
	}
	// 处理所在的楼层
	task := make(map[string][]int)

	var liveFloorComment model.LiveFloorComment
	var liveCacheQueue model.LiveCacheQueue

	floorRows := liveFloorComment.GetByCidSlice(cidSlice)

	// 合并同页的只生成一次
	var mapLiveQueue []map[string]interface{}
	for _, li := range floorRows {
		page := int(math.Ceil(float64(li.FloorNum) / 100))
		if tool.InArray(page, task[li.Filename]) {
			continue
		}
		task[li.Filename] = append(task[li.Filename], page)

		if liveCacheQueue.GetByFPage(li.Filename, page) == nil {
			mapLiveQueue = append(mapLiveQueue, map[string]interface{}{
				"filename":   li.Filename,
				"page":       page,
				"createTime": time.Now(),
			})
		}

	}

	// 批量写入写源队列
	if len(mapLiveQueue) > 0 {
		liveCacheQueue.BatchSaveData(mapLiveQueue)
	}
}

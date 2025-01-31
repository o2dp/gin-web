package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"strings"
	"time"
)

var (
	dictNameCache       = cache.New(24*time.Hour, 48*time.Hour)
	dictNameAndKeyCache = cache.New(24*time.Hour, 48*time.Hour)
)

// 获取指定字典名称且字典数据key的字典数据(不返回err)
func (s *MysqlService) GetDictDataByDictNameAndDictDataKeyNoErr(dictName, dictDataKey string) models.SysDictData {
	dict, err := s.GetDictDataByDictNameAndDictDataKey(dictName, dictDataKey)
	if err != nil || dict == nil {
		return models.SysDictData{}
	}
	return *dict
}

// 获取指定字典名称且字典数据key的字典数据
func (s *MysqlService) GetDictDataByDictNameAndDictDataKey(dictName, dictDataKey string) (*models.SysDictData, error) {
	cacheKey := fmt.Sprintf("%s_%s", dictName, dictDataKey)
	oldCache, ok := dictNameAndKeyCache.Get(cacheKey)
	if ok {
		c, _ := oldCache.(models.SysDictData)
		return &c, nil
	}
	var err error
	list := make([]models.SysDictData, 0)
	err = s.tx.
		Model(&models.SysDictData{}).
		Preload("Dict").
		Order("created_at DESC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	for _, data := range list {
		if data.Dict.Name == dictName && data.Key == dictDataKey {
			// 写入缓存
			dictNameAndKeyCache.Set(cacheKey, data, cache.DefaultExpiration)
			return &data, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

// 获取指定名称的字典数据
func (s *MysqlService) GetDictDatasByDictName(name string) ([]models.SysDictData, error) {
	cacheKey := name
	oldCache, ok := dictNameCache.Get(cacheKey)
	if ok {
		c, _ := oldCache.([]models.SysDictData)
		return c, nil
	}
	var err error
	list := make([]models.SysDictData, 0)
	err = s.tx.
		Model(&models.SysDictData{}).
		Preload("Dict").
		Order("sort").
		Find(&list).Error
	if err != nil {
		return list, err
	}
	newList := make([]models.SysDictData, 0)
	for _, data := range list {
		if data.Dict.Name == name {
			newList = append(newList, data)
		}
	}
	// 写入缓存
	dictNameCache.Set(cacheKey, newList, cache.DefaultExpiration)
	return newList, nil
}

// 获取所有字典
func (s *MysqlService) GetDicts(req *request.DictRequestStruct) ([]models.SysDict, error) {
	var err error
	list := make([]models.SysDict, 0)
	query := s.tx.
		Model(&models.SysDict{}).
		Preload("DictDatas").
		Order("created_at DESC")
	name := strings.TrimSpace(req.Name)
	if name != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	desc := strings.TrimSpace(req.Desc)
	if desc != "" {
		query = query.Where("desc = ?", desc)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	// 查询列表
	err = s.Find(query, &req.PageInfo, &list)
	return list, err
}

// 获取所有字典数据
func (s *MysqlService) GetDictDatas(req *request.DictDataRequestStruct) ([]models.SysDictData, error) {
	var err error
	list := make([]models.SysDictData, 0)
	query := s.tx.
		Model(&models.SysDictData{}).
		Preload("Dict").
		Order("created_at DESC")
	key := strings.TrimSpace(req.Key)
	if key != "" {
		query = query.Where("key LIKE ?", fmt.Sprintf("%%%s%%", key))
	}
	val := strings.TrimSpace(req.Val)
	if val != "" {
		query = query.Where("val LIKE ?", fmt.Sprintf("%%%s%%", val))
	}
	attr := strings.TrimSpace(req.Attr)
	if attr != "" {
		query = query.Where("attr = ?", attr)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.DictId != nil {
		query = query.Where("dict_id = ?", *req.DictId)
	}
	// 查询列表
	err = s.Find(query, &req.PageInfo, &list)
	return list, err
}

// 创建字典
func (s *MysqlService) CreateDict(req *request.CreateDictRequestStruct) (err error) {
	err = s.Create(req, new(models.SysDict))
	dictNameCache.Flush()
	dictNameAndKeyCache.Flush()
	return
}

// 更新字典
func (s *MysqlService) UpdateDictById(id uint, req request.UpdateDictRequestStruct) (err error) {
	err = s.UpdateById(id, req, new(models.SysDict))
	dictNameCache.Flush()
	dictNameAndKeyCache.Flush()
	return
}

// 批量删除字典
func (s *MysqlService) DeleteDictByIds(ids []uint) (err error) {
	err = s.DeleteByIds(ids, new(models.SysDict))
	dictNameCache.Flush()
	dictNameAndKeyCache.Flush()
	return
}

// 创建字典数据
func (s *MysqlService) CreateDictData(req *request.CreateDictDataRequestStruct) (err error) {
	err = s.Create(req, new(models.SysDictData))
	dictNameCache.Flush()
	dictNameAndKeyCache.Flush()
	return
}

// 更新字典数据
func (s *MysqlService) UpdateDictDataById(id uint, req request.UpdateDictDataRequestStruct) (err error) {
	err = s.UpdateById(id, req, new(models.SysDictData))
	dictNameCache.Flush()
	dictNameAndKeyCache.Flush()
	return
}

// 批量删除字典数据
func (s *MysqlService) DeleteDictDataByIds(ids []uint) (err error) {
	err = s.DeleteByIds(ids, new(models.SysDictData))
	dictNameCache.Flush()
	dictNameAndKeyCache.Flush()
	return
}

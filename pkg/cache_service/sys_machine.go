package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"strings"
)

func (s *RedisService) GetMachines(req *request.MachineRequestStruct) ([]models.SysMachine, error) {
	if !global.Conf.System.UseRedis || !global.Conf.System.UseRedisService {
		// 不使用redis
		return s.mysql.GetMachines(req)
	}
	var err error
	list := make([]models.SysMachine, 0)
	query := s.redis.
		Table(new(models.SysMachine).TableName()).
		Order("created_at DESC")
	host := strings.TrimSpace(req.Host)
	if host != "" {
		query = query.Where("host", "contains", host)
	}
	loginName := strings.TrimSpace(req.LoginName)
	if loginName != "" {
		query = query.Where("login_name", "contains", loginName)
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator", "contains", creator)
	}
	if req.Status != nil {
		query = query.Where("status", "=", *req.Status)
	}
	// 查询列表
	err = s.Find(query, &req.PageInfo, &list)
	return list, err
}

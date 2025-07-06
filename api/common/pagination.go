package common

// ===============================
// 通用分页响应结构
// ===============================

// PaginationResponse 通用分页响应结构
type PaginationResponse[T any] struct {
	List  []T   `json:"list"`  // 数据列表
	Total int64 `json:"total"` // 总记录数
	Page  int   `json:"page"`  // 当前页码
	Size  int   `json:"size"`  // 每页大小
}

// PaginationRequest 通用分页请求结构
type PaginationRequest struct {
	Page      int    `json:"page" d:"1" v:"min:1#页码最小为1"`                // 页码
	Size      int    `json:"size" d:"10" v:"between:1,100#每页数量为1-100"`   // 每页大小
	SortBy    string `json:"sort_by,omitempty"`                          // 排序字段
	SortOrder string `json:"sort_order" d:"desc" v:"in:asc,desc#排序方向无效"` // 排序方向
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	CurrentPage  int   `json:"current_page"`  // 当前页码
	PageSize     int   `json:"page_size"`     // 每页大小
	TotalPages   int   `json:"total_pages"`   // 总页数
	TotalRecords int64 `json:"total_records"` // 总记录数
	HasNext      bool  `json:"has_next"`      // 是否有下一页
	HasPrev      bool  `json:"has_prev"`      // 是否有上一页
}

// NewPaginationResponse 创建分页响应
func NewPaginationResponse[T any](list []T, total int64, page, size int) PaginationResponse[T] {
	return PaginationResponse[T]{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	}
}

// GetPaginationInfo 获取分页信息
func GetPaginationInfo(total int64, page, size int) PaginationInfo {
	totalPages := int((total + int64(size) - 1) / int64(size))

	return PaginationInfo{
		CurrentPage:  page,
		PageSize:     size,
		TotalPages:   totalPages,
		TotalRecords: total,
		HasNext:      page < totalPages,
		HasPrev:      page > 1,
	}
}

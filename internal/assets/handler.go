package assets

import( "github.com/gin-gonic/gin"
	"fmt"
	"strconv"
	"strings"
)
type Handler struct {
	Service *Service
	Repo    *Repository
}

func (h *Handler) GetAssets(c *gin.Context) {
	ctx := c.Request.Context()

	status := c.Query("status")
	atype := c.Query("type")
	category := c.Query("category")
	search := c.Query("search")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query := `SELECT id, asset_code, name, serial_no, asset_type, category, status 
	          FROM asset_inventory WHERE is_deleted = FALSE`

	var args []interface{}
	var conditions []string
	i := 1

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status=$%d", i))
		args = append(args, status)
		i++
	}
	if atype != "" {
		conditions = append(conditions, fmt.Sprintf("asset_type=$%d", i))
		args = append(args, atype)
		i++
	}
	if category != "" {
		conditions = append(conditions, fmt.Sprintf("category=$%d", i))
		args = append(args, category)
		i++
	}
	if search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR serial_no ILIKE $%d)", i, i))
		args = append(args, "%"+search+"%")
		i++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, limit, offset)

	rows, err := h.Repo.DB.Query(ctx, query, args...)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var list []AssetListDTO

	for rows.Next() {
		var a AssetListDTO
		rows.Scan(&a.ID, &a.Code, &a.Name, &a.SerialNo, &a.Type, &a.Category, &a.Status)
		list = append(list, a)
	}

	// count query
	countQ := `SELECT COUNT(*) FROM asset_inventory WHERE is_deleted=FALSE`
	if len(conditions) > 0 {
		countQ += " AND " + strings.Join(conditions, " AND ")
	}

	var total int
	h.Repo.DB.QueryRow(ctx, countQ, args[:len(args)-2]...).Scan(&total)

	c.JSON(200, gin.H{
		"success": true,
		"data":    list,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}
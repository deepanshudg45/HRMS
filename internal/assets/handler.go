package assets

import( 
	"github.com/gofiber/fiber/v2"
	"fmt"
	"strconv"
	"strings"
	"context"
)
type Handler struct {
	Service *Service
	Repo    *Repository
}

func (h *Handler) GetAssets(c *fiber.Ctx) error {
	ctx := context.Background() 

	status := c.Query("status")
	atype := c.Query("type")
	category := c.Query("category")
	search := c.Query("search")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
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
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer rows.Close()

	var list []AssetListDTO
	fmt.Println("Handler DB:", h.Repo.DB)
	for rows.Next() {
		var a AssetListDTO
		if err := rows.Scan(&a.ID, &a.Code, &a.Name, &a.SerialNo, &a.Type, &a.Category, &a.Status); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		list = append(list, a)
	}

	// count query
	countQ := `SELECT COUNT(*) FROM asset_inventory WHERE is_deleted=FALSE`
	if len(conditions) > 0 {
		countQ += " AND " + strings.Join(conditions, " AND ")
	}

	var total int
	err = h.Repo.DB.QueryRow(ctx, countQ, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    list,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

package serve

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	guardian "github.com/hyperupcall/redpanda/server/guardian"
	"github.com/hyperupcall/redpanda/server/store"
)

func hasError(ctx *gin.Context, err error) bool {
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return true
	}

	return false
}

func returnSuccess(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func Serve(store *store.Store) {
	r := gin.Default()
	g := guardian.New(store)

	r.POST("/api/action/apply", func(ctx *gin.Context) {
		type Schema struct {
			Transaction string `json:"transaction" binding:"required"`
		}
		var data Schema
		if err := ctx.BindJSON(&data); err != nil {
			return
		}

		content, err := g.ActionApply(data.Transaction)
		if hasError(ctx, err) {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"contents": content})
	})

	r.POST("/api/action/refresh", func(ctx *gin.Context) {
		type Schema struct {
			Transaction string `json:"transaction" binding:"required"`
		}
		var data Schema
		if err := ctx.BindJSON(&data); err != nil {
			return
		}

		content, err := g.ActionRefresh(data.Transaction)
		if hasError(ctx, err) {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"contents": content})
	})

	r.POST("/api/action/commit", func(ctx *gin.Context) {

		type Schema struct {
			Transaction   string `json:"transaction" binding:"required"`
			CommitMessage string `json:commitMessage" binding:"required"`
		}
		var data Schema
		if err := ctx.BindJSON(&data); err != nil {
			return
		}

		err := g.ActionCommit(data.Transaction)
		if hasError(ctx, err) {
			return
		}

		returnSuccess(ctx)
	})

	r.POST("/api/action/push", func(ctx *gin.Context) {
		type Schema struct {
			Transaction   string `json:"transaction" binding:"required"`
			CommitMessage string `json:commitMessage" binding:"required"`
		}
		var data Schema
		if err := ctx.BindJSON(&data); err != nil {
			return
		}

		err := g.ActionPush(data.Transaction)
		if hasError(ctx, err) {
			return
		}

		returnSuccess(ctx)
	})

	r.POST("/api/transformer/add", func(c *gin.Context) {
		type Schema struct {
			Transaction string `json:"transaction" binding:"required"`
			Type        string `json:"type" binding:"required"`
			Transformer string `json:"transformer" binding:"required"`
			Content     string `json:"content" binding:"required"`
		}
		var data Schema

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := store.TransformerAdd(data.Transaction, data.Type, data.Transformer, data.Content); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
		return
	})

	r.POST("/api/transformer/remove", func(c *gin.Context) {
		type Schema struct {
			Transaction string `json:"transaction" binding:"required"`
			Transformer string `json:"transformer" binding:"required"`
		}
		var data Schema

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := store.TransformerRemove(data.Transaction, data.Transformer); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
		return
	})

	r.POST("/api/transformer/edit", func(c *gin.Context) {
		type Schema struct {
			Transaction string `json:"transaction" binding:"required"`
			Transformer string `json:"transformer" binding:"required"`
			NewContent  string `json:"newContent" binding:"required"`
		}
		var data Schema

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := store.TransformerEdit(data.Transaction, data.Transformer, data.NewContent); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
		return
	})

	r.POST("/api/transformer/order", func(c *gin.Context) {
		type Schema struct {
			Transaction string `json:"transaction" binding:"required"`
			Order       string `json:"order" binding:"required"`
		}
		var data Schema

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := store.TransformerOrder(data.Transaction, data.Order); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
		return
	})

	r.POST("/api/repo/add", func(c *gin.Context) {
		type Schema struct {
			Transaction string `json:"transaction" binding:"required"`
			Repo        string `json:"repo" binding:"required"`
		}
		var data Schema

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := store.RepoAdd(data.Transaction, data.Repo); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
		return
	})

	r.POST("/api/repo/remove", func(c *gin.Context) {
		type Schema struct {
			Transaction string `json:"transaction" binding:"required"`
			Repo        string `json: "repo" binding: "required"`
		}
		var data Schema

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := store.RepoRemove(data.Transaction, data.Repo); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
		return
	})

	r.POST("/api/transaction/get", func(c *gin.Context) {
		type Schema struct {
			Name string `json: "name" binding:"required"`
		}
		var data Schema
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := store.TransactionGet(data.Name)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": result})
	})

	r.POST("/api/transaction/add", func(c *gin.Context) {
		type Schema struct {
			Name string `json: "name" binding:"required"`
		}
		var data Schema
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := store.TransactionAdd(data.Name); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	})

	r.POST("/api/transaction/remove", func(c *gin.Context) {
		type Schema struct {
			Name string `json: "name" binding:"required"`
		}
		var data Schema
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := store.TransactionRemove(data.Name); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	})

	r.POST("/api/transaction/rename", func(c *gin.Context) {
		type Schema struct {
			OldName string `json: "oldName" binding:"required"`
			NewName string `json: "newName" binding:"required"`
		}
		var data Schema
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := store.TransactionRename(data.OldName, data.NewName); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	})

	r.POST("/api/transaction/list", func(c *gin.Context) {
		type Schema struct{}
		var data Schema
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		transactions := store.TransactionList()
		c.JSON(http.StatusOK, gin.H{"transactions": transactions})
		return
	})

	if err := r.Run(":3000"); err != nil {
		log.Fatalln(err)
	}
}

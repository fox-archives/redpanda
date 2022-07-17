package serve

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyperupcall/redpanda/server/manager"
	"github.com/hyperupcall/redpanda/server/store"
)

func Serve(store *store.Store) {
	r := gin.Default()
	m := manager.New(store)

	r.POST("/api/step/initialize", func(c *gin.Context) {
		type Schema struct {
		}
		var data Schema
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := m.Initialize()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	r.POST("/api/step/idempotent-apply", func(c *gin.Context) {
		type Schema struct {
			Transaction string `json:"transaction" binding:"required"`
		}
		var data Schema
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := m.IdempotentApply(data.Transaction)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	r.POST("/api/step/diff", func(c *gin.Context) {
		type Schema struct {
			Transaction string `json:"transaction" binding:"required"`
		}
		var data Schema
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := m.Diff(data.Transaction)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"transactions": transactions})
		return
	})

	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}

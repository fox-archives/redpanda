package serve

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyperupcall/redpanda/server/store"
)

func Serve(store store.Store) {
	r := gin.Default()

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

	// r.POST("/api/user/list-repos", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"repos": []string{"a", "b"},
	// 	})
	// })

	// r.POST("/api/repos/list", func(c *gin.Context) {
	// 	repos := store.RepoList()

	// 	c.JSON(http.StatusOK, gin.H{
	// 		"repos": repos,
	// 	})
	// })

	// r.POST("/api/repos/add", func(c *gin.Context) {
	// 	type Schema struct {
	// 		Name string `json:"name" binding:"required"`
	// 	}
	// 	var data Schema

	// 	if err := c.ShouldBindJSON(&data); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	names := strings.Split(data.Name, ",")
	// 	if err := store.RepoAdd(names); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	c.Status(http.StatusOK)
	// 	return
	// })

	// r.POST("/api/repos/remove", func(c *gin.Context) {
	// 	type Schema struct {
	// 		Name string `json: "name" binding: "required"`
	// 	}
	// 	var data Schema

	// 	if err := c.ShouldBindJSON(&data); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	names := strings.Split(data.Name, ",")
	// 	if err := store.RepoRemove(names); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	c.Status(http.StatusOK)
	// 	return
	// })

	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}

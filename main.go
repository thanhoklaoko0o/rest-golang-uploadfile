package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	init()
	r := gin.Default()
	r.Static("static_file", "./assets")
	r.Use(aCustomMiddleware())

	r.GET("/ping", getPing)
	r.POST("/ping", postPing)
	r.GET("/detail/:id", getDetail)

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			//Khởi tạo middleware chỉ dùng cho group có prefix là v1
			v1.Use(groupV1Middleware())
			v1.GET("/ping", func(context *gin.Context) {
				context.JSON(http.StatusOK, gin.H{
					"message": "Ping",
				})
			})
			v1.GET("/pong", func(context *gin.Context) {
				context.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})
		}
		v2 := api.Group("/v2")
		{
			v2.GET("/a", func(context *gin.Context) {
				context.JSON(http.StatusOK, gin.H{
					"message": "a",
				})
			})
			v2.GET("/b", func(context *gin.Context) {
				context.JSON(http.StatusOK, gin.H{
					"message": "b",
				})
			})
		}
	}
	// rest to db
	r.POST("/todo", createTodo)
	r.GET("/todo", fetchAllTodo)
	r.GET("/todo/:id", fetchSingleTodo)
	r.PUT("/todo/:id", updateTodo)
	r.DELETE("/todo:id", deleteTodo)
	// upload file

	//upload single file
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.POST("/upload", func(context *gin.Context) {
		// single file
		file, _ := context.FormFile("file")
		log.Println(file.Filename)

		// Upload the file to specific dst.
		context.SaveUploadedFile(file, "./assets/upload/"+file.Filename)
		// context.SaveUploadedFile(file, "./assets/upload/"+uuid.New().String()+filepath.Ext(file.Filename))

		context.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	//upload multiple file
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.POST("/upload_multiple_file", func(context *gin.Context) {
		// Multipart form
		form, _ := context.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			// Upload the file to specific dst.
			context.SaveUploadedFile(file, "./assets/upload/"+uuid.New().String()+filepath.Ext(file.Filename))
		}
		context.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	r.Run(":333")
}

func groupV1Middleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		log.Println("I'm in a group v1 middleware")
		context.Next()
	}
}

func aCustomMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		log.Println("I'm in a global middleware")
		// jwt
		// auth
		if true {
			context.Next()
		}
	}
}

// func aCustomMiddleware(context *gin.Context) {
// 	// jwt
// 	// auth
// 	log.Println("This is global middleware")
// 	context.Next()
// }

func getDetail(context *gin.Context) {
	id := context.Param("id")
	context.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func postPing(context *gin.Context) {
	address := context.DefaultPostForm("addr", "VietNam")
	context.JSON(http.StatusOK, gin.H{
		"message": "Hello from " + address + " POST ping",
	})
}

func getPing(context *gin.Context) {
	name := context.DefaultQuery("name", "guest")
	var data = map[string]interface{}{
		"message": "Hello " + name + "from GET ping",
	}
	// context.String(http.StatusOK, "Ping")
	//context.JSON(http.StatusOK, gin.H{"message": "ok..........."})
	context.JSON(http.StatusOK, data)
}

var db *gorm.DBfunc init() {
	//open a db connection
	var err error
	db, err = gorm.Open("mysql", "root:sa123456@/golangdb?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
	 panic("failed to connect database")
	} 
	//Migrate the schema
	db.AutoMigrate(&todoModel{})
   }
   type (
	// todoModel describes a todoModel type
	todoModel struct {
	 gorm.Model
	 Title     string `json:"title"`
	 Completed int    `json:"completed"`
	}
   // transformedTodo represents a formatted todo
	transformedTodo struct {
	 ID        uint   `json:"id"`
	 Title     string `json:"title"`
	 Completed bool   `json:"completed"`
	}
   )

   // createTodo add a new todo
func createTodo(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := todoModel{Title: c.PostForm("title"), Completed: completed}
	db.Save(&todo)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": todo.ID})
   }

   // fetchAllTodo fetch all todos
func fetchAllTodo(c *gin.Context) {
	var todos []todoModel
	var _todos []transformedTododb.Find(&todos)if len(todos) <= 0 {
	 c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
	 return
	}
   //transforms the todos for building a good response
	for _, item := range todos {
	 completed := false
	 if item.Completed == 1 {
	  completed = true
	 } else {
	  completed = false
	 }
	 _todos = append(_todos, transformedTodo{ID: item.ID, Title: item.Title, Completed: completed})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todos})
   }
   // fetchSingleTodo fetch a single todo
   func fetchSingleTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")db.First(&todo, todoID)if todo.ID == 0 {
	 c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
	 return
	}
	completed := false
	if todo.Completed == 1 {
	 completed = true
	} else {
	 completed = false
	}
	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
   }
   // updateTodo update a todo
   func updateTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")db.First(&todo, todoID)if todo.ID == 0 {
	 c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
	 return
	}
	db.Model(&todo).Update("title", c.PostForm("title"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Update("completed", completed)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!"})
   }
   // deleteTodo remove a todo
   func deleteTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")db.First(&todo, todoID)if todo.ID == 0 {
	 c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
	 return
	}
	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!"})
   }

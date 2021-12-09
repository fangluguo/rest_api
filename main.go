package main

import (
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Maintainer struct {
	Name  string `json:"name" yaml:"name" binding:"required"`
	Email string `json:"email" yaml:"email" binding:"required,isValidEmail"`
}

type AppMeta struct {
	Title       string       `json:"title" yaml:"title" binding:"required"`
	Version     string       `json:"version" yaml:"version" binding:"required"`
	Maintainers []Maintainer `json:"maintainers" yaml:"maintainers" binding:"required"`
	Company     string       `json:"company" yaml:"company" binding:"required"`
	Website     string       `json:"website" yaml:"website" binding:"required"`
	Source      string       `json:"source" yaml:"source" binding:"required"`
	License     string       `json:"license" yaml:"license" binding:"required"`
	Description string       `json:"description" yaml:"description" binding:"required"`
}

type MaintainerQuery struct {
	Name  string `json:"name" yaml:"name"`
	Email string `json:"email" yaml:"email"`
}

type AppMetaQuery struct {
	Title       string          `json:"title" yaml:"title"`
	Version     string          `json:"version" yaml:"version"`
	Maintainer  MaintainerQuery `json:"maintainer" yaml:"maintainer"`
	Company     string          `json:"company" yaml:"company"`
	Website     string          `json:"website" yaml:"website"`
	Source      string          `json:"source" yaml:"source"`
	License     string          `json:"license" yaml:"license"`
	Description string          `json:"description" yaml:"description"`
}

var appMetas []AppMeta

var isValidEmail validator.Func = func(fieldLevel validator.FieldLevel) bool {
	email := fieldLevel.Field().String()
	log.Println("email: ", email)
	_, err := mail.ParseAddress(email)
	return err == nil
}

func main() {
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("isValidEmail", isValidEmail)
	} else {
		log.Println("Cannot register email validator.")
	}

	router.POST("/appmetas", createAppMeta)
	router.POST("/appmetas:query", queryAppMetas)

	router.Run("localhost:8080")
}

// TODO: Impement string length limit validation.
// TODO: Decide if to report error for non-existent field to avoid spelling error.
// TODO: Investigate general yaml/json query.
// TODO: More fancy query: expression for different and, or, not, in operations.
// TODO: Inveistgate why the validator doesn't work on the second level fields.
// 		 It works if the validator is on the first level fields.
// TODO: queryAppMetas

func createAppMeta(c *gin.Context) {
	var appMeta AppMeta
	if err := c.BindYAML(&appMeta); err != nil {
		c.YAML(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	appMetas = append(appMetas, appMeta)
	c.YAML(http.StatusCreated, appMeta)
}

// Query options:
// * Returns all if no fields are provided.
// * An AppMeta will be in the result if all provided query values are substrings
//   of the given AppMeta's corresponding fields.
// * For the maintainer array, if one maintainer is matched, it is counted as a match for the query's
//   Maintainer field.
// * Only a subset of fields are implemented.
func queryAppMetas(c *gin.Context) {
	var query AppMetaQuery
	// TODO: BindYAML cannot handle empty or non-existent query data.
	// Current solution is to provide an empty yaml: {}.
	if err := c.BindYAML(&query); err != nil {
		c.YAML(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	log.Println(query)
	var matchedAppMetas []AppMeta
	for _, a := range appMetas {
		isMatched := true
		if query.Title != "" && !strings.Contains(a.Title, query.Title) {
			isMatched = false
		}
		if query.Version != "" && !strings.Contains(a.Version, query.Version) {
			isMatched = false
		}
		if query.Maintainer != (MaintainerQuery{}) {
			isFound := false
			for _, maintainer := range a.Maintainers {
				if query.Maintainer.Email != "" && strings.Contains(maintainer.Email, query.Maintainer.Email) {
					isFound = true
					break
				}
			}
			if !isFound {
				isMatched = false
			}
		}
		if isMatched {
			matchedAppMetas = append(matchedAppMetas, a)
		}
	}
	c.YAML(http.StatusOK, matchedAppMetas)
}

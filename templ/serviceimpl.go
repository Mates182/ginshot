package templates

import (
	"fmt"

	"github.com/mates182/ginshot/formatter"
	"github.com/mates182/ginshot/models"
)

func GetServiceImplTemplate(config *models.ProjectConfig, dbName string, serviceName string, model string, id string, crudType int) string {
	var logic string
	proyectName := formatter.ToPascalCase(config.ProjectName)
	switch crudType {
	case 1:
		logic = GetCreateLogicTemplate(config, model, id)
	case 2:
		logic = GetReadLogicTemplate(config, model, id)
	case 3:
		logic = GetUpdateLogicTemplate(config, model, id)
	case 4:
		logic = GetDeleteLogicTemplate(config, model, id)
	case 5:
		logic = GetListLogicTemplate(config, model, id)
	default:
		logic = fmt.Sprintf("response := responses.%sResponse{}", proyectName)
	}

	return fmt.Sprintf(`package service

import (
	requests "%s/data/requests"
	responses "%s/data/responses"
	"net/http"
	`+func() string {
		if dbName != "" {
			if dbName == "mongo" {
				return `"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	` + func() string {
					if crudType == 2 || crudType == 4 {
						return ""
					}
					return `"` + config.ProjectName + "/models" + `"`
				}() + `
	"context"`
			} else if dbName == "redis" {
				return `"github.com/go-redis/redis/v8"`
			}
		}
		return ""
	}()+`
)
type %sServiceImpl struct {
	// Add Components
	`+func() string {
		if dbName != "" {
			return "DBClient *" + dbName + ".Client"
		}
		return ""
	}()+`
}

func New%sServiceImpl(`+func() string {
		if dbName != "" {
			return "dbClient *" + dbName + ".Client"
		}
		return ""
	}()+`) %sService {
	return &%sServiceImpl{
		// Add Components
		`+func() string {
		if dbName != "" {
			return "DBClient: dbClient,"
		}
		return ""
	}()+`
	}
}

func (service *%sServiceImpl) %sHandler(request requests.%sRequest) (int, responses.%sResponse) {
	%s
	return http.StatusOK, response
}
`, config.ProjectName, config.ProjectName, serviceName, serviceName, serviceName, serviceName, serviceName, serviceName, serviceName, serviceName, logic)
}

func GetCreateLogicTemplate(config *models.ProjectConfig, model string, id string) string {
	proyectName := formatter.ToPascalCase(config.ProjectName)
	return fmt.Sprintf(`mongoCollection := service.DBClient.Database("%s").Collection("%s")
	if request.%s.%s != "" {
		existing := mongoCollection.FindOne(context.Background(), bson.M{"%s": request.%s.%s})
		if existing.Err() == nil {
			return http.StatusConflict, responses.%sResponse{Message: "%s with the same %s already exists"}
		}
	}

	result, err := mongoCollection.InsertOne(context.Background(), request.%s)
	if err != nil {
		return http.StatusInternalServerError, responses.%sResponse{Message: "Error inserting %s"}
	}

	fmt.Println(result.InsertedID)

	response := responses.%sResponse{Message: "%s created successfully", %s: request.%s}`, config.Database.Name, config.Database.Table, model, id, id, model, id, proyectName, model, id, model, proyectName, model, proyectName, model, model, model)
}

func GetDeleteLogicTemplate(config *models.ProjectConfig, model string, id string) string {
	proyectName := formatter.ToPascalCase(config.ProjectName)
	return fmt.Sprintf(`mongoCollection := service.DBClient.Database("%s").Collection("%s")

	filter := bson.M{"%s": request.%s.%s}
	result, err := mongoCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return http.StatusInternalServerError, responses.%sResponse{Message: "Error deleting %s"}
	}
	if result.DeletedCount == 0 {
		return http.StatusNotFound, responses.%sResponse{Message: "%s not found"}
	}

	response := responses.%sResponse{Message: "%s deleted successfully", %s: request.%s}`, config.Database.Name, config.Database.Table, id, model, id, proyectName, model, proyectName, model, proyectName, model, model, model)
}
func GetReadLogicTemplate(config *models.ProjectConfig, model string, id string) string {
	proyectName := formatter.ToPascalCase(config.ProjectName)
	return fmt.Sprintf(`mongoCollection := service.DBClient.Database("%s").Collection("%s")

	var %s models.%s
	err := mongoCollection.FindOne(context.Background(), bson.M{"%s": request.%s.%s}).Decode(&%s)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, responses.%sResponse{Message: "%s not found"}
	}
	if err != nil {
		return http.StatusInternalServerError, responses.%sResponse{Message: "Error fetching %s"}
	}

	response := responses.%sResponse{Message: "%s retrieved successfully", %s: %s}`, config.Database.Name, config.Database.Table, model, model, id, model, id, model, proyectName, model, proyectName, model, proyectName, model, model, model)
}
func GetListLogicTemplate(config *models.ProjectConfig, model string, id string) string {
	proyectName := formatter.ToPascalCase(config.ProjectName)
	return fmt.Sprintf(`mongoCollection := service.DBClient.Database("%s").Collection("%s")

	cursor, err := mongoCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return http.StatusInternalServerError, responses.%sResponse{Message: "Error fetching %s"}
	}
	defer cursor.Close(context.Background())

	var %s []models.%s
	if err := cursor.All(context.Background(), &%s); err != nil {
		return http.StatusInternalServerError, responses.%sResponse{Message: "Error decoding %s"}
	}

	response := responses.%sResponse{Message: "All %s retrieved successfully", %s: %s}`, config.Database.Name, config.Database.Table,
		proyectName, model,
		model, model,
		model,
		proyectName, model,
		proyectName, model, model, model)
}
func GetUpdateLogicTemplate(config *models.ProjectConfig, model string, id string) string {
	proyectName := formatter.ToPascalCase(config.ProjectName)
	return fmt.Sprintf(`mongoCollection := service.DBClient.Database("%s").Collection("%s")
	filter := bson.M{"%s": request.%s.%s}
	update := bson.M{"$set": request.%s}
	result, err := mongoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return http.StatusInternalServerError, responses.%sResponse{Message: "Error updating %s"}
	}
	if result.MatchedCount == 0 {
		return http.StatusNotFound, responses.%sResponse{Message: "%s not found"}
	}return http.StatusInternalServerError, responses.%sResponse{Message: "Error fetching %s"}
	}

	response := responses.%sResponse{Message: "%s updated successfully", %s: %s}`, config.Database.Name, config.Database.Table,
		id, model, id,
		model,
		proyectName, model,
		proyectName, model,
		proyectName, model,
		proyectName, model, model, model)

}

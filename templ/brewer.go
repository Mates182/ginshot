package templates

import (
	"fmt"

	"github.com/mates182/ginshot/formatter"
)

func GetBrewerTemplate(Model, id string) string {
	model := formatter.ToLowerCase(Model)
	return fmt.Sprintf(`{
  "root": {
    "domain": "%s",
    "master-docker-compose": true,
    "gitignore": true,
    "database": {
      "type": "mongo",
      "name": "%s",
      "collection": "%s",
      "model": "%s",
      "id": "%s"
    }
  },
  "general": {
    "port": 80,
    "dockerfile": true,
    "docker_compose": true,
    "models": {
      "%s": {
        "ID": "string",
        "Message": "string",
        "Greeting": "Greeting"
      },
      "Greeting": {
        "Hello": "string"
      }
    },
    "messages": {
      "%sTopicMessage": {
        "Service": "string",
        "Message": "string",
        "%s": "models.%s"
      }
    }
  },
  "services": {
    "create_%s": {
      "type": "create",
      "requests": {
        "Create%sRequest": {
          "%s": "models.%s"
        }
      },
      "responses": {
        "Create%sResponse": {
          "%s": "models.%s",
          "Message": "string"
        }
      },
      "cors": {
        "methods": ["GET", "POST"]
      }
    },
    "get_%s_by_id": {
      "type": "get",
      "requests": {
        "Get%sByIdRequest": {
          "ID": "string"
        }
      },
      "responses": {
        "Get%sByIdResponse": {
          "%s": "models.%s",
          "Message": "string"
        }
      },
      "cors": {
        "methods": ["GET"]
      }
    },
    "list_%s": {
      "type": "list",
      "responses": {
        "List%sResponse": {
          "Message": "string",
          "%s": "[]models.%s"
        }
      },
      "cors": {
        "methods": ["GET"]
      }
    },
    "update_%s": {
      "type": "update",
      "requests": {
        "Update%sRequest": {
          "%s": "models.%s"
        }
      },
      "responses": {
        "Update%sResponse": {
          "%s": "models.%s",
          "Message": "string"
        }
      },
      "cors": {
        "methods": ["GET", "PUT"]
      }
    },
    "delete_%s": {
      "type": "delete",
      "requests": {
        "Delete%sRequest": {
          "ID": "string"
        }
      },
      "responses": {
        "Delete%sResponse": {
          "ID": "string",
          "Message": "string"
        }
      },
      "cors": {
        "methods": ["GET", "DELETE"]
      }
    }
  }
}
`, model, model, model, Model, id, Model, Model, Model, Model,
		model, Model, Model, Model, Model, Model, Model,
		model, Model, Model, Model, Model,
		model, Model, Model, Model,
		model, Model, Model, Model, Model, Model, Model,
		model, Model, Model)
}

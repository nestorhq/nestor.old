package resources

import (
	"errors"

	"github.com/nestorhq/nestor/internal/config"
)

// reference:
// https://github.com/pcarion/audiencefm_000/blob/7b9f21351f45df6f3a7b59e96d48371a8530e8fb/nestor/src/aws/policy.js

// PolicyStatement define an aws policy statement
type PolicyStatement struct {
	Effect   string
	Action   []string
	Resource string
}

// GetPolicyStatementsForLambda returns an array of policy statements
func (res *Resources) GetPolicyStatementsForLambda(permissions []config.LambdaPermission) ([]PolicyStatement, error) {
	var statements []PolicyStatement
	for _, permission := range permissions {
		resourceID := permission.ResourceID
		resource := res.findresourceByID(resourceID)
		if resource == nil {
			return nil, errors.New("Unknown resourceId:" + resourceID)
		}
		switch resource.resourceType {
		case s3Bucket:
			var actions []string
			for _, action := range permission.Actions {
				switch action.Operation {
				case "read":
					actions = append(actions, "s3:GetObject")
				case "write":
					actions = append(actions, "s3:PutObject")
				case "delete":
					actions = append(actions, "s3:DeleteObject")
				default:
					return nil, errors.New("Invalid operation on bucket:" + action.Operation)
				}
			}
			statements = append(statements, PolicyStatement{
				Effect:   "Allow",
				Resource: resource.awsID,
				Action:   actions,
			})

		case dynamoDbTable:
			var actions []string
			for _, action := range permission.Actions {
				switch action.Operation {
				case "read":
					actions = append(actions, "dynamodb:GetItem")
				case "query":
					actions = append(actions, "dynamodb:Query")
				case "write":
					actions = append(actions, "dynamodb:PutItem")
					actions = append(actions, "dynamodb:UpdateItem")
				case "delete":
					actions = append(actions, "dynamodb:DeleteItem")
				default:
					return nil, errors.New("Invalid operation on dynamo table:" + action.Operation)
				}
			}
			statements = append(statements, PolicyStatement{
				Effect:   "Allow",
				Resource: resource.awsID,
				Action:   actions,
			})

		default:
			return nil, errors.New("no policy can be set on:" + resourceID)
		}
	}
	return statements, nil
}

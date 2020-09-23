package items

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/southworks/gnalog/demo/repository/config"
	"github.com/southworks/gnalog/demo/repository/kafka"

	//Blank import
	_ "github.com/lib/pq"
)

//Server structure, holds every method
type Server struct{}

type item struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type itemKey struct {
	ID string `json:"id"`
}

type itemUpdate struct {
	Description string `json:":d"`
}

//TestGRPC Always responds with OK
func (s *Server) TestGRPC(ctx context.Context, request *TestGRPCRequest) (*TestGRPCResponse, error) {
	return &TestGRPCResponse{Ok: true}, nil
}

func dynamodbConnect() (*dynamodb.DynamoDB, error) {
	c, err := config.Get()
	config := &aws.Config{
		Region:   aws.String(c.Dynamodb.Region),
		Endpoint: aws.String(c.Dynamodb.Endpoint),
	}
	sess := session.Must(session.NewSession(config))
	return dynamodb.New(sess), err
}

func postgreSQLConnect() (*sql.DB, error) {
	c, err := config.Get()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", c.Postgresql.Host, c.Postgresql.Port, c.Postgresql.User, c.Postgresql.Password, c.Postgresql.Dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	return db, err
}

//CreateItem evaluates relational, connects and inserts to correspondant database
func (s *Server) CreateItem(ctx context.Context, request *CreateItemRequest) (*Result, error) {
	id := request.GetId()
	description := request.GetDescription()
	relational := request.GetRelational()
	c, err := config.Get()
	if err != nil {
		return nil, err
	}
	p := kafka.Publisher{}
	p.Connect("1", c.Kafka.Brokers, "logs")
	if relational {
		db, err := postgreSQLConnect()
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		defer db.Close()
		sqlStatement := `
			INSERT INTO item (id, description)
			VALUES ($1, $2)
		`
		_, err = db.Exec(sqlStatement, id, description)
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
	} else {
		svc, err := dynamodbConnect()
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		av, err := dynamodbattribute.MarshalMap(item{ID: id, Description: description})
		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("Items"),
		}
		_, err = svc.PutItem(input)
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
	}
	message := "CREATED"
	go p.Publish(context.Background(), nil, []byte(message))
	return &Result{Error: false, Message: message, Id: id}, nil
}

//ReadItem evaluates relational, connects and reads from correspondant database
func (s *Server) ReadItem(ctx context.Context, request *ReadItemRequest) (*Result, error) {
	id := request.GetId()
	i := item{}
	relational := request.GetRelational()
	c, err := config.Get()
	if err != nil {
		return nil, err
	}
	p := kafka.Publisher{}
	p.Connect("1", c.Kafka.Brokers, "logs")
	if relational {
		db, err := postgreSQLConnect()
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		defer db.Close()
		sqlStatement := `
			SELECT id, description
			FROM item 
			WHERE id = $1`
		row := db.QueryRow(sqlStatement, id)
		err = row.Scan(&i.ID, &i.Description)
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
	} else {
		svc, err := dynamodbConnect()
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		key, err := dynamodbattribute.MarshalMap(itemKey{ID: id})
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		input := &dynamodb.GetItemInput{
			Key:       key,
			TableName: aws.String("Items"),
		}
		result, err := svc.GetItem(input)
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		err = dynamodbattribute.UnmarshalMap(result.Item, &i)
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
	}
	var message string
	if i.ID != "" {
		message = "FOUND - ID: " + i.ID + " - Description: " + i.Description
		go p.Publish(context.Background(), nil, []byte(message))
		return &Result{Error: false, Message: "FOUND", Id: i.ID, Description: i.Description}, nil
	} else {
		message = "NOT FOUND"
	}
	go p.Publish(context.Background(), nil, []byte(message))
	return &Result{Error: true, Message: "NOT FOUND"}, nil
}

//UpdateItem evaluates relational, connects and updates to correspondant database
func (s *Server) UpdateItem(ctx context.Context, request *UpdateItemRequest) (*Result, error) {
	id := request.GetId()
	description := request.GetDescription()
	relational := request.GetRelational()
	c, err := config.Get()
	if err != nil {
		return nil, err
	}
	p := kafka.Publisher{}
	p.Connect("1", c.Kafka.Brokers, "logs")
	if relational {
		db, err := postgreSQLConnect()
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		defer db.Close()
		sqlStatement := `
			UPDATE item
			SET description = $2
			WHERE id = $1
		`
		_, err = db.Exec(sqlStatement, id, description)
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
	} else {
		svc, err := dynamodbConnect()
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		key, err := dynamodbattribute.MarshalMap(itemKey{ID: id})
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		update, err := dynamodbattribute.MarshalMap(itemUpdate{Description: description})
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		input := &dynamodb.UpdateItemInput{
			Key:                       key,
			TableName:                 aws.String("Items"),
			UpdateExpression:          aws.String("set description=:d"),
			ExpressionAttributeValues: update,
			ReturnValues:              aws.String("NONE"),
		}
		_, err = svc.UpdateItem(input)
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
	}
	message := "UPDATED"
	go p.Publish(context.Background(), nil, []byte(message))
	return &Result{Error: false, Message: message}, nil
}

//DeleteItem evaluates relational, connects and deletes to correspondant database
func (s *Server) DeleteItem(ctx context.Context, request *DeleteItemRequest) (*Result, error) {
	id := request.GetId()
	relational := request.GetRelational()
	c, err := config.Get()
	if err != nil {
		return nil, err
	}
	p := kafka.Publisher{}
	p.Connect("1", c.Kafka.Brokers, "logs")
	if relational {
		db, err := postgreSQLConnect()
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		defer db.Close()
		sqlStatement := `
			DELETE FROM item
			WHERE id = $1;		
		`
		_, err = db.Exec(sqlStatement, id)
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
	} else {
		svc, err := dynamodbConnect()
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		id := request.GetId()
		key, err := dynamodbattribute.MarshalMap(itemKey{ID: id})
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		input := &dynamodb.DeleteItemInput{
			Key:       key,
			TableName: aws.String("Items"),
		}
		_, err = svc.DeleteItem(input)
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
	}
	message := "DELETED"
	go p.Publish(context.Background(), nil, []byte(message))
	return &Result{Error: false, Message: message}, nil
}

//ListIds evaluates relational, connects and reads from correspondant database
func (s *Server) ListIds(ctx context.Context, request *ListIdsRequest) (*ListIdsResponse, error) {
	ids := []string{}
	relational := request.GetRelational()
	c, err := config.Get()
	if err != nil {
		return nil, err
	}
	p := kafka.Publisher{}
	p.Connect("1", c.Kafka.Brokers, "logs")
	if relational {
		db, err := postgreSQLConnect()
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		defer db.Close()
		rows, err := db.Query("SELECT id FROM item")
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var id string
			err = rows.Scan(&id)
			if err != nil {
				go p.Publish(context.Background(), nil, []byte(err.Error()))
				return nil, err
			}
			ids = append(ids, id)
		}
		err = rows.Err()
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
	} else {
		svc, err := dynamodbConnect()
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		proj := expression.NamesList(expression.Name("id"), expression.Name("description"))
		expr, err := expression.NewBuilder().WithProjection(proj).Build()
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		params := &dynamodb.ScanInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			ProjectionExpression:      expr.Projection(),
			TableName:                 aws.String("Items"),
		}
		result, err := svc.Scan(params)
		if err != nil {
			go p.Publish(context.Background(), nil, []byte(err.Error()))
			return nil, err
		}
		for _, element := range result.Items {
			i := item{}
			err = dynamodbattribute.UnmarshalMap(element, &i)
			if err != nil {
				go p.Publish(context.Background(), nil, []byte(err.Error()))
				return nil, err
			}
			ids = append(ids, i.ID)
		}
	}
	var message string
	if len(ids) > 0 {
		message = "OK"
	} else {
		message = "EMPTY"
	}
	go p.Publish(context.Background(), nil, []byte(message))
	return &ListIdsResponse{Error: false, Message: message, Ids: ids}, nil
}

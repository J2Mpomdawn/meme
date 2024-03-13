package service

import (
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"meme/model"
)

// get streaming configuration
func GetStreamConf() model.StreamConf {
	//if already acquired
	if StreamConf.Country != "" {
		return StreamConf
	}

	//create DynamoDB client
	svc := new_svc()

	//get session parameter
	res, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("meme_mori_memo"),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(os.Getenv("AWS_Memo")),
			},
		},
	})

	if err != nil {
		LogPrintln("red", "GetItem", err)
	}

	return model.StreamConf{
		Country: *res.Item["gvg_country"].N,
		World:   *res.Item["gvg_world"].N,
		Group:   *res.Item["gvg_group"].N,
		Class:   *res.Item["gvg_class"].N,
		Block:   *res.Item["gvg_block"].N,
		Castle:  *res.Item["gvg_castle"].N,
		Status:  *res.Item["status"].BOOL,
	}
}

// set streaming configuration
func SetStreamConf(country int, world int, group int, class int, block int, castle int, status bool) error {
	//create DynamoDB client
	svc := new_svc()

	//get session parameter
	_, err := svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("meme_mori_memo"),
		Item: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(os.Getenv("AWS_Memo")),
			},
			"gvg_country": {
				N: aws.String(strconv.Itoa(country)),
			},
			"gvg_world": {
				N: aws.String(strconv.Itoa(world)),
			},
			"gvg_group": {
				N: aws.String(strconv.Itoa(group)),
			},
			"gvg_class": {
				N: aws.String(strconv.Itoa(class)),
			},
			"gvg_block": {
				N: aws.String(strconv.Itoa(block)),
			},
			"gvg_castle": {
				N: aws.String(strconv.Itoa(castle)),
			},
			"status": {
				BOOL: aws.Bool(status),
			},
		},
	})

	StreamConf = model.StreamConf{
		Country: strconv.Itoa(country),
		World:   strconv.Itoa(world),
		Group:   strconv.Itoa(group),
		Class:   strconv.Itoa(class),
		Block:   strconv.Itoa(block),
		Castle:  strconv.Itoa(castle),
		Status:  status,
	}

	return err
}

// object for dynamoDb operation
func new_svc() *dynamodb.DynamoDB {
	//get session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_Region")),
	}))

	//create DynamoDB client
	return dynamodb.New(sess)
}

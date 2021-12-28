package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/labstack/echo/v4"
)

type Customer struct {
	Cid     string `json:"cid" bson:"cid"`
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
}
type Productdetails struct {
	Pid    string `json:"pid" bson:"pid"`
	Cid    string `json:"cid" bson:"cid"`
	Amount int    `json:"amount" bson:"amount"`
}

//var customers = make(map[string]Customer)
var collection *mongo.Collection
var collection1 *mongo.Collection
var ctx = context.TODO()

//var client mongo.Client

func main() {
	e := echo.New()

	c := Customer{}
	p := Productdetails{}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	collection = client.Database("User").Collection("Customer")
	collection1 = client.Database("User").Collection("Productdetails")

	//db := client.Database("User")

	e.GET("customers/:id", c.findById)
	e.POST("/customers", c.add)
	e.PUT("/customers/:id", c.update)
	e.DELETE("/customers/:id", c.delete)
	e.POST("/productdetails", p.padd)
	e.Logger.Fatal(e.Start(":8080"))
}

//Product Details adding
func (c Productdetails) padd(context echo.Context) error {
	p := &Productdetails{}
	if err := context.Bind(&p); err != nil {
		log.Fatal(err)
	}
	x := p.doProductAdd()

	return context.JSON(http.StatusAccepted, x)
}
func (c Productdetails) doProductAdd() string {

	res, err := collection1.InsertOne(context.Background(), c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.InsertedID.(primitive.ObjectID).Timestamp())
	return "Product Added"
}

//posting
func (c Customer) add(context echo.Context) error {
	u := &Customer{}
	if err := context.Bind(&u); err != nil {
		log.Fatal(err)
	}
	x := u.doAdd()

	return context.JSON(http.StatusAccepted, x)
}

func (c Customer) doAdd() string {

	res, err := collection.InsertOne(context.Background(), c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.InsertedID.(primitive.ObjectID).Timestamp())
	return "Customer Added"
}

//Delete

func (c Customer) delete(context echo.Context) error {
	id := context.Param("id")
	return context.JSON(http.StatusAccepted, c.doDelete(id))
}

func (c Customer) doDelete(id string) string {
	_, err := collection.DeleteOne(ctx, bson.D{{"cid", id}})
	if err != nil {
		log.Fatal(err)
	}
	return "Customer Deleted"
}

// getting
func (c Customer) findById(context echo.Context) error {
	id := context.Param("id")
	return context.JSON(http.StatusAccepted, c.getCustomerByid(id))

}

func (c Customer) getCustomerByid(id string) Customer {
	var customer Customer
	err := collection.FindOne(ctx, bson.D{{"cid", id}}).Decode(&customer)
	if err != nil {
		log.Fatal(err)
	}
	return customer
}
func (c Customer) update(context echo.Context) error {
	u := &Customer{}

	if err := context.Bind(&u); err != nil {
		log.Fatal(err)
	}
	return context.JSON(http.StatusAccepted, u.doUpdate())
}

func (c Customer) doUpdate() string {
	filter := bson.D{{"cid", c.Cid}}
	update := bson.D{{"$set", &c}}
	_, err := collection.UpdateOne(
		ctx,
		filter,
		update,
	)
	if err != nil {
		log.Fatal(err)
	}
	return "Updated"
}

/*
//put
func (c Customer) update(context echo.Context) error {
	u := &Customer{}

	if err := context.Bind(&u); err != nil {
		log.Fatal(err)
	}
	ID := context.Param("id")
	u.ID = ID
	filter := bson.D{{"id", c.ID}}
	update := bson.D{{"$set", *u}}
	_, err := collection.UpdateOne(
		ctx,
		filter,
		update,
	)
	if err != nil {
		log.Fatal(err)
	}
	return context.JSON(http.StatusOK, "Updataed")
}
*/
//	return context.JSON(http.StatusAccepted, u.doUpdate())

/*
func (c Customer) doUpdate() string {
filter:=bson.D{{"id",c.ID}}
update:=bson.D{{"$set",*c}}
_,err:=collection.UpdateOne(
	ctx,
	filter,
	update,
)
if err!=nil{
	log.Fatal(err)
}
return "Updated"
}

*/

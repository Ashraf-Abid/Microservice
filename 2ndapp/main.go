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

type Product struct {
	Pid     string `json:"pid" bson:"pid"`
	Pname   string `json:"pname" bson:"pname"`
	Pamount int    `json:"pamount" bson:"pamount"`
}

//var customers = make(map[string]Customer)
var collection *mongo.Collection
var ctx = context.TODO()

//var client mongo.Client

func main() {
	e := echo.New()

	p := Product{}

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
	collection = client.Database("Products").Collection("details")

	//db := client.Database("User")

	e.GET("products/:id", p.findById)
	e.POST("/products", p.add)
	e.PUT("/products/:id", p.update)
	e.DELETE("/products/:id", p.delete)
	//	e.PUT("/products/:id/:amount", p.padd)
	e.Logger.Fatal(e.Start(":8000"))

}

/*//Adding from Api
func (c Product) padd(context echo.Context) error {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"status": status,
			},
		},
		{
			"$group": bson.M{
				"_id":   "",
				"total": bson.M{"$sum": "$quantity"},
			},
		},
	}
	result := []bson.M{}
	err = db.C("product").Pipe(pipeline).All(&result)
	return result[0]["total"].(float64), nil
	p := &Product{}
	if err := context.Bind(&p); err != nil {
		log.Fatal(err)
	}
	id := context.Param("id")
	amount := context.Param("amount")
	fmt.Println(id, amount)
	amount = amount[1:]
	id = id[1:]

	a, _ := strconv.Atoi(amount)
	fmt.Println(id, amount, a)

	x := p.doProductAdd(id, a)

	return context.JSON(http.StatusAccepted, x)
}
func (c Product) doProductAdd(id string, need int) string {
	//filter := bson.D{{"pid", c.Pid}}
	//fmt.Println(filter)

	var product Product
	//available:=bson.D{{"amount",c.Pamount}}
	err := collection.FindOne(ctx, bson.D{{"pid", id}}).Decode(&product)
	fmt.Println(product)

	if err != nil {
		fmt.Println("filter")
		log.Fatal(err)
	}

	availabe := product.Pamount - need
	//fmt.Println("Available")
	fmt.Println(availabe)
	if availabe < 0 {
		return "Not available"
	} else {
		//fmt.Println(availabe)
		product.Pamount = availabe
		filter := bson.D{{"pid", product.Pid}}
		update := bson.D{{"$set", product}}
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

}
*/
//posting
func (c Product) add(context echo.Context) error {
	u := &Product{}
	if err := context.Bind(&u); err != nil {
		log.Fatal(err)
	}
	x := u.doAdd()

	return context.JSON(http.StatusAccepted, x)
}

func (c Product) doAdd() string {

	res, err := collection.InsertOne(context.Background(), c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.InsertedID.(primitive.ObjectID).Timestamp())
	return "Product Added"
}

//Delete

func (c Product) delete(context echo.Context) error {
	id := context.Param("id")
	return context.JSON(http.StatusAccepted, c.doDelete(id))
}

func (c Product) doDelete(id string) string {
	_, err := collection.DeleteOne(ctx, bson.D{{"pid", id}})
	if err != nil {
		log.Fatal(err)
	}
	return "Customer Deleted"
}

// getting
func (c Product) findById(context echo.Context) error {
	id := context.Param("id")
	fmt.Println(id)
	p := c.getCustomerByid(id)
	fmt.Println(p)
	//fmt.Println(p.)
	return context.JSON(http.StatusAccepted, p)

}

func (c Product) getCustomerByid(id string) Product {
	var product Product
	err := collection.FindOne(ctx, bson.D{{"pid", id}}).Decode(&product)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(product)
	return product
}
func (c Product) update(context echo.Context) error {
	u := &Product{}

	if err := context.Bind(&u); err != nil {
		log.Fatal(err)
	}
	return context.JSON(http.StatusAccepted, u.doUpdate())
}

func (c Product) doUpdate() string {
	filter := bson.D{{"cid", c.Pid}}
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

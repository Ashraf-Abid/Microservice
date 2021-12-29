package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
type Response struct {
	Pid     string `json:"pid" bson:"pid"`
	Pname   string `json:"pname" bson:"pname"`
	Pamount int    `json:"pamount" bson:"pamount"`
}

type newstruct struct {
	//_id   string `json:"" bson:""`
	Id    interface{} `json:"_id" bson:"_id"`
	Count int         `json:"count" bson:"count"`
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
	//q:=newstruct{}

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
	e.POST("/setProduct", p.saveproduct)

	e.Logger.Fatal(e.Start(":8080"))
}
func (p Productdetails) saveproduct(context echo.Context) error {
	u := Productdetails{}
	q := newstruct{}
	if err := context.Bind(&u); err != nil {
		log.Fatal(err)
	}

	id := u.Pid
	fmt.Println("AMOUNT", u.Amount)

	product, err := getProduct(id)
	fmt.Println(product.Pamount)
	//id:=product.Pid
	//fmt.Println(product, id)

	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"pid": product.Pid,
			},
		},
		{
			"$group": bson.M{
				"_id":   bson.M{"pid": "$pid"},
				"count": bson.M{"$sum": "$amount"},
			},
		},
	}
	//fmt.Println(pipeline)

	var result []bson.M

	cur, err := collection1.Aggregate(ctx, pipeline)
	// fmt.Println(cur)
	if err != nil {
		log.Println("[ERROR]", err)
	}
	_ = cur.All(ctx, &result)
	fmt.Println(result[0])
	bsonBytes, _ := bson.Marshal(result[0])
	bson.Unmarshal(bsonBytes, &q)
	fmt.Println(q.Count)
	/*if product.Pamount-q.Count > 0 {
		fmt.Println("Available")
	} else {
		fmt.Println("Not Availabe")
	}*/
	availabe := product.Pamount - q.Count
	fmt.Println(availabe)
	if availabe >= 0 {
		fmt.Println("Available")
		if availabe-u.Amount >= 0 {
			res, _ := collection1.InsertOne(ctx, u)
			fmt.Println(res)
			return context.JSON(http.StatusAccepted, "successfully store ")
		} else {
			fmt.Println("No availbe")
		}

	} else {
		fmt.Println("No available")
	}

	return context.JSON(http.StatusAccepted, "Not accepeted ")

}

//code by shoron bhai
func getProduct(id string) (*Response, error) {
	//fmt.Println("Here I am")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8000/products/"+id, nil)
	if req != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	fmt.Println(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in connecting service: %s", err.Error()))
	}
	defer resp.Body.Close()
	if resp.Body != nil {
		jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error in getting response body: %s", err.Error()))
		}

		response := Response{}
		err = json.Unmarshal([]byte(jsonDataFromHttp), &response) // here!
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error in parsing response body: %s", err.Error()))
		}

		// if resp.StatusCode == http.StatusOK {

		// }
		//fmt.Println("Here I am")
		fmt.Println(&response)

		return &response, nil

	} else {
		return nil, errors.New("something went wrong")
	}
}

/*
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
*/
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

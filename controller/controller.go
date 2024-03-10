package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/shubash/pipo/moddel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// var Email string = "shubashkarn"

// var MailPassword string = "wvwmpibufwgrikgs"

// const connectingstring =

const dbname = "pipo"
const colname = "user"

var Collection *mongo.Collection

func init() {
	godotenv.Load(".env")
	
	ClientOptions := options.Client().ApplyURI(os.Getenv("URI"))

	client, err := mongo.Connect(context.TODO(), ClientOptions)
	if err != nil {
		panic(err.Error())

	}
	fmt.Println("coonecting to mongo db")
	Collection = client.Database(dbname).Collection(colname)
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome to pipo website")
}
func Singup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	var user moddel.USER

	_ = json.NewDecoder(r.Body).Decode(&user)
	count, err := Collection.CountDocuments(context.Background(), bson.M{"email": user.Email})
	// count, err := Collection.CountDocuments(context.Background(), bson.M{"email": user.Email})
	if err != nil {
		log.Println(err)
		return
	}
	if count != 0 {
		fmt.Println("this email is used allready ")
		return
	}
	num, err := Collection.CountDocuments(context.Background(), bson.M{"phone": user.Phone})
	if err != nil {
		log.Println(err)
	}
	if num != 0 {
		fmt.Println("this number is use to make a new id ")
	}

	password := Hashpassword(user.Password)
	user.Password = password

	insertoneuser(user)
	json.NewEncoder(w).Encode(user)
}
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json-x-www-from-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var user moddel.USER

	var founduser moddel.USER

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
	}
	err1 := Collection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&founduser)
	if err1 != nil {
		// log.Println(err1)
		log.Println("email is incorect")
		// json.NewEncoder(w).Encode("email incorrect")
		return

	}
	err2 := bcrypt.CompareHashAndPassword([]byte(founduser.Password), []byte(user.Password))
	if err2 != nil {
		log.Println("password does not match")
		return
	}

	json.NewEncoder(w).Encode(founduser)
	fmt.Println("user is on")
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := getalluser()
	json.NewEncoder(w).Encode(user)

}
func DeleteAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json-x-www-from-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	count := deletealluser()
	json.NewEncoder(w).Encode(count)
}

func insertoneuser(user moddel.USER) {
	insert, err := Collection.InsertOne(context.Background(), user)
	if err != nil {
		panic(err)
	}
	fmt.Println("inserted one user in our database", insert.InsertedID)
}
func getalluser() []primitive.M {
	cur, err := Collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		panic(err)
	}
	var users []primitive.M
	for cur.Next(context.Background()) {
		var user bson.M
		err := cur.Decode(&user)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	defer cur.Close(context.Background())
	return users
}

func Hashpassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

// func verifypassword(userPassword string, proviedPassword string) error {
// 	err := bcrypt.CompareHashAndPassword([]byte(proviedPassword), []byte(userPassword))

// 	if err != nil {

// 		panic("password is wrong")
// 		return err
// 	}
// 	return nil

// }

//	func VerifyPassword(userPassword string, proviedPassword string) (bool, string) {
//		err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(userPassword))
//		check := true
//		msg := ""
//		if err != nil {
//			msg = fmt.Sprintln("password is wrong")
//			check = false
//		}
//		return check, msg
//	}
func deletealluser() int64 {
	deleteresult, err := Collection.DeleteMany(context.Background(), bson.M{}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("number of user is delete is ", deleteresult.DeletedCount)
	return (deleteresult.DeletedCount)
}

// func Sendemail(toemail string, message string) {
// 	from := Email
// 	password := MailPassword

// 	smtpHost := "smtp.gmail.com"
// 	smtpPort := "587"

// 	//receiver
// 	to := []string{toemail}

// 	//message

// 	// subject := "subject:forget password"
// 	// body := "the welcome from here"
// 	// message := []byte(subject + body)

// 	//auth

// 	auth := smtp.PlainAuth("", from, password, smtpHost)

// 	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}
// 	fmt.Println("email is send ")
// }

func Updatepassword(userid string, newPassword string) {
	Id, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		log.Println("this name user is not persent")
		return
	}
	filter := bson.M{"_id": Id}
	hashpassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println("hasing problem ")
	}
	update := bson.M{"$set": bson.M{"password": string(hashpassword)}}
	result, err := Collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Println("miss")
		return
	}
	fmt.Println("password has been change ", result.ModifiedCount)
}
func UPDATEPASSWORD(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json-x-www-from-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")
	params := mux.Vars(r)
	newPassword := r.FormValue("newPassword")
	// newPassword := bcrypt.CompareHashAndPassword(userpas)
	Updatepassword(params["id"], newPassword)
}

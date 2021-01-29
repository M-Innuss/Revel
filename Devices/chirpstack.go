package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/M-Innuss/Revel/User"
	"github.com/brocaar/chirpstack-api/go/as/external/api"

	"github.com/brocaar/lorawan"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
)

const (
	// postgres = "sqlite3"
	// memStr1     = "file::memory:?cache=shared"
	postgres = "postgres"
	memStr1  = "user=postgres dbname=postgres sslmode=disable password=postgres"
)

// configuration
var (
	// This must point to the API interface
	server = "localhost:8080"

	// The DevEUI for which we want to enqueue the downlink
	devEUI = lorawan.EUI64{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}

	// The API token (retrieved using the web-interface)
	apiToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcGlfa2V5X2lkIjoiNGM2OTU4ZTAtM2NmMS00NmQzLTlmYzMtNGUyZjZmNzMxM2Q2IiwiYXVkIjoiYXMiLCJpc3MiOiJhcyIsIm5iZiI6MTYxMTcwNTc3NSwic3ViIjoiYXBpX2tleSJ9.hp-sUJZQBeJh-xzrYIuwuYq-5R8mLICFCXEE0D5ZdFM"
)

type APIToken string

func (a APIToken) GetRequestMetadata(ctx context.Context, url ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", a),
	}, nil
}

func (a APIToken) RequireTransportSecurity() bool {
	return false
}

type Person struct {
	UserName string
}

// func about_page(w http.ResponseWriter, r *http.Request) {
// 	fmt.Printf("Main function is working!\n ")
// }

// const AddForm = `
// <form method="POST" action="/ListUserRequest()">
// URL: <input type="text" name="url">
// <input type="submit" value="ListUserRequest()">
// </form>
// `

//-------------------------------Main Function------------------------

func main() {
	fmt.Printf("Main function is working!\n ")

}

// ----------------------- Flash - session ---------------------------
// 	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {

// 		cookie, err := req.Cookie("Session-id")
// 		if err != nil {
// 			id := (uuid.NewV4())
// 			cookie = &http.Cookie{
// 				Name:  "session-id",
// 				Value: id.String(),
// 			}
// 			http.SetCookie(res, cookie)
// 		}
// 		fmt.Println(cookie)
// 		http.ListenAndServe(":9000", nil)

// 	})

// 	// }
// }

//Database Configuration

//Get deviceid from a specific user (email)

func GetDeviceIdFromDB(email string) interface{} {

	// var (
	// 	ctx context.Context
	// )

	type Account struct {
		IdNumber uint
		Email    string
		DeviceId string
	}

	db, err := sql.Open(postgres, memStr1)
	if err != nil {
		log.Fatalf("error opening DB (%s)", err)
	}

	userSql := "SELECT deviceid FROM account WHERE email = $1 LIMIT 10"

	//
	rows, err := db.Query(userSql, email)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	type DeviceInfo struct {
		DeviceId string
	}

	var device DeviceInfo
	var deviceSlice []DeviceInfo

	defer rows.Close()
	for rows.Next() {

		err = rows.Scan(&device.DeviceId)
		if err != nil {
			// handle this error
			panic(err)
		}
		//fmt.Printf("DeviceId of User:  %s \n", device.DeviceId)

		deviceSlice = append(deviceSlice, device)
		//fmt.Println(deviceSlice)
	}
	fmt.Println(deviceSlice)
	return deviceSlice
}

//Create a new table

func CreateTable() {

	//Open/Run Database
	// db, err := sql.Open(postgres, memStr1)
	db, err := sql.Open(postgres, memStr1)
	if err != nil {
		log.Fatalf("error opening DB (%s)", err)
	}

	log.Printf("Creating new table")
	if _, crErr := User.CreateAccountTable(db); crErr != nil {
		log.Fatalf("Error creating table (%s)", crErr)
	}
	log.Printf("Created")
}

//Insert new Account

func InsertNewUserDB(IdNumberInput uint, EmailInput string, DeviceIdInput string) {

	db, err := sql.Open(postgres, memStr1)
	if err != nil {
		log.Fatalf("error opening DB (%s)", err)
	}

	me := User.Account{IdNumber: IdNumberInput, Email: EmailInput, DeviceId: DeviceIdInput}

	if SelectUserDB(IdNumberInput, EmailInput, DeviceIdInput) != me {
		log.Printf("Inserting %+v into the DB", me)
		if _, insErr := User.InsertAccount(db, me); insErr != nil {
			log.Fatalf("Error inserting new Account into the DB (%s)", insErr)
		}
		log.Printf("User saved")
	}
	log.Printf("User already is connected to this device")
}

//Select Account

func SelectUserDB(IdNumberInput uint, EmailInput string, DeviceIdInput string) User.Account {

	db, err := sql.Open(postgres, memStr1)
	if err != nil {
		log.Fatalf("error opening DB (%s)", err)
	}

	log.Printf("Selecting Account from the DB")
	selectedMe := User.Account{}

	if err := User.SelectAccount(db, IdNumberInput, EmailInput, DeviceIdInput, &selectedMe); err != nil {
		//log.Fatalf("Error selecting Account from the DB (%s)", err)
		log.Printf("Error selecting Account from the DB (%s)", err)

	}
	//log.Printf("%+v", selectedMe)
	return selectedMe

}

// //Update Account

func UpdateUserDB(IdNumberInput uint, EmailInput string, DeviceIdInput string) {

	db, err := sql.Open(postgres, memStr1)
	if err != nil {
		log.Fatalf("error opening DB (%s)", err)
	}

	log.Printf("Updating Account in the DB: successful")
	updatedMe := User.Account{
		IdNumber: 1,
		Email:    "martinsinnuss@gmail.com",
		DeviceId: "001a",
	}
	selectedMe := User.Account{}

	if err := User.UpdateAccount(db, selectedMe.IdNumber, selectedMe.Email, selectedMe.DeviceId, updatedMe); err != nil {
		log.Fatalf("Error updating Account in the DB (%s)", err)
	}

}

// //Delete Account

func DeleteUserDB(IdNumberInput uint, EmailInput string, DeviceIdInput string) string {

	db, err := sql.Open(postgres, memStr1)
	if err != nil {
		log.Fatalf("error opening DB (%s)", err)
	}
	log.Printf("Deleting Account from the DB")
	selectedMe := User.Account{}
	if delErr := User.DeleteAccount(db, selectedMe.IdNumber, selectedMe.Email, selectedMe.DeviceId); delErr != nil {
		log.Fatalf("Error deleting Account from the DB (%s)", delErr)
	}
	log.Printf("User is Deleted")
	return EmailInput
}

//----------------------Create device -----------------------------------
func CreateNewDevice(deviceid int64, name string) {
	// fmt.Printf("%v\n", user2)
	// fmt.Printf("%v\n", user3)

	//Converting (deviceid int64) to (deviceidstring str)
	deviceidstring := strconv.FormatInt(deviceid, 10)

	// Define the new user with IsAdmin(bool), IsActive(bool), Email(String)
	DeviceProfile1 := &api.DeviceProfile{
		Id:   deviceidstring,
		Name: name,
	}

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(APIToken(apiToken)),
		grpc.WithInsecure(), // remove this when using TLS
	}
	// connect to the gRPC server
	conn, err := grpc.Dial(server, dialOpts...)
	if err != nil {
		panic(err)
	}

	// define the DeviceQueueService client
	NewDevice := api.NewDeviceProfileServiceClient(conn)
	// Save with CreateUserRequest the new user(User) with Password(String)
	resp, err := NewDevice.Create(context.Background(), &api.CreateDeviceProfileRequest{
		DeviceProfile: DeviceProfile1,
	})

	if err != nil {
		panic(err)

	}

	// Print out the new user's Id
	fmt.Printf("The id number of the saved device: %+v \n ", resp.Id)

}

// ------------------------Get Data Device-------------------------------
func GetDeviceData(DeviceId string) {

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(APIToken(apiToken)),
		grpc.WithInsecure(), // remove this when using TLS
	}
	// connect to the gRPC server
	conn, err := grpc.Dial(server, dialOpts...)
	if err != nil {
		panic(err)
	}

	// define the DeviceQueueService client
	NewDevice := api.NewDeviceServiceClient(conn)
	// Save with CreateUserRequest the new user(User) with Password(String)
	resp, err := NewDevice.Get(context.Background(), &api.GetDeviceRequest{
		DevEui: DeviceId,
	})

	if err != nil {
		panic(err)

	}
	fmt.Printf("Information about the device with Id: %+v \n ", resp.Device.DevEui)
	fmt.Printf("Device's name: %+v \n ", resp.Device.Name)
	fmt.Printf("Batterystatus: %+v \n ", resp.DeviceStatusBattery)
	fmt.Printf("Location: %+v \n ", resp.Location)
	fmt.Printf("Description: %+v \n ", resp.Device.Description)

}

// ---------------------------GetDeviceBattery-------------------------------------

func GetDeviceBattery(DeviceId string) uint32 {

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(APIToken(apiToken)),
		grpc.WithInsecure(), // remove this when using TLS
	}
	// connect to the gRPC server
	conn, err := grpc.Dial(server, dialOpts...)
	if err != nil {
		panic(err)
	}

	// define the DeviceQueueService client
	NewDevice := api.NewDeviceServiceClient(conn)
	// Save with CreateUserRequest the new user(User) with Password(String)
	resp, err := NewDevice.Get(context.Background(), &api.GetDeviceRequest{
		DevEui: DeviceId,
	})

	if err != nil {
		panic(err)

	}

	fmt.Printf("Batterystatus: %+v \n ", resp.DeviceStatusBattery)
	return resp.DeviceStatusBattery
}

//-----------------------Get device name------------------------------

func GetDeviceName(DeviceId string) string {

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(APIToken(apiToken)),
		grpc.WithInsecure(), // remove this when using TLS
	}
	// connect to the gRPC server
	conn, err := grpc.Dial(server, dialOpts...)
	if err != nil {
		panic(err)
	}

	// define the DeviceQueueService client
	NewDevice := api.NewDeviceServiceClient(conn)
	// Save with CreateUserRequest the new user(User) with Password(String)
	resp, err := NewDevice.Get(context.Background(), &api.GetDeviceRequest{
		DevEui: DeviceId,
	})

	if err != nil {
		panic(err)

	}

	fmt.Printf("Device's name: %+v \n ", resp.Device.Name)
	return resp.Device.Name
}

// ------------------------Activate Device-------------------------------

func ActiveDevice(DeviceId string) {

	device := &api.DeviceActivation{
		DevEui: DeviceId,
	}

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(APIToken(apiToken)),
		grpc.WithInsecure(), // remove this when using TLS
	}
	// connect to the gRPC server
	conn, err := grpc.Dial(server, dialOpts...)
	if err != nil {
		panic(err)
	}

	// define the DeviceQueueService client
	NewDevice := api.NewDeviceServiceClient(conn)
	// Save with CreateUserRequest the new user(User) with Password(String)
	resp, err := NewDevice.Activate(context.Background(), &api.ActivateDeviceRequest{
		DeviceActivation: device,
	})

	if err != nil {
		panic(err)

	}

	// Print out the new user's Id
	fmt.Printf("Activate: %+v \n ", resp)

}

// ------------------------Delete Device-------------------------------
func DeleteDevice(DeviceId string) {

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(APIToken(apiToken)),
		grpc.WithInsecure(), // remove this when using TLS
	}
	// connect to the gRPC server
	conn, err := grpc.Dial(server, dialOpts...)
	if err != nil {
		panic(err)
	}

	NewDevice := api.NewDeviceServiceClient(conn)

	resp, err := NewDevice.Delete(context.Background(), &api.DeleteDeviceRequest{
		DevEui: DeviceId,
	})

	if err != nil {
		panic(err)

	}

	// Print out the new user's Id
	fmt.Printf("The device has been successfully deleted!: %+v \n ", resp)

}

// ------------------------Update Device-------------------------------

func UpdateDeviceProfileName(id string, name string) {

	deviceprofile := &api.DeviceProfile{
		Id:   id,
		Name: name,
	}

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(APIToken(apiToken)),
		grpc.WithInsecure(), // remove this when using TLS
	}
	// connect to the gRPC server
	conn, err := grpc.Dial(server, dialOpts...)
	if err != nil {
		panic(err)
	}

	NewDevice := api.NewDeviceProfileServiceClient(conn)

	resp, err := NewDevice.Update(context.Background(), &api.UpdateDeviceProfileRequest{
		DeviceProfile: deviceprofile,
	})

	if err != nil {
		panic(err)

	}

	// Print out the new user's Id
	fmt.Printf("The name of device profile has been successfully updated!: %+v \n ", resp)

}

// ------------------------List Device Data-------------------------------

// func ListDeviceData() {

// 	dialOpts := []grpc.DialOption{
// 		grpc.WithBlock(),
// 		grpc.WithPerRPCCredentials(APIToken(apiToken)),
// 		grpc.WithInsecure(), // remove this when using TLS
// 	}
// 	// connect to the gRPC server
// 	conn, err := grpc.Dial(server, dialOpts...)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// define the DeviceQueueService client
// 	NewDevice := api.NewDeviceServiceClient(conn)
// 	// Save with CreateUserRequest the new user(User) with Password(String)
// 	resp, err := NewDevice.List(context.Background(), &api.ListDeviceRequest{
// 		Limit:  10,
// 		Offset: 0,
// 	})

// 	if err != nil {
// 		panic(err)

// 	}

// 	// Print out
// 	fmt.Printf("Number of saved devices: %+v \n ", resp.TotalCount)

// }

// ------------------------Save new user-------------------------------

func SaveNewUser(Email string, Password string, Password2 string) {
	// fmt.Printf("%v\n", user2)
	// fmt.Printf("%v\n", user3)

	// Define the new user with IsAdmin(bool), IsActive(bool), Email(String)
	user := &api.User{
		IsAdmin:  false,
		IsActive: true,
		Email:    Email,
	}

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(APIToken(apiToken)),
		grpc.WithInsecure(), // remove this when using TLS
	}
	// connect to the gRPC server
	conn, err := grpc.Dial(server, dialOpts...)

	if err != nil || Password != Password2 {
		panic(err)

	}

	// define the Service client
	userClient := api.NewUserServiceClient(conn)
	// Save with CreateUserRequest the new user(User) with Password(String)
	resp, err := userClient.Create(context.Background(), &api.CreateUserRequest{
		User:     user,
		Password: Password,
	})

	if err != nil {
		panic(err)
	}

	// Print out the new user's Id
	fmt.Printf("New user has been saved! The id number of the saved user: %+v \n ", resp.Id)

}

// ------------------------List alle saved users----------------------------
// Show all users with ListUserRequest

func ListUserRequest() {

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(APIToken(apiToken)),
		grpc.WithInsecure(), // remove this when using TLS
	}
	// connect to the gRPC server
	conn, err := grpc.Dial(server, dialOpts...)
	if err != nil {
		panic(err)
	}
	// define the DeviceQueueService client
	userClient := api.NewUserServiceClient(conn)
	resp, err := userClient.List(context.Background(), &api.ListUserRequest{

		Limit:  20,
		Offset: 0,
	})

	if err != nil {
		panic(err)
	}

	// List all users with ListUserResponse

	// fmt.Printf("Id: %+v \n ", resp.Result[0].Id)

	// fmt.Printf("Email: %+v \n ", resp.Result[0].Username)

	// s := fmt.Sprintf("%v\n", resp.TotalCount)

	/** converting the str1 variable into an int using Atoi method */

	// s1 := "2"

	fmt.Printf("Number of saved users: %+v \n ", resp.TotalCount)

	fmt.Println("All saved users: \n ")
	for i := 0; i < int(resp.TotalCount); i++ {
		fmt.Printf("Id: %v, Email: %+v \n", resp.Result[i].Id, resp.Result[i].Username)

	}

}

// func GetCookienName(Email string) int64 {

// 	dialOpts := []grpc.DialOption{
// 		grpc.WithBlock(),
// 		grpc.WithPerRPCCredentials(APIToken(apiToken)),
// 		grpc.WithInsecure(), // remove this when using TLS
// 	}

// }

func GetUserId(Email string) int64 {

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(APIToken(apiToken)),
		grpc.WithInsecure(), // remove this when using TLS
	}
	// connect to the gRPC server
	conn, err := grpc.Dial(server, dialOpts...)
	if err != nil {
		panic(err)
	}
	// define the DeviceQueueService client
	userClient := api.NewUserServiceClient(conn)
	resp, err := userClient.List(context.Background(), &api.ListUserRequest{

		Limit:  20,
		Offset: 0,
	})

	if err != nil {
		panic(err)
	}

	// List all users with ListUserResponse
	// fmt.Printf("Number of saved users: %+v \n ", resp.TotalCount)

	// fmt.Printf("Id: %+v \n ", resp.Result[0].Id)

	// fmt.Printf("Email: %+v \n ", resp.Result[0].Username)

	// s := fmt.Sprintf("%v\n", resp.TotalCount)

	/** converting the str1 variable into an int using Atoi method */

	// s1 := "2"

	for i := 0; i < int(resp.TotalCount); i++ {
		//fmt.Printf("Id: %v, Email: %+v \n", resp.Result[i].Id, resp.Result[i].Username)
		if Email == resp.Result[i].Username {
			fmt.Printf("Id of the registered user: %v \n", resp.Result[i].Id)
			return resp.Result[i].Id

		}

	}
	return 0

}

// ----------------------Delete user-----------------------------------------

func DeleteUser(id int64) {

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(APIToken(apiToken)),
		grpc.WithInsecure(), // remove this when using TLS
	}
	// connect to the gRPC server
	conn, err := grpc.Dial(server, dialOpts...)
	if err != nil {
		panic(err)
	}
	// define the DeviceQueueService client
	userClient := api.NewUserServiceClient(conn)

	// Write the id number of the user, that will be deleted
	resp, err := userClient.Delete(context.Background(), &api.DeleteUserRequest{
		Id: id,
	})

	if err != nil {
		panic(err)

	}
	fmt.Printf("The user has been successfully deleted! %+v \n ", resp)
}

// ----------------------Update user-----------------------------------------

func UpdateUser(id int64, email string, isActive bool, isAdmin bool) {

	user := &api.User{
		//Write the id number of the user, that his email, IsActive, IsActive can be updated
		Id: id,
		//Update the email
		Email: email,
		//Update the IsActive (bool)
		IsActive: isActive,
		//Update the IsAdmin(bool)
		IsAdmin: isAdmin,
	}

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(APIToken(apiToken)),
		grpc.WithInsecure(), // remove this when using TLS
	}
	// connect to the gRPC server
	conn, err := grpc.Dial(server, dialOpts...)
	if err != nil {
		panic(err)
	}
	// define the DeviceQueueService client
	userClient := api.NewUserServiceClient(conn)

	//Write the user, that will be updated
	resp, err := userClient.Update(context.Background(), &api.UpdateUserRequest{
		User: user,
	})

	if err != nil {
		panic(err)
	}
	fmt.Printf("The user with email "+email+" has been successfully updated! %+v \n ", resp)

}

//------------------------Get data for a particular user---------------------

func GetDataUser(id int64) {

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(APIToken(apiToken)),
		grpc.WithInsecure(), // remove this when using TLS
	}
	// connect to the gRPC server
	conn, err := grpc.Dial(server, dialOpts...)
	if err != nil {
		panic(err)
	}
	// define the DeviceQueueService client
	userClient := api.NewUserServiceClient(conn)

	//Write the id(int) from the user, from that the data will be shown
	resp, err := userClient.Get(context.Background(), &api.GetUserRequest{
		Id: id,
	})

	if err != nil {
		panic(err)
	}
	fmt.Printf("Here is the information about the user: %+v \n ", resp)

}

// ------------------------Update the password------------------

func UpdatePassword(id int64, password string) {

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(APIToken(apiToken)),
		grpc.WithInsecure(), // remove this when using TLS
	}
	// connect to the gRPC server
	conn, err := grpc.Dial(server, dialOpts...)
	if err != nil {
		panic(err)
	}
	// define the DeviceQueueService client
	userClient := api.NewUserServiceClient(conn)

	//Define the user thats password will be updated
	resp, err := userClient.UpdatePassword(context.Background(), &api.UpdateUserPasswordRequest{
		UserId:   id,
		Password: password,
	})

	if err != nil {
		panic(err)
	}
	fmt.Printf("The password has been successfully updated! Neues passwort: %+v \n ", resp)
	fmt.Printf("%s\n", password)

}

//---------------------------Login---------------------------

func LoginRequest(email string, password string) {

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(APIToken(apiToken)),
		grpc.WithInsecure(), // remove this when using TLS
	}
	// connect to the gRPC server
	conn, err := grpc.Dial(server, dialOpts...)
	if err != nil {
		panic(err)
	}
	// define the DeviceQueueService client
	userClient := api.NewInternalServiceClient(conn)

	resp, err := userClient.Login(context.Background(), &api.LoginRequest{
		Username: email,
		Password: password,
	})

	if err != nil {
		panic(err)

	}
	fmt.Printf("Congratulation "+email+" , you have successfully logged in! %+v \n ", resp)

}

// ------------------------Data base for users (postgresql) ---------------------------------

// func GetStuff() MyStruct {

// 	sql := "SELECT id, name from my_table "
// 	rows, err := app.DB.Query(sql)

// }

// func InitDB() {
// 	driver := revel.Config.StringDefault("db.driver", "mysql")
// 	connect_string := revel.Config.StringDefault("db.connect", "root:root@localhost/test")

// 	Db, err = sql.Open(driver, connect_string)

// }

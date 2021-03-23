package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/CaninoDev/gastro/server/api"
	"github.com/CaninoDev/gastro/server/internal/config"
)

var configYAML   = flag.String("c", "config.yml", "configure db")
var db *gorm.DB

func main() {
	flag.Parse()
	_, databaseC, _, _, err := config.Load(*configYAML)
	if err != nil {
		log.Fatalf("error parsing config.yml %v", err)
	}

	gormDB, err := Start(databaseC)
	if err != nil {
		log.Panic(err)
	}
	db = gormDB
	err = PopulateDB()
	if err != nil {
		log.Fatal(err)
	}
}

// Start returns a configured instance of db{}
func Start(cfg config.Database) (*gorm.DB, error) {

	return open(cfg, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		Logger: newLogger(),
	})
}

func newLogger() logger.Interface {
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 100 * time.Millisecond,
			Colorful:      true,
			LogLevel:      logger.Info,
		})
	return newLogger
}

func open(dbConf config.Database, gormCfg *gorm.Config) (*gorm.DB, error) {
	var dialect gorm.Dialector

	switch strings.ToLower(dbConf.Type) {
	case "mysql":
		dialect = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&allowOldPasswords=1",
			dbConf.User, dbConf.Pass, dbConf.Host, dbConf.Port, dbConf.Name))
	case "postgres":
		dialect = postgres.Open(fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
			dbConf.User, dbConf.Pass, dbConf.Host, dbConf.Port, dbConf.Name))
	case "sqlite3":
		dialect = sqlite.Open(fmt.Sprintf("file:%s?_auth&_auth_user=%s&_auth_pass=%s",
			dbConf.Host, dbConf.User, dbConf.Name))
	default:
		return &gorm.DB{}, fmt.Errorf("%s is unsupported", dbConf.Type)
	}



	db, err := gorm.Open(dialect, gormCfg)
	if err != nil {
		return db, err
	}

	return db, err
}

func StrPtr(s string) *string {
	return &s
}

// PopulateDatabase populates the db with sample data.
func PopulateDB () error {
	// Drop all Tables
	db.Migrator().DropTable(&api.Section{})
	db.Migrator().DropTable(&api.Item{})
	db.Migrator().DropTable(&api.User{})
	db.Migrator().DropTable(&api.Account{})
	fmt.Print("Old tables have been deleted...")

	// Migrate model over to db
	err := db.AutoMigrate(&api.Section{}, &api.Item{}, &api.User{}, &api.Account{})
	if err != nil {
		return fmt.Errorf("error migrating scheme to db: %v", err)
	}

	log.Println("New tables are  migrated successfully!")

	bagel := api.Item{Title: "Bagel", Description: StrPtr("Your choice of H&H bagel"), ListOrder: 1, Price: 395, Active: true}
	bagelwcreamcheese := api.Item{Title: "Bagel w/ Cream Cheese", Description: StrPtr("Toasted H&H Bagel with your choice of cream cheese."), ListOrder: 2, Price: 595, Active: true}
	bagelwlox := api.Item{Title: "Bagel with Lox", Description: StrPtr("Your choice of H&H bagels and Atlantic smoked lox."), ListOrder: 3, Price: 995, Active: true}

	bagels := api.Section{Title: "Bagels", ListOrder: 2, Items: []api.Item{bagel, bagelwcreamcheese, bagelwlox}}

	waffle := api.Item{Title: "Waffles", Description: StrPtr("5 slices of thick homemade waffles with light cream."), ListOrder: 1, Price: 775, Active: true}
	montecristowaffle := api.Item{Title: "Monte Cristo Waffle Sandwich", Description: StrPtr("Whole wheat Waffle sandwich with swiss cheese, raspberry jam, honey baked, and ham."), ListOrder: 2, Price: 1125, Active: true}
	spicywaffle := api.Item{Title: "Southern Waffle Sandwich", Description: StrPtr("Belgian waffles with spicy maple syrup, cheddar cheese, butter and cinnamon."), ListOrder: 3, Price: 1035, Active: true}

	waffles := api.Section{Title: "Waffles", ListOrder: 3, Items: []api.Item{waffle, montecristowaffle, spicywaffle}}

	eggswbacon := api.Item{Title: "Eggs (Any style)", Description: StrPtr("Fresh farm eggs cooked your way with thick applewood smoked bacon."), Price: 725, Active: true}
	eggsbenedict := api.Item{Title: "Eggs Benedict", Description: StrPtr("Eggs Benedict with homemade Hollandaise sauce."), Price: 775, Active: true}
	countryomelettes := api.Item{Title: "Country Omelettes", Description: StrPtr("Omelette stuffed with chorizo, green peppers, onion, and manchego cheese."), Price: 1295, ListOrder: 1, Active: true}
	classicomelettes := api.Item{Title: "Classic Omelettes", Description: StrPtr("Omelette with select herbs and freshly ground black pepper."), Price: 1095, ListOrder: 2, Active: true}
	westernomelette := api.Item{Title: "Western Omelettes", Description: StrPtr("Omelette with green bell pepper, red bell pepper, and monteray jack."), Price: 995, ListOrder: 3, Active: true}

	eggs := api.Section{Title: "Eggs & Omelettes", ListOrder: 4, Items: []api.Item{eggsbenedict, eggswbacon, countryomelettes, classicomelettes, westernomelette}}

	breakfast := &api.Section{Title: "Breakfast", ListOrder: 1, SubSections: []api.Section{bagels, eggs, waffles}}
	db.Create(&breakfast)

	sunomonosalad := api.Item{Title: "Sunomono Salad", Description: StrPtr("Thin rice noodles, shrimp, crab, soy sauce and rice vinegar."), Price: 395, ListOrder: 1, Active: true}
	cobbsalad := api.Item{Title: "Cobb Salad", Description: StrPtr("Blue cheese, grilled chicken breasts, red wine vinegar, eggs, and bacon."), Price: 645, Active: true}
	handpies := api.Item{Title: "Korean Beef Hand Pies", Description: StrPtr("Beef short rubes, rice noodles, hoisin sauce, chili sauce, and soy sauce."), Price: 795, Active: true}
	bruschetta := api.Item{Title: "Bruchetta", Description: StrPtr("Toasted baguettes with goat cheese, brown sugar, and cherry tomatoes."), Price: 495, Active: true}

	starters := api.Section{Title: "Starters", ListOrder: 9, Items: []api.Item{sunomonosalad, cobbsalad, handpies, bruschetta}}

	classictunamelt := api.Item{Title: "Classic Tuna Melt", Description: StrPtr("Sourdough bread, fresh tuna, red onions, dill pickles, celery and butter."), Price: 895, Active: true}
	hamdandcheesesandwich := api.Item{Title: "Perfect Ham and Cheese Sandwich", Description: StrPtr("Sourdough bread, swiss cheese, ham, honey, mustard, mayonnaise, and pickle."), Price: 995, Active: true}
	lemonchickenwraps := api.Item{Title: "Lemon Chicken Wrap", Description: StrPtr("Pita bread, grilled chicken breast, greek yogurt, garlic, Sriracha sauce, paprika."), Price: 1095, Active: true}
	hasselbacktomatoclubs := api.Item{Title: "Hassel Back Tomato Club", Description: StrPtr("Bibb lettuce leaves, ripe avocados, swiss cheese, plum tomatoes, and turkey"), Price: 1095, Active: true}

	sandwiches := api.Section{Title: "Sandwiches", ListOrder: 7, Items: []api.Item{classictunamelt, hamdandcheesesandwich, lemonchickenwraps, hasselbacktomatoclubs}}

	carrotgingersoup := api.Item{Title: "Carrot Ginger Soup", Description: StrPtr("Carrot ginger soup with coconut milk, apple cider vinegar, and maple syrup"), Price: 895, Active: true}
	classicchickensoup := api.Item{Title: "Homemade Chicken Soup", Description: StrPtr("Chicken soup with Israeli couscous."), Price: 695, Active: true}

	soups := api.Section{Title: "Soups", ListOrder: 6, Items: []api.Item{carrotgingersoup, classicchickensoup}}

	lunch := &api.Section{Title: "Lunch", ListOrder: 5, SubSections: []api.Section{soups, sandwiches}}
	db.Create(&lunch)

	entrees := api.Section{Title: "Entrées", ListOrder: 10}
	desserts := api.Section{Title: "Desserts", ListOrder: 11}

	dinner := &api.Section{Title: "Dinner", ListOrder: 8, SubSections: []api.Section{starters, entrees, desserts}}
	db.Create(&dinner)

	return nil
}




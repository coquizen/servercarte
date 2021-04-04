package gormDB

import (
	"fmt"
	"github.com/CaninoDev/gastro/server/domain/account"
	"github.com/CaninoDev/gastro/server/domain/menu"
	"github.com/CaninoDev/gastro/server/domain/user"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/CaninoDev/gastro/server/internal/logger"
)

func StrPtr(s string) *string {
	return &s
}

func seedDB(db *gorm.DB) error {
	// Drop all Tables
	var migrator = db.Migrator()
	if err := migrator.DropTable(&menu.Section{}, &menu.Item{}, &user.User{}, &account.Account{}); err != nil {
		return err
	}

	logger.Info.Println("Accounts tables have been dropped...")
	logger.Info.Println("Now migrating...")
	// Migrate model over to db
	err := db.AutoMigrate(&menu.Section{}, &menu.Item{}, &user.User{}, &account.Account{})
	if err != nil {
		return fmt.Errorf("error migrating scheme to db: %v", err)
	}
	logger.Info.Println("Successful")

	logger.Info.Println("Seeding...")

	bagel := menu.Item{Title: "Bagel", Description: StrPtr("Your choice of H&H bagel"), Type: 0, ListOrder: 1,
		Price:  395,
		Active: true}
	bagelwcreamcheese := menu.Item{Title: "Bagel w/ Cream Cheese", Description: StrPtr("Toasted H&H Bagel with your choice of cream cheese."), Type: 0, ListOrder: 2, Price: 595, Active: true}
	bagelwlox := menu.Item{Title: "Bagel with Lox", Description: StrPtr("Your choice of H&H bagels and Atlantic smoked lox."), Type: 0, ListOrder: 3, Price: 995, Active: true}

	bagels := menu.Section{Title: "Bagels", ListOrder: 2, Type: 1, Items: []menu.Item{bagel, bagelwcreamcheese, bagelwlox}}

	waffle := menu.Item{Title: "Waffles", Description: StrPtr("5 slices of thick homemade waffles with light cream."), Type: 0, ListOrder: 1, Price: 775, Active: true}
	montecristowaffle := menu.Item{Title: "Monte Cristo Waffle Sandwich", Description: StrPtr("Whole wheat Waffle sandwich with swiss cheese, raspberry jam, honey baked, and ham."), Type: 0, ListOrder: 2, Price: 1125, Active: true}
	spicywaffle := menu.Item{Title: "Southern Waffle Sandwich", Description: StrPtr("Belgian waffles with spicy maple syrup, cheddar cheese, butter and cinnamon."), Type: 0, ListOrder: 3, Price: 1035, Active: true}

	waffles := menu.Section{Title: "Waffles", ListOrder: 3, Type: 1, Items: []menu.Item{waffle, montecristowaffle, spicywaffle}}

	eggswbacon := menu.Item{Title: "Eggs (Any style)", Description: StrPtr("Fresh farm eggs cooked your way with thick applewood smoked bacon."), Type: 0, Price: 725, Active: true}
	eggsbenedict := menu.Item{Title: "Eggs Benedict", Description: StrPtr("Eggs Benedict with homemade Hollandaise sauce."), Type: 0, Price: 775, Active: true}
	countryomelettes := menu.Item{Title: "Country Omelettes", Description: StrPtr("Omelette stuffed with chorizo, green peppers, onion, and manchego cheese."), Type: 0, Price: 1295, ListOrder: 1, Active: true}
	classicomelettes := menu.Item{Title: "Classic Omelettes", Description: StrPtr("Omelette with select herbs and freshly ground black pepper."), Type: 0, Price: 1095, ListOrder: 2, Active: true}
	westernomelette := menu.Item{Title: "Western Omelettes", Description: StrPtr("Omelette with green bell pepper, red bell pepper, and monteray jack."), Type: 0, Price: 995, ListOrder: 3, Active: true}

	eggs := menu.Section{Title: "Eggs & Omelettes", ListOrder: 4, Type: 1, Items: []menu.Item{eggsbenedict, eggswbacon, countryomelettes, classicomelettes, westernomelette}}

	breakfast := &menu.Section{Title: "Breakfast", ListOrder: 1, Type: 0, SubSections: []menu.Section{bagels, eggs, waffles}}
	db.Create(&breakfast)

	sunomonosalad := menu.Item{Title: "Sunomono Salad", Description: StrPtr("Thin rice noodles, shrimp, crab, soy sauce and rice vinegar."), Type: 0, Price: 395, ListOrder: 1, Active: true}
	cobbsalad := menu.Item{Title: "Cobb Salad", Description: StrPtr("Blue cheese, grilled chicken breasts, red wine vinegar, eggs, and bacon."), Type: 0, Price: 645, Active: true}
	handpies := menu.Item{Title: "Korean Beef Hand Pies", Description: StrPtr("Beef short rubes, rice noodles, hoisin sauce, chili sauce, and soy sauce."), Type: 0, Price: 795, Active: true}
	bruschetta := menu.Item{Title: "Bruchetta", Description: StrPtr("Toasted baguettes with goat cheese, brown sugar, and cherry tomatoes."), Type: 0, Price: 495, Active: true}

	starters := menu.Section{Title: "Starters", ListOrder: 9, Items: []menu.Item{sunomonosalad, cobbsalad, handpies, bruschetta}}

	classictunamelt := menu.Item{Title: "Classic Tuna Melt", Description: StrPtr("Sourdough bread, fresh tuna, red onions, dill pickles, celery and butter."), Type: 0, Price: 895, Active: true}
	hamdandcheesesandwich := menu.Item{Title: "Perfect Ham and Cheese Sandwich", Description: StrPtr("Sourdough bread, swiss cheese, ham, honey, mustard, mayonnaise, and pickle."), Type: 0, Price: 995, Active: true}
	lemonchickenwraps := menu.Item{Title: "Lemon Chicken Wrap", Description: StrPtr("Pita bread, grilled chicken breast, greek yogurt, garlic, Sriracha sauce, paprika."), Type: 0, Price: 1095, Active: true}
	hasselbacktomatoclubs := menu.Item{Title: "Hassel Back Tomato Club", Description: StrPtr("Bibb lettuce leaves, ripe avocados, swiss cheese, plum tomatoes, and turkey"), Type: 0, Price: 1095, Active: true}

	sandwiches := menu.Section{Title: "Sandwiches", ListOrder: 7, Type: 1, Items: []menu.Item{classictunamelt, hamdandcheesesandwich, lemonchickenwraps, hasselbacktomatoclubs}}

	carrotgingersoup := menu.Item{Title: "Carrot Ginger Soup", Description: StrPtr("Carrot ginger soup with coconut milk, apple cider vinegar, and maple syrup"), Type: 0, Price: 895, Active: true}
	classicchickensoup := menu.Item{Title: "Homemade Chicken Soup", Description: StrPtr("Chicken soup with Israeli couscous."), Type: 0, Price: 695, Active: true}

	soups := menu.Section{Title: "Soups", ListOrder: 6, Type: 1, Items: []menu.Item{carrotgingersoup, classicchickensoup}}

	lunch := &menu.Section{Title: "Lunch", ListOrder: 5, Type: 0, SubSections: []menu.Section{soups, sandwiches}}
	db.Create(&lunch)

	entrees := menu.Section{Title: "Entrées", Type: 1, ListOrder: 10}
	desserts := menu.Section{Title: "Desserts", Type: 1, ListOrder: 11}

	dinner := &menu.Section{Title: "Dinner", Type: 0, ListOrder: 8, SubSections: []menu.Section{starters, entrees, desserts}}
	db.Create(&dinner)
	logger.Info.Println("Successful")

	// Now create admin user/account
	logger.Info.Println("Creating default admin user/account...")
	password := "administrator"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	db.Create(&account.Account{
		Username: "admin",
		Password: string(hashedPassword),
		Role:     account.Admin,
		User: user.User{
			FirstName:       "Auguste",
			LastName:        "Gusteau",
			Address1:        "31 Rue Cambon",
			ZipCode:         75001,
			Email:           "admin@Gusteaus.com",
			TelephoneNumber: "7185550193",
		},
	})

	logger.Info.Println("Successful")
	logger.Info.Println("Creating employee account/user...")

	password = "employee"
	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	db.Create(&account.Account{
		Username: "employee",
		Password: string(hashedPassword),
		Role:     account.Employee,
		User: user.User{
			FirstName:       "Remy",
			LastName:        "Ratatouille",
			Address1:        "10 Rue Egout",
			ZipCode:         75002,
			Email:           "remy@pixar.com",
			TelephoneNumber: "7185550192",
		},
	})

	logger.Info.Println("Successful")
	logger.Info.Println("Creating guest account/user...")
	password = "guest"
	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	db.Create(&account.Account{
		Username: "guest",
		Password: string(hashedPassword),
		Role:     account.Guest,
		User: user.User{
			FirstName:       "Anton",
			LastName:        "Ego",
			Address1:        "99 Tour D'Ivoire",
			ZipCode:         75003,
			Email:           "anton@divoire.com",
			TelephoneNumber: "7185550194",
		},
	})

	logger.Info.Println("Successful")

	return nil
}

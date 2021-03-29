package gormDB

import (
	"fmt"

	"github.com/CaninoDev/gastro/server/api/account"
	"github.com/CaninoDev/gastro/server/api/menu"
	"github.com/CaninoDev/gastro/server/api/user"

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

	logger.Info.Println("List tables have been dropped...")
	logger.Info.Println("Now migrating...")
	// Migrate model over to db
	err := db.AutoMigrate(&menu.Section{}, &menu.Item{}, &user.User{}, &account.Account{})
	if err != nil {
		return fmt.Errorf("error migrating scheme to db: %v", err)
	}
	logger.Info.Println("Successful")

	logger.Info.Println("Seeding...")
	bagel := menu.Item{Title: "Bagel", Description: StrPtr("Your choice of H&H bagel"), ListOrder: 1, Price: 395, Active: true}
	bagelwcreamcheese := menu.Item{Title: "Bagel w/ Cream Cheese", Description: StrPtr("Toasted H&H Bagel with your choice of cream cheese."), ListOrder: 2, Price: 595, Active: true}
	bagelwlox := menu.Item{Title: "Bagel with Lox", Description: StrPtr("Your choice of H&H bagels and Atlantic smoked lox."), ListOrder: 3, Price: 995, Active: true}

	bagels := menu.Section{Title: "Bagels", ListOrder: 2, Items: []menu.Item{bagel, bagelwcreamcheese, bagelwlox}}

	waffle := menu.Item{Title: "Waffles", Description: StrPtr("5 slices of thick homemade waffles with light cream."), ListOrder: 1, Price: 775, Active: true}
	montecristowaffle := menu.Item{Title: "Monte Cristo Waffle Sandwich", Description: StrPtr("Whole wheat Waffle sandwich with swiss cheese, raspberry jam, honey baked, and ham."), ListOrder: 2, Price: 1125, Active: true}
	spicywaffle := menu.Item{Title: "Southern Waffle Sandwich", Description: StrPtr("Belgian waffles with spicy maple syrup, cheddar cheese, butter and cinnamon."), ListOrder: 3, Price: 1035, Active: true}

	waffles := menu.Section{Title: "Waffles", ListOrder: 3, Items: []menu.Item{waffle, montecristowaffle, spicywaffle}}

	eggswbacon := menu.Item{Title: "Eggs (Any style)", Description: StrPtr("Fresh farm eggs cooked your way with thick applewood smoked bacon."), Price: 725, Active: true}
	eggsbenedict := menu.Item{Title: "Eggs Benedict", Description: StrPtr("Eggs Benedict with homemade Hollandaise sauce."), Price: 775, Active: true}
	countryomelettes := menu.Item{Title: "Country Omelettes", Description: StrPtr("Omelette stuffed with chorizo, green peppers, onion, and manchego cheese."), Price: 1295, ListOrder: 1, Active: true}
	classicomelettes := menu.Item{Title: "Classic Omelettes", Description: StrPtr("Omelette with select herbs and freshly ground black pepper."), Price: 1095, ListOrder: 2, Active: true}
	westernomelette := menu.Item{Title: "Western Omelettes", Description: StrPtr("Omelette with green bell pepper, red bell pepper, and monteray jack."), Price: 995, ListOrder: 3, Active: true}

	eggs := menu.Section{Title: "Eggs & Omelettes", ListOrder: 4, Items: []menu.Item{eggsbenedict, eggswbacon, countryomelettes, classicomelettes, westernomelette}}

	breakfast := &menu.Section{Title: "Breakfast", ListOrder: 1, SubSections: []menu.Section{bagels, eggs, waffles}}
	db.Create(&breakfast)

	sunomonosalad := menu.Item{Title: "Sunomono Salad", Description: StrPtr("Thin rice noodles, shrimp, crab, soy sauce and rice vinegar."), Price: 395, ListOrder: 1, Active: true}
	cobbsalad := menu.Item{Title: "Cobb Salad", Description: StrPtr("Blue cheese, grilled chicken breasts, red wine vinegar, eggs, and bacon."), Price: 645, Active: true}
	handpies := menu.Item{Title: "Korean Beef Hand Pies", Description: StrPtr("Beef short rubes, rice noodles, hoisin sauce, chili sauce, and soy sauce."), Price: 795, Active: true}
	bruschetta := menu.Item{Title: "Bruchetta", Description: StrPtr("Toasted baguettes with goat cheese, brown sugar, and cherry tomatoes."), Price: 495, Active: true}

	starters := menu.Section{Title: "Starters", ListOrder: 9, Items: []menu.Item{sunomonosalad, cobbsalad, handpies, bruschetta}}

	classictunamelt := menu.Item{Title: "Classic Tuna Melt", Description: StrPtr("Sourdough bread, fresh tuna, red onions, dill pickles, celery and butter."), Price: 895, Active: true}
	hamdandcheesesandwich := menu.Item{Title: "Perfect Ham and Cheese Sandwich", Description: StrPtr("Sourdough bread, swiss cheese, ham, honey, mustard, mayonnaise, and pickle."), Price: 995, Active: true}
	lemonchickenwraps := menu.Item{Title: "Lemon Chicken Wrap", Description: StrPtr("Pita bread, grilled chicken breast, greek yogurt, garlic, Sriracha sauce, paprika."), Price: 1095, Active: true}
	hasselbacktomatoclubs := menu.Item{Title: "Hassel Back Tomato Club", Description: StrPtr("Bibb lettuce leaves, ripe avocados, swiss cheese, plum tomatoes, and turkey"), Price: 1095, Active: true}

	sandwiches := menu.Section{Title: "Sandwiches", ListOrder: 7, Items: []menu.Item{classictunamelt, hamdandcheesesandwich, lemonchickenwraps, hasselbacktomatoclubs}}

	carrotgingersoup := menu.Item{Title: "Carrot Ginger Soup", Description: StrPtr("Carrot ginger soup with coconut milk, apple cider vinegar, and maple syrup"), Price: 895, Active: true}
	classicchickensoup := menu.Item{Title: "Homemade Chicken Soup", Description: StrPtr("Chicken soup with Israeli couscous."), Price: 695, Active: true}

	soups := menu.Section{Title: "Soups", ListOrder: 6, Items: []menu.Item{carrotgingersoup, classicchickensoup}}

	lunch := &menu.Section{Title: "Lunch", ListOrder: 5, SubSections: []menu.Section{soups, sandwiches}}
	db.Create(&lunch)

	entrees := menu.Section{Title: "Entrées", ListOrder: 10}
	desserts := menu.Section{Title: "Desserts", ListOrder: 11}

	dinner := &menu.Section{Title: "Dinner", ListOrder: 8, SubSections: []menu.Section{starters, entrees, desserts}}
	db.Create(&dinner)
	logger.Info.Println("Successful")

	// Now create admin user/account
	logger.Info.Println("Creating default admin user/account...")
	adminUser := user.User{
		FirstName: "Auguste",
		LastName:  "Gusteau",
		Address1:  "31 Rue Cambon",
		Address2:  "",
		ZipCode:   75001,
		Email:     "admin@Gusteaus.com",
	}
	db.Create(&adminUser)
	password := "administrator"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	adminAccount := &account.Account{
		Username: "admin",
		Password: string(hashedPassword),
		Role:     account.Admin,
	}
	db.Create(&adminAccount)

	logger.Info.Println("Successful")
	logger.Info.Println("Creating employee account/user...")
	employeeUser := user.User{
		FirstName: "Remy",
		LastName:  "Ratatouille",
		Address1:  "10 Rue Egout",
		Address2:  "",
		ZipCode:   75002,
		Email:     "remy@pixar.com",
	}
	db.Create(&employeeUser)

	password = "employee"
	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	employeeAccount := &account.Account{
		Username: "employee",
		Password: string(hashedPassword),
		Role:     account.Employee,


	}
	db.Create(&employeeAccount)


	logger.Info.Println("Successful")
	logger.Info.Println("Creating guest account/user...")
	guestUser := user.User{
		FirstName: "Anton",
		LastName:  "Ego",
		Address1:  "99 Tour D'Ivoire",
		Address2:  "",
		ZipCode:   75003,
		Email:     "anton@divoire.com",
	}

	db.Create(&guestUser)

	password = "guest"
	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	guestAccount := &account.Account{
		Username: "guest",
		Password: string(hashedPassword),
		Role:     account.Guest,
	}

	db.Create(&guestAccount)

	guestAccount.User = guestUser
	adminAccount.User = adminUser
	employeeAccount.User = employeeUser

	if err := db.Session(&gorm.Session{FullSaveAssociations: true, AllowGlobalUpdate: true}).Save(&guestAccount).Error; err != nil {
		return err
	}
	if err := db.Session(&gorm.Session{FullSaveAssociations: true,
		AllowGlobalUpdate: true}).Save(&adminAccount).Error; err != nil {
		return err
	}
	if err := db.Session(&gorm.Session{FullSaveAssociations: true,
		AllowGlobalUpdate: true}).Save(&employeeAccount).Error; err != nil {
		return err
	}
	//
	// // if db.Session(&gorm.Session{FullSaveAssociations: true, AllowGlobalUpdate: true}).Model(&guestAccount).Association("User").Append(&guestUser); err != nil {
	// // 	return err
	// // }
	// if db.Session(&gorm.Session{FullSaveAssociations: true, AllowGlobalUpdate: true}).Model(&employeeAccount).Association("User").Append(&employeeUser); err != nil {
	// 	return err
	// }
	//
	//
	// if db.Session(&gorm.Session{FullSaveAssociations: true, AllowGlobalUpdate: true}).Model(&adminAccount).Association("User").Append(&adminUser); err != nil {
	// 	return err
	// }

	logger.Info.Println("Successful")

	return nil
}

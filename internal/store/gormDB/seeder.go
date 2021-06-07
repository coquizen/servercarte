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

//goland:noinspection SpellCheckingInspection
func seedDB(db *gorm.DB) error {
	// Drop all Tables
	var migrator = db.Migrator()
	if err := migrator.DropTable(&menu.Section{}, &menu.Item{}, &user.User{}, &account.Account{}); err != nil {
		return err
	}

	logger.Info.Println("All tables have been dropped...")
	logger.Info.Println("Now migrating...")
	// Migrate model over to db
	err := db.AutoMigrate(&menu.Section{}, &menu.Item{}, &user.User{}, &account.Account{})
	if err != nil {
		return fmt.Errorf("error migrating scheme to db: %v", err)
	}
	logger.Info.Println("Successful")

	logger.Info.Println("Seeding...")

	everything := &menu.Item{Title: "Everything", Description: StrPtr("With onions, sesame seeds, & poppy seeds"), ListOrder : 1, Type: menu.AddOn, Price: 0, Active: true}
	sesame := &menu.Item{Title: "Sesame Seed", ListOrder : 2, Type: menu.AddOn, Price: 0, Active: true}
	poppy := &menu.Item{Title: "Poppy Seed", ListOrder : 3, Type: menu.AddOn, Price: 0, Active: true}
	plain := &menu.Item{Title: "Plain", ListOrder : 4, Type: menu.AddOn, Price: 0, Active: true}
	chocolate := &menu.Item{Title: "Chocolate Chip", Description: StrPtr("With Godiva chocolate chips"), ListOrder: 5, Type: menu.AddOn, Price: 0, Active: true}
	onion := &menu.Item{Title: "Onion", ListOrder : 6, Type: menu.AddOn, Price: 0, Active: true}

	bagelcontainer := &menu.Section{Title: "Bagel Container", ListOrder: 0, Type: menu.Container, Items: []menu.Item{*everything, *chocolate, *sesame, *poppy, *plain, *onion}}

	plaincreamcheese := &menu.Item{Title: "Plain Cream Cheese", ListOrder: 1, Price: 50, Type: menu.AddOn, Active: true}
	scallioncreamcheese := &menu.Item{Title: "Scallion Cream Cheese", ListOrder: 2, Price: 100, Type: menu.AddOn, Active: true}
	lightplaincreamcheese := &menu.Item{Title: "Light Cream Cheese", ListOrder: 3, Price: 25, Type: menu.AddOn, Active: true}
	garliccreamcheese := &menu.Item{Title: "Garlic Cream Cheese", ListOrder: 4, Price: 75, Type: menu.AddOn, Active: true}
	onionandchivecreamcheese := &menu.Item{Title: "Onion & Chive Cream Cheese", ListOrder: 5, Price: 150, Type: menu.AddOn, Active: true}
	smokedsalmoncreamcheese := &menu.Item{Title: "Smoked Salmon Cream Cheese", ListOrder: 6, Price: 50, Type: menu.AddOn, Active: true}

	bagelcondimentcontainer := &menu.Section{Title: "Bagel Condiments",
		ListOrder: 0, Type: menu.Container, Active: true, Visible: false,
		Items: []menu.Item{*plaincreamcheese, *scallioncreamcheese, *lightplaincreamcheese, *onionandchivecreamcheese,*garliccreamcheese, *smokedsalmoncreamcheese }}

	bagel := &menu.Item{Title: "Bagel", Description: StrPtr("Your choice of bagel."), Type: menu.Plate, ListOrder: 1, Price: 395, Active: true, AddOns: *bagelcontainer}
	bagelwcreamcheese := &menu.Item{Title: "Bagel w/ Cream Cheese", Description: StrPtr("Toasted H&H Bagel with your choice of cream cheese."), Type: menu.Plate, ListOrder: 2, Price: 595, Active: true, AddOns: *bagelcontainer, Condiments: *bagelcondimentcontainer}
	bagelwlox := &menu.Item{Title: "Bagel with Lox", Description: StrPtr("Your choice of H&H bagels and Atlantic smoked lox."), Type: menu.Plate, ListOrder: 3, Price: 995, Active: true, AddOns: *bagelcontainer, Condiments: *bagelcondimentcontainer}

	bagels := menu.Section{Title: "Bagels", ListOrder: 1, Type: menu.Category, Items: []menu.Item{*bagelwcreamcheese, *bagelwlox, *bagel}}

	maplesyrup := &menu.Item{Title: "Canadian Maple Syrup", ListOrder: 1, Type: menu.Condiment, Price: 0, Active: true}
	butter := &menu.Item{Title: "Butter", ListOrder: 2, Type:menu.Condiment, Price: 0, Active: true}

	wafflecondiments := &menu.Section{Title: "Waffle Condiments", ListOrder: 0, Type: menu.Container, Visible: false, Items: []menu.Item{*maplesyrup, *butter}}

	waffle := &menu.Item{Title: "Waffles", Description: StrPtr("5 slices of thick homemade waffles with light cream."), Type: menu.Plate, ListOrder: 1, Price: 775, Active: true, Condiments: *wafflecondiments	}
	montecristowaffle := &menu.Item{Title: "Monte Cristo Waffle Sandwich", Description: StrPtr("Whole wheat Waffle sandwich with swiss cheese, raspberry jam, honey baked, and ham."), Type: menu.Plate, ListOrder: 2, Price: 1125, Active: true, Condiments: *wafflecondiments}
	spicywaffle := &menu.Item{Title: "Southern Waffle Sandwich", Description: StrPtr("Belgian waffles with spicy maple syrup, cheddar cheese, butter and cinnamon."), Type: menu.Plate, ListOrder: 3, Price: 1035, Active: true, Condiments: *wafflecondiments}

	waffles := menu.Section{Title: "Waffles", ListOrder: 2, Type: menu.Category, Items: []menu.Item{*waffle, *montecristowaffle, *spicywaffle}}

	eggswbacon := &menu.Item{Title: "Eggs (Any style)", Description: StrPtr("Fresh farm eggs cooked your way with thick Applewood smoked bacon."), ListOrder: 1, Type: menu.Plate, Price: 725, Active: true}
	eggsbenedict := &menu.Item{Title: "Eggs Benedict", Description: StrPtr("Eggs Benedict with homemade Hollandaise sauce."), ListOrder: 2, Type: menu.Plate, Price: 775, Active: true}
	countryomelettes := &menu.Item{Title: "Country Omelettes", Description: StrPtr("Omelette stuffed with chorizo, green peppers, onion, and Manchego cheese."), ListOrder: 3, Type: menu.Plate, Price: 1295, Active: true}
	classicomelettes := &menu.Item{Title: "Classic Omelettes", Description: StrPtr("Omelette with select herbs and freshly ground black pepper."), ListOrder: 4, Type: menu.Plate, Price: 1095, Active: true}
	westernomelette := &menu.Item{Title: "Western Omelettes", Description: StrPtr("Omelette with green bell pepper, red bell pepper, and Monterrey jack."), ListOrder: 5,  Type: menu.Plate, Price: 995, Active: true}

	eggs := menu.Section{Title: "Eggs & Omelettes", ListOrder: 3, Type: menu.Category, Items: []menu.Item{*eggsbenedict, *eggswbacon, *countryomelettes, *classicomelettes, *westernomelette}}
	//_ = db.Model(&eggs).Association("Items").Append([]menu.Item{*eggsbenedict, *eggswbacon, *countryomelettes, *classicomelettes, *westernomelette})

	turkeybacon := &menu.Item{Title: "Turkey bacon", ListOrder: 3, Type: menu.Side, Active: true, Price: 200}
	chickensausage := &menu.Item{Title: "Chicken Sausage", Type: menu.Side, ListOrder: 2, Active: true, Price: 200}
	potatolatke := &menu.Item{Title: "Potatoe Latke", Type: menu.Side, ListOrder: 1, Active: true, Price: 200}
	breakfastsidecontainer := &menu.Section{Title: "Sides", Items: []menu.Item{*turkeybacon, *chickensausage, *potatolatke}, Type: menu.Container, Active: true, Visible: true, ListOrder: 0}
	breakfast := &menu.Section{Title: "Breakfast", ListOrder: 1, Type: menu.Meal, SubSections: []menu.Section{bagels, eggs, waffles, *breakfastsidecontainer}, Active: true, Visible: true}
	db.Create(&breakfast)

	sunomonosalad := menu.Item{Title: "Sunomono Salad", Description: StrPtr("Thin rice noodles, shrimp, crab, soy sauce and rice vinegar."), Type: menu.Plate, Price: 395, ListOrder: 1, Active: true}
	cobbsalad := menu.Item{Title: "Cobb Salad", Description: StrPtr("Blue cheese, grilled chicken breasts, red wine vinegar, eggs, and bacon."), Type: menu.Plate, Price: 645, ListOrder: 2, Active: true}

	salads := menu.Section{Title: "Salads", Type: menu.Category, Items: []menu.Item{sunomonosalad, cobbsalad}, Active: true, Visible: true, ListOrder: 1}

	classictunamelt := menu.Item{Title: "Classic Tuna Melt", Description: StrPtr("Sourdough bread, fresh tuna, red onions, dill pickles, celery and butter."), Type: menu.Plate, ListOrder: 1,Price: 895, Active: true}
	hamdandcheesesandwich := menu.Item{Title: "Perfect Ham and Cheese Sandwich", Description: StrPtr("Sourdough bread, swiss cheese, ham, honey, mustard, mayonnaise, and pickle."), ListOrder: 2, Type: menu.Plate, Price: 995, Active: true}
	lemonchickenwraps := menu.Item{Title: "Lemon Chicken Wrap", Description: StrPtr("Pita bread, grilled chicken breast, greek yogurt, garlic, Sriracha sauce, paprika."), Type: menu.Plate, ListOrder: 3, Price: 1095, Active: true}
	hasselbacktomatoclubs := menu.Item{Title: "Hassel Back Tomato Club", Description: StrPtr("Bibb lettuce leaves, ripe avocados, swiss cheese, plum tomatoes, and turkey"), Type: menu.Plate, ListOrder: 4, Price: 1095, Active: true}

	sandwiches := menu.Section{Title: "Sandwiches", ListOrder: 2, Type: menu.Category, Items: []menu.Item{classictunamelt, hamdandcheesesandwich, lemonchickenwraps, hasselbacktomatoclubs}}

	carrotgingersoup := menu.Item{Title: "Carrot Ginger Soup", Description: StrPtr("Carrot ginger soup with coconut milk, apple cider vinegar, and maple syrup"), Type: menu.Plate, ListOrder: 1, Price: 895, Active: true}
	classicchickensoup := menu.Item{Title: "Homemade Chicken Soup", Description: StrPtr("Chicken soup with Israeli couscous."), Type: menu.Plate, ListOrder: 2 ,Price: 695, Active: true}

	soups := menu.Section{Title: "Soups", ListOrder: 3, Type: menu.Category, Items: []menu.Item{carrotgingersoup, classicchickensoup}}

	lunch := &menu.Section{Title: "Lunch", ListOrder: 2, Type: menu.Meal, SubSections: []menu.Section{soups, salads, sandwiches}}
	db.Create(&lunch)

	handpies := menu.Item{Title: "Korean Beef Hand Pies", Description: StrPtr("Beef short rubes, rice noodles, hoisin sauce, chili sauce, and soy sauce."), Type: menu.Plate, ListOrder: 1, Price: 795, Active: true}
	bruschetta := menu.Item{Title: "Bruschetta", Description: StrPtr("Toasted baguettes with goat cheese, brown sugar, and cherry tomatoes."), Type: menu.Plate, ListOrder: 2, Price: 495, Active: true}
	grilledpolenta := menu.Item{Title: "Grilled Polenta w/ Wild Mushrooms", Type: menu.Plate, Price: 695, Active: true}

	starters := menu.Section{Title: "Starters", ListOrder: 1, Type: menu.Category, Items: []menu.Item{handpies, bruschetta, grilledpolenta}}

	ossobuco := menu.Item{Title: "Osso Buco", Description: StrPtr("Braised veal shank served over risotto milanese"), ListOrder: 3, Type: menu.Plate, Price: 2895, Active: true}
	shortribtortelloni := menu.Item{Title: "Short Ribs Tortelloni", Description: StrPtr("Tortelloni stuffed with braised beef short ribs"), ListOrder: 1, Type: menu.Plate, Price: 2495, Active: true}
	gnocchi := menu.Item{Title: "Gnocchi Castelmagno", Description: StrPtr("Handmade Kale Potato Gnocchi in Castelmagno cream sauce"), ListOrder: 2, Type: menu.Plate, Price: 1895, Active: true}
	entrees := menu.Section{Title: "Entrées", Type: menu.Category, ListOrder: 2, Items: []menu.Item{ossobuco, shortribtortelloni, gnocchi}}

	dinner := &menu.Section{Title: "Dinner", Type: menu.Meal, ListOrder: 3, SubSections: []menu.Section{starters, entrees}, Active: true, Visible: true}
	db.Create(&dinner)

	logger.Info.Println("Successful")

	logger.Info.Println("Seeding Desserts menu")

	chocolatemousse := menu.Item{Title: "Chocolate Mousse & Whipped Cream", ListOrder: 1, Type: menu.Plate, Price: 1250, Active: true}
	tiramisu := menu.Item{Title: "Tiramisu", ListOrder: 2, Type: menu.Plate, Price: 1000, Active: true}
	cheesecake := menu.Item{Title: "Housemade Cheesecake with Fresh Strawberries", Description: StrPtr("Warm chocolate sauce"), Type: menu.Plate, ListOrder: 3, Price: 1295, Active: true}
	gastrotartufo := menu.Item{Title: "Gastro's Tartufo", Description: StrPtr("Van Leeuwen Chocolate & Vanilla Bean Ice Cream rolled in Chocolate Chunks"), Type: menu.Plate, ListOrder: 4, Price: 1225, Active: true}

	desserts := menu.Section{Title: "Desserts", Type: menu.Meal, ListOrder: 4, Items: []menu.Item{chocolatemousse,tiramisu, cheesecake, gastrotartufo}}
	db.Create(&desserts)
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
